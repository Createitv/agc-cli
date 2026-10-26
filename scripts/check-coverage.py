"""Keep each coverage run isolated; publish the complete profile atomically."""
import argparse
import os
import pathlib
import re
import shlex
import subprocess
import tempfile


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--go', default='go')
    parser.add_argument('--minimum', type=float)
    args = parser.parse_args()
    command = shlex.split(args.go)
    with tempfile.TemporaryDirectory(prefix='agc-coverage-') as work:
        profile = pathlib.Path(work)/'coverage.out'
        result = subprocess.run([*command,'test','./...','-coverprofile='+str(profile)])
        if result.returncode: return result.returncode
        result = subprocess.run([*command,'tool','cover','-func='+str(profile)],capture_output=True,text=True)
        print(result.stdout,end='')
        if result.returncode:
            print(result.stderr,end='')
            return result.returncode
        matches = re.findall(r'^total:\s+.*?([\d.]+)%\s*$',result.stdout,re.MULTILINE)
        if len(matches)!=1: raise RuntimeError('Missing unambiguous total coverage')
        total = float(matches[0])
        # 兼容已有本地报告路径；独立校验完成后才保存，读者不会看到半写文件。
        fd, temporary = tempfile.mkstemp(prefix='.coverage-',dir='.')
        try:
            with os.fdopen(fd,'wb') as output: output.write(profile.read_bytes())
            os.replace(temporary,'coverage.out')
        finally:
            pathlib.Path(temporary).unlink(missing_ok=True)
        if args.minimum is not None:
            if total < args.minimum:
                print(f'coverage {total:.1f}% is below {args.minimum:.1f}%')
                return 1
            print(f'coverage {total:.1f}% >= {args.minimum:.1f}%')
    return 0


if __name__=='__main__': raise SystemExit(main())
