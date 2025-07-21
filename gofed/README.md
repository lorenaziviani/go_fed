# Gofed

Gofed é uma aplicação demonstrativa que simula um ambiente federado com microsserviços em Go utilizando GraphQL. Cada microsserviço possui seu próprio schema e expõe parte do domínio (ex: usuários, produtos). A federação é feita via Apollo Gateway (Node.js).

## Objetivo

Demonstrar padrões de microsserviços federados, resoluções concorrentes em Go (WaitGroup, context.Context, canais), simulação de problemas de performance e mitigação com paralelismo e cache.

## Tech Stack

- Go 1.18+
- GraphQL (gqlgen)
- Apollo Gateway (Node.js)
- Docker (opcional)
- Draw.io (diagramas)

## Estrutura

```
gofed/
├── services/
│   ├── users/
│   └── products/
├── gateway/
├── docs/
│   └── diagrama.drawio
├── README.md
├── Makefile
├── .env.example
└── go.work
```

## Serviço users

Para rodar localmente:

```bash
cd gofed/services/users
go run main.go
```

Para rodar via Docker:

```bash
cd gofed/services/users
docker build -t gofed-users .
docker run -p 8081:8081 gofed-users
```
