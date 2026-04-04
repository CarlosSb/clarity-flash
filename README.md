# AulaFlash

Grava sua aula ou reuniao e transforma em resumo profissional + flashcards inteligentes de forma discreta e automatica.

## Visão

O melhor assistente brasileiro para quem assiste aulas e reunioes (Zoom, Google Meet, Teams) e quer de forma automatica:

- **Resumo profissional** claro, com action items, decisoes e destaques
- **Flashcards inteligentes** para revisao rapida
- **Modo Quiz** para testar conhecimento

Serve para **estudantes** e **profissionais em home office**.

## Stack

| Camada | Tecnologia |
|---|---|
| Frontend | Vue 3 (Composition API) + Vite + Tailwind CSS + Pinia |
| Backend | Go (Gin) |
| Banco | PostgreSQL 15+ |
| Captura de audio | Chrome Extension (Manifest V3) com `chrome.tabCapture` |
| STT | Groq Whisper (prioridade) / fallback Whisper local |
| LLM | Hugging Face Inference (Llama 3.1 8B / Qwen2.5 7B) ou Ollama |
| Real-time | WebSocket |

## Estrutura do Projeto

```
clarityflash/
├── frontend/              # Vue 3 SPA
│   ├── src/
│   │   ├── components/    # Bento grid, flashcards, quiz, layout
│   │   ├── views/         # Home, Detalhe, Quiz
│   │   ├── store/         # Pinia stores
│   │   ├── services/      # API clients
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
│   │   ├── config/        # Environment config
│   │   ├── domain/        # Entities, repository interfaces, services
│   │   ├── handler/       # HTTP handlers
│   │   ├── middleware/    # Auth, CORS
│   │   ├── service/       # Business logic orchestrator
│   │   ├── repository/    # PostgreSQL implementation
│   │   ├── cache/         # In-memory cache
│   │   ├── websocket/     # WS handler
│   │   └── worker/        # Background worker
│   ├── pkg/               # Reusable packages
│   │   ├── audio/         # Audio processing
│   │   ├── stt/           # Groq Whisper client
│   │   ├── llm/           # Hugging Face / Ollama client
│   │   └── storage/       # Local file storage
│   ├── migrations/        # SQL migrations
│   ├── Makefile
│   ├── go.mod
│   └── go.sum
├── extension/             # Chrome Extension (MV3)
│   ├── manifest.json
│   └── src/
│       ├── background/    # Service worker - capture logic
│       ├── popup/         # Popup UI (Vue component)
│       ├── content/       # Content script
│       └── icons/         # Extension icons
├── docs/
└── .env.example
```

## Prerequisitos

- Go 1.21+
- Node.js 20+
- PostgreSQL 15+
- Chrome/Brave para carregar a extensao
- Chave de API: Groq (STT) + Hugging Face ou Ollama local (LLM)

## Setup Rapido

### 1. Banco de dados

```bash
createdb aulaflash
# rode as migrations
psql aulaflash < backend/migrations/001_initial.sql
```

### 2. Backend

```bash
cd backend
cp config/.env.example config/.env   # ajuste as chaves de API
make run
# ou
go run cmd/server/main.go
```

Servidor sobe em `http://localhost:8080`.

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

App disponivel em `http://localhost:5173`.

### 4. Chrome Extension

1. Abra `chrome://extensions/`
2. Ative **Modo desenvolvedor**
3. Clique em **Carregar sem compactacao**
4. Selecione a pasta `extension/`
5. Clique no icone da extensao para gravar audio da aba ativa

## Fluxo de Uso

1. Abra uma reuniao/aula no navegador
2. Clique no icone do AulaFlash e inicie a gravacao
3. Ao terminar, o audio e enviado ao backend automaticamente
4. O backend processa: transcricao -> resumo -> flashcards
5. Acesse a interface web para ver resultados e fazer quiz

## Roadmap

| Fase | Entregas |
|---|---|
| **MVP** | Gravacao, transcricao, resumo, flashcards, interface basica, exportacao CSV |
| **v1.0** | Assistente em tempo real, cache inteligente, modo Estudante/Profissional, dark mode completo |
| **v2.0** | Mapa mental automatico, OCR de slides, assistente avancado |

## Licenca

MIT
