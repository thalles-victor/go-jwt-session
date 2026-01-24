# 📚 API Documentation - JWT Session

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Autenticação](#autenticação)
3. [Sessões](#sessões)
4. [Recuperação de Senha](#recuperação-de-senha)
5. [Códigos de Erro](#códigos-de-erro)

---

## 🎯 Visão Geral

A API JWT Session fornece endpoints para autenticação de usuários, gerenciamento de sessões e recuperação de senha.

**Base URL:** `http://localhost:8080`

**Versão da API:** `v1`

### Formato de Resposta

Todas as respostas seguem o formato JSON:

```json
{
  "message": "mensagem de sucesso",
  "data": { ... }
}
```

### Tratamento de Erros

Erros retornam o seguinte formato:

```json
{
  "code": "ERROR_CODE",
  "error": "mensagem de erro"
}
```

---

## 🔐 Autenticação

### POST `/v1/auth/sign-up`

Cria uma nova conta de usuário.

**Request Body:**
```json
{
  "name": "João Silva",
  "email": "joao@example.com",
  "password": "senha123"
}
```

**Response 200:**
```json
{
  "message": "usuário cadastrado com sucesso",
  "user": {
    "id": "uuid",
    "name": "João Silva",
    "email": "joao@example.com",
    "created_at": "2025-01-24T10:00:00Z"
  },
  "access_token": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires": "1h"
  }
}
```

**Erros:**
- `400` - `INVALID_PAYLOAD` - Payload inválido
- `409` - `USER_ALREADY_EXISTS` - Usuário já cadastrado
- `500` - `INTERNAL_ERROR` - Erro interno do servidor

---

### POST `/v1/auth/sign-in`

Autentica um usuário existente.

**Request Body:**
```json
{
  "email": "joao@example.com",
  "password": "senha123"
}
```

**Response 200:**
```json
{
  "message": "usuário logado com sucesso",
  "user": {
    "id": "uuid",
    "name": "João Silva",
    "email": "joao@example.com",
    "created_at": "2025-01-24T10:00:00Z"
  },
  "access_token": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires": "1h"
  }
}
```

**Erros:**
- `400` - `INVALID_PAYLOAD` - Payload inválido
- `401` - `INVALID_PASSWORD` - Senha inválida
- `404` - `USER_NOT_FOUND` - Usuário não encontrado
- `500` - `INTERNAL_ERROR` - Erro interno do servidor

---

## 🔑 Sessões

### POST `/v1/auth/session/sign-up`

Cria uma nova conta de usuário com sessão.

**Request Body:**
```json
{
  "name": "João Silva",
  "email": "joao@example.com",
  "password": "senha123"
}
```

**Response 200:**
```json
{
  "message": "usuário cadastrado com sucesso",
  "user": {
    "id": "uuid",
    "name": "João Silva",
    "email": "joao@example.com",
    "created_at": "2025-01-24T10:00:00Z"
  },
  "access_token": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires": "1h"
  },
  "session": {
    "id": "uuid",
    "user_id": "uuid",
    "browser": "Mozilla/5.0...",
    "ip": "192.168.1.1",
    "created_at": "2025-01-24T10:00:00Z",
    "expires_at": "2025-01-25T10:00:00Z"
  }
}
```

**Headers:**
- `User-Agent` - Navegador do cliente (automático)
- `X-Forwarded-For` - IP do cliente (automático)

**Erros:**
- `400` - `INVALID_PAYLOAD` - Payload inválido
- `409` - `USER_ALREADY_EXISTS` - Usuário já cadastrado
- `500` - `INTERNAL_ERROR` - Erro interno do servidor

---

### POST `/v1/auth/session/sign-in`

Autentica um usuário e cria uma sessão.

**Request Body:**
```json
{
  "email": "joao@example.com",
  "password": "senha123"
}
```

**Response 200:**
```json
{
  "message": "usuário logado com sucesso",
  "user": {
    "id": "uuid",
    "name": "João Silva",
    "email": "joao@example.com",
    "created_at": "2025-01-24T10:00:00Z"
  },
  "access_token": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires": "1h"
  },
  "session": {
    "id": "uuid",
    "user_id": "uuid",
    "browser": "Mozilla/5.0...",
    "ip": "192.168.1.1",
    "created_at": "2025-01-24T10:00:00Z",
    "expires_at": "2025-01-25T10:00:00Z"
  }
}
```

**Erros:**
- `400` - `INVALID_PAYLOAD` - Payload inválido
- `401` - `INVALID_PASSWORD` - Senha inválida
- `404` - `USER_NOT_FOUND` - Usuário não encontrado
- `500` - `INTERNAL_ERROR` - Erro interno do servidor

---

### GET `/v1/auth/session/`

Valida a sessão atual do usuário.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response 200:**
```json
{
  "message": "sessão válida",
  "user_id": "uuid"
}
```

**Erros:**
- `401` - Token inválido ou expirado
- `401` - Sessão não encontrada ou expirada

---

### GET `/v1/auth/session/all`

Lista todas as sessões do usuário autenticado.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response 200:**
```json
{
  "message": "sessões encontradas com sucesso",
  "sessions": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "browser": "Mozilla/5.0...",
      "ip": "192.168.1.1",
      "created_at": "2025-01-24T10:00:00Z",
      "expires_at": "2025-01-25T10:00:00Z"
    }
  ]
}
```

**Erros:**
- `401` - Token inválido ou expirado
- `404` - `USER_NOT_FOUND` - Usuário não encontrado
- `500` - `INTERNAL_ERROR` - Erro interno do servidor

---

### DELETE `/v1/auth/session/:sessionId`

Deleta uma sessão específica do usuário.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `sessionId` - ID da sessão a ser deletada

**Response 200:**
```json
{
  "message": "sessão deletada com sucesso"
}
```

**Erros:**
- `401` - Token inválido ou expirado
- `403` - `SESSION_NOT_OWNED` - Sessão não pertence ao usuário
- `404` - `SESSION_NOT_FOUND` - Sessão não encontrada
- `404` - `USER_NOT_FOUND` - Usuário não encontrado
- `500` - `INTERNAL_ERROR` - Erro interno do servidor

---

## 🔄 Recuperação de Senha

### POST `/v1/auth/recovery/request/:email`

Solicita um código de recuperação de senha.

**Path Parameters:**
- `email` - Email do usuário

**Response 200:**
```json
{
  "message": "código de recuperação enviado com sucesso"
}
```

**Erros:**
- `404` - `USER_NOT_FOUND` - Usuário não encontrado
- `500` - `INTERNAL_ERROR` - Erro interno do servidor

**Nota:** Um email será enviado para o usuário com o código de recuperação.

---

### POST `/v1/auth/recovery/change-password`

Altera a senha usando o código de recuperação.

**Request Body:**
```json
{
  "email": "joao@example.com",
  "new_password": "novaSenha123",
  "code": "abc123xyz"
}
```

**Response 200:**
```json
{
  "message": "senha alterada com sucesso"
}
```

**Erros:**
- `400` - `INVALID_PAYLOAD` - Payload inválido
- `401` - `INVALID_RECOVERY_CODE` - Código de recuperação inválido
- `404` - `USER_NOT_FOUND` - Usuário não encontrado
- `404` - `RECOVERY_NOT_FOUND` - Dados de recuperação não encontrados
- `406` - `RECOVERY_EXPIRED` - Código de recuperação expirado
- `406` - `RECOVERY_ATTEMPTS_EXCEEDED` - Número de tentativas excedido
- `500` - `INTERNAL_ERROR` - Erro interno do servidor

**Nota:** Após alterar a senha, todas as sessões do usuário serão invalidadas.

---

## ⚠️ Códigos de Erro

### Códigos de Erro Comuns

| Código | Status HTTP | Descrição |
|--------|-------------|-----------|
| `INVALID_PAYLOAD` | 400 | Payload da requisição inválido |
| `INVALID_PASSWORD` | 401 | Senha informada é inválida |
| `INVALID_TOKEN` | 401 | Token JWT inválido ou expirado |
| `INVALID_RECOVERY_CODE` | 401 | Código de recuperação inválido |
| `USER_NOT_FOUND` | 404 | Usuário não encontrado |
| `SESSION_NOT_FOUND` | 404 | Sessão não encontrada |
| `RECOVERY_NOT_FOUND` | 404 | Dados de recuperação não encontrados |
| `USER_ALREADY_EXISTS` | 409 | Usuário já cadastrado |
| `RECOVERY_EXPIRED` | 406 | Código de recuperação expirado |
| `RECOVERY_ATTEMPTS_EXCEEDED` | 406 | Número de tentativas excedido |
| `SESSION_NOT_OWNED` | 403 | Sessão não pertence ao usuário |
| `INTERNAL_ERROR` | 500 | Erro interno do servidor |

---

## 🔒 Autenticação

### JWT Token

Para endpoints protegidos, inclua o token JWT no header:

```
Authorization: Bearer <jwt_token>
```

O token JWT contém:
- `sub`: ID do usuário ou sessão
- `exp`: Data de expiração (1 hora)

### Validação de Sessão

Para endpoints que requerem sessão (`/v1/auth/session/*`), o token JWT deve conter o ID da sessão no campo `sub`. A sessão será validada no banco de dados e deve estar ativa e não expirada.

---

## 📝 Exemplos de Uso

### Exemplo: Fluxo Completo de Autenticação

1. **Criar conta:**
```bash
curl -X POST http://localhost:8080/v1/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com",
    "password": "senha123"
  }'
```

2. **Fazer login:**
```bash
curl -X POST http://localhost:8080/v1/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao@example.com",
    "password": "senha123"
  }'
```

3. **Validar sessão:**
```bash
curl -X GET http://localhost:8080/v1/auth/session/ \
  -H "Authorization: Bearer <jwt_token>"
```

### Exemplo: Recuperação de Senha

1. **Solicitar código:**
```bash
curl -X POST http://localhost:8080/v1/auth/recovery/request/joao@example.com
```

2. **Alterar senha:**
```bash
curl -X POST http://localhost:8080/v1/auth/recovery/change-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao@example.com",
    "new_password": "novaSenha123",
    "code": "abc123xyz"
  }'
```

---

## 📚 Referências

- [Fiber Framework](https://docs.gofiber.io/)
- [JWT.io](https://jwt.io/)
