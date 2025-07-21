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
│   ├── users/          # Serviço de usuários (porta 8081)
│   └── products/       # Serviço de produtos (porta 8082)
├── gateway/            # Apollo Gateway (porta 4000)
├── docs/
│   └── arquitecture.drawio
├── README.md
├── Makefile
├── env.example
└── go.work
```

## Serviço Users

### Funcionalidades

- **Query `users`**: Retorna todos os usuários
- **Query `user(id: ID!)`**: Retorna um usuário específico por ID
- **Dados mock**: Alice e Bob pré-cadastrados

### Como executar

#### Localmente:

```bash
cd gofed/services/users
go run main.go
```

#### Via Docker:

```bash
cd gofed/services/users
docker build -t gofed-users .
docker run -p 8081:8081 gofed-users
```

#### Via Makefile:

```bash
make run-users
```

### Testando o serviço

1. **Acesse o GraphQL Playground**: http://localhost:8081/
2. **Teste a query `users`**:
   ```graphql
   query {
     users {
       id
       name
       email
     }
   }
   ```
3. **Teste a query `user`**:
   ```graphql
   query {
     user(id: "1") {
       id
       name
       email
     }
   }
   ```

### Endpoints

- **GraphQL Playground**: http://localhost:8081/
- **GraphQL Query**: http://localhost:8081/query

## Próximos passos

- [ ] Implementar serviço products
- [ ] Configurar Apollo Gateway
- [ ] Implementar federação GraphQL
- [ ] Adicionar resoluções concorrentes
- [ ] Implementar cache e otimizações
