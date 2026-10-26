"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { main, resolveBinary, targetFor } = require("../bin/agc.js");

const cases = [
  ["darwin", "x64", ["darwin-x64", "agc"]],
  ["darwin", "arm64", ["darwin-arm64", "agc"]],
  ["linux", "x64", ["linux-x64", "agc"]],
  ["linux", "arm64", ["linux-arm64", "agc"]],
  ["win32", "x64", ["win32-x64", "agc.exe"]],
];

for (const [platform, arch, expected] of cases) {
  test(`resolves ${platform}/${arch}`, () => {
    assert.deepEqual(targetFor(platform, arch), expected);
  });
}

test("rejects unsupported targets", () => {
  assert.throws(() => targetFor("win32", "arm64"), /does not provide a binary/);
  assert.throws(() => targetFor("freebsd", "x64"), /does not provide a binary/);
});

test("resolves an existing embedded binary", () => {
  assert.match(resolveBinary("linux", "x64", () => true), /vendor\/linux-x64\/agc$/);
});

test("rejects an incomplete package", () => {
  assert.throws(() => resolveBinary("linux", "x64", () => false), /package is incomplete/);
});

test("executes the embedded binary and forwards arguments", () => {
  let invocation;
  const status = main({
    platform: "linux", arch: "x64", argv: ["version"], existsSync: () => true,
    spawn: (binary, argv, options) => { invocation = { binary, argv, options }; return { status: 0, signal: null }; },
  });
  assert.equal(status, 0);
  assert.deepEqual(invocation.argv, ["version"]);
  assert.equal(invocation.options.stdio, "inherit");
});

test("propagates child exit status", () => {
  assert.equal(main({ platform: "linux", arch: "x64", existsSync: () => true, spawn: () => ({ status: 23, signal: null }) }), 23);
});

test("returns failure when the child cannot start", () => {
  assert.equal(main({ platform: "linux", arch: "x64", existsSync: () => true, spawn: () => ({ error: new Error("denied") }) }), 1);
});

test("forwards child signals to the launcher process", () => {
  let forwarded;
  const status = main({
    platform: "linux", arch: "x64", existsSync: () => true,
    spawn: () => ({ status: null, signal: "SIGTERM" }),
    kill: (pid, signal) => { forwarded = { pid, signal }; },
  });
  assert.equal(status, 1);
  assert.deepEqual(forwarded, { pid: process.pid, signal: "SIGTERM" });
});

test("treats a null child status as failure", () => {
  assert.equal(main({ platform: "linux", arch: "x64", existsSync: () => true, spawn: () => ({ status: null, signal: null }) }), 1);
});
