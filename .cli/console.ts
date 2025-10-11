import { color } from "./config"

export async function runc(command: string, description: string) {
  try {
    if (description) {
      console.log(`${color.blue}${description}${color.reset}`)
    }
    console.log(`${color.cyan}Ejecutando: ${color.green}${command}${color.reset}`)

    const output = await Bun.$`${command.split(" ")}`.text()

    console.log(`${color.green}✓${color.reset} ${output}\n`)
  } catch (error) {
    console.error(`${color.red}✗ ${command}:${color.reset}`)
    console.error(`${color.yellow}${error}${color.reset}`)
    process.exit(1)
  }
}

export function inputc(question: string): string {
  const answer = prompt(`${color.bold}${color.cyan}${question}${color.reset}`, "")
  if (!answer) {
    return ""
  }
  return answer
}
