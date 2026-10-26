"""Real OpenSSL regression: a CA-first chain must select the CSR's leaf."""
import importlib.util
import pathlib
import subprocess
import tempfile
import unittest

spec = importlib.util.spec_from_file_location('package', pathlib.Path(__file__).with_name('live-package.py'))
m = importlib.util.module_from_spec(spec);spec.loader.exec_module(m)

class FingerprintTests(unittest.TestCase):
    def run_command(self, command):
        return subprocess.run(command, capture_output=True, check=True).stdout

    def test_ca_first_chain_selects_csr_key_and_rejects_unrelated_chain(self):
        with tempfile.TemporaryDirectory() as work:
            w = pathlib.Path(work)
            for name in ('root', 'leaf'):
                self.run_command(['openssl', 'req', '-newkey', 'ec', '-pkeyopt', 'ec_paramgen_curve:prime256v1', '-nodes', '-keyout', str(w/(name+'.key')), '-out', str(w/(name+'.csr')), '-subj', '/CN=Offline-'+name])
                self.run_command(['openssl', 'req', '-x509', '-key', str(w/(name+'.key')), '-in', str(w/(name+'.csr')), '-out', str(w/(name+'.pem')), '-days', '1'])
            chain = w/'chain.cer';chain.write_bytes((w/'root.pem').read_bytes()+(w/'leaf.pem').read_bytes())
            expected = self.run_command(['openssl', 'x509', '-in', str(w/'leaf.pem'), '-noout', '-fingerprint', '-sha256']).decode().split('=')[1].strip().replace(':', '')
            legacy = self.run_command(['openssl', 'x509', '-in', str(chain), '-noout', '-fingerprint', '-sha256']).decode().split('=')[1].strip().replace(':', '')
            self.assertNotEqual(legacy, expected)  # 原算法取到链首 CA，复现错误。
            self.assertEqual(m.certificate_fingerprint_for_csr(chain, w/'leaf.csr'), expected)
            self.assertEqual(m.lifecycle.certificate_fingerprint_for_csr(chain, w/'leaf.csr'), expected)
            with self.assertRaises(RuntimeError):m.certificate_fingerprint_for_csr(w/'root.pem', w/'leaf.csr')
            chain.write_bytes((w/'leaf.pem').read_bytes()*2)
            with self.assertRaises(RuntimeError):m.certificate_fingerprint_for_csr(chain, w/'leaf.csr')

if __name__ == '__main__':unittest.main()
