# tomlights

Go backend (go-atlas), MongoDB, email sender (SMTP / Resend), Vue (Nuxt 3) frontend.

## Run locally

```bash
cp .env.example .env
cp web/.env.example web/.env

make mongo-up      # MongoDB in docker
make run           # backend: gRPC :7779, HTTP :9080
make web-install
make web-dev       # http://localhost:3000
```

## Run in Docker

```bash
cp .env.example .env
make compose-up    # http://localhost:3000
```

Configuration lives in `.env` (see `.env.example`) and `configs/*.yaml`. Run `make help` to see all commands.

Development rules: [docs/SERVICE_DEVELOPMENT_STANDARD.md](docs/SERVICE_DEVELOPMENT_STANDARD.md), [docs/PROTO_DEVELOPMENT_STANDARD.md](docs/PROTO_DEVELOPMENT_STANDARD.md).
