import pathlib
import unittest


ROOT = pathlib.Path(__file__).resolve().parents[1]
NPM_WORKFLOW = ROOT / ".github" / "workflows" / "npm.yml"


class NPMWorkflowTests(unittest.TestCase):
    def test_setup_node_does_not_use_boolean_cache_value(self):
        workflow = NPM_WORKFLOW.read_text(encoding="utf-8")

        self.assertNotIn(
            "cache: false",
            workflow,
            "actions/setup-node only accepts a package-manager name for cache; "
            "omit cache to disable it",
        )

    def test_missing_version_does_not_turn_npm_json_error_into_integrity(self):
        workflow = NPM_WORKFLOW.read_text(encoding="utf-8")

        self.assertNotIn(
            'dist.integrity --json',
            workflow,
            "npm view --json prints an error object to stdout for a missing "
            "version, which must not be mistaken for an existing integrity",
        )
        self.assertIn(
            'case "$published_integrity" in',
            workflow,
            "the workflow must only treat an sha512 integrity as published",
        )

    def test_registry_verification_waits_for_new_package_processing(self):
        workflow = NPM_WORKFLOW.read_text(encoding="utf-8")

        self.assertIn(
            "for attempt in {1..30}",
            workflow,
            "npm says newly published packages may take a few minutes to become available",
        )
        self.assertIn(
            'npm view "@createitv/agc-cli@$version" version',
            workflow,
            "the retry must wait for the exact version before installing",
        )


if __name__ == "__main__":
    unittest.main()
