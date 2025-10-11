// configuracion de variables

// code: ubicacion de los archivos de codigo
// dev: ubicacion de los archivos compilados para el desarrollo
// build:
export const dir = {
  code: {
    ts: "internal/ui/ts",
    js: "internal/ui/ts",
    css: "internal/ui/css",
    wasm: "internal/ui/wasm",
    qtpl: "internal/ui/pages",
  },
  dev: {
    js: "public/js",
    ts: "public/js",
    css: "public/css",
    wasm: "public/wasm",
    qtpl: "internal/ui/view",
  },
  build: {
    js: "dist/public/js",
    ts: "dist/public/js",
    css: "dist/public/css",
    wasm: "dist/public/wasm",
    qtpl: "internal/ui/view",
  },
}

// Colores para la consola
export const color = {
  red: "\x1b[31m",
  green: "\x1b[32m",
  yellow: "\x1b[33m",
  blue: "\x1b[34m",
  magenta: "\x1b[35m",
  cyan: "\x1b[36m",
  reset: "\x1b[0m",
  bold: "\x1b[1m",
}
