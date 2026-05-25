#!/usr/bin/env node

const { spawnSync } = require("child_process");
const path = require("path");

const binary = path.join(
  __dirname,
  "bin",
  "ddx" + (process.platform === "win32" ? ".exe" : "")
);

const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });

process.exit(result.status ?? 1);
