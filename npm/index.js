#!/usr/bin/env node

const { spawnSync } = require("child_process");
const path = require("path");

const binary = path.join(
  __dirname,
  "bin",
  "ddx" + (process.platform === "win32" ? ".exe" : "")
);

function ddx(filePath, options = {}) {
  const args = [];
  for (const filter of options.filters ?? []) {
    args.push("-f", filter);
  }
  args.push(filePath);

  const result = spawnSync(binary, args, { encoding: "utf8" });
  if (result.status !== 0) {
    throw new Error(result.stderr?.trim() || `ddx exited with code ${result.status}`);
  }
  return result.stdout;
}

if (require.main === module) {
  const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });
  process.exit(result.status ?? 1);
} else {
  module.exports = ddx;
}
