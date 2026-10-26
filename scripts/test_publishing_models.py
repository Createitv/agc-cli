import importlib.util
from pathlib import Path
import unittest

spec = importlib.util.spec_from_file_location('models', Path(__file__).with_name('audit-publishing-models.py'))
models = importlib.util.module_from_spec(spec)
spec.loader.exec_module(models)


class ModelAuditTests(unittest.TestCase):
    def test_missing_nested_required_field_is_detected(self):
        doc = {'sourceUrl': 'https://official/model', 'tables': [{'rows': [['tagId', 'M', 'String(64)', '标签ID']]}]}
        result = models.audit_model(doc, [], 'body.tagQualInfos[]', {}, 'qualification-review-create')
        self.assertEqual(result[0]['field'], 'body.tagQualInfos[].tagId')
        self.assertEqual(result[0]['code'], 'missing_field')

    def test_shared_name_resolves_exact_source_link(self):
        doc = {'sourceUrl': 'https://official/user-info', 'tables': [{'rows': [['baseInfo', 'O', 'AgreementBaseInfo', '基础信息']]}], 'links': [
            {'title': 'AgreementBaseInfo', 'sourceUrl': 'https://official/user-base'},
            {'title': 'AgreementBaseInfo（隐私政策协议）', 'sourceUrl': 'https://official/privacy-base'}]}
        base = {'sourceUrl': 'https://official/user-base', 'tables': [{'rows': [['name', 'M', 'String(100)', '协议名称']]}]}
        fields = [{'name': 'baseInfo', 'type': 'object', 'properties': [{'name': 'name', 'type': 'string', 'required': True, 'maxLength': 100}]}]
        self.assertEqual(models.audit_model(doc, fields, 'body.agreementInfos[]', {'https://official/user-base': base}, 'user-agreement-create'), [])

    def test_integer_width_is_not_a_bound(self):
        doc = {'sourceUrl': 'https://official/model', 'tables': [{'rows': [['size', 'M', 'Integer(64)', '大小']]}]}
        self.assertEqual(models.audit_model(doc, [{'name': 'size', 'required': True, 'type': 'integer'}], 'body.file', {}, 'fixture'), [])

    def test_optional_permissions_array_retains_explicit_interval(self):
        doc = {'sourceUrl': 'https://official/model', 'tables': [{'rows': [['permissions', 'O', 'List<String(100)>', '数组长度为1~200。']]}]}
        fields = [{'name': 'permissions', 'type': 'array', 'itemsType': 'string', 'itemsMaxLength': 100, 'maxItems': 200}]
        issues = models.audit_model(doc, fields, 'body.permissions[]', {}, 'fixture')
        self.assertEqual([item['code'] for item in issues], ['minItems'])


if __name__ == '__main__':
    unittest.main()
