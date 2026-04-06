# ClarityFlash

Grava sua aula ou reunião e transforma em resumo claro + flashcards inteligentes de forma discreta e automática.

## Visão Geral

ClarityFlash é um assistente inteligente que escuta aulas e reuniões (Zoom, Google Meet, Teams) via Chrome Extension e entrega automaticamente:

- **Resumo profissional** com destaques, decisões, action items e conceitos-chave
- **10-15 flashcards** de estudo gerados por IA com frente/verso/nível de dificuldade
- **Modo Quiz** para revisão interativa
- **Exportação** para CSV (compatível com Anki) e texto simples

O produto atende tanto **estudantes** quanto **profissionais em home office**, com foco em simplicidade, discrição e qualidade em português brasileiro.

## Stack Tecnológico

### Frontend
| Tecnologia | Versão |
|---|---|
| Vue | 3.4.27 |
| Vite | 5.2.12 |
| TypeScript | 5.4.5 |
| Tailwind CSS | 3.4.4 |
| Pinia | 2.1.7 |
| Vue Router | 4.3.2 |
| Axios | 1.7.2 |

### Backend
| Tecnologia | Versão |
|---|---|
| Go | 1.21 |
| lib/pq (PostgreSQL driver) | 1.10.9 |
| gorilla/websocket | 1.5.1 |
| godotenv | 1.5.1 |

### IA / Modelos
| Função | Provedor | Modelo |
|---|---|---|
| **Transcrição (STT)** | Groq Whisper | `whisper-large-v3` |
| **Geração de Conteúdo (LLM)** | Hugging Face Inference API | `meta-llama/llama-3.1-8b-instruct` |
| **LLM Alternativo** | Ollama (local) | Configurável via `.env` |

- **Transcrição**: Groq Whisper Large V3 — converte áudio em texto em português (PT-BR)
- **Resumo**: Llama 3.1 8B Instruct — gera resumos profissionais com highlights, decisões, action items e conceitos-chave
- **Flashcards**: Llama 3.1 8B Instruct — cria 10-15 cards com estrutura front/back/difficulty
- Prompts em português brasileiro com formato de saída JSON forçado
- Modelos configuráveis via `.env`: `GROQ_MODEL` e `LLM_MODEL`

### Infraestrutura
- **Banco de dados**: PostgreSQL 15+ (atualmente Neon DB — cloud-hosted)
- **Chrome Extension**: Manifest V3 com `tabCapture` API
- **Comunicação real-time**: WebSocket (gorilla/websocket)

## Pré-requisitos

- Node.js 20+
- Go 1.21+
- PostgreSQL 15+ (ou conta Neon DB)
- Chave de API Groq (obtenha em https://console.groq.com)
- Chave de API Hugging Face (ou Ollama local como alternativa)
- Chrome/Brave browser para a extensão

## Como Executar

### 1. Banco de Dados

```bash
# Criar banco local
createdb aulaflash

# Executar migrações
psql aulaflash < backend/migrations/001_initial.sql
psql aulaflash < backend/migrations/002_add_auth.sql

# Ou via Makefile
make -C backend migrate
```

### 2. Backend

```bash
cd backend
cp ../.env.example ../.env
# Edite .env com suas chaves de API
make run
# ou: go run cmd/server/main.go
```

O servidor inicia em `http://localhost:8081`.

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

A aplicação fica disponível em `http://localhost:5173`.

### 4. Chrome Extension

1. Abra `chrome://extensions/` (ou `brave://extensions/`)
2. Ative o **Modo desenvolvedor**
3. Clique em **Carregar sem compactação**
4. Selecione a pasta `extension/`
5. Clique no ícone da extensão para iniciar a gravação da aba ativa

### Comandos Disponíveis (Makefile)

| Comando | Descrição |
|---|---|
| `make run` | Executa o servidor em modo desenvolvimento |
| `make build` | Compila o binário em `bin/aulaflash` |
| `make migrate` | Executa as migrações do banco |
| `make test` | Roda os testes |
| `make clean` | Remove binários e uploads temporários |

## Fluxo de Uso

1. Abra uma reunião/aula no navegador
2. Clique no ícone do ClarityFlash e inicie a gravação
3. Ao terminar, o áudio é enviado ao backend automaticamente
4. O backend processa: **transcrição → resumo → flashcards**
5. Acesse a interface web para ver resultados, estudar com flashcards e fazer quiz

## Estrutura do Projeto

```
clarity-flash/
├── frontend/              # Vue 3 SPA
│   ├── src/
│   │   ├── components/    # Bento grid, flashcards, quiz, layout
│   │   ├── views/         # Home, SessionDetail, Quiz
│   │   ├── store/         # Pinia stores
│   │   ├── services/      # API clients (axios)
│   │   ├── composables/   # Reusable logic
│   │   ├── router/        # Vue Router
│   │   └── styles/        # Global CSS + Tailwind
│   ├── index.html
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── package.json
├── backend/
│   ├── cmd/               # Entry points (server, worker, migrate)
│   ├── internal/          # Private application code
│   │   ├── api/           # Router and routes
│   │   ├── auth/          # JWT token service
│   │   ├── config/        # Environment config
│   │   ├── domain/        # Entities, repository interfaces, services
│   │   ├── handler/       # HTTP handlers (session, auth, export, health)
│   │   ├── middleware/    # Auth (JWT + fallback), CORS
│   │   ├── service/       # Business logic (Processor, AuthService)
│   │   ├── repository/    # PostgreSQL implementation
│   │   ├── cache/         # In-memory cache
│   │   ├── websocket/     # WS hub for real-time events
│   │   ├── worker/        # Background worker
│   │   └── logger/        # Logging
│   ├── pkg/               # Reusable packages
│   │   ├── audio/         # Audio processing (WAV conversion, validation)
│   │   ├── stt/           # Groq Whisper client
│   │   ├── llm/           # Hugging Face / Ollama client
│   │   └── storage/       # Local file storage
│   ├── migrations/        # SQL migrations (001_initial, 002_add_auth)
│   ├── Makefile
│   ├── go.mod
│   └── go.sum
├── extension/             # Chrome Extension (MV3)
│   ├── manifest.json
│   └── src/
│       ├── background/    # Service worker - capture logic
│       ├── popup/         # Popup UI
│       ├── content/       # Content script
│       └── icons/         # Extension icons
├── docs/
├── .env.example
└── README.md
```

## API Endpoints

### Saúde
| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/health` | Health check |

### WebSocket
| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/ws` | WebSocket para eventos em tempo real |
| `GET` | `/ws-stream` | WebSocket stream de áudio (frames binários, extensão) |

### Autenticação
| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/api/auth/register` | Registro de usuário |
| `POST` | `/api/auth/login` | Login de usuário |

### Sessões
| Método | Rota | Auth | Descrição |
|---|---|---|---|
| `POST` | `/api/sessions/upload` | Sim | Upload de áudio para processamento |
| `POST` | `/api/sessions/stream-init` | Não | Cria sessão vazia para streaming |
| `PATCH` | `/api/sessions/{id}/audio-chunk` | Não | Recebe chunk de áudio em streaming |
| `POST` | `/api/sessions/{id}/audio-complete` | Não | Sinaliza fim do streaming e inicia processamento |
| `POST` | `/api/sessions/{id}/upload-complete` | Não | Upload completo de áudio (fallback offline) |
| `GET` | `/api/sessions/{id}` | Sim | Detalhes de uma sessão |
| `GET` | `/api/sessions` | Sim | Lista sessões do usuário |
| `DELETE` | `/api/sessions/{id}` | Sim | Deleta uma sessão |

### Exportação
| Método | Rota | Auth | Descrição |
|---|---|---|---|
| `GET` | `/api/export/{id}/csv` | Sim | Exporta flashcards em CSV (compatível Anki) |
| `GET` | `/api/export/{id}/txt` | Sim | Exporta flashcards em texto simples |

> **Auth**: Rotas marcadas como "Sim" requerem JWT token no header `Authorization: Bearer <token>` ou fallback via header `X-User-ID`.

## Status do Desenvolvimento

### MVP: ~80% Completo

- ✅ Chrome Extension com captura de áudio (tabCapture)
- ✅ Backend API em Go com Clean Architecture
- ✅ Upload de áudio (arquivo completo e streaming)
- ✅ Transcrição com Groq Whisper Large V3
- ✅ Geração de resumo com Llama 3.1 8B (Hugging Face)
- ✅ Geração de flashcards (10-15 cards por sessão)
- ✅ PostgreSQL com migrações (users, sessions, flashcards)
- ✅ WebSocket para atualizações em tempo real
- ✅ Frontend Vue 3 com Tailwind CSS
- ✅ Sistema básico de autenticação (JWT + fallback)
- ✅ Exportação (CSV para Anki, TXT)
- ⚠️ Componentes de Flashcard (FlipCard, DeckList) — necessita testes
- ⚠️ Modo Quiz — necessita testes
- ❌ Assistente IA em tempo real (planejado para v1.0)
- ❌ Sistema de cache inteligente (planejado para v1.0)
- ❌ Modos Estudante/Profissional (planejado para v1.0)
- ❌ Geração de mapa mental (planejado para v2.0)
- ❌ OCR de slides (planejado para v2.0)

## Habilidades Necessárias para Desenvolvimento

- **Go** — desenvolvimento backend, Clean Architecture
- **Vue 3 + TypeScript** — desenvolvimento frontend
- **PostgreSQL** — gerenciamento de banco de dados e migrações SQL
- **WebSockets** — comunicação em tempo real
- **Chrome Extension (Manifest V3)** — captura de áudio com tabCapture API
- **Integração com IA/LLM** — Groq API, Hugging Face Inference, prompt engineering
- **Processamento de áudio** — conversão WAV, streaming, validação

## Configuração do .env

```env
# Backend
SERVER_PORT=8081
DATABASE_URL=postgres://postgres:postgres@localhost:5432/aulaflash?sslmode=disable
UPLOAD_DIR=/tmp/aulaflash-uploads

# Groq (STT)
GROQ_API_KEY=sua_chave_aqui
GROQ_MODEL=whisper-large-v3

# LLM (HuggingFace ou Ollama)
HUGGING_FACE_TOKEN=sua_chave_aqui
LLM_MODEL=meta-llama/llama-3.1-8b-instruct

# Ollama (alternativa local)
USE_OLLAMA=false
OLLAMA_URL=http://localhost:11434

# Auth (opcional)
API_KEY=sua_chave_aqui
```

## Licença

MIT
