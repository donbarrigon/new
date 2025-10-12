import { spawn } from "child_process"
import { existsSync, watch } from "fs"
import { copyFile, mkdir, readdir, readFile, writeFile } from "fs/promises"
import { extname, join } from "path"
import { color, dir } from "./config"

export async function devMode() {
  await setupDirectories([dir.dev.js, dir.dev.css, dir.dev.wasm, dir.dev.qtpl])
  await compileCSS(true)
  await compileScripts(true)
}

// funcion para crear las carpetas antes de compilar
async function setupDirectories(dirs: string[]) {
  for (const dir of dirs) {
    if (!existsSync(dir)) {
      await mkdir(dir, { recursive: true })
    }
  }
}

async function compileCSS(isDev: boolean = true) {
  try {
    if (!existsSync(dir.code.css)) {
      console.log(`⚠️  ${color.red}Directorio ${dir.code.css} no existe${color.red}`)
      return
    }

    const files = await getFilesRecursively(dir.code.css, [".css"])

    for (const file of files) {
      const relativePath = file.replace(`${dir.code.css}/`, "")
      const outputPath = join(isDev ? dir.dev.css : dir.build.css, relativePath)
      const outputDir = outputPath.substring(0, outputPath.lastIndexOf("/"))

      if (!existsSync(outputDir)) {
        await mkdir(outputDir, { recursive: true })
      }

      let content = await readFile(file, "utf-8")

      if (!isDev) {
        // Minificar CSS en producción
        content = content
          .replace(/\/\*[\s\S]*?\*\//g, "") // Remover comentarios
          .replace(/\s+/g, " ") // Colapsar espacios
          .replace(/;\s*}/g, "}") // Remover último semicolon
          .replace(/\s*{\s*/g, "{") // Limpiar llaves
          .replace(/;\s*/g, ";") // Limpiar semicolons
          .trim()
      }

      await writeFile(outputPath, content)
      console.log(`${color.green}✓ Copiado: ${file} -> ${color.reset}${outputPath}`)
    }
  } catch (error) {
    console.error(`${color.red}✗ Error copiando CSS:${color.reset}`, error)
  }
}

async function compileScripts(isDev: boolean = true) {
  try {
    if (!existsSync(dir.code.js)) {
      console.log(`⚠️  ${color.yellow}Directorio ${dir.code.js} no existe${color.reset}`)
      return
    }

    const files = await getFilesRecursively(dir.code.js, [".js", ".ts"])

    for (const file of files) {
      const relativePath = file.replace(`${dir.code.js}/`, "")
      const outputPath = join(isDev ? dir.dev.js : dir.build.js, relativePath.replace(/\.[jt]s$/, ".js"))
      const outputDir = outputPath.substring(0, outputPath.lastIndexOf("/"))

      if (!existsSync(outputDir)) {
        await mkdir(outputDir, { recursive: true })
      }

      const result = await Bun.build({
        entrypoints: [file],
        target: "browser",
        minify: !isDev,
        sourcemap: isDev ? "inline" : "none",
        outdir: isDev ? dir.dev.js : dir.build.js,
        naming: relativePath.replace(/\.[jt]s$/, ".js"),
      })

      if (result.success) {
        console.log(`${color.green}✓ Compilado: ${file} -> ${color.reset}${outputPath}`)
      } else {
        console.error(`${color.red}✗ Error compilando ${file}:${color.reset}`, result.logs)
      }
    }
  } catch (error) {
    console.error(`${color.red}✗ Error en compilación de scripts:${color.reset}`, error)
  }
}

// funcion optener
async function getFilesRecursively(dir: string, extensions: string[]): Promise<string[]> {
  const files = []
  const extensionsArray = Array.isArray(extensions) ? extensions : [extensions]

  try {
    const items = await readdir(dir, { withFileTypes: true })

    for (const item of items) {
      const fullPath = join(dir, item.name)

      if (item.isDirectory()) {
        const subFiles = await getFilesRecursively(fullPath, extensions)
        files.push(...subFiles)
      } else if (extensionsArray.some((ext) => item.name.endsWith(ext))) {
        files.push(fullPath)
      }
    }
  } catch (error) {
    console.warn(`${color.yellow}Directorio no existe [${dir}] ${extensions}${color.reset}`)
    console.warn(error)
  }

  return files
}
