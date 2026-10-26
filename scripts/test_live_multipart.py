import importlib.util
import pathlib
import json
import hashlib
import subprocess
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('multipart', pathlib.Path(__file__).with_name('live-multipart.py'))
m = importlib.util.module_from_spec(spec)
spec.loader.exec_module(m)


class UploadTests(unittest.TestCase):
    def invoke(self, data, exit_code=0):
        with tempfile.TemporaryDirectory() as work:
            path = pathlib.Path(work)/'part'
            path.write_bytes(b'owned-part-content')
            result = subprocess.CompletedProcess([], exit_code, json.dumps({'data': data}), 'signed-secret')
            with patch.object(m.subprocess, 'run', return_value=result):
                return m.upload_part(pathlib.Path('/fixture/agc'), 'explicit-profile', {'method': 'PUT', 'url': 'https://example.com/signed-secret', 'headers': {'Authorization': 'signed-secret'}}, path)

    def test_incomplete_address_map_rejected_before_any_part_upload(self):
        valid = {'method': 'PUT', 'url': 'https://example.com/signed-secret', 'partObjectId': 'owned-part', 'headers': {}}
        m.validate_addresses({'additionalProp1': valid}, {'additionalProp1'})
        for invalid in [dict(valid, partObjectId=''), dict(valid, method='GET'), dict(valid, url='http://example.com'), dict(valid, headers=[]), None]:
            with self.assertRaises(RuntimeError) as caught:
                m.validate_addresses({'additionalProp1': invalid}, {'additionalProp1'})
            self.assertNotIn('signed-secret', str(caught.exception))
        with self.assertRaises(RuntimeError):
            m.validate_addresses({'additionalProp2': valid}, {'additionalProp1'})

    def test_exact_part_bytes(self):
        receipt = self.invoke({'dryRun': False, 'statusCode': 200, 'etag': '"'+hashlib.md5(b'owned-part-content').hexdigest()+'"'})
        self.assertEqual(receipt['sha256'], hashlib.sha256(b'owned-part-content').hexdigest())

    def test_false_success_rejected(self):
        for data in [{'dryRun': True, 'statusCode': 200, 'etag': 'wrong'}, {'dryRun': False, 'statusCode': 200, 'etag': 'wrong'}, {'dryRun': False, 'statusCode': 403, 'etag': hashlib.md5(b'owned-part-content').hexdigest()}]:
            with self.assertRaises(RuntimeError): self.invoke(data)

    def test_failure_does_not_leak_signed_credentials(self):
        with self.assertRaises(RuntimeError) as caught: self.invoke({}, 1)
        self.assertNotIn('signed-secret', str(caught.exception))


if __name__ == '__main__': unittest.main()
