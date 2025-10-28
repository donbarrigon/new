# new 0.8.1

Este es un punto de partida pa’ arrancar mis proyectos con **Go** y **MongoDB**, usando una estructura **MVC** monolitica con una ui **SSR** o una **API**, todo montado como a mí me gusta: sencillo y sin exceso de dependencias.

La idea es que no me toque arrancar de cero cada vez que vaya a montar un proyecto, sino que ya tenga la base lista pa’ empezar a meterle funcionalidad de una. 🚀

este proyecto nacio por que cuando empece a aprender go, y no me gustaron los frameworks que habia en el momento, todos inspirados en algo que no me gusta "express" no es que sean malos solo que me gustan los frameworks con structura.

## 📥 Instalación

### 🚀 Iniciar el proyecto (recomendado)

descarga el proyecto e instala las dependencias

```bash
make project
```

## 📥 Recursos nesesarios

### ***Go***
El lenguaje de programación backend.
Sitio oficial: [https://go.dev/](https://go.dev/)

### ***Bun***
runtime/entorno JavaScript / TypeScript.
Sitio oficial: [https://bun.sh/](https://bun.sh/)
debe instalarlo para que funcione el cli.

### ***Alpine.js***
Su nuevo y ligero marco de JavaScript.
Simple, ligero y poderoso como el infierno.
Sitio oficial: [https://alpinejs.dev/](https://alpinejs.dev/)
**exelente motor de plantillas**

### ***Mongo db***
motor de base de datos
Sitio oficial: [https://www.mongodb.com/](https://www.mongodb.com/)

### ***make*** 
Es el asistente cli para facilitar el trabajo en el proyecto.
- Compila el ts, js y las Templates.
- Minifica el css y js en produccion.
- Crea codigo base para el backend y templates para el front.
- Tiene modo desarrollador que compila y reinicia el server al realizar cambios.



## 🛠️ Comandos de desarrollo

Inicia el server en modo desarrollo:

```bash
make run
```

Construir para producción:

```bash
make build
```


## 🧩 Comandos para las templates

usa snake case en minuscula para nombrar el modelo o o entidad
Sintaxis make <template> [<dominio>].<nombre>

```bash
make model mi_entidad
make migration mi_entidad
make repository mi_entidad
make resource mi_entidad
make seed mi_entidad

make view mi_entidad
make page mi_entidad
make component mi_entidad
make ts mi_entidad
make js mi_entidad
make css mi_entidad
make wasm mi_entidad

make controller mi_entidad
make middleware mi_entidad
make policy mi_entidad
make route mi_entidad
make service mi_entidad
make validator mi_entidad
```

# 🧩 Comandos para multiples templates

```bash
make model seed migration bill
```
puede combinarlos y crear tantos como nesesite de una sola vez
Sintaxis make <template> <template> <template> [<dominio>].<nombre>


```bash
make model seed migration dashboard.bill
```
o dentro del dominio

# 🧩 Multiples templates segun responsabilidad

```bash
make db mi_entidad
```
model, seed y migration

```bash
make ui mi_entidad
```
model, seed y migration

```bash
make handler mi_entidad
```
controller, middleware, policy, route, service y validator

```bash
make mvc mi_entidad
```
crea todo


Este proyecto fue creado con ❤️ por Don Barrigon
Distribuido bajo la [MIT License](./LICENSE).

This project uses `bun init` in bun v1.2.21. [Bun](https://bun.com) is a fast all-in-one JavaScript runtime.
