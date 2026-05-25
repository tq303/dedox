#!/usr/bin/env node

const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");

const VERSION = require("./package.json").version;

const PLATFORM_MAP = {
  "darwin-arm64": "ddx-darwin-arm64",
  "darwin-x64": "ddx-darwin-amd64",
  "linux-x64": "ddx-linux-amd64",
  "win32-x64": "ddx-windows-amd64.exe",
};

const key = `${process.platform}-${process.arch}`;
const binaryName = PLATFORM_MAP[key];

if (!binaryName) {
  console.error(`ddx: unsupported platform ${key}`);
  process.exit(1);
}

const url = `https://github.com/tq303/ddx/releases/download/v${VERSION}/${binaryName}`;
const dest = path.join(__dirname, "bin", "ddx" + (process.platform === "win32" ? ".exe" : ""));

fs.mkdirSync(path.join(__dirname, "bin"), { recursive: true });

console.log(`ddx: downloading binary for ${key}...`);

function download(url, dest, cb) {
  const file = fs.createWriteStream(dest);
  https.get(url, (res) => {
    if (res.statusCode === 302 || res.statusCode === 301) {
      return download(res.headers.location, dest, cb);
    }
    if (res.statusCode !== 200) {
      cb(new Error(`download failed: ${res.statusCode}`));
      return;
    }
    res.pipe(file);
    file.on("finish", () => file.close(cb));
  }).on("error", cb);
}

download(url, dest, (err) => {
  if (err) {
    console.error(`ddx: failed to download binary: ${err.message}`);
    process.exit(1);
  }
  fs.chmodSync(dest, 0o755);
  console.log("ddx: ready");
});
