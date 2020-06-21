# Cliente GRPC con Golang

go-client-grpc es una aplicación realizado para el proyecto 2 del curso de Sistemas Operativos. Este cliente se conecta con el servidor [python-server-grpc](https://github.com/NeftXx/python-server-grpc)

## Construir imagen

Primero es loguearse con docker login, ingresando su usuario y contraseña de DockerHub:

```bash
docker login
```

Luego debemos construir la imagen, con el siguiente comando:

```bash
docker build -t go-client-grpc .
```

Luego tagueamos la imagen local con nuestro usuario y el repositorio de docker hub.

> Nota: \$1 es el tu nombre de usuario de docker hub

```bash
docker tag go-client-grpc $1/go-client-grpc
```

Por último, subimos la imagen al repositorio

```bash
docker push $1/go-client-grpc
```

Si deseas puedes usar el script [build-docker.sh](build-docker.sh) dandole permisos de ejecución, ejecutandoló y mandando de parámetro el nombre del usuario del DockerHub (Debes estar en la carpeta del Dockerfile).

```bash
chmod 777 build-docker.sh
build-docker.sh usuarioDockerHub
```

## Como usar la imagen

El siguiente ejemplo, crea un servidor en http://localhost:4000.

```bash
docker run -d \
  --name go-client \
  -p 4000:4000 \
  neftxx/go-client-grpc
```

### Configuración

Las variables de entorno se pasan al comando de ejecución para configurar un contenedor.

| Nombre   | Valor por defecto | Descripción             |
| -------- | ----------------- | ----------------------- |
| URL_GRPC | "localhost"       | Url de un servidor GRPC |

#### Ejemplo

```bash
docker run -d \
  --name go-client \
  -p 4000:4000 \
  -e URL_GRPC="IP:9000" \
  neftxx/go-client-grpc
```
