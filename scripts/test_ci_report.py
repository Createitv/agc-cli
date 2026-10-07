import importlib.util
import json
import tempfile
import unittest
from pathlib import Path

spec = importlib.util.spec_from_file_location("ci_report", Path(__file__).with_name("ci-report.py"))
report = importlib.util.module_from_spec(spec)
spec.loader.exec_module(report)


def events(*items):
    return "\n".join(json.dumps(item) for item in items)


class ReportTests(unittest.TestCase):
    def test_pass_skip_and_subtest_failure(self):
        data = events(
            {"Action": "run", "Package": "example/pkg", "Test": "TestParent"},
            {"Action": "fail", "Package": "example/pkg", "Test": "TestParent/child", "Elapsed": 0.1},
            {"Action": "fail", "Package": "example/pkg", "Test": "TestParent"},
            {"Action": "skip", "Package": "example/pkg", "Test": "TestSkip"},
            {"Action": "pass", "Package": "example/pkg", "Test": "TestPass"},
            {"Action": "fail", "Package": "example/pkg"},
        )
        result = report.parse(data, 1)
        self.assertEqual(result["status"], "failed")
        self.assertEqual(result["tests"], 4)
        self.assertEqual(result["failures"], 2)
        self.assertEqual(result["skipped"], 1)

    def test_build_error_is_reported(self):
        data = events({"Action": "build-output", "ImportPath": "example/pkg", "Output": "compile error\n"},
                      {"Action": "build-fail", "ImportPath": "example/pkg"})
        self.assertEqual(report.parse(data, 1)["failures"], 1)

    def test_failed_package_without_test(self):
        result = report.parse(events({"Action": "fail", "Package": "example/pkg"}), 1)
        self.assertEqual(result["failures"], 1)

    def test_missing_and_truncated_data(self):
        for data in ("", '{"Action":', events({"Action": "run", "Package": "example/pkg", "Test": "TestLost"})):
            with self.subTest(data=data):
                self.assertEqual(report.parse(data, 0)["status"], "failed")

    def test_exit_code_overrides_passing_events(self):
        self.assertEqual(report.parse(events({"Action": "pass", "Package": "example/pkg"}), 2)["status"], "failed")

    def test_artifacts_and_known_metadata_only(self):
        import xml.etree.ElementTree as ET
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp)
            result = report.parse(events({"Action": "pass", "Package": "example/pkg", "Test": "TestOK"},
                                         {"Action": "pass", "Package": "example/pkg"}), 0)
            report.write_reports(result, output, {"GITHUB_SHA": "abc", "GITHUB_SERVER_URL": "https://github.com",
                "GITHUB_REPOSITORY": "owner/repo", "GITHUB_RUN_ID": "42", "SECRET_TOKEN": "never-print"}, "total: (statements) 75.0%\n")
            self.assertEqual(ET.parse(output / "junit.xml").getroot().tag, "testsuites")
            text = (output / "result.json").read_text()
            self.assertNotIn("never-print", text)
            self.assertIn("https://github.com/owner/repo/actions/runs/42", text)
            self.assertIn("75.0%", (output / "summary.md").read_text())


if __name__ == "__main__":
    unittest.main()
