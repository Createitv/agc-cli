#!/usr/bin/env python3
"""Build the npm package from checksum-verified GoReleaser archives."""

import argparse
import hashlib
import json
from pathlib import Path
import re
import shutil
import stat
import tarfile
import zipfile


ROOT = Path(__file__).resolve().parents[1]
PACKAGE_SOURCE = ROOT / "packaging" / "npm"
TARGETS = (
    ("darwin", "amd64", "darwin-x64", "agc", "tar.gz"),
    ("darwin", "arm64", "darwin-arm64", "agc", "tar.gz"),
    ("linux", "amd64", "linux-x64", "agc", "tar.gz"),
    ("linux", "arm64", "linux-arm64", "agc", "tar.gz"),
    ("windows", "amd64", "win32-x64", "agc.exe", "zip"),
)


def read_checksums(path):
    checksums = {}
    for line in path.read_text(encoding="utf-8").splitlines():
        parts = line.split()
        if len(parts) == 2 and re.fullmatch(r"[0-9a-fA-F]{64}", parts[0]):
            checksums[parts[1].lstrip("*")] = parts[0].lower()
    return checksums


def verify_archive(path, checksums):
    expected = checksums.get(path.name)
    if not expected:
        raise ValueError(f"checksum missing for {path.name}")
    observed = hashlib.sha256(path.read_bytes()).hexdigest()
    if observed != expected:
        raise ValueError(f"checksum mismatch for {path.name}")


def archive_binary(path, binary_name):
    if path.suffix == ".zip":
        with zipfile.ZipFile(path) as archive:
            matches = [name for name in archive.namelist() if Path(name).name == binary_name and not name.endswith("/")]
            if len(matches) != 1:
                raise ValueError(f"expected one {binary_name} in {path.name}, found {len(matches)}")
            return archive.read(matches[0])
    with tarfile.open(path, "r:gz") as archive:
        matches = [member for member in archive.getmembers() if Path(member.name).name == binary_name and member.isfile()]
        if len(matches) != 1:
            raise ValueError(f"expected one {binary_name} in {path.name}, found {len(matches)}")
        source = archive.extractfile(matches[0])
        if source is None:
            raise ValueError(f"cannot read {binary_name} from {path.name}")
        return source.read()


def build_package(version, release_dir, output_dir):
    if not re.fullmatch(r"\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?", version):
        raise ValueError("version must be a semantic version without a leading v")
    if output_dir.exists() and any(output_dir.iterdir()):
        raise ValueError(f"output directory is not empty: {output_dir}")
    output_dir.mkdir(parents=True, exist_ok=True)
    checksums = read_checksums(release_dir / "checksums.txt")

    for goos, goarch, npm_target, binary_name, extension in TARGETS:
        archive_name = f"agc-cli_{version}_{goos}_{goarch}.{extension}"
        archive = release_dir / archive_name
        if not archive.is_file():
            raise ValueError(f"release archive missing: {archive_name}")
        verify_archive(archive, checksums)
        destination = output_dir / "vendor" / npm_target / binary_name
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_bytes(archive_binary(archive, binary_name))
        if goos != "windows":
            destination.chmod(destination.stat().st_mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)

    shutil.copytree(PACKAGE_SOURCE / "bin", output_dir / "bin")
    (output_dir / "bin" / "agc.js").chmod(0o755)
    shutil.copy2(PACKAGE_SOURCE / "README.md", output_dir / "README.md")
    shutil.copy2(ROOT / "LICENSE", output_dir / "LICENSE")
    package = {
        "name": "@createitv/agc-cli",
        "version": version,
        "description": "AppGallery Connect command center for CI, scripts, and agents",
        "license": "MIT",
        "homepage": "https://agccli.app/",
        "repository": {"type": "git", "url": "git+https://github.com/Createitv/agc-cli.git"},
        "bugs": {"url": "https://github.com/Createitv/agc-cli/issues"},
        "keywords": ["appgallery-connect", "huawei", "harmonyos", "cli", "cicd"],
        "bin": {"agc": "bin/agc.js"},
        "files": ["bin", "vendor", "README.md", "LICENSE"],
        "engines": {"node": ">=18"},
        "publishConfig": {"access": "public"},
    }
    (output_dir / "package.json").write_text(json.dumps(package, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--version", required=True)
    parser.add_argument("--release-dir", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    build_package(args.version, args.release_dir.resolve(), args.output.resolve())
    print(args.output.resolve())


if __name__ == "__main__":
    main()
