#!/usr/bin/env bun

import { rmSync } from "node:fs"

// ─── validar el nombre del proyecto ─────────────────────────────────────────
const raw = Bun.argv[2] ?? "my-forge-app"
const appName = raw.toLowerCase()

if (!/^[a-z0-9-]+$/.test(appName)) {
  console.error(`✗ Nombre inválido: "${appName}"`)
  console.error("  Solo se permiten letras minúsculas, números y guiones medios (-)")
  process.exit(1)
}

console.log(`\n  forge — creando proyecto "${appName}"...\n`)

// ─── clonar el template ──────────────────────────────────────────────────────
const clone = await Bun.spawn(
  ["git", "clone", "--depth", "1", "https://github.com/donbarrigon/create-ironforge", appName],
  { stdout: "inherit", stderr: "inherit" },
).exited

if (clone !== 0) {
  console.error("✗ Error al clonar el template")
  process.exit(1)
}

process.chdir(appName)

// ─── eliminar .git ───────────────────────────────────────────────────────────
rmSync(".git", { recursive: true, force: true })

// ─── renombrar en Cargo.toml ─────────────────────────────────────────────────
const cargoPath = "Cargo.toml"
const cargoContent = await Bun.file(cargoPath).text()
await Bun.write(cargoPath, cargoContent.replace(`name = "create-forge"`, `name = "${appName}"`))

// ─── actualizar package.json ─────────────────────────────────────────────────
const pkgPath = "package.json"
const pkg = await Bun.file(pkgPath).json()
pkg.name = appName
delete pkg.bin
await Bun.write(pkgPath, `${JSON.stringify(pkg, null, 2)}\n`)

// ─── crear env.json desde env.example.json ───────────────────────────────────
const envPath = "env.example.json"
const env = await Bun.file(envPath).json()
env.app.name = appName
await Bun.write("env.json", `${JSON.stringify(env, null, 2)}\n`)

// ─── eliminar index.ts ───────────────────────────────────────────────────────
rmSync("index.ts", { force: true })

// ─── bun install ─────────────────────────────────────────────────────────────
console.log("  instalando dependencias...\n")
const install = await Bun.spawn(["bun", "install"], { stdout: "inherit", stderr: "inherit" }).exited

if (install !== 0) {
  console.error("✗ Error en bun install")
  process.exit(1)
}

// ─── abrir vscode ────────────────────────────────────────────────────────────
await Bun.spawn(["code", "."], { stdout: "inherit", stderr: "inherit" }).exited

// ─── mensaje de éxito ────────────────────────────────────────────────────────
const green = "\x1b[92m"
const bold = "\x1b[1m"
const dim = "\x1b[2m"
const r = "\x1b[0m"

console.log(`
${bold}${green}  ███████╗ ██████╗ ██████╗  ██████╗ ███████╗${r}
${bold}${green}  ██╔════╝██╔═══██╗██╔══██╗██╔════╝ ██╔════╝${r}
${bold}${green}  █████╗  ██║   ██║██████╔╝██║  ███╗█████╗  ${r}
${bold}${green}  ██╔══╝  ██║   ██║██╔══██╗██║   ██║██╔══╝  ${r}
${bold}${green}  ██║     ╚██████╔╝██║  ██║╚██████╔╝███████╗${r}
${bold}${green}  ╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝${r}
${dim}                                       v0.1.0${r}

  ✓ Proyecto "${appName}" creado exitosamente

  Para empezar:
    cd ${appName}
    cargo run
`)
