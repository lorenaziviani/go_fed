# 🚀 Gofed - Federation de Microsserviços com GraphQL

O **Gofed** é uma aplicação demonstrativa que simula um ambiente federado com microsserviços em Go usando GraphQL. Cada microsserviço possui seu schema e expõe parte do domínio (ex: usuários, produtos). A federação é feita via Apollo Gateway.

## 🎯 Objetivo

Demonstrar a implementação de **GraphQL Federation** com microsserviços em Go, incluindo:

- ✅ **Resoluções concorrentes** com WaitGroup, context.Context, canais
- ✅ **Benchmarks e race detection** para validação de performance e segurança
- ✅ **Semáforo customizado** para controle de backpressure e limitação de concorrência
- ✅ **Métricas Prometheus** com contadores de requisições, latência e cache
- ✅ **Request tracing** com TraceID em logs e headers
- Simulação de problemas de performance mitigados com paralelismo e cache
- Federation com Apollo Gateway e diretivas `@key`

## 🛠️ Tech Stack

- **Go 1.24.3**: Linguagem principal para microsserviços
- **GraphQL**: API query language
- **gqlgen**: Biblioteca Go para GraphQL
- **Apollo Gateway (Node.js)**: Para GraphQL federation
- **Docker & Docker Compose**: Containerização e orquestração
- **Federation v2.0**: Com diretivas `@key` para referências cruzadas
- **Concurrency Patterns**: WaitGroup, Channels, Context, Semaphore
- **Testing**: Race detection, Benchmarks, Unit tests
- **Backpressure Control**: Semáforo customizado com chan struct{}
- **Observability**: Prometheus metrics, Request tracing, Structured logging

## 📁 Estrutura do Monorepo

```
gofed/
├── services/
│   ├── users/          # Serviço de usuários (porta 8081)
│   └── products/       # Serviço de produtos (porta 8082)
├── shared/
│   └── metrics/        # Pacote compartilhado de métricas
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

# Testar semáforo e backpressure
./scripts/test-semaphore.sh
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

#### Query com Semáforo - Produtos com Backpressure

```graphql
query {
  productsWithSemaphore(ids: ["1", "2", "3", "4", "5"]) {
    id
    name
    description
    price
    category
    owner {
      id
      name
      email
    }
  }
}
```

#### Estatísticas do Semáforo

```graphql
query {
  semaphoreStats {
    max
    current
    available
    usage
  }
}
```

#### Produtos por Categoria

```graphql
query {
  productsByCategory(category: "Electronics") {
    id
    name
    category
    price
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

# Query com semáforo - produtos com backpressure
curl -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"}'

# Estatísticas do semáforo
curl -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ semaphoreStats { max current available usage } }"}'

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
| **Users Metrics**      | `http://localhost:8081/metrics` | Prometheus metrics   |
| **Products Metrics**   | `http://localhost:8082/metrics` | Prometheus metrics   |

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

### Semáforo Customizado (Products Service)

- **Limitação**: Máximo 3 resoluções concorrentes
- **Implementação**: `chan struct{}` com mutex
- **Backpressure**: Controle automático de sobrecarga
- **Monitoramento**: Estatísticas em tempo real
- **Latência**: 200ms por produto (simulando carga)

### Performance

- **Query Concorrente (5 usuários)**: ~0.16s
- **Query Concorrente (8 usuários)**: ~0.16s
- **Query com Semáforo (5 produtos)**: ~0.4s (backpressure)
- **Queries Sequenciais**: ~0.09s cada (0.45s total para 5)

### Exemplo de Implementação

```go
// Semáforo customizado
type Semaphore struct {
    permits chan struct{}
    mu      sync.RWMutex
    max     int
    current int
}

// Resolução com semáforo
func (r *Resolver) ProductsWithSemaphore(ctx context.Context, ids []string) ([]*model.Product, error) {
    // Adquirir permissão do semáforo
    if err := r.semaphore.Acquire(ctx); err != nil {
        return nil, err
    }
    defer r.semaphore.Release()

    // Processar produto...
}
```

## 🔒 Semáforo e Backpressure

### Implementação do Semáforo

```go
// Semáforo customizado usando chan struct{}
type Semaphore struct {
    permits chan struct{}
    mu      sync.RWMutex
    max     int
    current int
}

// Métodos principais
func (s *Semaphore) Acquire(ctx context.Context) error
func (s *Semaphore) Release()
func (s *Semaphore) Stats() map[string]int
```

### Características do Semáforo

- **Limitação**: Máximo 3 resoluções simultâneas
- **Thread-safe**: Mutex para operações concorrentes
- **Context-aware**: Suporte a cancelamento e timeout
- **Estatísticas**: Monitoramento em tempo real
- **Backpressure**: Controle automático de carga

### Monitoramento

```graphql
query {
  semaphoreStats {
    max # Máximo de permissões (3)
    current # Permissões em uso
    available # Permissões disponíveis
    usage # Percentual de uso (%)
  }
}
```

### Comandos de Teste

```bash
# Testar semáforo
make test-semaphore

# Verificar estatísticas
make test-semaphore-stats

# Testar performance
make test-semaphore-performance
```

## 📊 Métricas e Observabilidade

### Métricas Prometheus

Cada serviço expõe métricas em `/metrics` com os seguintes indicadores:

#### Contadores de Requisições

- `graphql_requests_total` - Total de requisições por serviço, endpoint e tipo de operação
- `graphql_errors_total` - Total de erros por serviço e tipo
- `cache_hits_total` - Hits no cache por serviço
- `cache_misses_total` - Misses no cache por serviço

#### Histogramas de Latência

- `graphql_request_duration_seconds` - Duração das requisições em segundos

#### Gauges de Estado

- `graphql_active_requests` - Número de requisições ativas
- `semaphore_current` - Goroutines atuais usando o semáforo
- `semaphore_max` - Máximo de goroutines permitidas

### Request Tracing

- **TraceID**: Gerado automaticamente ou recebido via header `X-Trace-ID`
- **Logs Estruturados**: Incluem TraceID em todas as entradas
- **Headers de Resposta**: TraceID retornado em `X-Trace-ID`

### Exemplo de Métricas

```bash
# Ver métricas do Users Service
curl http://localhost:8081/metrics

# Ver métricas do Products Service
curl http://localhost:8082/metrics

# Exemplo de saída
# HELP graphql_requests_total Total de requisições GraphQL por serviço e endpoint
# TYPE graphql_requests_total counter
graphql_requests_total{service="users",endpoint="/query",operation_type="query"} 42

# HELP graphql_request_duration_seconds Duração das requisições GraphQL em segundos
# TYPE graphql_request_duration_seconds histogram
graphql_request_duration_seconds_bucket{service="users",endpoint="/query",operation_type="query",le="0.1"} 35
```

### Middleware de Observabilidade

```go
// Chain de middleware: Trace -> Metrics -> Logging
handlerWithMiddleware := metrics.TraceMiddleware(
    metrics.MetricsMiddleware("users")(
        middleware.LoggingMiddleware(logger)(mux),
    ),
)
```

### Logs com TraceID

```json
{
  "level": "info",
  "msg": "Request started",
  "method": "POST",
  "path": "/query",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "time": "2024-01-15T10:30:00Z"
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

# Semáforo e Backpressure
SEMAPHORE_MAX_CONCURRENT=3
SEMAPHORE_TIMEOUT=10s
BACKPRESSURE_ENABLED=true
```

## 📈 Próximos Passos

- [x] **Resoluções concorrentes** (WaitGroup, context.Context, channels) ✅
- [x] **Benchmarks e race detection** (go test -race, go test -bench) ✅
- [x] **Semáforo customizado** (backpressure, limitação de concorrência) ✅
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
4. **Products Service**: Gerencia dados de produtos (com semáforo e backpressure)
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
