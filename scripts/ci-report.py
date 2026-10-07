#!/usr/bin/env python3
"""Convert go test -json output to version-linked, portable CI reports."""
import argparse
import json
import os
import re
from pathlib import Path
import xml.etree.ElementTree as ET


def parse(stream, exit_code):
    tests, packages, diagnostics = {}, {}, []
    valid_events = 0
    for line_number, line in enumerate(stream.splitlines(), 1):
        if not line.strip():
            continue
        try:
            event = json.loads(line)
            if not isinstance(event, dict) or not isinstance(event.get("Action"), str):
                raise ValueError("missing event action")
        except (ValueError, TypeError):
            diagnostics.append(f"Invalid JSON event at line {line_number}")
            continue
        valid_events += 1
        package = event.get("Package") or event.get("ImportPath") or "unknown"
        entry = packages.setdefault(package, {"status": "incomplete", "output": []})
        action, name = event["Action"], event.get("Test")
        if action in ("output", "build-output"):
            output = event.get("Output", "")
            if name:
                tests.setdefault((package, name), {"status": "incomplete", "elapsed": 0, "output": []})["output"].append(output)
            else:
                entry["output"].append(output)
        elif action in ("run", "pass", "fail", "skip") and name:
            test = tests.setdefault((package, name), {"status": "incomplete", "elapsed": 0, "output": []})
            if action != "run":
                test["status"] = action
                test["elapsed"] = event.get("Elapsed", 0)
        elif action in ("pass", "fail", "skip", "build-fail"):
            entry["status"] = "fail" if action == "build-fail" else action
    cases = []
    for (package, name), test in sorted(tests.items()):
        if test["status"] == "incomplete":
            test["status"] = "fail"
            test["output"].append("Test stream ended before a terminal event")
        cases.append(dict(package=package, name=name, **test))
    for package, entry in sorted(packages.items()):
        has_failure = any(case["package"] == package and case["status"] == "fail" for case in cases)
        if entry["status"] in ("incomplete", "fail") and not has_failure:
            cases.append({"package": package, "name": "[package/build]", "status": "fail", "elapsed": 0,
                          "output": entry["output"] or ["Package failed or stream ended before a terminal event"]})
    if not valid_events:
        diagnostics.append("No Go test events received")
    if exit_code != 0 and not any(case["status"] == "fail" for case in cases):
        diagnostics.append(f"go test exited with code {exit_code}")
    for diagnostic in diagnostics:
        cases.append({"package": "ci", "name": "[report] " + diagnostic, "status": "fail", "elapsed": 0, "output": [diagnostic]})
    failures = sum(case["status"] == "fail" for case in cases)
    return {"status": "failed" if failures or exit_code else "passed", "exitCode": exit_code,
            "tests": len(cases), "failures": failures, "skipped": sum(case["status"] == "skip" for case in cases),
            "packages": len(packages), "diagnostics": diagnostics, "cases": cases}


def write_reports(result, output_dir, environment, coverage=""):
    output_dir.mkdir(parents=True, exist_ok=True)
    metadata = {"sha": environment.get("GITHUB_SHA", ""), "runId": environment.get("GITHUB_RUN_ID", ""),
                "runAttempt": environment.get("GITHUB_RUN_ATTEMPT", ""), "ref": environment.get("GITHUB_REF", "")}
    server = environment.get("GITHUB_SERVER_URL", "https://github.com").rstrip("/")
    repository, run_id = environment.get("GITHUB_REPOSITORY", ""), metadata["runId"]
    metadata["runURL"] = f"{server}/{repository}/actions/runs/{run_id}" if repository and run_id else ""
    matches = re.findall(r"^total:\s+.*?([\d.]+%)\s*$", coverage, re.MULTILINE)
    total = matches[-1] if matches else ""
    artifact = dict(result, **metadata, coverageTotal=total)
    (output_dir / "result.json").write_text(json.dumps(artifact, indent=2) + "\n", encoding="utf-8")
    root = ET.Element("testsuites", tests=str(result["tests"]), failures=str(result["failures"]), skipped=str(result["skipped"]))
    suites = {}
    for case in result["cases"]:
        suite = suites.setdefault(case["package"], [])
        suite.append(case)
    for package, cases in suites.items():
        suite = ET.SubElement(root, "testsuite", name=package, tests=str(len(cases)),
                              failures=str(sum(c["status"] == "fail" for c in cases)),
                              skipped=str(sum(c["status"] == "skip" for c in cases)))
        for case in cases:
            node = ET.SubElement(suite, "testcase", classname=package, name=case["name"], time=str(case["elapsed"]))
            if case["status"] == "fail":
                ET.SubElement(node, "failure", message="Go test failure").text = "".join(case["output"])
            elif case["status"] == "skip":
                ET.SubElement(node, "skipped")
    ET.indent(root)
    ET.ElementTree(root).write(output_dir / "junit.xml", encoding="utf-8", xml_declaration=True)
    lines = [f"## Go test: {result['status']}", "", f"Tests: {result['tests']} · Failures: {result['failures']} · Skipped: {result['skipped']}"]
    if metadata["sha"]:
        lines += ["", f"Commit: `{metadata['sha']}`"]
    if metadata["runURL"]:
        lines += ["", f"[Workflow run]({metadata['runURL']})"]
    if total:
        lines += ["", f"Statement coverage: **{total}**"]
    failures = [case for case in result["cases"] if case["status"] == "fail"]
    if failures:
        lines += ["", "Failed tests:", ""]
        lines += [f"- `{case['package']}: {case['name']}`" for case in failures[:30]]
        if len(failures) > 30:
            lines.append(f"- {len(failures) - 30} more; see JUnit report.")
    (output_dir / "summary.md").write_text("\n".join(lines) + "\n", encoding="utf-8")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--input", required=True, type=Path)
    parser.add_argument("--output-dir", type=Path, default=Path("generatedtestReports"))
    parser.add_argument("--exit-code", required=True, type=int)
    parser.add_argument("--coverage-summary", type=Path)
    args = parser.parse_args()
    stream = args.input.read_text(encoding="utf-8", errors="replace") if args.input.exists() else ""
    coverage = args.coverage_summary.read_text(encoding="utf-8", errors="replace") if args.coverage_summary and args.coverage_summary.exists() else ""
    write_reports(parse(stream, args.exit_code), args.output_dir, os.environ, coverage)


if __name__ == "__main__":
    main()
