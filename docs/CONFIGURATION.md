# ⚙️ Configuração do Projeto - JWT Session

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Variáveis de Ambiente](#variáveis-de-ambiente)
3. [Configuração do Banco de Dados](#configuração-do-banco-de-dados)
4. [Configuração de Email](#configuração-de-email)
5. [Configuração JWT](#configuração-jwt)
6. [Configuração de Logs](#configuração-de-logs)
7. [Exemplo de Arquivo .env](#exemplo-de-arquivo-env)
8. [Setup Inicial](#setup-inicial)

---

## 🎯 Visão Geral

O projeto utiliza variáveis de ambiente para configuração. Todas as configurações são carregadas do arquivo `.env` na raiz do projeto usando a biblioteca `godotenv`.

---

## 🔧 Variáveis de Ambiente

### Variáveis Obrigatórias

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `POSTGRES_USER` | Usuário do PostgreSQL | `postgres` |
| `POSTGRES_PASSWORD` | Senha do PostgreSQL | `mypassword` |
| `DB_HOST` | Host do banco de dados | `localhost` |
| `POSTGRES_PORT` | Porta do PostgreSQL | `5432` |
| `POSTGRES_DB` | Nome do banco de dados | `jwt_session` |
| `JWT_SEC_KEY` | Chave secreta para assinatura JWT | `your-secret-key-here` |
| `MAIL_HOST` | Host do servidor SMTP | `smtp.gmail.com` |
| `MAIL_USER` | Usuário do email | `your-email@gmail.com` |
| `MAIL_PASS` | Senha do email | `your-app-password` |
| `MAIL_PORT` | Porta do servidor SMTP | `587` |

### Variáveis Opcionais

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `LOKI_CONNECTION` | URL de conexão do Grafana Loki | `http://localhost:3100` |

---

## 🗄️ Configuração do Banco de Dados

### PostgreSQL

O projeto utiliza PostgreSQL como banco de dados. A string de conexão é montada automaticamente a partir das variáveis de ambiente.

**Formato da String de Conexão:**
```
postgres://{POSTGRES_USER}:{POSTGRES_PASSWORD}@{DB_HOST}:{POSTGRES_PORT}/{POSTGRES_DB}?sslmode=disable
```

### Exemplo de Configuração

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=mypassword
DB_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=jwt_session
```

### Migrações

Execute as migrações SQL na ordem:

1. `migrations/create-user.sql` - Tabela de usuários
2. `migrations/create-session.sql` - Tabela de sessões
3. `migrations/create-recovery.sql` - Tabela de recuperação

**Exemplo:**
```bash
psql -U postgres -d jwt_session -f migrations/create-user.sql
psql -U postgres -d jwt_session -f migrations/create-session.sql
psql -U postgres -d jwt_session -f migrations/create-recovery.sql
```

---

## 📧 Configuração de Email

### SMTP

O projeto utiliza SMTP para envio de emails. Configure as credenciais do seu provedor de email.

### Gmail

Para usar Gmail, você precisa:

1. Ativar a verificação em duas etapas
2. Gerar uma senha de app:
   - Acesse: https://myaccount.google.com/apppasswords
   - Gere uma senha de app
   - Use essa senha no `MAIL_PASS`

**Configuração:**
```env
MAIL_HOST=smtp.gmail.com
MAIL_USER=your-email@gmail.com
MAIL_PASS=your-app-password
MAIL_PORT=587
```

### Outros Provedores

#### Outlook/Hotmail
```env
MAIL_HOST=smtp-mail.outlook.com
MAIL_PORT=587
```

#### SendGrid
```env
MAIL_HOST=smtp.sendgrid.net
MAIL_USER=apikey
MAIL_PASS=your-sendgrid-api-key
MAIL_PORT=587
```

#### Mailgun
```env
MAIL_HOST=smtp.mailgun.org
MAIL_PORT=587
```

---

## 🔐 Configuração JWT

### Chave Secreta

A chave secreta JWT (`JWT_SEC_KEY`) é usada para assinar e validar tokens. **Nunca compartilhe esta chave publicamente**.

**Recomendações:**
- Use uma chave forte (mínimo 32 caracteres)
- Gere uma chave aleatória para produção
- Não commite a chave no repositório

**Gerar uma chave segura:**
```bash
# Linux/Mac
openssl rand -base64 32

# Ou usando Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

**Exemplo:**
```env
JWT_SEC_KEY=your-very-secret-key-minimum-32-characters-long
```

### Expiração do Token

A expiração padrão do token JWT é de **1 hora**. Para alterar, edite o arquivo:
`internal/shared/jwt/jwt.go`

```go
"exp": time.Now().Add(time.Hour * 1).Unix(), // Altere aqui
```

---

## 📊 Configuração de Logs

### Grafana Loki (Opcional)

O projeto suporta envio de logs para Grafana Loki. Se não configurado, os logs serão apenas salvos localmente.

**Configuração:**
```env
LOKI_CONNECTION=http://localhost:3100
```

**Formato da URL:**
- Com autenticação: `http://user:password@localhost:3100`
- Sem autenticação: `http://localhost:3100`

### Logs Locais

Os logs são salvos automaticamente no arquivo `app.log` na raiz do projeto.

**Níveis de Log:**
- `INFO` - Informações gerais
- `WARN` - Avisos
- `ERROR` - Erros
- `DEBUG` - Debug (desenvolvimento)

---

## 📝 Exemplo de Arquivo .env

Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:

```env
# ============================================
# Database Configuration
# ============================================
POSTGRES_USER=postgres
POSTGRES_PASSWORD=mypassword
DB_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=jwt_session

# ============================================
# JWT Configuration
# ============================================
JWT_SEC_KEY=your-very-secret-key-minimum-32-characters-long-change-in-production

# ============================================
# Email Configuration (SMTP)
# ============================================
MAIL_HOST=smtp.gmail.com
MAIL_USER=your-email@gmail.com
MAIL_PASS=your-app-password
MAIL_PORT=587

# ============================================
# Logging Configuration (Optional)
# ============================================
LOKI_CONNECTION=http://localhost:3100
```

### .env.example

É recomendado criar um arquivo `.env.example` no repositório (sem valores sensíveis) para servir como template:

```env
# Database Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
DB_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=jwt_session

# JWT Configuration
JWT_SEC_KEY=your-secret-key-here

# Email Configuration
MAIL_HOST=smtp.gmail.com
MAIL_USER=your-email@gmail.com
MAIL_PASS=your-app-password
MAIL_PORT=587

# Logging Configuration (Optional)
LOKI_CONNECTION=http://localhost:3100
```

---

## 🚀 Setup Inicial

### 1. Clonar o Repositório

```bash
git clone <repository-url>
cd go-jwt-session
```

### 2. Instalar Dependências

```bash
go mod download
```

### 3. Configurar Banco de Dados

```bash
# Criar banco de dados
createdb jwt_session

# Executar migrações
psql -U postgres -d jwt_session -f migrations/create-user.sql
psql -U postgres -d jwt_session -f migrations/create-session.sql
psql -U postgres -d jwt_session -f migrations/create-recovery.sql
```

### 4. Criar Arquivo .env

```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

### 5. Executar a Aplicação

```bash
# Desenvolvimento
go run cmd/server/main.go

# Ou compilar e executar
go build -o jwt-session cmd/server/main.go
./jwt-session
```

---

## 🔒 Segurança

### Boas Práticas

1. **Nunca commite o arquivo `.env`**
   - Adicione `.env` ao `.gitignore`
   - Use `.env.example` como template

2. **Use chaves fortes em produção**
   - JWT_SEC_KEY: mínimo 32 caracteres aleatórios
   - POSTGRES_PASSWORD: senha forte

3. **Proteja as credenciais**
   - Use variáveis de ambiente do sistema em produção
   - Considere usar um gerenciador de secrets (Vault, AWS Secrets Manager, etc.)

4. **Configuração de Email**
   - Use senhas de app (não a senha principal)
   - Considere usar serviços de email transacional (SendGrid, Mailgun, etc.)

### .gitignore

Certifique-se de que seu `.gitignore` inclui:

```
.env
.env.local
app.log
*.log
```

---

## 🐳 Docker

### Docker Compose (Exemplo)

Crie um arquivo `docker-compose.yml` para facilitar o desenvolvimento:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: mypassword
      POSTGRES_DB: jwt_session
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  app:
    build: .
    ports:
      - "8080:8080"
    env_file:
      - .env
    depends_on:
      - postgres

volumes:
  postgres_data:
```

---

## 🧪 Ambiente de Desenvolvimento vs Produção

### Desenvolvimento

- Use banco de dados local
- Logs detalhados (DEBUG)
- Chaves de desenvolvimento (não use em produção)

### Produção

- Use banco de dados gerenciado
- Logs apenas ERROR e WARN
- Chaves fortes e únicas
- Variáveis de ambiente do sistema (não arquivo .env)
- SSL/TLS habilitado
- Rate limiting configurado

---

## 📚 Referências

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Gmail App Passwords](https://support.google.com/accounts/answer/185833)
- [JWT.io](https://jwt.io/)
- [Grafana Loki](https://grafana.com/docs/loki/latest/)

---

## ❓ Troubleshooting

### Erro: "error when loading .env file"

**Solução:** Certifique-se de que o arquivo `.env` existe na raiz do projeto.

### Erro: "error when connect in the database"

**Solução:** 
- Verifique se o PostgreSQL está rodando
- Confirme as credenciais no `.env`
- Teste a conexão: `psql -U postgres -h localhost`

### Erro: "MAIL_HOST not declared"

**Solução:** Adicione todas as variáveis de email obrigatórias no `.env`.

### Erro ao enviar email

**Solução:**
- Verifique as credenciais SMTP
- Para Gmail, use senha de app (não a senha principal)
- Verifique se a porta está correta (587 para TLS)
