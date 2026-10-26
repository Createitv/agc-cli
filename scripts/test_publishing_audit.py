"""Independent regressions for documented table/example conflicts."""
import importlib.util
import unittest
from pathlib import Path

spec = importlib.util.spec_from_file_location('audit', Path(__file__).with_name('audit-publishing-contracts.py'))
audit = importlib.util.module_from_spec(spec)
spec.loader.exec_module(audit)


class AuditTests(unittest.TestCase):
    def test_shared_transport_keeps_agreement_document_identity(self):
        ref = {'checkedAt': 'fixture', 'endpoints': [{
            'name': '查询隐私政策协议列表', 'method': 'POST',
            'path': '/api/publish/v2/agreement/list', 'sourceUrl': 'https://official/privacy',
            'parameters': [],
        }]}
        local = [
            {'id': 'user-agreement-list', 'method': 'POST', 'path': '/api/publish/v2/agreement/list', 'sourceUrl': 'https://official/user'},
            {'id': 'privacy-agreement-list', 'method': 'POST', 'path': '/api/publish/v2/agreement/list', 'sourceUrl': 'https://official/privacy'},
        ]
        self.assertEqual(audit.audit(ref, local)['endpoints'][0]['command'], 'privacy-agreement-list')

    def check(self, endpoint, parameter, actual):
        ref = {'checkedAt': 'fixture', 'endpoints': [{'name': endpoint, 'method': 'POST', 'path': '/fixture', 'parameters': [parameter]}]}
        local = [{'id': 'fixture', 'method': 'POST', 'path': '/fixture', 'parameters': [actual]}]
        return audit.audit(ref, local)['endpoints'][0]

    def test_identity_alternative(self):
        p = {'in': 'header', 'name': 'appId', 'required': True, 'type': 'String(32)'}
        result = self.check('查询应用所有版本类型的版本列表', p, {**p, 'required': False, 'type': 'string'})
        self.assertEqual(result['parameterIssues'], [])
        self.assertTrue(result['reviewedExceptions'])

    def test_agreement_example_array(self):
        p = {'in': 'body', 'name': 'agreementInfos', 'required': False, 'type': 'AgreementInfo'}
        result = self.check('新建用户协议', p, {**p, 'type': 'array'})
        self.assertEqual(result['parameterIssues'], [])
        self.assertTrue(result['reviewedExceptions'])
        # A table-shaped object must still be reported as wrong.
        bad = self.check('新建用户协议', p, {**p, 'type': 'object'})
        self.assertTrue(bad['parameterIssues'])

    def test_no_blanket_requiredness_exemption(self):
        p = {'in': 'header', 'name': 'appId', 'required': True, 'type': 'String(32)'}
        result = self.check('新建用户协议', p, {**p, 'required': False, 'type': 'string'})
        self.assertTrue(result['parameterIssues'])


if __name__ == '__main__':
    unittest.main()
