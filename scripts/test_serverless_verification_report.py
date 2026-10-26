import importlib.util,pathlib,unittest
p=pathlib.Path(__file__).with_name('serverless-verification-report.py');s=importlib.util.spec_from_file_location('slreport',p);m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class EvidenceTests(unittest.TestCase):
 def test_self_reported_success_cannot_replace_real_pass_events(self):
  endpoints=[{'familyId':'storage','id':'list-files','sourceUrl':'https://official'}]
  r=m.matrix(endpoints,[{'Action':'output','Test':'TestEveryEndpointHTTPTransport/storage/list-files','Output':'passed'}]);self.assertFalse(r['offlineComplete']);self.assertEqual(r['liveVerified'],0)
 def test_every_adapter_and_whole_package_pass_is_required(self):
  endpoints=[{'familyId':'storage','id':'list-files','sourceUrl':'https://official'}]
  events=[{'Action':'pass','Test':x+'/storage/list-files'} for x in ('TestEveryEndpointHTTPTransport','TestEveryServerlessOperationREST','TestEveryRegisteredEndpointCanDryRun')]
  self.assertFalse(m.matrix(endpoints,events)['offlineComplete'])
  events +=[{'Action':'pass','Package':'github.com/Createitv/agc-cli/'+p} for p in ('cmd/agc/command','pkg/server')]
  r=m.matrix(endpoints,events);self.assertTrue(r['offlineComplete']);self.assertEqual(r['operations'][0]['live']['status'],'not-run')
if __name__=='__main__':unittest.main()
