#!/usr/bin/env node

"use strict";

const fs = require("node:fs");
const path = require("node:path");
const { spawnSync } = require("node:child_process");

const targets = {
  "darwin-x64": ["darwin-x64", "agc"],
  "darwin-arm64": ["darwin-arm64", "agc"],
  "linux-x64": ["linux-x64", "agc"],
  "linux-arm64": ["linux-arm64", "agc"],
  "win32-x64": ["win32-x64", "agc.exe"],
};

function targetFor(platform = process.platform, arch = process.arch) {
  const key = `${platform}-${arch}`;
  const target = targets[key];
  if (!target) {
    throw new Error(
      `agc-cli does not provide a binary for ${platform}/${arch}. ` +
        "Supported targets: macOS x64/arm64, Linux x64/arm64, and Windows x64.",
    );
  }
  return target;
}

function resolveBinary(platform = process.platform, arch = process.arch, existsSync = fs.existsSync) {
  const target = targetFor(platform, arch);
  const binary = path.join(__dirname, "..", "vendor", target[0], target[1]);
  if (!existsSync(binary)) {
    throw new Error(`agc-cli package is incomplete: missing ${path.relative(path.join(__dirname, ".."), binary)}`);
  }
  return binary;
}

function main({
  platform = process.platform,
  arch = process.arch,
  argv = process.argv.slice(2),
  existsSync = fs.existsSync,
  spawn = spawnSync,
  kill = process.kill.bind(process),
} = {}) {
  let binary;
  try {
    binary = resolveBinary(platform, arch, existsSync);
  } catch (error) {
    console.error(error.message);
    return 1;
  }
  const result = spawn(binary, argv, { stdio: "inherit" });
  if (result.error) {
    console.error(`failed to start agc: ${result.error.message}`);
    return 1;
  }
  if (result.signal) {
    kill(process.pid, result.signal);
    return 1;
  }
  return result.status === null ? 1 : result.status;
}

module.exports = { main, resolveBinary, targetFor };

if (require.main === module) {
  process.exitCode = main();
}
