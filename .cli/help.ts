import { color } from "./config.ts"
export function showHelp() {
  console.log(`${color.bold}${color.cyan}make es el asistente para el proyecto${color.reset}

${color.bold}esta herramienta está inspirada en artisan de Laravel.${color.reset}
si estas viendo esto es por que el comando no existe o quieres ver el help

${color.bold}sintaxis:${color.reset}
bun make <comando> [argumentos]

${color.bold}Comandos disponibles:${color.reset}

${color.magenta}=== comandos para iniciar el proyecto ===${color.reset}
${color.green}init${color.reset}       
                Inicializa o renombra el proyecto.
                    -- elimina el historial de Git existente
                    -- instala las dependencias
                    -- renombra el package.json y el go.mod
                    -- renombra los imports de Go
                    -- inicializa un nuevo repositorio Git
                    -- realiza un commit inicial

${color.green}init fork${color.reset}  
                Igual que ${color.green}init${color.reset}, pero además configura un fork del proyecto actual.
                    -- realiza los mismos pasos que init
                    -- agrega un nuevo remoto (fork) en Git
                    -- conserva el historial del repositorio original

${color.green}merge upstream${color.reset}  
                Sincroniza el proyecto con el repositorio upstream.
                    -- obtiene los últimos cambios del repositorio original
                    -- fusiona las actualizaciones en la rama actual


${color.magenta}=== comandos para correr el proyecto ===${color.reset}

${color.green}bun run dev${color.reset}
                Corre el servidor de desarrollo.
                    -- compila las diferentes partes del proyecto
                    -- escucha los cambios y reinicia el servidor

${color.green}bun run build${color.reset}
                Compila el proyecto en la carpeta dist.


${color.magenta}=== comandos para desarrollo ===${color.reset}

${color.green}make${color.reset}
                Crea un nuevo recurso en el proyecto basado en una plantilla.
                ${color.red}Sintaxis:${color.reset}
                    ${color.green}bun make${color.reset} <${color.blue}recurso${color.reset}> <${color.yellow}nombre${color.reset}>

${color.green}recursos:${color.reset}

    ${color.blue}model${color.reset}         -- Define la estructura de datos principal y sus métodos asociados.
    ${color.blue}migration${color.reset}     -- Crea un archivo de migración para la base de datos.
    ${color.blue}repository${color.reset}    -- Implementa la capa de acceso a datos (DAO) del modelo.
    ${color.blue}resource${color.reset}      -- Genera un serializador para formatear la salida de datos (JSON/API).
    ${color.blue}seed${color.reset}          -- Crea un seeder para poblar datos iniciales o de prueba.

    ${color.blue}view${color.reset}          -- Genera una plantilla HTML basada en QuickTemplate.
    ${color.blue}page${color.reset}          -- Crea una página completa (vista con layout y controladores asociados).
    ${color.blue}component${color.reset}     -- Define un componente reutilizable de interfaz (HTML/TS/CSS).
    ${color.blue}ts${color.reset}            -- Crea un módulo TypeScript asociado al recurso.
    ${color.blue}js${color.reset}            -- Genera un script JavaScript estándar.
    ${color.blue}css${color.reset}           -- Crea una hoja de estilos base o específica del recurso.
    ${color.blue}wasm${color.reset}          -- Inicializa un módulo WebAssembly integrado con Go.

    ${color.blue}controller${color.reset}    -- Define la lógica HTTP del recurso (Create, Read, Update, Delete).
    ${color.blue}middleware${color.reset}    -- Crea un middleware para interceptar y procesar peticiones.
    ${color.blue}policy${color.reset}        -- Define las reglas de autorización y control de acceso.
    ${color.blue}route${color.reset}         -- Registra rutas HTTP y las asocia con controladores.
    ${color.blue}service${color.reset}       -- Implementa la lógica de negocio independiente del controlador.
    ${color.blue}validator${color.reset}     -- Define las reglas de validación para entradas de usuario o datos.


    ${color.blue}db${color.reset}          Crea todos los recursos de la base de datos
                  -- model
                  -- migration
                  -- seed

    ${color.blue}handler${color.reset}     Crea todos los recursos del handler
                  -- controller
                  -- middleware
                  -- policy
                  -- route
                  -- service
                  -- validator

    ${color.blue}ui${color.reset}          Crea todos los recursos de la interfaz
                  -- page
                  -- ts
                  -- css

    ${color.blue}mvc${color.reset}         Crea todos los recursos de la arquitectura MVC

    
${color.yellow}nombre${color.reset}
                ${color.red}Reglas de sintaxis${color.reset}:
                    -- usa snake_case en singular y minusculas sin acentos
                      Ej: ${color.yellow}user_profile${color.reset}  
                          -- model:      UserProfile
                          -- colleccion: user_profiles
                          -- controller: UserProfileCreate
                          -- view:       WriteUserProfilePage o StreamUserProfilePage
                          -- route:      /users/profiles, /users/profiles/create, /users/profiles/store ...
                    -- si requires una letras en mayusculas el las respeta colocala que la respeta
                    -- OjO los numeros y mayusculas no los pluraliza

                    -- para raritos que trabajan con DDD + MVC, pueden usar el punto (.)
                       la sintaxis es: <dominio>.<nombre>
                       te crea en internal una carpeta con el dominio y la estructura mvc dentro de ella
                       para que tabajes con el dominio aislado de los demas
                       Ej: ${color.yellow}billing.invoice${color.reset}
                           -- model:      Invoice
                           -- colección:  invoices
                           -- controller: InvoiceEdit, InvoiceUpdate
                           -- view:       WriteInvoicePage
                           -- route:      /billing/invoices, /billing/invoices/create ...

${color.red}Ejemplo${color.reset}
${color.yellow}bun make controller route movie_category${color.reset}


${color.red}Ejemplo usando dominio${color.reset}
${color.yellow}bun make model migration seed dashboard.movie_category${color.reset}
crea los archivos
/internal/dashboard/database/model/movie_category.go 
/internal/dashboard/database/migration/20250101_create_movie_categories_collection.go 
/internal/dashboard/database/seed/movie_categories.go
y en los templates
// model:     type MovieCategory
// migracion: coleccion de movie_categories func MovieCategoriesUp() func MovieCategoriesDown()
// seed:      func MovieCategories()
`)
}
