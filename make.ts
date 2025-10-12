import { showHelp } from "./.cli/help"
import { init } from "./.cli/init"

// make es el asistente cli para el proyecto
// usa la carpeta .cli para los comandos
// la carpeta .cli se llama asi por que la extencion de vscode Meterial Icons Theme le pone un icono lindo
// si se pregunta por que ise el cli en ts y no en go es por que me permite modificalo en cualquier momento sin compilar.
// ademas lo puedo personalizar facil para cualquier proyecto.
// y tambien es mas facil ya que bun trae compilador transpilador y modo de desarrollo integrado.

const command = process.argv[2]

if (!command) {
  showHelp()
  process.exit(0)
} else if (command === "help" || command === "h" || command === "-h" || command === "--help") {
  showHelp()
  process.exit(0)
} else if (command === "version" || command === "v" || command === "-v" || command === "--version") {
  console.log("version 0.8.1")
  process.exit(0)
} else if (command === "init") {
  await init()
  process.exit(0)
} else if (command === "dev") {
  process.exit(0)
}
