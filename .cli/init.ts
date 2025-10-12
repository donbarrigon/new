import { existsSync, readdirSync, readFileSync, renameSync, writeFileSync } from "node:fs"
import { basename, dirname, join } from "node:path"
import { color } from "./config"
import { inputc, runc } from "./console"
import { showHelp } from "./help"

export async function init() {
  const arg = process.argv[3]
  if (!arg) {
    await newProject()
    return
  }
  if (arg && arg === "fork") {
    await fork()
    return
  }
  showHelp()
  console.log(`${color.red}Parse tiene un error de sintaxis${color.reset}`)
  console.log(`talvez quiso decir`)
  console.log(`${color.green}bun make init${color.reset}`)
  console.log(`${color.green}bun make init fork${color.reset}`)
  process.exit(2)
}

async function newProject() {
  await initProject()
  console.log(`\n${color.magenta}Iniciando nuevo proyecto...${color.reset}`)

  // Eliminar historial de Git
  await runc(["rm", "-rf", ".git"], "📁 Eliminando historial de Git existente")

  // Instalar dependencias
  await runc(["bun", "install"], "📦 Instalando dependencias")

  // Inicializar nuevo repositorio Git
  await runc(["git", "init"], "🔧 Inicializando nuevo repositorio Git")
  await runc(["git", "add", "."], "📦 Agregando archivos al staging")
  await runc(["git", "commit", "-m", "feat: initial commit from donbarrigon/new"], "💾 Realizando commit inicial")
  try {
    runc(["code", "."], "Abriendo editor")
  } catch (e) {
    console.error(`${color.red}✗ no tienes vs code instalado:${color.reset} `)
  }

  console.log(`${color.bold}${color.green}🎉 Proyecto inicializado!${color.reset}`)
}

async function fork() {
  await initProject()
  console.log(`\n${color.bold}Configurando fork${color.reset}\n`)
  // Instalar dependencias
  await runc(["bun", "install"], "📦 Instalando dependencias")

  // Renombrar origin a upstream
  await runc(["git", "remote", "rename", "origin", "upstream"], "🔄 Renombrando origin a upstream")

  // Commit
  await runc(["git", "add", "."], "📦 Agregando cambios al staging")
  await runc(["git", "commit", "-m", "feat: initial commit from donbarrigon/new"], "💾 Realizando commit inicial")

  try {
    runc(["code", "."], "Abriendo editor")
  } catch (e) {
    console.error(`${color.red}✗ no tienes vs code instalado:${color.reset} `)
  }
  console.log(`${color.bold}${color.green}🎉 Fork configurado.!${color.reset}`)
}

async function initProject() {
  const projectName = inputc("Nombre del proyecto:", "gituser/app-name")
  if (!validateProjectName(projectName)) {
    console.error(
      `${color.red}Formato incorrecto. Usa el formato: usuarioDeGit/nombreDelProyecto sin espacios${color.reset}`,
    )
    process.exit(2)
  }

  if (!existsSync("package.json")) {
    console.error(`${color.red}✗ No se encontró package.json en el directorio actual${color.reset}`)
    process.exit(1)
  }

  console.log(`\n${color.bold}Iniciando configuración del proyecto: ${projectName}${color.reset}\n`)

  const projectNameOnly = projectName.split("/")[1]
  if (!projectNameOnly) {
    console.error(
      `${color.red}Error al obtener el nombre del proyecto. Usa el formato: usuarioDeGit/nombreDelProyecto${color.reset}`,
    )
    process.exit(2)
  }
  // Actualizar package.json
  updatePackageJson(projectNameOnly)

  // Actualizar go.mod
  await updateModule(projectName)

  // Renombrar la carpeta proyecto
  renameProjectDir(projectNameOnly)
}

function validateProjectName(name: string) {
  const regex = /^[a-zA-Z0-9._-]+\/[a-zA-Z0-9._-]+$/
  return regex.test(name)
}

function updatePackageJson(projectName: string) {
  try {
    const packageJson = JSON.parse(readFileSync("package.json", "utf8"))

    packageJson.name = projectName

    writeFileSync("package.json", JSON.stringify(packageJson, null, 2))
    console.log(`${color.green}✓ package.json actualizado - name: ${projectName}${color.reset}\n`)
  } catch (error) {
    console.error(`${color.red}✗ Error actualizando package.json: ${error}${color.reset}`)
    process.exit(1)
  }
}

export function renameProjectDir(projectName: string) {
  const currentDir = process.cwd() // Ej: /home/user/projects/new
  const parentDir = dirname(currentDir) // Ej: /home/user/projects
  const currentName = basename(currentDir) // Ej: new

  const newPath = join(parentDir, projectName)

  // Evita colisión si ya existe
  if (existsSync(newPath)) {
    console.error(`${color.red}✗ Ya existe una carpeta con el nombre${color.reset} '${projectName}'`)
    process.exit(1)
  }

  console.log(`🔄 Renombrando proyecto: '${currentName}' → '${projectName}'`)
  renameSync(currentDir, newPath)

  // Cambiar el directorio actual
  process.chdir(newPath)
  console.log(`✓ Proyecto renombrado y movido a '${newPath}'`)
}

/**
 * Actualiza el nombre del módulo en go.mod y todos los archivos .go
 * @param newModule - Nuevo nombre del módulo (ej: "github.com/usuario/proyecto")
 */
async function updateModule(newModule: string): Promise<void> {
  const oldModule = "donbarrigon/new"
  const oldImportPath = "donbarrigon/new/internal/"
  const newImportPath = `${newModule}/internal/`

  console.log(`🔄 Actualizando módulo de "${oldModule}" a "${newModule}"`)

  try {
    // 1. Actualizar go.mod
    console.log("📝 Actualizando go.mod...")
    const goModContent = readFileSync("go.mod", "utf-8")
    const updatedGoMod = goModContent.replace(`module ${oldModule}`, `module ${newModule}`)
    writeFileSync("go.mod", updatedGoMod, "utf-8")
    console.log(`${color.green}✓ go.mod actualizado${color.reset}`)

    // 2. Actualizar todos los archivos .go en internal/
    console.log("📝 Actualizando archivos .go en internal/...")
    const filesUpdated = updateGoFiles("internal", oldImportPath, newImportPath)
    console.log(`${color.green}✓ ${filesUpdated} archivos actualizados${color.reset}`)

    // 3. Ejecutar go mod tidy para limpiar dependencias
    await runc(["go", "mod", "tidy"], "🧹 Ejecutando go mod tidy")
    console.log(`${color.green}✓ Dependencias actualizadas${color.reset}`)

    console.log(
      `\n${color.green}✓ Módulo actualizado exitosamente a${color.reset} ${color.magenta}"${newModule}${color.reset}"`,
    )
  } catch (error) {
    console.error(`${color.red}✗ Error al actualizar el módulo:${color.reset}`, error)
    throw error
  }
}

/**
 * Actualiza recursivamente todos los archivos .go en un directorio
 */
function updateGoFiles(dir: string, oldPath: string, newPath: string): number {
  let count = 0
  const entries = readdirSync(dir, { withFileTypes: true })

  for (const entry of entries) {
    const fullPath = join(dir, entry.name)

    if (entry.isDirectory()) {
      // Recursión para subdirectorios
      count += updateGoFiles(fullPath, oldPath, newPath)
    } else if (entry.isFile() && entry.name.endsWith(".go")) {
      // Actualizar archivo .go
      const content = readFileSync(fullPath, "utf-8")
      if (content.includes(oldPath)) {
        const updatedContent = content.replaceAll(oldPath, newPath)
        writeFileSync(fullPath, updatedContent, "utf-8")
        count++
        console.log(`  ✓ ${fullPath}`)
      }
    }
  }

  return count
}
