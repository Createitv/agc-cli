import importlib.util
import pathlib
import unittest
import tempfile
import types
from unittest.mock import patch

spec = importlib.util.spec_from_file_location("live_lifecycle", pathlib.Path(__file__).with_name("live-lifecycle.py"))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class LifecycleTests(unittest.TestCase):
    def test_json_string_ret_preserves_code_without_embedded_private_data(self):
        evidence = module.safe_evidence({"ret": '{"code":123,"msg":"PRIVATE account"}'})
        self.assertEqual(evidence["businessCode"], 123)
        self.assertNotIn("PRIVATE", str(evidence))
        self.assertIsNone(module.safe_evidence({"ret": "invalid JSON"})["businessCode"])

    def test_failure_records_code_origin_without_raw_response(self):
        with tempfile.TemporaryDirectory() as tmp:
            args = types.SimpleNamespace(binary=pathlib.Path("offline-cli"), profile="default", app_id="dedicated", output_dir=pathlib.Path(tmp))
            flow = module.Flow(args, "error-proof", "offline")
            response = types.SimpleNamespace(returncode=1, stderr="agc endpoint returned HTTP 200 (2): denied", stdout='{"error":{"code":"2","codeField":"businessCode","resultCodes":{"rtnCode":"0","businessCode":"2"},"rawBody":{"invitationCode":"SECRET"}}}')
            with patch.object(module.subprocess, "run", return_value=response):
                with self.assertRaises(RuntimeError): flow.call("testing", "test-api-create-open-test-link")
            step = flow.ledger["steps"][0]
            self.assertEqual(step["codeField"], "businessCode")
            self.assertEqual(step["resultCodes"], {"rtnCode": "0", "businessCode": "2"})
            self.assertNotIn("SECRET", flow.output.read_text())

    def test_json_only_failure_preserves_http_status_without_private_message(self):
        with tempfile.TemporaryDirectory() as tmp:
            args = types.SimpleNamespace(binary=pathlib.Path("offline-cli"), profile="default", app_id="dedicated", output_dir=pathlib.Path(tmp))
            flow = module.Flow(args, "json-error-proof", "offline")
            response = types.SimpleNamespace(returncode=1, stderr="", stdout='{"error":{"code":"991","message":"agc endpoint returned HTTP 408 (991): PRIVATE", "codeField":"rtnCode","resultCodes":{"rtnCode":"991"}}}')
            with patch.object(module.subprocess, "run", return_value=response):
                with self.assertRaises(RuntimeError): flow.call("domains", "domain-api-precheck")
            step = flow.ledger["steps"][0]
            self.assertEqual(step.get("httpStatus"), 408)
            self.assertEqual(step["status"], "http-error")
            self.assertEqual(step["businessCode"], "991")
            self.assertNotIn("PRIVATE", flow.output.read_text())

    def test_step_can_select_and_record_separate_compile_profile(self):
        with tempfile.TemporaryDirectory() as tmp:
            args = types.SimpleNamespace(binary=pathlib.Path("offline-cli"), profile="default", app_id="dedicated", output_dir=pathlib.Path(tmp))
            flow = module.Flow(args, "profile-proof", "offline")
            response = types.SimpleNamespace(returncode=0, stderr="", stdout='{"data":{"dryRun":false,"statusCode":200,"body":{"ret":{"code":0}}}}')
            with patch.object(module.subprocess, "run", return_value=response) as run:
                flow.call("publishing", "query-aabfile", query={"appId": "dedicated", "pkgIds": "owned-package"}, profile="api-client")
            self.assertEqual(run.call_args.args[0][2], "api-client")
            self.assertEqual(flow.ledger["steps"][0]["profile"], "api-client")
            self.assertEqual(args.profile, "default")

    def test_pms_success_uses_nested_error_result(self):
        self.assertEqual(module.safe_evidence({"error": {"errorCode": 0}, "result": "owned-group-id"})["businessCode"], 0)

    def test_cleanup_refuses_resources_not_created_by_run(self):
        ledger = {"owned": {"groupId": "mine"}}
        self.assertTrue(module.is_owned(ledger, "groupId", "mine"))
        self.assertFalse(module.is_owned(ledger, "groupId", "existing-business-group"))
        self.assertFalse(module.is_owned(ledger, "groupId", None))

    def test_evidence_does_not_persist_invitation_secrets_or_accounts(self):
        data = {"rtnCode": 0, "groupId": "g", "invitationCode": "SECRET", "testerInfo": [{"hwAccount": "PRIVATE"}]}
        evidence = module.safe_evidence(data)
        self.assertEqual(evidence["businessCode"], 0)
        self.assertEqual(evidence["groupId"], "g")
        self.assertNotIn("SECRET", str(evidence))
        self.assertNotIn("PRIVATE", str(evidence))

    def test_rename_readback_checks_both_identity_and_name(self):
        groups = [{"groupId": "mine", "groupName": "renamed"}, {"groupId": "other", "groupName": "desired"}]
        self.assertTrue(module.group_matches(groups, "mine", "renamed"))
        self.assertFalse(module.group_matches(groups, "mine", "desired"))

    def test_business_subcode_is_not_hidden_by_success_code(self):
        self.assertEqual(module.safe_evidence({"rtnCode": 0, "businessCode": 100})["businessCode"], 100)

    def test_provisioning_requires_dedicated_app_identity(self):
        fixture = {"appId": "dedicated", "readbackVerified": True}
        self.assertTrue(module.is_dedicated_app("dedicated", fixture))
        self.assertFalse(module.is_dedicated_app("business", fixture))

    def test_mutation_requires_successful_readback_and_cleanup(self):
        ledger = {"cleanupVerified": True, "steps": [
            {"family": "testing", "endpoint": "test-api-edit-test-group", "status": "passed"},
            {"family": "testing", "endpoint": "test-api-get-test-grouplist", "status": "passed", "assertion": "same group ID has updated name"},
            {"family": "testing", "endpoint": "test-api-stop-invite-code", "status": "business-error"},
        ]}
        features = module.verified_features(ledger)
        self.assertIn("testing/test-api-edit-test-group", features)
        self.assertNotIn("testing/test-api-stop-invite-code", features)
        ledger["cleanupVerified"] = False
        self.assertEqual(module.verified_features(ledger), {})


if __name__ == "__main__":
    unittest.main()
