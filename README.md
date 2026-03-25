# myapp

Proyecto base en Go siguiendo los últimos estándares y buenas prácticas.

## Stack

| Área       | Tecnología                                |
| ---------- | ----------------------------------------- |
| Lenguaje   | Go 1.23                                   |
| HTTP       | `net/http` stdlib (Go 1.22+ enhanced mux) |
| Logging    | `log/slog` stdlib                         |
| Config     | Variables de entorno + `.env`             |
| Contenedor | Docker multi-stage + distroless           |
| CI/CD      | GitHub Actions                            |

## Estructura del proyecto

```
.
├── cmd/
│   └── api/
│       └── main.go            # Punto de entrada
├── internal/
│   ├── config/
│   │   └── config.go          # Configuración desde env vars
│   ├── handler/
│   │   ├── handler.go         # Base de handlers + registro de rutas
│   │   ├── health.go          # Handlers de health/readiness/info
│   │   └── health_test.go     # Tests de handlers
│   ├── middleware/
│   │   └── middleware.go      # RequestID, Logger, Recover
│   └── server/
│       └── server.go          # Servidor HTTP + graceful shutdown
├── pkg/
│   └── logger/
│       └── logger.go          # Wrapper de slog
├── .air.toml                  # Hot-reload (air)
├── .env.example               # Variables de entorno de ejemplo
├── .gitignore
├── .golangci.yml              # Configuración de linter
├── .github/
│   └── workflows/
│       └── ci.yml             # Test → Lint → Build
├── docker-compose.yml
├── Dockerfile                 # Multi-stage + distroless
├── go.mod
├── Makefile
└── README.md
```

## Inicio rápido

```bash
# 1. Copiar variables de entorno
cp .env.example .env

# 2. Descargar dependencias
make tidy

# 3. Ejecutar
make run
```

El servidor arranca en `http://localhost:8080`.

## Endpoints

| Método | Ruta                     | Descripción                  |
| ------ | ------------------------ | ---------------------------- |
| `GET`  | `/health`                | Liveness probe               |
| `GET`  | `/ready`                 | Readiness probe              |
| `GET`  | `/api/v1/info`           | Información de la aplicación |
| `GET`  | `/`                      | Frontend HTML Pokédex        |
| `GET`  | `/api/v1/pokemon`        | Lista paginada de Pokémon    |
| `GET`  | `/api/v1/pokemon/{name}` | Detalle de un Pokémon        |

## Comandos disponibles

```bash
make help          # Listar todos los targets
make build         # Compilar binario en ./bin/
make run           # Ejecutar con go run
make dev           # Hot-reload con air
make test          # Ejecutar tests con race detector
make test-cover    # Tests + reporte de cobertura HTML
make lint          # Ejecutar golangci-lint
make tidy          # go mod tidy + verify
make clean         # Limpiar artefactos
make docker-build  # Construir imagen Docker
make docker-up     # Levantar con Docker Compose
```

## Variables de entorno

| Variable                  | Valor por defecto | Descripción                                  |
| ------------------------- | ----------------- | -------------------------------------------- |
| `APP_NAME`                | `myapp`           | Nombre de la aplicación                      |
| `APP_ENV`                 | `development`     | Entorno (`development` / `production`)       |
| `APP_VERSION`             | `0.1.0`           | Versión                                      |
| `SERVER_HOST`             | `0.0.0.0`         | Host del servidor                            |
| `SERVER_PORT`             | `8080`            | Puerto del servidor                          |
| `SERVER_READ_TIMEOUT`     | `10s`             | Timeout de lectura                           |
| `SERVER_WRITE_TIMEOUT`    | `10s`             | Timeout de escritura                         |
| `SERVER_SHUTDOWN_TIMEOUT` | `30s`             | Timeout de graceful shutdown                 |
| `LOG_LEVEL`               | `info`            | Nivel de log (`debug`/`info`/`warn`/`error`) |

## Herramientas opcionales

```bash
# Hot-reload
go install github.com/air-verse/air@latest

# Linter
brew install golangci-lint
```

## Buenas prácticas incluidas

- **Arquitectura**: separación `cmd` / `internal` / `pkg`
- **Configuración**: solo variables de entorno, sin hardcoding
- **Logging**: `log/slog` estructurado (JSON en prod, texto en dev)
- **Middleware**: RequestID, Logger, Recover (anti-panic)
- **Graceful shutdown**: espera conexiones activas antes de cerrar
- **Docker**: imagen multi-stage + distroless (mínima superficie de ataque)
- **Tests**: race detector activado por defecto
- **CI**: pipeline de test → lint → build en GitHub Actions
