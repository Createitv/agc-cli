import importlib.util
import pathlib
import unittest

spec = importlib.util.spec_from_file_location("api_suite", pathlib.Path(__file__).with_name("api-suite.py"))
suite = importlib.util.module_from_spec(spec)
spec.loader.exec_module(suite)


class SuiteTests(unittest.TestCase):
    def test_http_success_does_not_hide_business_failure(self):
        self.assertEqual(suite.response_status(200, {"rtnCode": 0, "ret": {"code": 12}}), "business-error")
        self.assertEqual(suite.response_status(401, {"code": 0}), "http-error")
        self.assertEqual(suite.response_status(200, {"groups": []}), "response-needs-review")
        self.assertEqual(suite.response_status(200, {"rtnCode": "0"}), "verified-read")

    def test_group_member_fixture_uses_matching_application(self):
        group = {"appId": "app", "groupId": "group"}
        fixture = suite.read_fixture({"familyId": "testing", "id": "test-api-query-test-user"}, {}, "app", "project", group)
        self.assertEqual(fixture["headers"], {"appId": "app", "groupId": "group"})
        self.assertIsNone(suite.read_fixture({"familyId": "testing", "id": "test-api-query-test-user"}, {}, "other-app", "project", group))

    def test_prior_success_cannot_enable_a_mutation(self):
        fixture = suite.read_fixture({"familyId": "publishing", "id": "app-submit", "method": "POST"}, {"status": "verified-read", "testedQuery": {}}, "app", "project", {})
        self.assertIsNone(fixture)

    def test_official_table_preserves_parameter_locations_and_types(self):
        html = '<h4>请求参数</h4><h5>Header</h5><table><tr><th>参数名称</th><th>必选(M)/可选(O)</th><th>类型</th><th>参数说明</th></tr><tr><td>appId</td><td>M</td><td>String(32)</td><td>应用ID</td></tr></table><h5>Body</h5><table><tr><th>参数名称</th><th>必选(M)/可选(O)</th><th>类型</th><th>参数说明</th></tr><tr><td>groupType</td><td>O</td><td>Integer</td><td>0: 外部</td></tr></table><h4>响应参数</h4><table><tr><td>groupId</td><td>O</td><td>String</td><td>ID</td></tr></table>'
        params = suite.extract_parameters(html)
        self.assertEqual([(p["name"], p["in"], p["required"], p["type"]) for p in params], [("appId", "header", True, "String(32)"), ("groupType", "body", False, "Integer")])
        self.assertEqual(suite.extract_parameters(html.replace("Header", "[h2]Header").replace("Body", "[h2]Body")), params)

    def test_comment_read_supplies_bounded_millisecond_window(self):
        endpoint = {"familyId": "comments", "id": "com-rating-harmonyos", "method": "GET"}
        fixture = suite.read_fixture(endpoint, {"testedQuery": {"appId": "old"}}, "app", "project", {})
        self.assertEqual(fixture["query"]["countries"], "CN")
        self.assertGreater(int(fixture["query"]["endTime"]), int(fixture["query"]["beginTime"]))

    def test_live_cli_timeout_is_reported_without_output(self):
        import subprocess
        from unittest.mock import patch
        with patch.object(suite.subprocess, "run", side_effect=subprocess.TimeoutExpired("cli", 1, output="SECRET")):
            result = suite.invoke_cli("binary", {"familyId": "testing", "id": "query"}, {"headers": {}}, "profile", False)
        self.assertEqual(result["status"], "timeout")
        self.assertNotIn("SECRET", str(result))

    def test_cli_error_retains_codes_without_raw_body(self):
        from unittest.mock import patch
        import subprocess
        proc = subprocess.CompletedProcess([], 1, "", "Error: agc endpoint returned HTTP 200 (123): SECRET")
        with patch.object(suite.subprocess, "run", return_value=proc):
            result = suite.invoke_cli("binary", {"familyId": "testing", "id": "query"}, {}, "profile", False)
        self.assertEqual(result["status"], "business-error")
        self.assertEqual(result["businessCode"], "123")
        self.assertNotIn("SECRET", str(result))


if __name__ == "__main__":
    unittest.main()
