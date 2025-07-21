# 🚀 Gofed - Federation de Microsserviços com GraphQL

O **Gofed** é uma aplicação demonstrativa que simula um ambiente federado com microsserviços em Go usando GraphQL. Cada microsserviço possui seu schema e expõe parte do domínio (ex: usuários, produtos). A federação é feita via Apollo Gateway.

## 🎯 Objetivo

Demonstrar a implementação de **GraphQL Federation** com microsserviços em Go, incluindo:

- ✅ **Resoluções concorrentes** com WaitGroup, context.Context, canais
- ✅ **Benchmarks e race detection** para validação de performance e segurança
- Simulação de problemas de performance mitigados com paralelismo e cache
- Federation com Apollo Gateway e diretivas `@key`

## 🛠️ Tech Stack

- **Go 1.24.3**: Linguagem principal para microsserviços
- **GraphQL**: API query language
- **gqlgen**: Biblioteca Go para GraphQL
- **Apollo Gateway (Node.js)**: Para GraphQL federation
- **Docker & Docker Compose**: Containerização e orquestração
- **Federation v2.0**: Com diretivas `@key` para referências cruzadas
- **Concurrency Patterns**: WaitGroup, Channels, Context
- **Testing**: Race detection, Benchmarks, Unit tests

## 📁 Estrutura do Monorepo

```
gofed/
├── services/
│   ├── users/          # Serviço de usuários (porta 8081)
│   └── products/       # Serviço de produtos (porta 8082)
├── gateway/            # Apollo Federation Gateway (porta 4000)
├── docs/              # Documentação e diagramas
├── examples/          # Exemplos de queries GraphQL
├── scripts/           # Scripts de teste e automação
├── docker-compose.yml # Orquestração dos serviços
├── Makefile          # Comandos de automação
└── env.example       # Variáveis de ambiente
```

## 🚀 Como Executar

### Opção 1: Docker Compose (Recomendado)

```bash
# Construir e subir todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar serviços
docker-compose down
```

### Opção 2: Localmente

```bash
# Terminal 1: Users Service
make run-users

# Terminal 2: Products Service
make run-products

# Terminal 3: Gateway
make run-gateway
```

### Opção 3: Comandos Makefile

```bash
# Executar todos os serviços (instruções)
make run-all

# Construir imagens Docker
make docker-build

# Subir com Docker Compose
make docker-up

# Parar Docker Compose
make docker-down
```

## 🧪 Testando a Federation

### 1. Script de Testes Automatizado

```bash
# Executar todos os testes
./scripts/test-queries.sh
```

### 2. Queries de Exemplo

#### Query Básica de Usuários

```graphql
query {
  users {
    id
    name
    email
  }
}
```

#### Query Básica de Produtos

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

#### Query Federada - Produtos com Owner

```graphql
query {
  products {
    id
    name
    owner {
      id
      name
      email
    }
  }
}
```

#### Query Federada - Usuários e Produtos Juntos

```graphql
query {
  users {
    id
    name
  }
  products {
    id
    name
    owner {
      id
      name
    }
  }
}
```

#### Query Concorrente - Múltiplos Usuários (WaitGroup + Channels)

```graphql
query {
  usersByIds(ids: ["1", "2", "3", "4", "5"]) {
    id
    name
    email
  }
}
```

### 3. Testes com curl

```bash
# Query federada com owner
curl -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ products { id name owner { id name email } } }"}'

# Query concorrente - múltiplos usuários
curl -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ usersByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name email } }"}'

# Query complexa federada
curl -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ users { id name } products { id name price category owner { id name } } }"}'
```

## 🔗 Endpoints

| Serviço                | URL                             | Descrição            |
| ---------------------- | ------------------------------- | -------------------- |
| **Users Service**      | `http://localhost:8081/query`   | GraphQL endpoint     |
| **Users Health**       | `http://localhost:8081/healthz` | Health check         |
| **Products Service**   | `http://localhost:8082/query`   | GraphQL endpoint     |
| **Products Health**    | `http://localhost:8082/healthz` | Health check         |
| **Apollo Gateway**     | `http://localhost:4000/`        | Federation endpoint  |
| **GraphQL Playground** | `http://localhost:4000/`        | Interface interativa |

## 🔑 Federation Features

### Diretivas @key Implementadas

- **User**: `@key(fields: "id")` - Permite busca por ID
- **Product**: `@key(fields: "id")` - Permite busca por ID

### \_\_resolveReference

- **User.\_\_resolveReference**: Resolve referências por `id`
- **Product.\_\_resolveReference**: Resolve referências por `id`

### Queries Federadas Suportadas

✅ **Busca direta por ID**
✅ **Referências cruzadas entre serviços**
✅ **Queries combinadas de múltiplos serviços**
✅ **Resolução automática de entidades relacionadas**

## ⚡ Concurrency Features

### Resoluções Concorrentes Implementadas

- **WaitGroup**: Sincronização de goroutines
- **Channels**: Comunicação entre goroutines
- **Context**: Cancelamento e timeout
- **Timeout**: 5 segundos por query
- **Latência Simulada**: 100ms por usuário

### Performance

- **Query Concorrente (5 usuários)**: ~0.16s
- **Query Concorrente (8 usuários)**: ~0.16s
- **Queries Sequenciais**: ~0.09s cada (0.45s total para 5)

### Exemplo de Implementação

```go
func (r *Resolver) UsersByIds(ctx context.Context, ids []string) ([]*model.User, error) {
    // Contexto com timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // Canais para resultados e erros
    resultChan := make(chan *model.User, len(ids))
    errorChan := make(chan error, len(ids))

    // WaitGroup para sincronização
    var wg sync.WaitGroup

    // Goroutines para cada ID
    for _, id := range ids {
        wg.Add(1)
        go fetchUser(id, &wg, resultChan, errorChan, ctx)
    }

    // Coletar resultados
    // ...
}
```

## 📊 Benchmark e Race Detection

### Comandos de Teste

```bash
# Race detection
make test-race

# Performance benchmarks
make test-benchmark

# Benchmarks detalhados
make test-benchmark-detail

# Todos os testes
make test-all
```

### Resultados dos Benchmarks

#### Performance Comparativa

```
BenchmarkSequentialUserResolution-8    2   505109980 ns/op    48 B/op    0 allocs/op
BenchmarkConcurrentUserResolution-8    10  101074488 ns/op    2563 B/op  38 allocs/op
```

#### Análise de Performance

- **Sequencial (5 usuários)**: ~505ms
- **Concorrente (5 usuários)**: ~101ms
- **Melhoria**: **5x mais rápido**
- **Memory Overhead**: ~2.4KB vs 48B
- **Allocations**: 37 vs 0

#### Scalability por Tamanho

- **Small (2 usuários)**: ~101ms, 1.3KB, 22 allocs
- **Medium (8 usuários)**: ~101ms, 3.5KB, 50 allocs
- **Large (20 usuários)**: ~101ms, 7.2KB, 99 allocs

### Testes de Segurança

✅ **Race Detection**: Nenhuma race condition detectada
✅ **Timeout Handling**: Context timeout funcionando
✅ **Cancellation**: Context cancellation funcionando
✅ **Valid Data**: Dados válidos processados corretamente
✅ **Invalid Data**: IDs inválidos tratados adequadamente
✅ **Stress Test**: 50 goroutines simultâneas sem problemas

### Script de Benchmark

```bash
./scripts/benchmark.sh
```

Executa automaticamente:

- Race detection tests
- Performance benchmarks
- Unit tests
- Análise detalhada de performance

## 📊 Apollo Studio

Para análise avançada e debugging:

1. Acesse: https://studio.apollographql.com/
2. Conecte seu endpoint: `http://localhost:4000/`
3. Explore o schema federado
4. Analise performance e queries

## 🔧 Configuração

### Variáveis de Ambiente

Copie `env.example` para `.env` e ajuste conforme necessário:

```env
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

# GraphQL
GRAPHQL_PLAYGROUND_ENABLED=true
GRAPHQL_INTROSPECTION_ENABLED=true

# Logging
LOG_LEVEL=info

# Federation
FEDERATION_ENABLED=true
FEDERATION_VERSION=2

# Testes e Desenvolvimento
TEST_TIMEOUT=30s
DEBUG_MODE=false
BENCHMARK_ENABLED=true
RACE_DETECTION_ENABLED=true
```

## 📈 Próximos Passos

- [x] **Resoluções concorrentes** (WaitGroup, context.Context, channels) ✅
- [x] **Benchmarks e race detection** (go test -race, go test -bench) ✅
- [ ] **Cache e otimizações de performance**
- [ ] **Novos serviços** (orders, reviews) que referenciam users/products
- [ ] **Autenticação e autorização**
- [ ] **Métricas e monitoring**

## 🏗️ Arquitetura

![Arquitetura Gofed](docs/arquitecture.drawio)

### Componentes

1. **Frontend/Client**: Consome o GraphQL federado
2. **Apollo Gateway**: Orquestra e combina schemas
3. **Users Service**: Gerencia dados de usuários (com concorrência e testes)
4. **Products Service**: Gerencia dados de produtos
5. **Mock Data**: Dados de exemplo em memória

### Fluxo de Dados

1. Cliente envia query para Apollo Gateway
2. Gateway analisa e roteia para serviços apropriados
3. Serviços processam e retornam dados (concorrentemente)
4. Gateway combina resultados e retorna resposta unificada

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Push para a branch
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
