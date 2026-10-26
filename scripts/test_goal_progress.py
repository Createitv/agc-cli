import importlib.util
import pathlib
import unittest
import json
import copy

spec = importlib.util.spec_from_file_location("progress", pathlib.Path(__file__).with_name("goal-progress.py"))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class ProgressTests(unittest.TestCase):
    def test_domain_precheck_requires_real_code_file_readback_and_cleanup(self):
        baseline = {'originalBlocked': [{'family': 'domains', 'endpoint': 'domain-api-precheck'}]}
        ledger = {'flow': 'domain-precheck-custom-host', 'status': 'verified-precheck', 'appId': 'dedicated',
                  'owned': {'workerName': 'owned-worker', 'hostname': 'owned.example.test'},
                  'hostedDomain': 'https://owned.example.test', 'cleanupVerified': True,
                  'verificationFile': {'fileName': 'dedicated.txt', 'bytes': 32, 'sha256': 'a'*64,
                                       'sourceURL': 'https://connect-api.cloud.huawei.com/api/dms/domain-manage/v1/app/domain/verify-file', 'downloadHTTPStatus': 200},
                  'publicReadbacks': [{'httpStatus': 200, 'bytes': 32, 'sha256': 'a'*64, 'matchesDownloadedFile': True}],
                  'steps': [{'family': 'domains', 'endpoint': 'domain-api-precheck', 'status': 'passed', 'exitCode': 0, 'httpStatus': 200, 'businessCode': 0}]}
        self.assertEqual(module.collect(baseline, [('domain', ledger)])['verifiedOriginalLifecycles'], 1)
        for section, key, value in [('verificationFile', 'sha256', 'b'*64), ('verificationFile', 'fileName', 'other.txt'), ('verificationFile', 'downloadHTTPStatus', 403)]:
            broken=copy.deepcopy(ledger);broken[section][key]=value
            self.assertEqual(module.collect(baseline, [('domain', broken)])['verifiedOriginalLifecycles'], 0)
        for key,value in [('cleanupVerified',False), ('status','incomplete')]:
            broken=copy.deepcopy(ledger);broken[key]=value
            self.assertEqual(module.collect(baseline, [('domain', broken)])['verifiedOriginalLifecycles'], 0)
        broken=copy.deepcopy(ledger);broken['steps'][0]['businessCode']=203489282
        self.assertEqual(module.collect(baseline, [('domain', broken)])['verifiedOriginalLifecycles'], 0)
        for key, value in [('sha256', 'b'*64), ('bytes', 31), ('httpStatus', 403), ('matchesDownloadedFile', False)]:
            broken = copy.deepcopy(ledger)
            broken['publicReadbacks'][0][key] = value
            self.assertEqual(module.collect(baseline, [('domain', broken)])['verifiedOriginalLifecycles'], 0)
        for section, key, value in [('owned', 'hostname', 'unowned.example.test'), ('owned', 'workerName', ''), ('verificationFile', 'sourceURL', 'https://unowned.example.test/file')]:
            broken = copy.deepcopy(ledger)
            broken[section][key] = value
            self.assertEqual(module.collect(baseline, [('domain', broken)])['verifiedOriginalLifecycles'], 0)
        broken = copy.deepcopy(ledger)
        broken['publicReadbacks'] = []
        self.assertEqual(module.collect(baseline, [('domain', broken)])['verifiedOriginalLifecycles'], 0)

    def test_multipart_counts_only_parts_with_independent_receipts(self):
        baseline = {"originalBlocked": [{"family": "upload", "endpoint": e} for e in ("obbfile-init", "obbfile-uploadurl", "obbfile-upload", "obbfile-compose")]}
        # Offline assertions use a dedicated fixture, not mutable live evidence.
        ledger = {
            'flow': 'multipart-upload', 'status': 'verified-part-upload-and-compose',
            'owned': {'objectId': 'offline-object', 'nspUploadId': 'offline-session'},
            'sourceArtifact': {'bytes': 12, 'sha256': 'a'*64},
            'steps': [
                {'endpoint': 'obbfile-init', 'status': 'passed', 'httpStatus': 200, 'businessCode': 0},
                {'endpoint': 'obbfile-uploadurl', 'status': 'passed', 'httpStatus': 200, 'businessCode': 0},
                {'endpoint': 'obbfile-upload', 'status': 'passed', 'partCount': 2}],
            'partReceipts': {
                'part1': {'httpStatus': 200, 'bytes': 6, 'sha256': 'b'*64, 'expectedMD5': 'c'*32, 'etag': 'c'*32},
                'part2': {'httpStatus': 200, 'bytes': 6, 'sha256': 'd'*64, 'expectedMD5': 'e'*32, 'etag': 'e'*32}}}
        self.assertEqual(module.collect(baseline, [('multipart', ledger)])['verifiedOriginalLifecycles'], 3)
        broken = copy.deepcopy(ledger)
        next(iter(broken['partReceipts'].values()))['expectedMD5'] = '0'*32
        self.assertEqual(module.collect(baseline, [('multipart', broken)])['verifiedOriginalLifecycles'], 0)

    def test_compile_query_requires_owned_package_and_independent_state(self):
        baseline = {"originalBlocked": [{"family": "publishing", "endpoint": "query-aabfile"}]}
        source = {"runId": "source", "appId": "dedicated", "owned": {"pkgId": "mine"}, "steps": [{"family": "publishing", "endpoint": "package-compile-status-v3", "status": "passed", "packageId": "mine", "parseState": 2}]}
        read = {"flow": "legacy-compile-readback", "status": "verified-read", "exitCode": 0, "httpStatus": 200, "businessCode": 0, "appId": "dedicated", "packageId": "mine", "sourcePackageEvidence": "source", "expectedParseState": 2, "packageStates": [{"pkgId": "mine", "successStatus": 2}]}
        report = module.collect(baseline, [("source-file", source), ("read-file", read)])
        self.assertEqual(report["verifiedOriginalLifecycles"], 1)
        read["packageStates"][0]["pkgId"] = "unowned"
        self.assertEqual(module.collect(baseline, [("source-file", source), ("read-file", read)])["verifiedOriginalLifecycles"], 0)

    def test_retained_group_requires_ownership_and_subsequent_readback(self):
        baseline = {"originalBlocked": [{"family": "pms", "endpoint": "addproductgroup-harmonyosnext"}]}
        ledger = {"flow": "pms-group-create", "status": "verified-retained", "owned": {"groupId": "mine"}, "steps": [
            {"endpoint": "addproductgroup-harmonyosnext", "status": "passed", "businessCode": 0},
            {"endpoint": "getproductgroup-harmonyosnext", "status": "passed", "assertion": "dedicated empty subscription group read back"}]}
        report = module.collect(baseline, [("group", ledger)])
        self.assertEqual(report["verifiedOriginalLifecycles"], 1)
        self.assertFalse(report["completed"][0]["cleanupVerified"])
        ledger["steps"].reverse()
        self.assertEqual(module.collect(baseline, [("group", ledger)])["verifiedOriginalLifecycles"], 0)

    def test_upload_needs_address_and_content_assertion_and_reports_retention(self):
        baseline = {"originalBlocked": [{"family": "upload", "endpoint": "upload-file-new"}]}
        ledger = {"flow": "signed-file-upload", "status": "verified-upload", "objectId": "test-object", "sha256": "fixture", "contentLength": 68, "retention": "retained", "steps": [
            {"endpoint": "upload-url-new", "status": "passed", "businessCode": 0},
            {"endpoint": "upload-file-new", "status": "passed", "httpStatus": 200, "etag": "md5", "assertion": "server ETag matches uploaded fixture MD5"}]}
        report = module.collect(baseline, [("upload", ledger)])
        self.assertEqual(report["verifiedOriginalLifecycles"], 1)
        self.assertFalse(report["completed"][0]["cleanupVerified"])
        ledger["steps"][1]["assertion"] = None
        self.assertEqual(module.collect(baseline, [("upload", ledger)])["verifiedOriginalLifecycles"], 0)

    def test_helpers_duplicates_and_self_reported_success_do_not_inflate_progress(self):
        baseline = {"originalBlocked": [{"family": "publishing", "endpoint": "language-info-delete"}]}
        ledger = {"cleanupVerified": True, "verifiedFeatures": {"fake": {}}, "steps": [
            {"family": "publishing", "endpoint": "language-info-delete", "status": "passed"},
            {"family": "publishing", "endpoint": "app-info-query-v3", "status": "passed", "assertion": "temporary language absent after deletion"}]}
        report = module.collect(baseline, [("a", ledger), ("b", ledger)])
        self.assertEqual(report["verifiedOriginalLifecycles"], 1)
        self.assertEqual(report["remainingOriginal"], 0)
        ledger["cleanupVerified"] = False
        self.assertEqual(module.collect(baseline, [("c", ledger)])["verifiedOriginalLifecycles"], 0)


if __name__ == "__main__":
    unittest.main()
