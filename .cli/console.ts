import { color } from "./config"

export async function runc(command: string[], description: string) {
  try {
    if (description) {
      console.log(`${color.blue}${description}${color.reset}`)
    }
    console.log(`${color.cyan}Ejecutando: ${color.green}${command}${color.reset}`)

    const proc = Bun.spawn(command)
    const output = await new Response(proc.stdout).text()

    console.log(`${color.green}✓${color.reset} ${output}\n`)
  } catch (error) {
    console.error(`${color.red}✗ ${command}:${color.reset}`)
    console.error(`${color.yellow}${error}${color.reset}`)
    process.exit(1)
  }
}

export function inputc(question: string, defaultValue: string = ""): string {
  const answer = prompt(`${color.bold}${color.cyan}${question}${color.reset}`, defaultValue)
  if (!answer) {
    return ""
  }
  return answer
}
