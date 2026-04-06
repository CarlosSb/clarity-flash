# ClarityFlash - Agent Skills & Development Guide

## Visão Geral do Projeto

ClarityFlash é um assistente inteligente que grava aulas/reuniões via Chrome Extension e transforma em:
- Resumo profissional com destaques, decisões, action items e conceitos-chave
- 10-15 flashcards de estudo gerados por IA
- Modo Quiz para revisão interativa
- Exportação para CSV (Anki) e texto simples

## Stack Tecnológico

### Frontend
- Vue 3.4.27 + Vite 5.2.12 + TypeScript 5.4.5
- Tailwind CSS 3.4.4 + Pinia 2.1.7 + Vue Router 4.3.2

### Backend
- Go 1.21 + PostgreSQL (lib/pq) + WebSocket (gorilla/websocket)

### IA / Modelos
| Função | Provedor | Modelo | Config |
|---|---|---|---|
| Transcrição (STT) | Groq Whisper | `whisper-large-v3` | `GROQ_MODEL` |
| Geração de Conteúdo (LLM) | Hugging Face | `meta-llama/llama-3.1-8b-instruct` | `LLM_MODEL` |
| LLM Alternativo | Ollama local | Configurável | `USE_OLLAMA=true` |

### Infraestrutura
- PostgreSQL 15+ (Neon DB cloud)
- Chrome Extension (Manifest V3, tabCapture API)

## Skills Necessárias para Desenvolvimento

1. **Go** - Backend development, Clean Architecture, HTTP handlers
2. **Vue 3 + TypeScript** - Frontend development, Composition API, Pinia
3. **PostgreSQL** - Database management, SQL migrations, JSONB
4. **WebSockets** - Real-time communication (gorilla/websocket)
5. **Chrome Extension (Manifest V3)** - Audio capture with tabCapture API
6. **AI/LLM Integration** - Groq API, Hugging Face Inference, prompt engineering
7. **Audio Processing** - WAV conversion, streaming, format validation

## Como Executar

### Pré-requisitos
- Node.js 20+
- Go 1.21+
- PostgreSQL 15+ (ou Neon DB)
- Groq API key (https://console.groq.com)
- Chrome/Brave browser

### Setup Rápido

```bash
# 1. Banco de dados
createdb clarityflash
psql clarityflash < backend/migrations/001_initial.sql
psql clarityflash < backend/migrations/002_add_auth.sql

# 2. Backend
cd backend
cp ../.env.example ../.env
# Edite .env com suas chaves de API
make run  # Servidor em http://localhost:8081

# 3. Frontend
cd frontend
npm install
npm run dev  # App em http://localhost:5173

# 4. Chrome Extension
# chrome://extensions/ -> Modo desenvolvedor -> Carregar sem compactação -> extension/
```

### CLI do Projeto (Makefile)

| Comando | Descrição |
|---|---|
| `make -C backend run` | Executa servidor em modo desenvolvimento |
| `make -C backend build` | Compila binário em `backend/bin/clarityflash` |
| `make -C backend migrate` | Executa migrações do banco |
| `make -C backend test` | Roda testes |
| `make -C backend clean` | Remove binários e uploads temporários |

### CLI do Projeto (Makefile)

| Comando | Descrição |
|---|---|
| `make -C backend run` | Executa servidor em modo desenvolvimento |
| `make -C backend build` | Compila binário em `backend/bin/clarityflash` |
| `make -C backend migrate` | Executa migrações do banco |
| `make -C backend test` | Roda testes |
| `make -C backend clean` | Remove binários e uploads temporários |

## Status do Desenvolvimento (MVP ~80%)

- ✅ Chrome Extension com captura de áudio
- ✅ Backend API em Go
- ✅ Upload de áudio (file + streaming)
- ✅ Transcrição Groq Whisper Large V3
- ✅ Geração de resumo com Llama 3.1 8B (via Groq)
- ✅ Geração de flashcards (10-15 cards)
- ✅ PostgreSQL com migrações
- ✅ WebSocket real-time
- ✅ Vue 3 frontend
- ✅ Auth básico (JWT + fallback)
- ✅ Exportação CSV/TXT
- ✅ CLI do projeto (Makefile)
- ⚠️ Flashcard UI - testes pendentes
- ⚠️ Modo Quiz - testes pendentes
- ❌ Assistente IA tempo real (v1.0)
- ❌ Cache inteligente (v1.0)
- ❌ Modos Estudante/Profissional (v1.0)
- ❌ Mapa mental (v2.0)
- ❌ OCR de slides (v2.0)

## Comandos Úteis

```bash
# Backend
make -C backend run       # Executar servidor
make -C backend build     # Compilar binário
make -C backend migrate   # Executar migrações
make -C backend test      # Rodar testes
make -C backend clean     # Limpar temporários

# Frontend
cd frontend && npm run dev    # Dev server
cd frontend && npm run build  # Build produção
```

## Estrutura de Arquivos

```
clarity-flash/
├── frontend/          # Vue 3 SPA
├── backend/           # Go server
│   ├── cmd/           # Entry points
│   ├── internal/      # App code (api, handler, service, repository)
│   ├── pkg/           # Reusable packages (stt, llm, audio, storage)
│   └── migrations/    # SQL migrations
├── extension/         # Chrome Extension (MV3)
├── docs/              # Documentação técnica
└── .env.example       # Template de variáveis de ambiente
```

## API Endpoints Principais

| Método | Rota | Descrição |
|---|---|---|
| GET | /health | Health check |
| POST | /api/auth/register | Registro |
| POST | /api/auth/login | Login |
| POST | /api/sessions/upload | Upload áudio |
| GET | /api/sessions/{id} | Detalhes sessão |
| GET | /api/sessions | Lista sessões |
| GET | /api/export/{id}/csv | Export CSV |
| GET | /api/export/{id}/txt | Export TXT |

## Modelos de IA - Detalhes

### Transcrição (STT)
- **Modelo**: Groq Whisper `whisper-large-v3`
- **Endpoint**: `https://api.groq.com/openai/v1/audio/transcriptions`
- **Parâmetros**: language=`pt`, response_format=`json`

### Geração de Conteúdo (LLM)
- **Modelo**: `meta-llama/llama-3.1-8b-instruct`
- **Endpoint**: `https://api-inference.huggingface.co/models/{model}`
- **Parâmetros**: max_new_tokens=2048, temperature=0.3
- **Prompts**: Em português brasileiro, saída JSON forçada
- **Usado para**: Resumos estruturados e geração de flashcards

## Checklist para Novos Desenvolvedores

- [ ] Instalar Node.js 20+, Go 1.21+, PostgreSQL 15+
- [ ] Criar conta no Groq (https://console.groq.com) e obter API key
- [ ] Configurar `.env` com credenciais
- [ ] Executar migrações do banco
- [ ] Rodar backend (`make run`)
- [ ] Rodar frontend (`npm run dev`)
- [ ] Carregar Chrome Extension
- [ ] Testar fluxo completo: gravar → processar → visualizar → estudar
