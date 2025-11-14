# Neoway Backend - Sistema de Processamento de Pedidos

Sistema de processamento assíncrono de pedidos desenvolvido em Go, utilizando **Arquitetura Hexagonal (Ports & Adapters)**.

## Tecnologias

- **Go 1.25.4** - Linguagem de programação
- **Gin** - Framework web para API REST
- **Swagger/Swaggo** - Documentação interativa da API
- **MongoDB** - Banco de dados NoSQL
- **Apache Kafka** - Plataforma de streaming de eventos
- **Docker Compose** - Orquestração de containers

### Pré-requisitos
- Docker
- Docker Compose

### Iniciar o Sistema

```bash
# Subir todos os serviços
docker compose up -d

# Ver logs
docker compose logs -f
```

Aguarde até 30 segundos para todos os serviços iniciarem.

### Acessar a Documentação

**Swagger UI**: http://localhost:8080/doc

## API Endpoints

| Método | Rota      | Descrição |
|--------|-----------|-----------|
| GET | `/health` | Health check |
| POST | `/orders` | Criar novo pedido |
| GET | `/doc`  | Documentação Swagger |

### Exemplos

```bash
# Health Check
curl http://localhost:8080/health

# Criar Pedido
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"product": "Notebook Dell", "quantity": 2}'
```

## Arquitetura

Este projeto implementa Arquitetura Hexagonal com separação em camadas:

```
internal/
├── domain/              # Regras de negócio
│   ├── entity/         # Entidades com validações
│   └── port/           # Interfaces (contratos)
├── application/        # Casos de uso
│   └── usecase/
└── adapter/            # Implementações
    ├── input/http/     # HTTP handlers
    └── output/         # MongoDB, RabbitMQ
```

### Fluxo de Processamento
1. API recebe pedido e valida
2. Persiste no MongoDB (status: CRIADO)
3. Publica mensagem no tópico Kafka
4. Worker consome mensagem do tópico
5. Worker processa (2s delay)
6. Worker atualiza status para PROCESSADO

## Testes

```bash
# Executar testes
go test ./...
```

```bash
# Com cobertura
go test -cover ./...
```

## Monitoramento

### Kafka UI
- **URL**: http://localhost:8090
- Interface web para monitorar tópicos, mensagens e consumidores
- Visualizar partições e offsets

### Acessar MongoDB e vizualizar pedidos
```bash
docker exec -it neoway-mongodb mongosh neoway
db.orders.find().pretty()
```

## Variáveis de Ambiente

| Variável | Padrão |
|----------|---------|
| `MONGODB_URI` | `mongodb://mongodb:27017` |
| `MONGODB_DATABASE` | `neoway` |
| `KAFKA_BROKERS` | `localhost:9092` |
| `KAFKA_TOPIC` | `orders` |
| `KAFKA_GROUP_ID` | `neoway-workers` (worker apenas) |
| `PORT` | `8080` |

## Boas Práticas Implementadas

- ✅ Arquitetura Hexagonal (Ports & Adapters)
- ✅ Clean Code e SOLID
- ✅ Domain-Driven Design
- ✅ Testes unitários (100% cobertura core)
- ✅ Documentação Swagger/OpenAPI
- ✅ Docker multi-stage builds (~8MB)
- ✅ Graceful shutdown
- ✅ Context propagation
- ✅ Error handling robusto

---

**URLs Importantes**:
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- Kafka UI: http://localhost:8090
