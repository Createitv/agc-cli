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


if __name__ == "__main__":
    unittest.main()
