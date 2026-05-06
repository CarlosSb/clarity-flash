# Guia Didático do ClarityFlash

## Bem-vindo ao ClarityFlash! 🎓

Olá! Este documento é um guia completo para entender o projeto ClarityFlash. Mesmo que você seja iniciante em desenvolvimento, vamos explicar tudo passo a passo de forma clara e tranquila. Vamos começar?

## 📖 O que é o ClarityFlash?

Imagine que você está em uma aula ou reunião importante, mas não consegue anotar tudo. O **ClarityFlash** resolve isso!

**O que ele faz:**
- ✅ Grava áudio de aulas/reuniões automaticamente
- ✅ Converte fala em texto (usando IA)
- ✅ Gera um resumo profissional com pontos-chave
- ✅ Cria flashcards de estudo (10-15 cartões)
- ✅ Permite revisar o conteúdo com quizzes interativos
- ✅ Exporta para Anki (app de flashcards) ou texto simples

**Como funciona no dia a dia:**
1. Você clica no botão da extensão no Chrome
2. O áudio da aula/reunião é gravado
3. IA processa tudo automaticamente
4. Você recebe resumo + flashcards prontos para estudo

## 🏗️ Arquitetura do Sistema

Vamos entender como o sistema funciona por dentro. Imagine uma casa:

- **Frontend**: A "fachada" bonita que você vê (interface web)
- **Backend**: A "estrutura" que processa tudo nos bastidores (servidor Go)
- **Banco de Dados**: O "porão" onde guardamos as informações (PostgreSQL)
- **Extensão Chrome**: O "portão de entrada" que captura o áudio
- **IA**: Os "cérebros" que entendem fala e geram conteúdo

### Fluxo de Dados

```
🎤 Áudio → Extensão Chrome → Backend → IA (Transcrição) → IA (Resumo + Flashcards) → Banco → Frontend
```

## 🛠️ Tecnologias Usadas

Não se preocupe se você não conhece tudo! Vamos explicar cada uma:

### Frontend (Interface Web)
- **Vue 3**: Framework JavaScript moderno para criar interfaces web
- **TypeScript**: JavaScript com "superpoderes" de tipagem (menos erros!)
- **Tailwind CSS**: Sistema de estilos pronto para usar
- **Pinia**: Gerenciador de estado (como uma "memória" da aplicação)
- **Vite**: Ferramenta ultra-rápida para desenvolvimento

### Backend (Servidor)
- **Go 1.21**: Linguagem de programação rápida e confiável
- **PostgreSQL**: Banco de dados avançado para armazenar dados
- **WebSocket**: Tecnologia para comunicação em tempo real
- **Fiber**: Framework web moderno para Go (versão 3)

### IA e Processamento
- **Groq Whisper Large V3**: IA que converte fala em texto
- **Llama 3.1 8B**: IA que gera resumos e flashcards
- **FFmpeg**: Ferramenta para processar áudio (converter formatos)

### Infraestrutura
- **Chrome Extension (Manifest V3)**: Extensão do navegador para capturar áudio
- **Docker**: "Caixinha" que empacota o projeto para qualquer computador

## 📁 Estrutura do Projeto

Vamos ver como o código está organizado:

```
clarity-flash/
├── frontend/           # Interface web (Vue.js)
│   ├── src/
│   │   ├── components/ # Peças reutilizáveis da interface
│   │   ├── views/      # Páginas principais
│   │   ├── stores/     # Estado da aplicação (Pinia)
│   │   └── router/     # Configuração de rotas
│   └── package.json    # Dependências do frontend
│
├── backend/            # Servidor (Go)
│   ├── cmd/server/     # Ponto de entrada da aplicação
│   ├── internal/       # Código interno (privado)
│   │   ├── api/        # Configuração de rotas HTTP
│   │   ├── auth/       # Sistema de autenticação
│   │   ├── handler/    # Controladores HTTP (como "recepcionistas")
│   │   ├── middleware/ # Código que roda entre requisições
│   │   ├── service/    # Lógica de negócio (regras do projeto)
│   │   └── repository/ # Acesso ao banco de dados
│   ├── pkg/            # Pacotes reutilizáveis
│   │   ├── audio/      # Processamento de áudio
│   │   ├── llm/        # Integração com IA de texto
│   │   ├── stt/        # Integração com IA de fala
│   │   └── storage/    # Sistema de arquivos
│   └── migrations/     # Scripts para criar tabelas no banco
│
├── extension/          # Extensão do Chrome
│   └── manifest.json   # Configuração da extensão
│
└── docs/              # Documentação
```

## 🚀 Como Rodar o Projeto

Vamos configurar tudo passo a passo!

### Pré-requisitos
Antes de começar, você precisa ter instalado:
- **Node.js 20+** (para o frontend)
- **Go 1.21+** (para o backend)
- **PostgreSQL 15+** (banco de dados)
- **Google Chrome** (para testar a extensão)

### Passo 1: Configurar o Banco de Dados

```bash
# Criar banco de dados
createdb clarityflash

# Executar migrações (criar tabelas)
psql clarityflash < backend/migrations/001_initial.sql
psql clarityflash < backend/migrations/002_add_auth.sql
```

### Passo 2: Configurar o Backend

```bash
# Entrar na pasta do backend
cd backend

# Copiar arquivo de exemplo de configuração
cp ../.env.example ../.env

# Editar .env com suas chaves de API
# (você precisa de uma chave da Groq: https://console.groq.com)

# Instalar dependências
go mod tidy

# Executar servidor
go run cmd/server/main.go
```

O servidor vai rodar em `http://localhost:8081`

### Passo 3: Configurar o Frontend

```bash
# Entrar na pasta do frontend
cd frontend

# Instalar dependências
npm install

# Executar servidor de desenvolvimento
npm run dev
```

O frontend vai rodar em `http://localhost:5173`

### Passo 4: Configurar a Extensão do Chrome

```bash
# Abrir chrome://extensions/
# Ativar "Modo desenvolvedor"
# Clicar em "Carregar sem compactação"
# Selecionar a pasta extension/
```

## 🔄 Como Funciona o Fluxo de Processamento

Vamos entender o "coração" do sistema - como um áudio vira flashcards:

### 1. Upload do Áudio
```go
// O usuário faz upload via extensão Chrome
// Handler recebe o arquivo multipart
func (h *SessionHandler) FiberUpload(c fiber.Ctx) error {
    file, err := c.FormFile("audio")
    // ... validações ...
    return h.processor.Process(c.Context(), session, src, file)
}
```

### 2. Processamento Assíncrono
```go
// O processamento roda em background
func (p *Processor) Process(ctx context.Context, session *repository.Session, file multipart.File, header *multipart.FileHeader) error {
    // 1. Salvar arquivo
    // 2. Iniciar goroutine para processamento
    go func() {
        p.runPipeline(context.Background(), session, path)
    }()
    return nil
}
```

### 3. Pipeline de IA
```go
func (p *Processor) runPipeline(ctx context.Context, session *repository.Session, audioPath string) error {
    // 1. Validar áudio
    // 2. Converter para WAV
    // 3. Transcrever com Groq Whisper
    transcript, err := p.sttClient.Transcribe(wavPath)

    // 4. Gerar resumo com Llama
    summaryPrompt := model.SummaryPrompt(transcript)
    summaryJSON, err := p.llmClient.Generate(ctx, summaryPrompt)

    // 5. Gerar flashcards
    flashcardPrompt := model.FlashcardPrompt(transcript)
    flashcardJSON, err := p.llmClient.Generate(ctx, flashcardPrompt)

    // 6. Salvar no banco
    return p.sessionRepo.UpdateStatus(ctx, session.ID, "completed")
}
```

## 🎯 Conceitos Importantes

### API REST
O backend expõe uma API REST. Pense nela como um restaurante:
- **GET**: Pedir informação (como ver o cardápio)
- **POST**: Enviar dados (fazer um pedido)
- **DELETE**: Remover algo

Exemplos:
- `GET /health` - Verificar se o servidor está funcionando
- `POST /api/sessions/upload` - Enviar áudio para processamento
- `GET /api/sessions` - Listar sessões do usuário

### Middleware
Código que roda automaticamente entre cada requisição:
```go
func FiberJWTOrFallbackAuth(tokenService *auth.TokenService, fallbackHeader string) fiber.Handler {
    return func(c fiber.Ctx) error {
        // Verificar token JWT
        // Ou usar fallback para extensão
        // Colocar user_id no contexto
        return c.Next() // Continuar para o handler
    }
}
```

### WebSocket
Permite comunicação em tempo real:
```go
// Cliente se conecta: ws://localhost:8081/ws
// Servidor pode enviar atualizações de progresso
// "Processamento: 25% concluído..."
```

## 🧪 Testando o Sistema

### Teste Básico
1. Acesse `http://localhost:5173` (frontend)
2. Faça login ou registre-se
3. Use a extensão para "gravar" um áudio de teste
4. Veja o processamento em tempo real
5. Visualize o resumo e flashcards gerados

### Endpoints da API
```bash
# Health check
curl http://localhost:8081/health

# Listar sessões (precisa auth)
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8081/api/sessions

# Upload (via form-data com arquivo)
curl -X POST -F "audio=@test.wav" -F "user_id=test" http://localhost:8081/api/sessions/upload
```

## 📚 Conceitos de Programação Aprendidos

Este projeto ensina muitos conceitos importantes:

### Backend
- **Clean Architecture**: Separação clara entre camadas
- **Dependency Injection**: Injeção de dependências para teste
- **Goroutines**: Processamento assíncrono
- **Context**: Controle de cancelamento e timeout
- **Middleware Pattern**: Código reutilizável entre requisições

### Frontend
- **Component-based Architecture**: Interface em componentes
- **State Management**: Gerenciamento de estado com Pinia
- **Reactive Programming**: Interface que se atualiza automaticamente
- **Routing**: Navegação entre páginas

### Infraestrutura
- **API Design**: Como criar boas APIs REST
- **Authentication**: JWT e fallback para extensões
- **File Upload**: Processamento de arquivos multipart
- **WebSockets**: Comunicação bidirecional

## 🔧 Desenvolvimento e Contribuição

### Como Contribuir
1. **Fork** o projeto no GitHub
2. **Clone** sua cópia local
3. Crie uma **branch** para sua feature: `git checkout -b feature/nova-funcionalidade`
4. **Implemente** suas mudanças
5. **Teste** tudo
6. **Commit** com mensagens claras
7. **Push** e crie um Pull Request

### Comandos Úteis
```bash
# Backend
make -C backend run      # Executar servidor
make -C backend build    # Compilar binário
make -C backend test     # Rodar testes
make -C backend migrate  # Executar migrações

# Frontend
cd frontend && npm run dev    # Servidor de desenvolvimento
cd frontend && npm run build  # Build de produção
```

## 🎉 Conclusão

Parabéns! Agora você entende como o ClarityFlash funciona. Este projeto combina muitas tecnologias modernas:

- **IA**: Para processamento inteligente de áudio e texto
- **Web Development**: Frontend moderno e backend robusto
- **DevOps**: Containers, bancos de dados, APIs
- **Arquitetura**: Padrões de design e boas práticas

O mais importante: este projeto resolve um problema real das pessoas que querem aprender de forma mais eficiente. Cada linha de código contribui para ajudar estudantes e profissionais a reterem conhecimento melhor.

**Próximos passos para aprender mais:**
1. Rode o projeto localmente
2. Experimente modificar algo pequeno
3. Leia o código dos handlers e services
4. Implemente uma nova funcionalidade
5. Contribua com a comunidade!

Boa sorte na sua jornada de aprendizado! 🚀</content>
<parameter name="filePath">GUIDE.md