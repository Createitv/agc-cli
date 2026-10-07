#!/usr/bin/env python3
"""Check publisher access and read back package versions against release checksums."""
import argparse
import base64
import json
import os
import re
import sys
import urllib.request

REPOSITORIES = {"homebrew-tap": ("master", "agc-cli.rb"), "scoop-bucket": ("main", "agc-cli.json")}


def api(path):
    token = os.environ.get("GH_TOKEN", "")
    if not token:
        raise ValueError("TAP_GITHUB_TOKEN is required")
    request = urllib.request.Request("https://api.github.com/" + path,
                                    headers={"Authorization": "Bearer " + token, "Accept": "application/vnd.github+json", "User-Agent": "agc-package-verification"})
    with urllib.request.urlopen(request, timeout=30) as response:
        return json.load(response)


def verify(formula, manifest, tag, checksums):
    version = tag.removeprefix("v")
    if not re.search(r'^\s*version\s+"' + re.escape(version) + r'"\s*$', formula, re.M):
        raise ValueError("Homebrew formula version does not match release")
    if manifest.get("version") != version:
        raise ValueError("Scoop manifest version does not match release")
    pairs = re.findall(r'url\s+"([^"]+)"\s+sha256\s+"([a-f0-9]{64})"', formula)
    expected = {f"agc-cli_{version}_{platform}_{arch}.tar.gz" for platform in ["darwin", "linux"] for arch in ["amd64", "arm64"]}
    if {url.rsplit("/", 1)[-1] for url, _ in pairs} != expected:
        raise ValueError("Homebrew platform archives are incomplete")
    windows = manifest.get("architecture", {}).get("64bit", {})
    if windows.get("bin") not in ("agc.exe", ["agc.exe"]):
        raise ValueError("Scoop executable is incorrect")
    pairs.append((windows.get("url", ""), windows.get("hash", "")))
    prefix = f"https://github.com/Createitv/agc-cli/releases/download/{tag}/"
    for url, checksum in pairs:
        filename = url.rsplit("/", 1)[-1]
        if not url.startswith(prefix) or not checksum or checksums.get(filename) != checksum:
            raise ValueError("Package URL/checksum does not match release: " + filename)
    if pairs[-1][0] != prefix + f"agc-cli_{version}_windows_amd64.zip":
        raise ValueError("Scoop archive does not match release")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--check-access", action="store_true")
    parser.add_argument("--tag")
    args = parser.parse_args()
    if not args.check_access and not args.tag:
        parser.error("--check-access or --tag is required")
    if args.check_access:
        for repository, (branch, _) in REPOSITORIES.items():
            metadata = api("repos/Createitv/" + repository)
            if not metadata.get("permissions", {}).get("push"):
                raise ValueError("Publisher token needs contents write access: " + repository)
            if metadata.get("default_branch") != branch:
                raise ValueError("Publisher default branch changed: " + repository)
            print(repository + ": publisher write access verified")
    if args.tag:
        if not re.fullmatch(r"v\d+\.\d+\.\d+", args.tag):
            raise ValueError("Readback expects a stable vMAJOR.MINOR.PATCH release")
        contents = []
        for repository, (branch, path) in REPOSITORIES.items():
            response = api(f"repos/Createitv/{repository}/contents/{path}?ref={branch}")
            contents.append(base64.b64decode(response["content"]).decode())
        url = f"https://github.com/Createitv/agc-cli/releases/download/{args.tag}/checksums.txt"
        with urllib.request.urlopen(url, timeout=30) as response:
            checksum_lines = response.read().decode().splitlines()
        checksums = {parts[1].lstrip("*"): parts[0] for line in checksum_lines if len(parts := line.split()) == 2}
        verify(contents[0], json.loads(contents[1]), args.tag, checksums)
        print("Homebrew and Scoop versions, archive URLs and SHA256 match " + args.tag)


if __name__ == "__main__":
    try:
        main()
    except Exception as error:
        print("Package verification failed: " + str(error), file=sys.stderr)
        sys.exit(1)
