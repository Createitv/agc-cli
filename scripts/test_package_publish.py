import importlib.util
from pathlib import Path
import unittest

spec = importlib.util.spec_from_file_location('publisher', Path(__file__).with_name('verify-package-publish.py'))
publisher = importlib.util.module_from_spec(spec)
spec.loader.exec_module(publisher)


class PublicationReadbackTest(unittest.TestCase):
    def setUp(self):
        self.prefix = 'https://github.com/Createitv/agc-cli/releases/download/v1.2.3/'
        self.hash = 'a' * 64
        self.archives = [f'agc-cli_1.2.3_{os}_{arch}.tar.gz' for os in ['darwin', 'linux'] for arch in ['amd64', 'arm64']]
        self.formula = 'version "1.2.3"\n' + '\n'.join(f'url "{self.prefix}{name}"\nsha256 "{self.hash}"' for name in self.archives)
        self.windows = 'agc-cli_1.2.3_windows_amd64.zip'
        self.manifest = {'version': '1.2.3', 'architecture': {'64bit': {'bin': 'agc.exe', 'url': self.prefix + self.windows, 'hash': self.hash}}}
        self.checksums = dict.fromkeys(self.archives + [self.windows], self.hash)

    def test_matching_release_is_accepted(self):
        publisher.verify(self.formula, self.manifest, 'v1.2.3', self.checksums)
        self.manifest['architecture']['64bit']['bin'] = ['agc.exe']
        publisher.verify(self.formula, self.manifest, 'v1.2.3', self.checksums)

    def test_stale_package_version_is_rejected(self):
        with self.assertRaisesRegex(ValueError, 'Homebrew formula version'):
            publisher.verify(self.formula.replace('version "1.2.3"', 'version "1.2.2"'), self.manifest, 'v1.2.3', self.checksums)
        self.manifest['version'] = '1.2.2'
        with self.assertRaisesRegex(ValueError, 'Scoop manifest version'):
            publisher.verify(self.formula, self.manifest, 'v1.2.3', self.checksums)

    def test_wrong_archive_hash_is_rejected(self):
        self.checksums[self.windows] = 'b' * 64
        with self.assertRaisesRegex(ValueError, 'checksum'):
            publisher.verify(self.formula, self.manifest, 'v1.2.3', self.checksums)

    def test_missing_platform_and_wrong_executable_are_rejected(self):
        with self.assertRaisesRegex(ValueError, 'incomplete'):
            publisher.verify('version "1.2.3"', self.manifest, 'v1.2.3', self.checksums)
        self.manifest['architecture']['64bit']['bin'] = 'other.exe'
        with self.assertRaisesRegex(ValueError, 'executable'):
            publisher.verify(self.formula, self.manifest, 'v1.2.3', self.checksums)

    def test_archive_from_wrong_release_is_rejected(self):
        self.manifest['architecture']['64bit']['url'] = self.prefix.replace('v1.2.3', 'v1.2.2') + self.windows
        with self.assertRaisesRegex(ValueError, 'URL/checksum'):
            publisher.verify(self.formula, self.manifest, 'v1.2.3', self.checksums)


if __name__ == '__main__':
    unittest.main()
