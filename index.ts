#!/usr/bin/env bun

const appName = Bun.argv[2] ?? "my-forge-app"

// clona el template
await Bun.spawn(["git", "clone", "--depth", "1", "https://github.com/donbarrigon/create-forge", appName], {
  stdout: "inherit",
  stderr: "inherit",
}).exited

// entra a la carpeta y ejecuta blacksmith init
process.chdir(appName)
await Bun.spawn(["bun", "blacksmith.ts", "init"], { stdout: "inherit", stderr: "inherit" }).exited
