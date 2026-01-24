# 🏗️ Arquitetura do Projeto - JWT Session

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Estrutura de Diretórios](#estrutura-de-diretórios)
3. [Padrões Arquiteturais](#padrões-arquiteturais)
4. [Camadas da Aplicação](#camadas-da-aplicação)
5. [Inversão de Dependências](#inversão-de-dependências)
6. [Domain Errors](#domain-errors)
7. [Fluxo de Dados](#fluxo-de-dados)

---

## 🎯 Visão Geral

O projeto **JWT Session** segue os princípios de **Clean Architecture** e **SOLID**, com foco em:

- **Separação de Responsabilidades**: Cada camada tem uma responsabilidade bem definida
- **Inversão de Dependências**: Dependemos de abstrações (interfaces), não de implementações concretas
- **Testabilidade**: Todas as dependências podem ser mockadas
- **Baixo Acoplamento**: Módulos independentes e reutilizáveis
- **Alta Coesão**: Componentes relacionados agrupados logicamente

### Tecnologias Principais

- **Go 1.24+** - Linguagem principal
- **Fiber v2** - Framework web
- **PostgreSQL** - Banco de dados relacional
- **JWT** - Autenticação baseada em tokens
- **Docker** - Containerização

---

## 📁 Estrutura de Diretórios

```
go-jwt-session/
│
├── cmd/
│   └── server/
│       └── main.go                    # Entry point da aplicação
│
├── internal/
│   ├── database/                      # Configurações de banco de dados
│   │   └── database.go
│   │
│   ├── domains/                      # Domínios de negócio
│   │   └── auth/                     # Domínio de autenticação
│   │       ├── contracts.go          # Interfaces (Service, JWTService, MailService, etc.)
│   │       ├── errors.go             # Domain Errors
│   │       ├── adapters.go           # Adaptadores externos (JWT, Mail, Code, Date)
│   │       ├── service.go            # Lógica de negócio
│   │       ├── controller.go         # Handler HTTP
│   │       └── dto.go                # Data Transfer Objects
│   │
│   ├── models/                       # Estruturas de dados
│   │   └── models.go                 # User, Session, Recovery
│   │
│   ├── repositories/                 # Camada de acesso a dados
│   │   └── postgres/
│   │       ├── user_repository.go
│   │       ├── session_repository.go
│   │       └── recovery_repository.go
│   │
│   ├── routes/                       # Configuração de rotas
│   │   └── routes.go                 # Setup de rotas e injeção de dependências
│   │
│   └── shared/                       # Utilitários compartilhados
│       ├── jwt/
│       │   └── jwt.go                # Geração e validação de JWT
│       ├── mail/
│       │   └── mail.go                # Envio de emails
│       ├── code/
│       │   └── code.go                # Geração de códigos
│       ├── date/
│       │   └── date.go                # Operações com datas
│       ├── config/
│       │   └── config.go              # Configurações e variáveis de ambiente
│       ├── logger/
│       │   └── logger.go              # Sistema de logs
│       └── middlewares/
│           └── middleware.go         # Middlewares HTTP
│
├── docs/                             # Documentação
│   ├── API.md
│   ├── ARCHITECTURE.md
│   └── CONFIGURATION.md
│
├── migrations/                       # Migrações do banco de dados
│   ├── create-user.sql
│   ├── create-session.sql
│   └── create-recovery.sql
│
├── go.mod                            # Dependências Go
└── README.md                         # README principal
```

---

## 🎨 Padrões Arquiteturais

### 1. Domain-Driven Design (DDD)

Cada domínio de negócio é isolado em seu próprio pacote com:

- **Entities**: Modelos de dados (`models/`)
- **Services**: Lógica de negócio
- **Repositories**: Acesso a dados
- **DTOs**: Objetos de transferência
- **Domain Errors**: Erros específicos do domínio

### 2. Dependency Injection

As dependências são injetadas via construtores, facilitando testes e manutenção:

```go
func NewService(
    repos Repositories,
    jwtSvc JWTService,
    mailSvc MailService,
    codeSvc CodeService,
    dateSvc DateService,
) Service {
    return &service{
        repos:   repos,
        jwtSvc:  jwtSvc,
        mailSvc: mailSvc,
        codeSvc: codeSvc,
        dateSvc: dateSvc,
    }
}
```

### 3. Repository Pattern

Abstrações para acesso a dados, permitindo trocar implementações facilmente:

```go
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Create(ctx context.Context, user *models.User) (*models.User, error)
    UpdatePassword(ctx context.Context, userID string, newPassword string) error
}
```

### 4. Adapter Pattern

Wrappers para serviços externos, isolando a lógica de integração:

```go
type JWTAdapter struct{}

func (a *JWTAdapter) GenerateJwt(sub string) (string, error) {
    return jwt.GenerateJwt(sub)
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
func (c *Controller) SignIn(ctx *fiber.Ctx) error {
    var dto SignInDto
    if err := ctx.BodyParser(&dto); err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "code": "INVALID_PAYLOAD",
            "error": "invalid request payload",
        })
    }
    
    result, err := c.service.SignIn(ctx.Context(), dto)
    if err != nil {
        return handleError(ctx, err)
    }
    
    return ctx.Status(200).JSON(fiber.Map{
        "message": "usuário logado com sucesso",
        "user": result.User,
        "access_token": result.AccessToken,
    })
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
func (s *service) SignIn(ctx context.Context, dto SignInDto) (*SignInResponse, error) {
    user, err := s.repos.User.GetByEmail(ctx, dto.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, NewDomainError(...)
    }
    
    if user.Password != dto.Password {
        return nil, ErrInvalidPassword
    }
    
    jwtToken, err := s.jwtSvc.GenerateJwt(user.ID)
    if err != nil {
        return nil, NewDomainError(...)
    }
    
    return &SignInResponse{
        User: user,
        AccessToken: struct{...}{JWT: jwtToken, Expires: "1h"},
    }, nil
}
```

### 3. **Repositories Layer** (Data Access)

Responsável por:
- Acesso a bancos de dados
- Queries e comandos
- Mapeamento de modelos

**Exemplo:**
```go
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := r.db.GetContext(ctx, &user, 
        `SELECT id, name, email, password, created_at, updated_at 
         FROM "users" WHERE email = $1`, email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

### 4. **Adapters Layer** (External Services)

Responsável por:
- Comunicação com serviços externos
- Wrapping de bibliotecas externas
- Isolamento de dependências

**Exemplo:**
```go
type MailAdapter struct{}

func (a *MailAdapter) SendCreateAccount(name, email string) error {
    return mail.SendCreateAccount(name, email)
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
// internal/domains/auth/contracts.go
type Service interface {
    SignUp(ctx context.Context, dto SignUpDto) (*SignUpResponse, error)
    SignIn(ctx context.Context, dto SignInDto) (*SignInResponse, error)
    // ...
}

type JWTService interface {
    GenerateJwt(sub string) (string, error)
    ParseJWT(tokenAsString string) (string, error)
}
```

**2. Implementar Services com Dependências via Interface**

```go
type service struct {
    repos     Repositories
    jwtSvc    JWTService    // Depende de interface
    mailSvc   MailService   // Depende de interface
    codeSvc   CodeService   // Depende de interface
    dateSvc   DateService   // Depende de interface
}
```

**3. Criar Adapters para Implementações Concretas**

```go
type JWTAdapter struct{}

func (a *JWTAdapter) GenerateJwt(sub string) (string, error) {
    return jwt.GenerateJwt(sub)
}
```

**4. Inicializar com Implementações Reais**

```go
// internal/routes/routes.go
jwtAdapter := auth.NewJWTAdapter()
mailAdapter := auth.NewMailAdapter()
codeAdapter := auth.NewCodeAdapter()
dateAdapter := auth.NewDateAdapter()

authService := auth.NewService(
    auth.Repositories{
        User:     userRepo,
        Session:  sessionRepo,
        Recovery: recoveryRepo,
    },
    jwtAdapter,
    mailAdapter,
    codeAdapter,
    dateAdapter,
)
```

### Benefícios

1. **Testabilidade**: Fácil criar mocks das dependências
2. **Flexibilidade**: Trocar implementações sem alterar lógica de negócio
3. **Baixo Acoplamento**: Módulos não dependem de implementações concretas
4. **Manutenibilidade**: Mudanças isoladas em cada camada

---

## ⚠️ Domain Errors

O projeto utiliza o padrão **Domain Errors** para tratamento estruturado de erros.

### Estrutura

```go
type DomainError struct {
    Code       string // ex: USER_NOT_FOUND
    Message    string // mensagem amigável
    HTTPStatus int    // código HTTP sugerido
    Op         string // operação que falhou
    Err        error  // erro original (wrapping)
}
```

### Exemplo de Uso

```go
var (
    ErrUserNotFound = NewDomainError(
        "USER_NOT_FOUND",
        "usuário não encontrado",
        http.StatusNotFound,
        "auth.signin",
        nil,
    )
)

// No service
if err == sql.ErrNoRows {
    return nil, ErrUserNotFound
}
```

### Tratamento no Controller

```go
func handleError(ctx *fiber.Ctx, err error) error {
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
```

---

## 🔄 Fluxo de Dados

### Fluxo de uma Requisição

```
1. HTTP Request
   ↓
2. Controller (validação de entrada)
   ↓
3. Service (lógica de negócio)
   ↓
4. Repository (acesso a dados)
   ↓
5. Database
   ↓
6. Repository (retorna dados)
   ↓
7. Service (processa dados)
   ↓
8. Controller (formata resposta)
   ↓
9. HTTP Response
```

### Exemplo: Sign In

```
POST /v1/auth/sign-in
   ↓
Controller.SignIn()
   ↓
Service.SignIn()
   ├─→ Repository.GetByEmail()
   │     └─→ Database Query
   ├─→ Validação de senha
   └─→ JWTAdapter.GenerateJwt()
   ↓
Controller retorna resposta
```

---

## 🧪 Testabilidade

### Estratégia de Testes

1. **Unit Tests**: Testam serviços isoladamente com mocks
2. **Integration Tests**: Testam fluxos completos (futuro)
3. **Mocks Manuais**: Implementações fake para testes

### Exemplo de Mock

```go
type mockJWTService struct {
    generateFunc func(string) (string, error)
}

func (m *mockJWTService) GenerateJwt(sub string) (string, error) {
    return m.generateFunc(sub)
}

// No teste
mockJWT := &mockJWTService{
    generateFunc: func(sub string) (string, error) {
        return "mock-token", nil
    },
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
- `errors.go` - Domain Errors

**Funções:**
- Construtores: `NewService()`, `NewController()`, `NewAdapter()`
- Interfaces: `Service`, `Repository`, `JWTService`
- Structs privados: `service`, `controller`, `adapter`

### 3. Context Propagation

```go
// ✅ BOM: Passar context nas chamadas
func (s *service) SignIn(ctx context.Context, dto SignInDto) error {
    user, err := s.repos.User.GetByEmail(ctx, dto.Email)
    // ...
}

// ❌ RUIM: Criar novo context
func (s *service) SignIn(ctx context.Context, dto SignInDto) error {
    user, err := s.repos.User.GetByEmail(context.Background(), dto.Email)
    // ...
}
```

### 4. Tratamento de Erros

```go
// ✅ BOM: Usar Domain Errors
if err == sql.ErrNoRows {
    return nil, ErrUserNotFound
}

// ✅ BOM: Wrapping de erros com contexto
if err != nil {
    return NewDomainError(
        "INTERNAL_ERROR",
        "erro ao buscar usuário",
        500,
        "auth.signin",
        err,
    )
}
```

---

## 📚 Referências

- [Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Fiber Framework](https://docs.gofiber.io/)
