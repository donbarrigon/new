# new

Este es un punto de partida pa’ arrancar mis proyectos con **Go** y **MongoDB**, usando una estructura **MVC** bien organizada, con **SSR**, todo montado como a mí me gusta: sencillo y sin exceso de dependencias.

La idea es que no me toque arrancar de cero cada vez que vaya a montar un proyecto, sino que ya tenga la base lista pa’ empezar a meterle funcionalidad de una. 🚀

este proyecto nacio por que cuando empece a aprender go no me gustaron los frameworks que habia en el momento.
Todos inspirados en algo que no me gusta express no es que sean malos solo que tengo gustos exoticos.

## 📥 Instalación

### 🚀 Iniciar el proyecto (recomendado)

descarga el proyecto e instala las dependencias

```bash
git clone https://github.com/donbarrigon/new.git
cd new
bun make init
```

### 🌀 Iniciar el proyecto como fork (opcional)

crea un fork al proyecto original donbarrigon/new
descarga el proyecto e instala las dependencias

```bash
git clone https://github.com/donbarrigon/new.git
cd new
bun make init fork
```
Si iniciaste el proyecto como un fork, puedes descargar actualizaciones con:

```bash
bun helper merge upstream
```


## 📥 Recursos nesesarios

### ***Go***
El lenguaje de programación backend.
Sitio oficial: [https://go.dev/](https://go.dev/)

### ***Bun***
runtime/entorno JavaScript / TypeScript.
Sitio oficial: [https://bun.sh/](https://bun.sh/)
debe instalarlo para que funcione el cli

### ***Quicktemplate***
motor de plantillas para Go.
Repositorio oficial GitHub: [https://github.com/valyala/quicktemplate](https://github.com/valyala/quicktemplate)
**exelente motor de plantillas**

### ***Mongo db***
motor de base de datos
Sitio oficial: [https://www.mongodb.com/](https://www.mongodb.com/)

### ***make*** 
Es el asistente cli para facilitar el trabajo en el proyecto.
- Compila el ts, js con y las Quicktemplate.
- Minifica el css y js en produccion.
- Crea las templates.
- Tiene modo desarrollador que compila y reinicia el server al realizar cambios.


## 🛠️ Comandos de desarrollo

Inicia el server en modo desarrollo:

```bash
bun run dev
```

Construir para producción:

```bash
bun run build
```


## 🧩 Comandos para las templates

usa snake case en minuscula para nombrar el modelo o o entidad
Sintaxis bun make <template> [<dominio>].<nombre>

```bash
bun make model mi_entidad
bun make migration mi_entidad
bun make repository mi_entidad
bun make resource mi_entidad
bun make seed mi_entidad

bun make view mi_entidad
bun make page mi_entidad
bun make component mi_entidad
bun make ts mi_entidad
bun make js mi_entidad
bun make css mi_entidad
bun make wasm mi_entidad

bun make controller mi_entidad
bun make middleware mi_entidad
bun make policy mi_entidad
bun make route mi_entidad
bun make service mi_entidad
bun make validator mi_entidad
```

# 🧩 Comandos para multiples templates

```bash
bun make model seed migration bill
```
puede combinarlos y crear tantos como nesesite de una sola vez
Sintaxis bun make <template> <template> <template> [<dominio>].<nombre>


```bash
bun make model seed migration dashboard.bill
```
o dentro del dominio

# 🧩 Multiples templates segun responsabilidad

```bash
bun make db mi_entidad
```
model, seed y migration

```bash
bun make ui mi_entidad
```
model, seed y migration

```bash
bun make handler mi_entidad
```
controller, middleware, policy, route, service y validator

```bash
bun make mvc mi_entidad
```
crea todo


Creado con ❤️ por Don Barrigon
Distribuido bajo la [MIT License](./LICENSE).

This project was created using `bun init` in bun v1.2.21. [Bun](https://bun.com) is a fast all-in-one JavaScript runtime.
