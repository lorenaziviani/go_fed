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
├── docker-compose.yml
└── go.work
```

## Serviços

### Users Service (Porta 8081)

#### Funcionalidades

- **Query `users`**: Retorna todos os usuários
- **Query `user(id: ID!)`**: Retorna um usuário específico por ID
- **Dados mock**: Alice e Bob pré-cadastrados
- **Health Check**: Endpoint `/healthz` para monitoramento
- **Logging estruturado**: Logs em JSON com contexto completo
- **Federation Support**: Diretiva `@key(fields: "id")` e `__resolveReference`

#### Como executar

```bash
# Localmente
cd gofed/services/users
go run main.go

# Via Docker
cd gofed/services/users
docker build -t gofed-users .
docker run -p 8081:8081 gofed-users

# Via Makefile
make run-users
```

#### Testando o serviço

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

### Products Service (Porta 8082)

#### Funcionalidades

- **Query `products`**: Retorna todos os produtos
- **Query `product(id: ID!)`**: Retorna um produto específico por ID
- **Dados mock**: iPhone, MacBook, Nike, Coffee Maker
- **Health Check**: Endpoint `/healthz` para monitoramento
- **Logging estruturado**: Logs em JSON com contexto completo
- **Federation Support**: Diretiva `@key(fields: "id")` e `__resolveReference`

#### Como executar

```bash
# Localmente
cd gofed/services/products
go run main.go

# Via Docker
cd gofed/services/products
docker build -t gofed-products .
docker run -p 8082:8082 gofed-products

# Via Makefile
make run-products
```

#### Testando o serviço

1. **Acesse o GraphQL Playground**: http://localhost:8082/
2. **Teste a query `products`**:
   ```graphql
   query {
     products {
       id
       name
       description
       price
       category
     }
   }
   ```
3. **Teste a query `product`**:
   ```graphql
   query {
     product(id: "1") {
       id
       name
       description
       price
       category
     }
   }
   ```
4. **Teste o Health Check**:
   ```bash
   curl http://localhost:8082/healthz
   ```

### Apollo Gateway (Porta 4000)

#### Funcionalidades

- **Federation**: Combina schemas dos serviços users e products
- **Proxy**: Faz proxy para os serviços individuais
- **GraphQL Playground**: Interface para testar queries federadas
- **Configurável**: URLs dos serviços via variáveis de ambiente
- **Docker Compose**: Integração com todos os serviços

#### Como executar

```bash
# Localmente (com serviços rodando)
cd gofed/gateway
npm install
npm start

# Via Docker
cd gofed/gateway
docker build -t gofed-gateway .
docker run -p 4000:4000 gofed-gateway

# Via Docker Compose (todos os serviços)
cd gofed
docker-compose up
```

#### Testando o gateway

1. **Acesse o GraphQL Playground**: http://localhost:4000/
2. **Teste queries federadas**:
   ```graphql
   query {
     users {
       id
       name
       email
     }
     products {
       id
       name
       price
     }
   }
   ```
3. **Teste queries individuais**:
   ```graphql
   query {
     user(id: "1") {
       id
       name
     }
     product(id: "1") {
       id
       name
       price
     }
   }
   ```

## Federation Support

Ambos os serviços estão preparados para federation com Apollo Gateway:

### Users Service

```graphql
type User @key(fields: "id") {
  id: ID!
  name: String!
  email: String!
}
```

### Products Service

```graphql
type Product @key(fields: "id") {
  id: ID!
  name: String!
  description: String!
  price: Float!
  category: String!
}
```

### \_\_resolveReference

Cada serviço implementa a função `__resolveReference` que permite ao Apollo Gateway resolver referências federadas:

- **Users**: Resolve referências por `id` do usuário
- **Products**: Resolve referências por `id` do produto

### Gateway Federation

O Apollo Gateway combina os schemas e faz proxy para os serviços:

- **Schema Federado**: Combina User e Product types
- **Proxy Resolvers**: Encaminha queries para os serviços corretos
- **Configuração**: URLs configuráveis via variáveis de ambiente

## Endpoints

### Users Service

- **GraphQL Playground**: http://localhost:8081/
- **GraphQL Query**: http://localhost:8081/query
- **Health Check**: http://localhost:8081/healthz

### Products Service

- **GraphQL Playground**: http://localhost:8082/
- **GraphQL Query**: http://localhost:8082/query
- **Health Check**: http://localhost:8082/healthz

### Apollo Gateway

- **GraphQL Playground**: http://localhost:4000/
- **GraphQL Query**: http://localhost:4000/

## Logging

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

## Variáveis de Ambiente

```bash
# Users Service
USERS_SERVICE_PORT=8081
USERS_SERVICE_HOST=localhost

# Products Service
PRODUCTS_SERVICE_PORT=8082
PRODUCTS_SERVICE_HOST=localhost

# Gateway
GATEWAY_PORT=4000
GATEWAY_HOST=localhost
USERS_SERVICE_URL=http://localhost:8081/query
PRODUCTS_SERVICE_URL=http://localhost:8082/query

# Logging
LOG_LEVEL=info

# Federation
FEDERATION_ENABLED=true
FEDERATION_VERSION=2
```

## Docker Compose

Para executar todos os serviços juntos:

```bash
cd gofed
docker-compose up
```

Isso irá:

- Construir e executar o Users Service na porta 8081
- Construir e executar o Products Service na porta 8082
- Construir e executar o Apollo Gateway na porta 4000
- Configurar a rede entre os serviços
- Usar as variáveis de ambiente corretas para cada serviço

## Próximos passos

- [ ] Adicionar resoluções concorrentes
- [ ] Implementar cache e otimizações
- [ ] Adicionar mais serviços (orders, reviews, etc.)
- [ ] Implementar autenticação e autorização
