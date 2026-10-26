"""Two actual Make invocations must not consume one another's profiles."""
import os
import pathlib
import subprocess
import tempfile
import unittest
import concurrent.futures

ROOT = pathlib.Path(__file__).resolve().parents[1]

class CoverageTests(unittest.TestCase):
    def test_parallel_coverage_uses_owned_profiles(self):
        with tempfile.TemporaryDirectory() as work:
            directory = pathlib.Path(work)
            fake = directory/'go'
            fake.write_text('''#!/usr/bin/env python3
import sys,os,pathlib,time
args=sys.argv[1:]
if args[0]=='test':
 p=pathlib.Path(next(a.split('=',1)[1] for a in args if a.startswith('-coverprofile=')))
 p.write_text(os.environ['COVERAGE_TEST_OWNER'])
 time.sleep(.3)
elif args[:2]==['tool','cover']:
 p=pathlib.Path(next(a.split('=',1)[1] for a in args if a.startswith('-func=')))
 if p.read_text()!=os.environ['COVERAGE_TEST_OWNER']:sys.exit(7)
 print('total: (statements) 84.0%')
else:sys.exit(8)
''')
            fake.chmod(0o700)
            def run(owner):
                env=dict(os.environ,COVERAGE_TEST_OWNER=owner)
                return subprocess.run(['make','-f',str(ROOT/'Makefile'),'coverage-check','GO='+str(fake)],cwd=directory,env=env,capture_output=True,text=True)
            with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
                results=list(pool.map(run,['first','second']))
            self.assertEqual([r.returncode for r in results],[0,0],str([(r.stdout,r.stderr) for r in results]))

if __name__=='__main__':unittest.main()
