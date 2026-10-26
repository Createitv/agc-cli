"""Continue an owned initialized upload; two part uploads run concurrently.
Only the dedicated APP artifact is used. Signed URLs and headers stay in memory.
No package registration, submission or publication is performed.
"""
import argparse
import concurrent.futures
import datetime
import hashlib
import importlib.util
import json
import pathlib
import subprocess
import tempfile
import zipfile

ROOT = pathlib.Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location('lifecycle', pathlib.Path(__file__).with_name('live-lifecycle.py'))
lifecycle = importlib.util.module_from_spec(spec)
spec.loader.exec_module(lifecycle)


def validate_addresses(addresses, expected_keys):
    # 合并必须使用服务器返回的分片对象 ID；全部地址先核对，再并发上传。
    if not isinstance(addresses, dict) or set(addresses) != set(expected_keys):
        raise RuntimeError('Signed addresses do not cover exact owned APP parts')
    for address in addresses.values():
        if (not isinstance(address, dict)
                or address.get('method') != 'PUT'
                or not isinstance(address.get('url'), str)
                or not address['url'].startswith('https://')
                or not isinstance(address.get('partObjectId'), str)
                or not address['partObjectId'].strip()
                or len(address['partObjectId']) > 64
                or not isinstance(address.get('headers', {}), dict)):
            raise RuntimeError('Incomplete signed part address; response withheld')


def upload_part(binary, profile, address, path):
    if address.get('method') != 'PUT' or not address.get('url', '').startswith('https://'):
        raise RuntimeError('Invalid signed part upload method/address')
    command = [str(binary), '--profile', profile, 'upload', 'obbfile-upload', '--invoke', '--dry-run=false', '--param', 'uploadUrl='+address['url'], '--file', str(path)]
    for key, value in address.get('headers', {}).items():
        command += ['--header', key+'='+str(value)]
    result = subprocess.run(command, cwd=ROOT, capture_output=True, text=True, timeout=90)
    if result.returncode:
        raise RuntimeError('Part CLI upload failed; sensitive response withheld')
    data = json.loads(result.stdout)['data']
    content = path.read_bytes()
    expected = hashlib.md5(content).hexdigest()
    etag = data.get('etag', '').strip('"')
    if data.get('dryRun') is not False or data.get('statusCode') != 200 or etag.lower() != expected:
        raise RuntimeError('Part response does not verify exact uploaded bytes')
    return {'httpStatus': 200, 'etag': etag, 'expectedMD5': expected, 'bytes': len(content), 'sha256': hashlib.sha256(content).hexdigest()}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--ledger', required=True, type=pathlib.Path)
    parser.add_argument('--binary', required=True, type=pathlib.Path)
    args = parser.parse_args()
    ledger = json.loads(args.ledger.read_text())
    fixture = json.loads((ROOT/'docs/verification/dedicated-test-app.json').read_text())
    if ledger.get('flow') != 'multipart-upload' or ledger.get('status') != 'initialized' or not lifecycle.is_dedicated_app(ledger.get('appId'), fixture):
        raise SystemExit('Owned initialized dedicated multipart ledger required')
    for key in ('objectId', 'nspUploadId'):
        if not ledger.get('owned', {}).get(key): raise SystemExit('Missing owned multipart session')
    builds = [json.loads(p.read_text()) for p in args.ledger.parent.glob('*-dedicated-app-package-build.json')]
    builds = [b for b in builds if b.get('appId') == fixture['appId'] and b.get('status') == 'verified-unsigned-app']
    artifact = max(builds, key=lambda b: b['checkedAt'])['artifacts'][0]
    source = pathlib.Path(artifact['path'])
    content = source.read_bytes()
    if hashlib.sha256(content).hexdigest() != artifact['sha256']:
        raise SystemExit('Dedicated APP changed since verified build')
    with zipfile.ZipFile(source) as archive:
        if json.loads(archive.read('pack.info'))['summary']['app']['bundleName'] != fixture['packageName']:
            raise SystemExit('Dedicated APP bundle mismatch')
    chunk_size = ledger['partBounds']['nspPartMinSize']
    if type(chunk_size) is not int or not 0 < chunk_size <= 100*1024*1024 or len(content) <= chunk_size:
        raise SystemExit('Require real multi-part APP within documented part bounds')
    # Resume the same session: never create a replacement on observation timeout.
    flow = object.__new__(lifecycle.Flow)
    flow.args = argparse.Namespace(binary=args.binary.resolve(), profile=ledger['profile'], app_id=ledger['appId'])
    flow.output, flow.ledger = args.ledger, ledger
    ledger['sourceArtifact'] = {'sha256': artifact['sha256'], 'bytes': len(content), 'bundleName': fixture['packageName']}
    query = {key: ledger['owned'][key] for key in ('objectId', 'nspUploadId')}
    try:
        with tempfile.TemporaryDirectory(prefix='agc-multipart-parts-') as work:
            parts = {}
            for i, start in enumerate(range(0, len(content), chunk_size), 1):
                key = 'additionalProp'+str(i)
                path = pathlib.Path(work)/('part-'+str(i))
                path.write_bytes(content[start:start+chunk_size]); parts[key] = path
            request = {k: {'length': p.stat().st_size, 'sha256': hashlib.sha256(p.read_bytes()).hexdigest()} for k, p in parts.items()}
            addresses = flow.call('upload', 'obbfile-uploadurl', query=query, body=request,
                assert_fn=lambda b: set(b.get('uploadInfoMap', {})) == set(parts), assertion='signed addresses cover exact owned APP parts')['uploadInfoMap']
            validate_addresses(addresses, parts)
            with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
                pending = {key: pool.submit(upload_part, args.binary.resolve(), ledger['profile'], addresses[key], path) for key, path in parts.items()}
                receipts = {key: task.result() for key, task in pending.items()}
            ledger['partReceipts'] = receipts
            ledger['steps'].append({'family': 'upload', 'endpoint': 'obbfile-upload', 'profile': ledger['profile'], 'status': 'passed', 'checkedAt': datetime.datetime.now(datetime.timezone.utc).isoformat(), 'assertion': 'every owned part HTTP 200 and ETag matches independent local MD5', 'partCount': len(receipts)})
            flow.save()
            compose = {key: {'partObjectId': addresses[key]['partObjectId'], 'etag': receipt['etag']} for key, receipt in receipts.items()}
            flow.call('upload', 'obbfile-compose', query=query, body=compose, assertion='owned session merged using exact uploaded part receipts')
            ledger['status'] = 'verified-part-upload-and-compose'
            ledger['retention'] = 'Dedicated remote object retained, not linked to app versions; no delete API available'
            ledger['limitations'] = ['Final merged bytes not independently downloaded; compose business success does not prove final byte equality.']
    except Exception as error:
        ledger['status'] = 'incomplete'
        ledger['failure'] = str(error) if isinstance(error, RuntimeError) else type(error).__name__
    finally:
        flow.save()
        print(ledger['status'], str(flow.output))


if __name__ == '__main__': main()
