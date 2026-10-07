import importlib.util
import pathlib
import unittest

spec = importlib.util.spec_from_file_location("live_lifecycle", pathlib.Path(__file__).with_name("live-lifecycle.py"))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class LifecycleTests(unittest.TestCase):
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
