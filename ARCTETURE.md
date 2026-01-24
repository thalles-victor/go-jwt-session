# 📚 GoodTrack Tracker - Documentação Arquitetural

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Arquitetura do Projeto](#arquitetura-do-projeto)
3. [Estrutura de Diretórios](#estrutura-de-diretórios)
4. [Padrões Arquiteturais](#padrões-arquiteturais)
5. [Camadas da Aplicação](#camadas-da-aplicação)
6. [Inversão de Dependências](#inversão-de-dependências)
7. [Contratos (Interfaces)](#contratos-interfaces)
8. [Adapters (Adaptadores)](#adapters-adaptadores)
9. [Domain Errors - Tratamento de Erros](#domain-errors---tratamento-de-erros)
10. [Testing & Mocks](#testing--mocks)
11. [Exemplos Práticos](#exemplos-práticos)
12. [Boas Práticas](#boas-práticas)

---

## 🎯 Visão Geral

O **GoodTrack Tracker** é um sistema de tracking e processamento de eventos de conversão para múltiplas plataformas de checkout (Hotmart, Hubla, PerfectPay, Cakto, etc.) com integração ao Meta Pixel e APIs de conversão.

### Tecnologias Principais

- **Go 1.21+** - Linguagem principal
- **Fiber v2** - Framework web
- **PostgreSQL** - Banco de dados relacional
- **MongoDB** - Banco de dados NoSQL para eventos
- **Docker** - Containerização
- **Air** - Hot reload para desenvolvimento

---

## 🏗️ Arquitetura do Projeto

O projeto segue os princípios de **Clean Architecture** e **SOLID**, com foco em:

- **Separação de Responsabilidades**: Cada camada tem uma responsabilidade bem definida
- **Inversão de Dependências**: Dependemos de abstrações (interfaces), não de implementações concretas
- **Testabilidade**: Todas as dependências podem ser mockadas
- **Baixo Acoplamento**: Módulos independentes e reutilizáveis
- **Alta Coesão**: Componentes relacionados agrupados logicamente

### Diagrama Conceitual

```
┌─────────────────────────────────────────────────────────────┐
│                         HTTP Layer                          │
│                    (Controllers/Routes)                     │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                      Business Logic                         │
│                        (Services)                           │
└───────────────────────────┬─────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
┌───────▼──────┐   ┌────────▼────────┐   ┌─────▼──────┐
│ Repositories │   │    Adapters     │   │   Proxies  │
│  (DB Access) │   │ (External APIs) │   │ (Services) │
└──────────────┘   └─────────────────┘   └────────────┘
```

---

## 📁 Estrutura de Diretórios

```
goodtrack-tracker/
│
├── cmd/
│   └── server/
│       └── main.go                    # Entry point da aplicação
│
├── internal/
│   ├── database/                      # Configurações de banco de dados
│   │   └── database.go
│   │
│   ├── domains/                       # Domínios de negócio
│   │   ├── checkouts/                 # Processamento de checkouts
│   │   │   ├── hotmart/
│   │   │   │   ├── contracts.go       # Interfaces (Service, MetaAPI, ConversionAPI)
│   │   │   │   ├── adapters.go        # Adaptadores externos (Meta, Currency)
│   │   │   │   ├── service.go         # Lógica de negócio
│   │   │   │   ├── controller.go      # Handler HTTP
│   │   │   │   ├── dto.go             # Data Transfer Objects
│   │   │   │   └── service_test.go    # Testes unitários
│   │   │   │
│   │   │   ├── hubla/
│   │   │   ├── perfectpay/
│   │   │   ├── cakto/
│   │   │   ├── kiwify/
│   │   │   └── ...                    # Outros checkouts
│   │   │
│   │   ├── events/                    # Sistema de eventos
│   │   │   ├── view/                  # Evento de visualização
│   │   │   │   ├── view.contract.go   # Interfaces
│   │   │   │   ├── view.adapter.go    # Adaptadores
│   │   │   │   ├── view.service.go    # Lógica de negócio
│   │   │   │   ├── view.controller.go # Handler HTTP
│   │   │   │   └── view.dto.go        # DTOs
│   │   │   │
│   │   │   └── lead/                  # Evento de lead
│   │   │       └── ...
│   │   │
│   │   └── leona/                     # WhatsApp Tracking
│   │       └── ...
│   │
│   ├── ipgeolocation/                 # Geolocalização de IPs
│   │   ├── ipgeolocation.contract.go  # Interfaces
│   │   ├── ipgeolocation.adapter.go   # Adaptadores
│   │   ├── ipgolocation.service.go    # Lógica de negócio
│   │   ├── ipgeolocation.controller.go # Handler HTTP
│   │   └── ipgeolocation.dto.go       # DTOs
│   │
│   ├── models/                        # Estruturas de dados
│   │   ├── events.go
│   │   ├── sales.go
│   │   ├── webhook.go
│   │   └── ...
│   │
│   ├── repositories/                  # Camada de acesso a dados
│   │   ├── postgres/
│   │   │   ├── pixel_repository.go
│   │   │   ├── webhook_repository.go
│   │   │   └── ...
│   │   │
│   │   └── mongo/
│   │       ├── event_repository.go
│   │       ├── sales_repository.go
│   │       └── ...
│   │
│   ├── proxy/                         # Serviços externos
│   │   ├── meta/                      # Meta Conversion API
│   │   │   ├── meta.go
│   │   │   └── api.go
│   │   │
│   │   ├── currency/                  # API de conversão de moeda
│   │   │   ├── api.go
│   │   │   ├── adapter.go
│   │   │   └── contracts.go
│   │   │
│   │   └── goelocation/               # API de geolocalização
│   │       └── geolocation.go
│   │
│   ├── routes/                        # Configuração de rotas
│   │   ├── routes.go                  # Setup de rotas
│   │   └── controllers.go             # Inicialização de controllers
│   │
│   └── shared/                        # Utilitários compartilhados
│       ├── cache/
│       │   ├── cache.go               # Interface de cache
│       │   └── cache.mock.go          # Mock para testes
│       │
│       ├── logger/
│       │   └── logger.go              # Sistema de logs
│       │
│       ├── middlewares/
│       │   ├── request_id.go
│       │   └── request_logger.go
│       │
│       └── utilities/
│           ├── consts.go
│           ├── format.go
│           └── hasher.go
│
├── docs/                              # Documentação
│   └── README.md                      # Este arquivo
│
├── .air.toml                          # Configuração do Air (hot reload)
├── docker-compose.yml                 # Compose para desenvolvimento
├── Dockerfile                         # Build da imagem
├── go.mod                             # Dependências Go
└── README.md                          # README principal
```

---

## 🎨 Padrões Arquiteturais

### 1. Domain-Driven Design (DDD)

Cada domínio de negócio (checkout, events, leona) é isolado em seu próprio pacote com:

- **Entities**: Modelos de dados (`models/`)
- **Services**: Lógica de negócio
- **Repositories**: Acesso a dados
- **DTOs**: Objetos de transferência

### 2. Dependency Injection

As dependências são injetadas via construtores, facilitando testes e manutenção:

```go
func NewService(
    webhookRepo WebhookRepository,
    metaAPI MetaAPI,
    conversionAPI ConversionAPI,
) Service {
    return &service{
        webhookRepo:   webhookRepo,
        metaAPI:       metaAPI,
        conversionAPI: conversionAPI,
    }
}
```

### 3. Repository Pattern

Abstrações para acesso a dados, permitindo trocar implementações facilmente:

```go
type PixelRepository interface {
    GetPixelsByWebhookID(ctx context.Context, webhookID string) ([]models.MetaPixel, error)
    GetMetaByPixelValue(pixel string) (models.MetaPixel, error)
}
```

### 4. Adapter Pattern

Wrappers para serviços externos, isolando a lógica de integração:

```go
type MetaAdapter struct{}

func (m *MetaAdapter) SendEvent(ctx context.Context, event meta.MetaEvent) error {
    return meta.SendEvent(ctx, event, "HotmartAdapter: ")
}
```

---

## 🔄 Camadas da Aplicação

### 1. **Controllers Layer** (HTTP Handlers)

Responsável por:
- Receber requisições HTTP
- Validar dados de entrada
- Chamar serviços
- Retornar respostas HTTP

**Exemplo:**

```go
// internal/domains/checkouts/hotmart/controller.go
type Controller struct {
    service Service
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
    var dto HotmartDto
    if err := ctx.BodyParser(&dto); err != nil {
        return ctx.Status(400).JSON(fiber.Map{"error": "invalid payload"})
    }
    
    webhookID := ctx.Params("webhookID")
    hotToken := ctx.Get("X-Hotmart-Hottok")
    
    if err := c.service.Handler(ctx.Context(), webhookID, hotToken, ctx.Body(), dto); err != nil {
        return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return ctx.Status(200).JSON(fiber.Map{"message": "success"})
}
```

### 2. **Services Layer** (Business Logic)

Responsável por:
- Regras de negócio
- Orquestração de dependências
- Validações complexas
- Transformações de dados

**Exemplo:**

```go
// internal/domains/checkouts/hotmart/service.go
type service struct {
    webhookRepo   WebhookRepository
    metaAPI       MetaAPI
    conversionAPI ConversionAPI
}

func (s *service) Handler(ctx context.Context, webhookID string, hotToken string, rawBody []byte, dto HotmartDto) error {
    // 1. Validar webhook
    webhook, err := s.webhookRepo.GetWebhookByID(ctx, webhookID)
    if err != nil {
        return err
    }
    
    // 2. Validar token usando Domain Error
    if hotToken != *webhook.ValidationToken {
        return ErrInvalidToken // Domain Error definido em errors.go
    }
    
    // 3. Converter moeda se necessário
    if dto.Currency != "BRL" {
        result, err := s.conversionAPI.ConvertToBRL(ctx, valueInCents, dto.Currency)
        if err != nil {
            return NewDomainError(
                ErrConvertPriceToBRL.Code,
                ErrConvertPriceToBRL.Message,
                ErrConvertPriceToBRL.HTTPStatus,
                ErrConvertPriceToBRL.Op,
                err,
            )
        }
        // ...
    }
    
    // 4. Enviar evento para Meta
    s.metaAPI.SendEvent(ctx, metaEvent)
    
    return nil
}
```

### 3. **Repositories Layer** (Data Access)

Responsável por:
- Acesso a bancos de dados
- Queries e comandos
- Mapeamento de modelos

**Exemplo:**

```go
// internal/repositories/postgres/webhook_repository.go
type WebhookRepository interface {
    GetWebhookByID(ctx context.Context, id string) (*models.WebhookModel, error)
}

type webhookRepository struct {
    db *sqlx.DB
}

func (r *webhookRepository) GetWebhookByID(ctx context.Context, id string) (*models.WebhookModel, error) {
    var webhook models.WebhookModel
    err := r.db.GetContext(ctx, &webhook, "SELECT * FROM webhooks WHERE id = $1", id)
    return &webhook, err
}
```

### 4. **Proxies Layer** (External Services)

Responsável por:
- Comunicação com APIs externas
- Tratamento de erros de rede
- Retry logic
- Rate limiting

**Exemplo:**

```go
// internal/proxy/goelocation/geolocation.go
type GeoLocationHTTP struct {
    Client *http.Client
}

func (g *GeoLocationHTTP) GetGeoLocation(ip string, geoLocation *GeoLocationResponseData, logPrefix string) error {
    url := fmt.Sprintf("https://get.geojs.io/v1/ip/geo/%s.json", ip)
    
    resp, err := g.Client.Get(url)
    if err != nil {
        return fmt.Errorf("erro ao fazer requisição: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("requisição retornou status: %s", resp.Status)
    }
    
    return json.NewDecoder(resp.Body).Decode(&geoLocation)
}
```

---

## 🔀 Inversão de Dependências

A inversão de dependências é um dos pilares do projeto. Implementamos através de **interfaces (contratos)**.

### Princípio

> Módulos de alto nível não devem depender de módulos de baixo nível. Ambos devem depender de abstrações.

### Como Implementamos

**1. Definir Contratos (Interfaces)**

```go
// internal/domains/checkouts/hotmart/contracts.go
package hotmart

// Service define o contrato para o serviço de Hotmart
type Service interface {
    Handler(ctx context.Context, webhookID string, hotToken string, rawBody []byte, dto HotmartDto) error
}

// MetaAPI define o contrato para envio de eventos ao Meta
type MetaAPI interface {
    SendEvent(ctx context.Context, event meta.MetaEvent) error
}

// ConversionAPI define o contrato para conversão de moedas
type ConversionAPI interface {
    ConvertToBRL(ctx context.Context, valueInCents uint64, currency string) (currency.ConversionResult, error)
}
```

**2. Implementar Services com Dependências via Interface**

```go
// internal/domains/checkouts/hotmart/service.go
package hotmart

type service struct {
    webhookRepo   WebhookRepository   // Depende de interface
    metaAPI       MetaAPI             // Depende de interface
    conversionAPI ConversionAPI       // Depende de interface
}

// NewService recebe as dependências via parâmetros (Dependency Injection)
func NewService(
    webhookRepo WebhookRepository,
    metaAPI MetaAPI,
    conversionAPI ConversionAPI,
) Service {
    return &service{
        webhookRepo:   webhookRepo,
        metaAPI:       metaAPI,
        conversionAPI: conversionAPI,
    }
}
```

**3. Criar Adapters para Implementações Concretas**

```go
// internal/domains/checkouts/hotmart/adapters.go
package hotmart

// MetaAdapter implementa a interface MetaAPI
type MetaAdapter struct{}

func NewMetaAdapter() *MetaAdapter {
    return &MetaAdapter{}
}

// Implementação do contrato MetaAPI
func (m *MetaAdapter) SendEvent(ctx context.Context, event meta.MetaEvent) error {
    // Chama a implementação real do pacote meta
    return meta.SendEvent(ctx, event, "HotmartAdapter: ")
}

// ConversionAdapter implementa a interface ConversionAPI
type ConversionAdapter struct{}

func NewConversionAdapter() *ConversionAdapter {
    return &ConversionAdapter{}
}

// Implementação do contrato ConversionAPI
func (c *ConversionAdapter) ConvertToBRL(ctx context.Context, valueInCents uint64, curr string) (currency.ConversionResult, error) {
    // Chama a implementação real do pacote currency
    return currency.ConvertToBRL(ctx, valueInCents, curr)
}
```

**4. Inicializar com Implementações Reais**

```go
// internal/routes/controllers.go
func InitializeControllers(db *sqlx.DB) *DomainControllers {
    // Repositórios
    webhookRepo := postgresRepo.NewWebhookRepository(db)
    
    // Adapters (implementações das interfaces)
    metaAdapter := hotmart.NewMetaAdapter()
    conversionAdapter := hotmart.NewConversionAdapter()
    
    // Service recebe as implementações via interface
    hotmartService := hotmart.NewService(
        webhookRepo,      // WebhookRepository
        metaAdapter,      // MetaAPI
        conversionAdapter, // ConversionAPI
    )
    
    return &DomainControllers{
        Hotmart: hotmart.NewController(hotmartService),
    }
}
```

### Benefícios da Inversão de Dependências

1. **Testabilidade**: Fácil criar mocks das dependências
2. **Flexibilidade**: Trocar implementações sem alterar lógica de negócio
3. **Baixo Acoplamento**: Módulos não dependem de implementações concretas
4. **Manutenibilidade**: Mudanças isoladas em cada camada

---

## 📜 Contratos (Interfaces)

Contratos definem **o quê** fazer, não **como** fazer.

### Estrutura de um Arquivo de Contratos

```go
// internal/domains/checkouts/hotmart/contracts.go
package hotmart

import (
    "context"
    "webhook-worker/internal/proxy/currency"
    "webhook-worker/internal/proxy/meta"
)

// Service é a interface principal do domínio
type Service interface {
    Handler(ctx context.Context, webhookID string, hotToken string, rawBody []byte, dto HotmartDto) error
}

// MetaAPI define operações de envio ao Meta Pixel
type MetaAPI interface {
    SendEvent(ctx context.Context, event meta.MetaEvent) error
}

// ConversionAPI define operações de conversão de moeda
type ConversionAPI interface {
    ConvertToBRL(ctx context.Context, valueInCents uint64, currency string) (currency.ConversionResult, error)
}
```

### Exemplos de Contratos no Projeto

#### 1. Cache Contract

```go
// internal/shared/cache/cache.go
package cache

type CacheContract interface {
    Get(key string) (any, bool)
    Set(key string, value any, ttl time.Duration)
    Delete(key string)
}
```

#### 2. Event Repository Contract

```go
// internal/domains/events/view/view.contract.go
package view

type EventRepository interface {
    Create(ctx context.Context, event models.EventModel) error
}
```

#### 3. Geolocation Service Contract

```go
// internal/ipgeolocation/ipgeolocation.contract.go
package ipgeolocation

type GetLocationService interface {
    GetLocation(ip string, geoLocation *goelocation.GeoLocationResponseData, logPrefix string) error
}

type Service interface {
    GetIpGeoLocation(ip string) (*goelocation.GeoLocationResponseData, error)
}
```

### Quando Criar um Contrato?

✅ **SIM, criar contrato quando:**
- A funcionalidade será testada com mocks
- A implementação pode mudar no futuro
- A dependência é externa (APIs, banco de dados)
- Você quer desacoplar camadas

❌ **NÃO criar contrato quando:**
- Funções utilitárias simples
- Estruturas de dados (DTOs, models)
- Constantes e enums

---

## 🔌 Adapters (Adaptadores)

Adapters fazem a ponte entre **contratos** e **implementações concretas**.

### Propósito dos Adapters

1. **Isolar dependências externas**: APIs, serviços, bibliotecas
2. **Permitir substituição**: Trocar implementação sem afetar o domínio
3. **Facilitar testes**: Criar mocks do adapter em vez da implementação real

### Estrutura de um Adapter

```go
// internal/domains/checkouts/hotmart/adapters.go
package hotmart

import (
    "context"
    "webhook-worker/internal/proxy/currency"
    "webhook-worker/internal/proxy/meta"
)

// ========================================
// Meta Adapter
// ========================================

type MetaAdapter struct{}

func NewMetaAdapter() *MetaAdapter {
    return &MetaAdapter{}
}

// Implementação da interface MetaAPI
func (m *MetaAdapter) SendEvent(ctx context.Context, event meta.MetaEvent) error {
    return meta.SendEvent(ctx, event, "HotmartAdapter: ")
}

// ========================================
// Conversion Adapter
// ========================================

type ConversionAdapter struct{}

func NewConversionAdapter() *ConversionAdapter {
    return &ConversionAdapter{}
}

// Implementação da interface ConversionAPI
func (c *ConversionAdapter) ConvertToBRL(ctx context.Context, valueInCents uint64, curr string) (currency.ConversionResult, error) {
    return currency.ConvertToBRL(ctx, valueInCents, curr)
}
```

### Adapter com Estado (Stateful)

Alguns adapters podem armazenar configurações:

```go
// internal/ipgeolocation/ipgeolocation.adapter.go
package ipgeolocation

type geoLocationAdapter struct {
    geoHTTP *goelocation.GeoLocationHTTP  // Dependência injetada
}

func NewGeoLocationAdapter(geoHTTP *goelocation.GeoLocationHTTP) GetLocationService {
    return &geoLocationAdapter{
        geoHTTP: geoHTTP,
    }
}

func (a *geoLocationAdapter) GetLocation(ip string, geoLocation *goelocation.GeoLocationResponseData, logPrefix string) error {
    // Delega para a implementação real
    return a.geoHTTP.GetGeoLocation(ip, geoLocation, logPrefix)
}
```

### Padrão de Nomenclatura

- **Arquivo**: `adapters.go`
- **Struct**: `MetaAdapter`, `ConversionAdapter`, etc.
- **Constructor**: `NewMetaAdapter()`, `NewConversionAdapter()`, etc.
- **Métodos**: Implementam os métodos da interface correspondente

---

## ⚠️ Domain Errors - Tratamento de Erros

O projeto utiliza o padrão **Domain Errors** para tratamento estruturado de erros, seguindo os princípios de Clean Architecture. Este padrão permite que erros carreguem informações suficientes para definir status HTTP, códigos padronizados e contexto de debugging.

### O Problema

Em Go, erros são apenas valores que implementam a interface `error`. Um simples `errors.New("user not found")` não carrega informações suficientes para:

- Definir o status HTTP correto na resposta
- Enviar um código padronizado para o cliente
- Identificar a operação que falhou
- Manter o erro original para debugging

### A Solução: Domain Errors

A ideia é criar um tipo de erro rico que carregue toda informação necessária:

```go
// internal/domains/checkouts/perfectpay/errors.go
package perfectpay

import (
	"fmt"
	"net/http"
)

// DomainError representa um erro do domínio com informações estruturadas
type DomainError struct {
	Code       string // ex: WEBHOOK_VALIDATION_TOKEN_REQUIRED
	Message    string // mensagem amigável
	HTTPStatus int    // código HTTP sugerido
	Op         string // operação que falhou (ex: "perfectpay.handler")
	Err        error  // erro original (wrapping)
}

// Error implementa a interface error
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap retorna o erro original para permitir unwrapping
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is permite comparação de erros por código usando errors.Is
func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// NewDomainError cria um novo DomainError
func NewDomainError(code, message string, httpStatus int, op string, err error) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Op:         op,
		Err:        err,
	}
}
```

### Por que isso funciona?

1. **Code**: Identificador único e padronizado do erro (ex: `WEBHOOK_VALIDATION_TOKEN_REQUIRED`)
2. **Message**: Mensagem amigável que pode ser exibida ao usuário
3. **HTTPStatus**: O middleware sabe qual status retornar sem adivinhar
4. **Op**: Contexto da operação - útil para logs e debugging
5. **Err**: Mantém o erro original através de `Unwrap()`

### Definindo Erros do Domínio

Cada módulo do seu domínio deve ter seu arquivo `errors.go`:

```go
// internal/domains/checkouts/perfectpay/errors.go
package perfectpay

import "net/http"

var (
	// ErrWebhookValidationTokenRequired é retornado quando o webhook não tem token de validação
	ErrWebhookValidationTokenRequired = NewDomainError(
		"WEBHOOK_VALIDATION_TOKEN_REQUIRED",
		"webhook não tem um token de validação cadastrado",
		http.StatusBadRequest,
		"perfectpay.handler",
		nil,
	)

	// ErrInvalidToken é retornado quando o token informado é inválido
	ErrInvalidToken = NewDomainError(
		"INVALID_TOKEN",
		"token informado inválido",
		http.StatusUnauthorized,
		"perfectpay.handler",
		nil,
	)

	// ErrEventAlreadyOccurred é retornado quando o evento já foi processado
	ErrEventAlreadyOccurred = NewDomainError(
		"EVENT_ALREADY_OCCURRED",
		"o evento já ocorreu",
		http.StatusConflict,
		"perfectpay.handler",
		nil,
	)

	// ErrSaveSale é retornado quando há erro ao salvar os dados da venda
	ErrSaveSale = NewDomainError(
		"SAVE_SALE_ERROR",
		"erro ao salvar os dados da venda",
		http.StatusInternalServerError,
		"perfectpay.handler",
		nil,
	)
)
```

### Usando Domain Errors no Service

Agora suas validações ficam extremamente legíveis:

```go
// internal/domains/checkouts/perfectpay/service.go
func (s *service) Handler(ctx context.Context, webhookID string, token string, rawBody []byte, dto PerfectPayWebhookDto) error {
	webhook, err := s.webhookRepo.GetWebhookByID(ctx, webhookID)
	if err != nil {
		return err
	}

	// Validações limpas e claras
	if webhook.ValidationToken == nil {
		return ErrWebhookValidationTokenRequired
	}

	if token == "" {
		return ErrAuthenticationTokenRequired
	}

	if token != *webhook.ValidationToken {
		return ErrInvalidToken
	}

	// Verificar idempotência
	exists, err := s.salesRepo.ExistsByIdempotencyKey(ctx, idepotencyKey)
	if err != nil {
		return NewDomainError(
			ErrCheckIdempotency.Code,
			ErrCheckIdempotency.Message,
			ErrCheckIdempotency.HTTPStatus,
			ErrCheckIdempotency.Op,
			err, // Wrapping do erro original
		)
	}

	if exists {
		return ErrEventAlreadyOccurred
	}

	// ... resto da lógica
}
```

### Wrapping de Erros

Quando você precisa preservar o erro original (por exemplo, de um repositório), use `NewDomainError` com o erro original:

```go
// Erro de repositório com contexto do domínio
payloadID, err := s.payloadRepo.SavePayload(ctx, payload)
if err != nil {
	return NewDomainError(
		ErrSaveWebhookPayload.Code,
		ErrSaveWebhookPayload.Message,
		ErrSaveWebhookPayload.HTTPStatus,
		ErrSaveWebhookPayload.Op,
		err, // Preserva o erro original do repositório
	)
}
```

### Tratamento no Controller/Middleware

O controller ou middleware pode extrair o `DomainError` e responder de forma padronizada:

```go
// internal/domains/checkouts/perfectpay/controller.go
func (c *Controller) Handle(ctx *fiber.Ctx) error {
	var dto PerfectPayWebhookDto
	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"code":  "INVALID_PAYLOAD",
			"error": "invalid request payload",
		})
	}

	webhookID := ctx.Params("webhookID")
	token := ctx.Get("X-PerfectPay-Token")

	err := c.service.Handler(ctx.Context(), webhookID, token, ctx.Body(), dto)
	if err != nil {
		var domainErr *DomainError
		if errors.As(err, &domainErr) {
			// Extrai o DomainError e responde com status e código apropriados
			return ctx.Status(domainErr.HTTPStatus).JSON(fiber.Map{
				"code":  domainErr.Code,
				"error": domainErr.Message,
			})
		}

		// Fallback para erros não-domínio
		return ctx.Status(500).JSON(fiber.Map{
			"code":  "INTERNAL_ERROR",
			"error": "internal server error",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "success"})
}
```

### Comparação de Erros com errors.Is

O método `Is` permite comparações por código:

```go
// Em testes ou validações
if errors.Is(err, ErrInvalidToken) {
	// Tratar erro de token inválido
}

// Funciona mesmo quando o erro foi wrapped
if errors.Is(err, ErrSaveSale) {
	// Tratar erro ao salvar venda
}
```

### Exemplo Completo: Perfect Pay

Veja a implementação completa no domínio Perfect Pay:

```go
// internal/domains/checkouts/perfectpay/errors.go
package perfectpay

// ... estrutura DomainError ...

var (
	ErrWebhookValidationTokenRequired = NewDomainError(
		"WEBHOOK_VALIDATION_TOKEN_REQUIRED",
		"webhook não tem um token de validação cadastrado",
		http.StatusBadRequest,
		"perfectpay.handler",
		nil,
	)

	ErrInvalidToken = NewDomainError(
		"INVALID_TOKEN",
		"token informado inválido",
		http.StatusUnauthorized,
		"perfectpay.handler",
		nil,
	)

	ErrEventAlreadyOccurred = NewDomainError(
		"EVENT_ALREADY_OCCURRED",
		"o evento já ocorreu",
		http.StatusConflict,
		"perfectpay.handler",
		nil,
	)

	// ... outros erros ...
)
```

### Benefícios do Padrão Domain Errors

✅ **Separação de responsabilidades**: O domínio define O QUE é o erro, não COMO apresentá-lo  
✅ **Consistência**: Todos os erros seguem o mesmo padrão  
✅ **Debugging facilitado**: O campo `Op` mostra exatamente onde o erro ocorreu  
✅ **Composição**: Use `Unwrap()` para manter a cadeia de erros  
✅ **Type-safe**: `errors.Is` e `errors.As` funcionam perfeitamente  
✅ **Testável**: Compare erros de forma clara nos testes  
✅ **Status HTTP definido no domínio**: Não precisa decidir no controller

### Convenções

1. **Arquivo**: Cada domínio deve ter um arquivo `errors.go`
2. **Nomenclatura**: `Err` + `NomeDescritivo` (ex: `ErrInvalidToken`)
3. **Códigos**: UPPER_SNAKE_CASE (ex: `INVALID_TOKEN`)
4. **Mensagens**: Em português, amigáveis ao usuário
5. **Status HTTP**: Use os códigos apropriados (400, 401, 404, 409, 500, etc.)
6. **Op**: Formato `domínio.operacao` (ex: `perfectpay.handler`)

### Mapeamento de Status HTTP

| Tipo de Erro | Status HTTP | Exemplo |
|-------------|-------------|---------|
| Validação | 400 Bad Request | `ErrWebhookValidationTokenRequired` |
| Autenticação | 401 Unauthorized | `ErrInvalidToken` |
| Não encontrado | 404 Not Found | `ErrWebhookNotFound` |
| Conflito | 409 Conflict | `ErrEventAlreadyOccurred` |
| Erro interno | 500 Internal Server Error | `ErrSaveSale` |

### Exemplo de Resposta ao Cliente

Quando um `DomainError` é retornado, o cliente recebe uma resposta consistente:

```json
{
  "code": "INVALID_TOKEN",
  "error": "token informado inválido"
}
```

Com status HTTP `401 Unauthorized` - tudo definido no domínio, não no handler.

---

## 🧪 Testing & Mocks

O projeto está estruturado para ser altamente testável através de mocks.

### Estratégia de Testes

1. **Unit Tests**: Testam serviços isoladamente com mocks
2. **Integration Tests**: Testam fluxos completos (futuro)
3. **Mocks Manuais**: Implementações fake para testes

### Criando Mocks

#### 1. Mock Inline (Em arquivo de teste)

```go
// internal/domains/checkouts/hotmart/service_test.go
package hotmart

// Mock do WebhookRepository
type mockWebhookRepository struct {
    webhook *models.WebhookModel
    err     error
}

func (m *mockWebhookRepository) GetWebhookByID(ctx context.Context, id string) (*models.WebhookModel, error) {
    return m.webhook, m.err
}

// Mock do MetaAPI
type mockMetaAPI struct {
    lastEvent *meta.MetaEvent
    err       error
}

func (m *mockMetaAPI) SendEvent(ctx context.Context, event meta.MetaEvent) error {
    m.lastEvent = &event  // Armazena para verificação
    return m.err
}
```

#### 2. Mock com Funções (Flexível)

```go
// internal/shared/cache/cache.mock.go
package cache

type CacheMock struct {
    GetFunc func(key string) (any, bool)
    SetFunc func(key string, value any, ttl time.Duration)
}

func (m *CacheMock) Get(key string) (any, bool) {
    return m.GetFunc(key)
}

func (m *CacheMock) Set(key string, value any, ttl time.Duration) {
    m.SetFunc(key, value, ttl)
}
```

**Uso:**

```go
cacheMock := &cache.CacheMock{
    GetFunc: func(key string) (any, bool) {
        if key == "test-key" {
            return "test-value", true
        }
        return nil, false
    },
    SetFunc: func(key string, value any, ttl time.Duration) {
        // não faz nada
    },
}
```

### Exemplo de Teste Completo

```go
// internal/domains/checkouts/hotmart/service_test.go
func TestHandler_WithSCK(t *testing.T) {
    // Setup
    validToken := "valid-token"
    webhookID := "webhook-123"
    sck := "sck-12345"

    // Mocks
    webhookRepo := &mockWebhookRepository{
        webhook: &models.WebhookModel{
            ID:              webhookID,
            UserID:          "user-1",
            ValidationToken: &validToken,
        },
    }
    
    metaAPI := &mockMetaAPI{}
    conversionAPI := &mockConversionAPI{}

    // Criar service com mocks
    service := NewService(webhookRepo, metaAPI, conversionAPI)

    // Preparar input
    rawBody, dto := getTestPayload(&sck)

    // Executar
    err := service.Handler(context.Background(), webhookID, validToken, rawBody, dto)

    // Assertions
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }

    // Verificar que Meta Event foi enviado
    if metaAPI.lastEvent == nil {
        t.Fatal("Expected Meta Event to be sent, but it was nil")
    }

    if metaAPI.lastEvent.EventType != "Purchase" {
        t.Errorf("Expected event type Purchase, got %s", metaAPI.lastEvent.EventType)
    }
}
```

### Estrutura de Testes

```
domain/
├── contracts.go           # Interfaces
├── adapters.go            # Implementações reais
├── service.go             # Lógica de negócio
├── service_test.go        # Testes + Mocks inline
└── controller.go          # HTTP handlers
```

### Boas Práticas de Testes

1. **Nomear testes descritivamente**: `TestHandler_WithSCK`, `TestHandler_InvalidToken`
2. **Usar mocks para dependências externas**: APIs, banco de dados
3. **Testar casos de sucesso e erro**: Happy path e edge cases
4. **Verificar side effects**: Checagem se métodos foram chamados
5. **Usar table-driven tests** para múltiplos cenários similares

---

## 💡 Exemplos Práticos

### Exemplo 1: Criar um Novo Domínio de Checkout

Vamos criar um domínio para a plataforma "NewPay".

#### Passo 1: Criar estrutura de arquivos

```
internal/domains/checkouts/newpay/
├── contracts.go
├── errors.go          # Domain Errors
├── adapters.go
├── service.go
├── controller.go
├── dto.go
└── service_test.go
```

#### Passo 2: Definir contratos

```go
// contracts.go
package newpay

import (
    "context"
    "webhook-worker/internal/proxy/currency"
    "webhook-worker/internal/proxy/meta"
)

type Service interface {
    Handler(ctx context.Context, webhookID string, token string, rawBody []byte, dto NewPayDto) error
}

type MetaAPI interface {
    SendEvent(ctx context.Context, event meta.MetaEvent) error
}

type ConversionAPI interface {
    ConvertToBRL(ctx context.Context, valueInCents uint64, currency string) (currency.ConversionResult, error)
}
```

#### Passo 2.5: Criar Domain Errors

```go
// errors.go
package newpay

import (
	"fmt"
	"net/http"
)

// DomainError representa um erro do domínio com informações estruturadas
type DomainError struct {
	Code       string
	Message    string
	HTTPStatus int
	Op         string
	Err        error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

func NewDomainError(code, message string, httpStatus int, op string, err error) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Op:         op,
		Err:        err,
	}
}

var (
	ErrInvalidToken = NewDomainError(
		"INVALID_TOKEN",
		"token informado inválido",
		http.StatusUnauthorized,
		"newpay.handler",
		nil,
	)

	ErrWebhookNotFound = NewDomainError(
		"WEBHOOK_NOT_FOUND",
		"webhook não encontrado",
		http.StatusNotFound,
		"newpay.handler",
		nil,
	)
)
```

#### Passo 3: Criar adapters

```go
// adapters.go
package newpay

import (
    "context"
    "webhook-worker/internal/proxy/currency"
    "webhook-worker/internal/proxy/meta"
)

type MetaAdapter struct{}

func NewMetaAdapter() *MetaAdapter {
    return &MetaAdapter{}
}

func (m *MetaAdapter) SendEvent(ctx context.Context, event meta.MetaEvent) error {
    return meta.SendEvent(ctx, event, "NewPayAdapter: ")
}

type ConversionAdapter struct{}

func NewConversionAdapter() *ConversionAdapter {
    return &ConversionAdapter{}
}

func (c *ConversionAdapter) ConvertToBRL(ctx context.Context, valueInCents uint64, curr string) (currency.ConversionResult, error) {
    return currency.ConvertToBRL(ctx, valueInCents, curr)
}
```

#### Passo 4: Implementar service

```go
// service.go
package newpay

import (
    "context"
    "errors"
    postgresRepo "webhook-worker/internal/repositories/postgres"
    mongoRepo "webhook-worker/internal/repositories/mongo"
)

type service struct {
    webhookRepo   postgresRepo.WebhookRepository
    payloadRepo   mongoRepo.PayloadRepository
    pixelRepo     postgresRepo.PixelRepository
    eventRepo     mongoRepo.EventRepository
    salesRepo     mongoRepo.SalesRepository
    metaAPI       MetaAPI
    conversionAPI ConversionAPI
}

func NewService(
    webhookRepo postgresRepo.WebhookRepository,
    payloadRepo mongoRepo.PayloadRepository,
    pixelRepo postgresRepo.PixelRepository,
    eventRepo mongoRepo.EventRepository,
    salesRepo mongoRepo.SalesRepository,
    metaAPI MetaAPI,
    conversionAPI ConversionAPI,
) Service {
    return &service{
        webhookRepo:   webhookRepo,
        payloadRepo:   payloadRepo,
        pixelRepo:     pixelRepo,
        eventRepo:     eventRepo,
        salesRepo:     salesRepo,
        metaAPI:       metaAPI,
        conversionAPI: conversionAPI,
    }
}

func (s *service) Handler(ctx context.Context, webhookID string, token string, rawBody []byte, dto NewPayDto) error {
    // 1. Validar webhook
    webhook, err := s.webhookRepo.GetWebhookByID(ctx, webhookID)
    if err != nil {
        // Wrapping com Domain Error
        return NewDomainError(
            ErrWebhookNotFound.Code,
            ErrWebhookNotFound.Message,
            ErrWebhookNotFound.HTTPStatus,
            ErrWebhookNotFound.Op,
            err,
        )
    }

    // 2. Validar token usando Domain Error
    if token != *webhook.ValidationToken {
        return ErrInvalidToken
    }

    // 3. Processar lógica específica do NewPay
    // ... sua lógica aqui

    return nil
}
```

#### Passo 5: Criar controller

```go
// controller.go
package newpay

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
    service Service
}

func NewController(service Service) *Controller {
    return &Controller{service: service}
}

func (c *Controller) Handle(ctx *fiber.Ctx) error {
    var dto NewPayDto
    if err := ctx.BodyParser(&dto); err != nil {
        return ctx.Status(400).JSON(fiber.Map{"error": "invalid payload"})
    }

    webhookID := ctx.Params("webhookID")
    token := ctx.Get("X-NewPay-Token")

    err := c.service.Handler(ctx.Context(), webhookID, token, ctx.Body(), dto)
    if err != nil {
        // Extrair DomainError se existir
        var domainErr *DomainError
        if errors.As(err, &domainErr) {
            return ctx.Status(domainErr.HTTPStatus).JSON(fiber.Map{
                "code":  domainErr.Code,
                "error": domainErr.Message,
            })
        }

        // Fallback para erros não-domínio
        return ctx.Status(500).JSON(fiber.Map{
            "code":  "INTERNAL_ERROR",
            "error": "internal server error",
        })
    }

    return ctx.Status(200).JSON(fiber.Map{"message": "success"})
}
```

#### Passo 6: Definir DTOs

```go
// dto.go
package newpay

type NewPayDto struct {
    OrderID      string  `json:"order_id"`
    CustomerName string  `json:"customer_name"`
    Email        string  `json:"email"`
    Amount       float64 `json:"amount"`
    Currency     string  `json:"currency"`
    Status       string  `json:"status"`
}
```

#### Passo 7: Registrar no inicializador

```go
// internal/routes/controllers.go
func InitializeControllers(db *sqlx.DB) *DomainControllers {
    // ... outros domínios
    
    // Initialize NewPay Domain
    newPayMetaAdapter := newpay.NewMetaAdapter()
    newPayConversionAdapter := newpay.NewConversionAdapter()
    newPayService := newpay.NewService(
        webhookRepo,
        payloadRepo,
        pixelRepo,
        eventRepo,
        salesRepo,
        newPayMetaAdapter,
        newPayConversionAdapter,
    )

    return &DomainControllers{
        // ... outros controllers
        NewPay: newpay.NewController(newPayService),
    }
}
```

#### Passo 8: Adicionar rota

```go
// internal/routes/routes.go
func SetupRoutes(app *fiber.App, db *sqlx.DB) {
    // ... outras rotas
    
    app.Post("/webhook/:webhookID/newpay", controllers.NewPay.Handle)
}
```

### Exemplo 2: Criar um Mock para Testes

```go
// service_test.go
package newpay

import (
    "context"
    "testing"
)

type mockMetaAPI struct {
    called bool
    err    error
}

func (m *mockMetaAPI) SendEvent(ctx context.Context, event meta.MetaEvent) error {
    m.called = true
    return m.err
}

func TestHandler_Success(t *testing.T) {
    // Setup mocks
    metaMock := &mockMetaAPI{}
    
    service := NewService(
        &mockWebhookRepo{...},
        metaMock,
        &mockConversionAPI{},
    )
    
    // Execute
    err := service.Handler(context.Background(), "webhook-id", "token", []byte("{}"), NewPayDto{})
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    if !metaMock.called {
        t.Error("Expected Meta API to be called")
    }
}
```

---

## ✅ Boas Práticas

### 1. Organização de Código

- ✅ Um domínio por diretório
- ✅ Contratos, adapters, service, controller separados
- ✅ DTOs em arquivo próprio
- ✅ Testes junto ao código testado

### 2. Nomenclatura

**Arquivos:**
- `contracts.go` - Interfaces
- `adapters.go` - Adaptadores
- `service.go` - Lógica de negócio
- `controller.go` - HTTP handlers
- `dto.go` - Data Transfer Objects
- `service_test.go` - Testes

**Funções:**
- Construtores: `NewService()`, `NewController()`, `NewAdapter()`
- Interfaces: `Service`, `Repository`, `API`
- Structs privados: `service`, `controller`, `adapter`

### 3. Injeção de Dependências

```go
// ✅ BOM: Receber interfaces
func NewService(repo Repository, api ExternalAPI) Service {
    return &service{repo: repo, api: api}
}

// ❌ RUIM: Criar dependências dentro
func NewService() Service {
    repo := NewRepository()  // Hard coupling
    api := NewExternalAPI()  // Hard coupling
    return &service{repo: repo, api: api}
}
```

### 4. Tratamento de Erros

```go
// ✅ BOM: Usar Domain Errors
if webhook.ValidationToken == nil {
    return ErrWebhookValidationTokenRequired
}

if token != *webhook.ValidationToken {
    return ErrInvalidToken
}

// ✅ BOM: Wrapping de erros com contexto
if err != nil {
    return NewDomainError(
        ErrSaveSale.Code,
        ErrSaveSale.Message,
        ErrSaveSale.HTTPStatus,
        ErrSaveSale.Op,
        err, // Preserva o erro original
    )
}

// ❌ RUIM: Perder contexto
if err != nil {
    return err
}

// ❌ RUIM: Decidir status HTTP no service
if err != nil {
    return fmt.Errorf("erro: %w", err) // Sem informação de status HTTP
}
```

### 5. Context Propagation

```go
// ✅ BOM: Passar context nas chamadas
func (s *service) Handler(ctx context.Context, ...) error {
    webhook, err := s.webhookRepo.GetWebhookByID(ctx, id)
    // ...
}

// ❌ RUIM: Criar novo context
func (s *service) Handler(ctx context.Context, ...) error {
    webhook, err := s.webhookRepo.GetWebhookByID(context.Background(), id)
    // ...
}
```

### 6. Retornar Interfaces, Receber Structs

```go
// ✅ BOM: Retorna interface
func NewService(repo Repository) Service {
    return &service{repo: repo}
}

// ❌ RUIM: Retorna struct
func NewService(repo Repository) *service {
    return &service{repo: repo}
}
```

### 7. Logs Estruturados

```go
// ✅ BOM: Logs com contexto
logger.Info.Printf("[HotmartService] Processing webhook %s for user %s", webhookID, userID)

// ❌ RUIM: Logs vagos
logger.Info.Printf("Processing webhook")
```

### 8. Validação de Entrada

```go
// ✅ BOM: Validar no controller
func (c *Controller) Handle(ctx *fiber.Ctx) error {
    var dto HotmartDto
    if err := ctx.BodyParser(&dto); err != nil {
        return ctx.Status(400).JSON(fiber.Map{"error": "invalid payload"})
    }
    
    if dto.OrderID == "" {
        return ctx.Status(400).JSON(fiber.Map{"error": "order_id is required"})
    }
    
    return c.service.Handler(ctx.Context(), dto)
}
```

### 9. Idempotência

```go
// ✅ BOM: Verificar se já processado
idempotencyKey := utilities.HashSHA256(string(rawBody))
if exists := s.eventRepo.ExistsByKey(ctx, idempotencyKey); exists {
    return nil  // Já processado, retorna sucesso
}
```

### 10. Testes

```go
// ✅ BOM: Testar casos de sucesso e erro
func TestHandler_Success(t *testing.T) { ... }
func TestHandler_InvalidToken(t *testing.T) { ... }
func TestHandler_WebhookNotFound(t *testing.T) { ... }

// ✅ BOM: Usar table-driven tests
func TestHandler(t *testing.T) {
    tests := []struct {
        name    string
        input   HotmartDto
        wantErr bool
    }{
        {"valid webhook", validDto, false},
        {"invalid token", invalidDto, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := service.Handler(ctx, tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("wanted error: %v, got: %v", tt.wantErr, err)
            }
        })
    }
}
```

---

## 🚀 Executando o Projeto

### Pré-requisitos

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 14+
- MongoDB 6+

### Setup

```bash
# Clone o repositório
cd goodtrack-tracker

# Copie o arquivo de ambiente
cp .env.example .env

# Configure as variáveis de ambiente
# Edite o arquivo .env com suas credenciais

# Suba os bancos de dados
docker-compose up -d postgres mongo

# Instale o Air para hot reload
go install github.com/cosmtrek/air@latest

# Execute a aplicação
air
```

### Comandos Úteis

```bash
# Rodar testes
go test ./...

# Rodar testes com coverage
go test -cover ./...

# Rodar testes de um pacote específico
go test ./internal/domains/checkouts/hotmart/...

# Build para produção
go build -o webhook-worker ./cmd/server

# Rodar com Docker
docker build -t goodtrack-tracker .
docker run -p 8181:8181 goodtrack-tracker
```

---

## 📚 Referências

- [Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Dependency Injection in Go](https://github.com/google/wire)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Fiber Framework](https://docs.gofiber.io/)

---

## 📝 Conclusão

Esta arquitetura foi desenhada para ser:

- **Escalável**: Adicionar novos domínios é simples e rápido
- **Testável**: Todas as dependências são mockáveis
- **Manutenível**: Código organizado e desacoplado
- **Flexível**: Fácil trocar implementações sem quebrar o sistema

Seguindo os padrões documentados, você conseguirá:

1. Criar novos domínios de checkout rapidamente
2. Escrever testes unitários eficientes
3. Integrar novos serviços externos
4. Manter o código limpo e organizado

**Happy Coding!** 🚀
