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
- Logrus (logging estruturado)

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
- **Health Check**: Endpoint `/healthz` para monitoramento
- **Logging estruturado**: Logs em JSON com contexto completo

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
4. **Teste o Health Check**:
   ```bash
   curl http://localhost:8081/healthz
   ```

### Endpoints

- **GraphQL Playground**: http://localhost:8081/
- **GraphQL Query**: http://localhost:8081/query
- **Health Check**: http://localhost:8081/healthz

### Logging

O serviço utiliza logging estruturado em JSON com os seguintes campos:

- `timestamp`: Timestamp ISO 8601
- `level`: Nível do log (info, warn, error, fatal)
- `message`: Mensagem do log
- `service`: Nome do serviço
- `version`: Versão do serviço
- `method`: Método HTTP (para requisições)
- `path`: Caminho da requisição
- `status_code`: Código de status HTTP
- `duration`: Duração da requisição

### Variáveis de Ambiente

```bash
# Nível de log (debug, info, warn, error, fatal)
LOG_LEVEL=info
```

## Próximos passos

- [ ] Implementar serviço products
- [ ] Configurar Apollo Gateway
- [ ] Implementar federação GraphQL
- [ ] Adicionar resoluções concorrentes
- [ ] Implementar cache e otimizações
