import importlib.util
import pathlib
import unittest
import tempfile
import hashlib
import io
import zipfile
import json
from unittest.mock import patch
from subprocess import CompletedProcess

spec = importlib.util.spec_from_file_location("package_flow", pathlib.Path(__file__).with_name("live-package.py"))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class PackageStateTests(unittest.TestCase):
    def test_inner_hap_identity_cannot_hide_behind_outer_pack_info(self):
        with tempfile.TemporaryDirectory() as work:
            app = pathlib.Path(work) / 'fixture.app'
            for inner_bundle in ['dedicated.bundle', 'other.bundle']:
                hap = io.BytesIO()
                with zipfile.ZipFile(hap, 'w') as z:
                    z.writestr('module.json', json.dumps({'app': {'bundleName': inner_bundle, 'debug': False}, 'module': {'name': 'entry'}}))
                with zipfile.ZipFile(app, 'w') as z:
                    z.writestr('pack.info', json.dumps({'summary': {'app': {'bundleName': 'dedicated.bundle'}}}))
                    z.writestr('entry.hap', hap.getvalue())
                if inner_bundle == 'dedicated.bundle':
                    report = module.inspect_app_modules(app, 'dedicated.bundle')
                    self.assertEqual(report[0]['debug'], False)
                else:
                    with self.assertRaises(RuntimeError): module.inspect_app_modules(app, 'dedicated.bundle')

    def test_retention_copies_only_app_and_preserves_exact_bytes(self):
        with tempfile.TemporaryDirectory() as work:
            directory = pathlib.Path(work)
            source = directory / "signed.app"
            source.write_bytes(b"signed-fixture")
            (directory / "private.p12").write_bytes(b"private-key")
            report = module.retain_signed_artifact(source, directory / "artifacts")
            self.assertEqual(report["sha256"], hashlib.sha256(b"signed-fixture").hexdigest())
            self.assertEqual(list((directory / "artifacts").iterdir()), [directory / "artifacts/signed.app"])
            self.assertEqual((directory / "artifacts/signed.app").stat().st_mode & 0o777, 0o600)
            with self.assertRaises(RuntimeError):
                module.retain_signed_artifact(directory / "private.p12", directory / "artifacts")
            with self.assertRaises(FileExistsError):
                module.retain_signed_artifact(source, directory / "artifacts")

    def test_sdk_zero_exit_with_error_is_not_success(self):
        # 实测 SDK 统计工具在参数错误时退出 0，故必须检查错误标记。
        output = 'Ohos BundleTool [Error]: code:9132105 error:private-keystore-password'
        with patch.object(module.subprocess, 'run', return_value=CompletedProcess([], 0, '', output)):
            with self.assertRaises(RuntimeError) as caught:
                module.run_sdk([], 'Package verification')
        self.assertIn('9132105', str(caught.exception))
        self.assertNotIn('private-keystore-password', str(caught.exception))

    def test_sdk_nonzero_error_is_sanitized(self):
        with patch.object(module.subprocess, 'run', return_value=CompletedProcess([], 1, 'ERROR: 11011005 secret', '')):
            with self.assertRaisesRegex(RuntimeError, 'SDK code 11011005') as caught:
                module.run_sdk([], 'Signing')
        self.assertNotIn('secret', str(caught.exception))

    def test_sdk_success_message_is_accepted(self):
        with patch.object(module.subprocess, 'run', return_value=CompletedProcess([], 0, 'Verification succeeded', '')):
            module.run_sdk([], 'Signing')

    def test_state_is_for_exact_owned_package(self):
        self.assertEqual(module.package_state({"pkgStateList": [{"pkgId": "mine", "successStatus": 1}]}, "mine"), 1)
        with self.assertRaises(RuntimeError):
            module.package_state({"pkgStateList": [{"pkgId": "other", "successStatus": 0}]}, "mine")

    def test_boolean_and_duplicate_states_are_not_success(self):
        for states in [[{"pkgId": "mine", "successStatus": False}], [{"pkgId": "mine", "successStatus": 0}, {"pkgId": "mine", "successStatus": 2}]]:
            with self.assertRaises(RuntimeError):
                module.package_state({"pkgStateList": states}, "mine")


if __name__ == "__main__":
    unittest.main()
