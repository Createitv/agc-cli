import hashlib
import importlib.util
import io
import json
import stat
from pathlib import Path
import tarfile
import tempfile
import unittest
import zipfile


SPEC = importlib.util.spec_from_file_location("npm_package", Path(__file__).with_name("build-npm-package.py"))
npm_package = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(npm_package)


class NpmPackageTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)
        self.release = self.root / "release"
        self.release.mkdir()
        checksums = []
        for goos, goarch, _, binary_name, extension in npm_package.TARGETS:
            archive = self.release / f"agc-cli_1.2.3_{goos}_{goarch}.{extension}"
            payload = f"{goos}/{goarch}".encode()
            if extension == "zip":
                with zipfile.ZipFile(archive, "w") as bundle:
                    bundle.writestr(binary_name, payload)
            else:
                with tarfile.open(archive, "w:gz") as bundle:
                    info = tarfile.TarInfo(binary_name)
                    info.size = len(payload)
                    bundle.addfile(info, io.BytesIO(payload))
            checksums.append(f"{hashlib.sha256(archive.read_bytes()).hexdigest()}  {archive.name}")
        (self.release / "checksums.txt").write_text("\n".join(checksums) + "\n")

    def tearDown(self):
        self.temp.cleanup()

    def test_builds_complete_script_free_package(self):
        output = self.root / "package"
        npm_package.build_package("1.2.3", self.release, output)
        manifest = json.loads((output / "package.json").read_text())
        self.assertEqual(manifest["name"], "@createitv/agc-cli")
        self.assertEqual(manifest["version"], "1.2.3")
        self.assertEqual(manifest["bin"], {"agc": "bin/agc.js"})
        self.assertNotIn("scripts", manifest)
        for _, _, npm_target, binary_name, _ in npm_package.TARGETS:
            self.assertTrue((output / "vendor" / npm_target / binary_name).is_file())
        self.assertTrue((output / "vendor" / "linux-x64" / "agc").stat().st_mode & stat.S_IXUSR)

    def test_rejects_checksum_mismatch(self):
        archive = self.release / "agc-cli_1.2.3_linux_amd64.tar.gz"
        archive.write_bytes(archive.read_bytes() + b"tampered")
        with self.assertRaisesRegex(ValueError, "checksum mismatch"):
            npm_package.build_package("1.2.3", self.release, self.root / "package")

    def test_rejects_missing_platform(self):
        (self.release / "agc-cli_1.2.3_windows_amd64.zip").unlink()
        with self.assertRaisesRegex(ValueError, "release archive missing"):
            npm_package.build_package("1.2.3", self.release, self.root / "package")

    def test_rejects_invalid_version(self):
        with self.assertRaisesRegex(ValueError, "semantic version"):
            npm_package.build_package("v1.2.3", self.release, self.root / "package")

    def test_rejects_non_empty_output_directory(self):
        output = self.root / "package"
        output.mkdir()
        (output / "existing").write_text("keep")
        with self.assertRaisesRegex(ValueError, "not empty"):
            npm_package.build_package("1.2.3", self.release, output)

    def test_rejects_checksum_missing_from_manifest(self):
        (self.release / "checksums.txt").write_text("")
        with self.assertRaisesRegex(ValueError, "checksum missing"):
            npm_package.build_package("1.2.3", self.release, self.root / "package")


if __name__ == "__main__":
    unittest.main()
