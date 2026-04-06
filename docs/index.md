# ClarityFlash — Documentação Técnica Completa

## 1. Visão Geral do Produto

**Nome:** ClarityFlash  
**Tagline:** Grava sua aula ou reunião e transforma em resumo claro + flashcards inteligentes de forma discreta e automática.

**Objetivo principal:**  
Criar um assistente inteligente que escuta aulas e reuniões (Zoom, Google Meet, Teams) via Chrome Extension e entrega automaticamente:
- Resumo profissional claro e acionável com destaques, decisões, action items e conceitos-chave
- 10-15 flashcards inteligentes para revisão rápida com níveis de dificuldade
- Modo Quiz para estudo interativo
- Exportação para CSV (compatível com Anki) e texto simples

O produto atende tanto **estudantes** quanto **profissionais em home office**, com foco em simplicidade, discrição e qualidade em português brasileiro.

---

## 2. Regras de Negócio

- A gravação deve ser **100% discreta** (sem bot visível na reunião para outros participantes)
- O app tem **dupla utilidade**: acadêmica (estudantes e professores) e profissional (reuniões de trabalho)
- Todo processamento começa gratuito (utilizando free tiers sempre que possível)
- O usuário tem total controle sobre o que é gravado e armazenado
- Conteúdo similar de uma mesma aula/reunião pode ser reutilizado por múltiplos usuários (cache inteligente — planejado)
- **Privacidade**: nunca armazenar áudio original após o processamento; apenas texto processado
- O áudio original é deletado automaticamente após o pipeline de processamento

---

## 3. Stack Tecnológico Detalhado

### 3.1 Frontend

| Tecnologia | Versão | Uso |
|---|---|---|
| Vue | 3.4.27 | Framework reativo principal |
| Vite | 5.2.12 | Build tool e dev server |
| TypeScript | 5.4.5 | Tipagem estática |
| Tailwind CSS | 3.4.4 | Utility-first CSS framework |
| Pinia | 2.1.7 | State management |
| Vue Router | 4.3.2 | Roteamento client-side |
| Axios | 1.7.2 | HTTP client para API |

### 3.2 Backend

| Tecnologia | Versão | Uso |
|---|---|---|
| Go | 1.21 | Linguagem principal do servidor |
| lib/pq | 1.10.9 | Driver PostgreSQL puro |
| gorilla/websocket | 1.5.1 | Servidor WebSocket |
| godotenv | 1.5.1 | Carregamento de variáveis de ambiente |

### 3.3 IA / Modelos

| Função | Provedor | Modelo | Configuração |
|---|---|---|---|
| **Transcrição (STT)** | Groq Whisper | `whisper-large-v3` | `GROQ_MODEL` |
| **Geração de Conteúdo (LLM)** | Hugging Face Inference API | `meta-llama/llama-3.1-8b-instruct` | `LLM_MODEL` |
| **LLM Alternativo** | Ollama (local) | Configurável | `USE_OLLAMA=true` |

**Detalhes dos Modelos:**

- **Transcrição**: Groq Whisper Large V3 — converte áudio em texto em português brasileiro (PT-BR). O parâmetro `language` é fixado como `"pt"` na requisição.
- **Resumo**: Llama 3.1 8B Instruct via Hugging Face Inference API — gera resumos profissionais com estrutura JSON contendo: título, descrição, highlights, decisões, action items e conceitos-chave.
- **Flashcards**: Llama 3.1 8B Instruct — cria 10-15 cards com estrutura `{front, back, difficulty}` onde difficulty varia de 1 (fácil) a 3 (difícil).
- **Prompts** são em português brasileiro e exigem formato de saída JSON puro (sem markdown adicional).
- **Fallback**: Ollama local pode ser ativado definindo `USE_OLLAMA=true` no `.env`.

### 3.4 Infraestrutura

- **Banco de dados**: PostgreSQL 15+ (atualmente Neon DB — cloud-hosted com SSL)
- **Chrome Extension**: Manifest V3 com `chrome.tabCapture` API para captura de áudio da aba
- **Comunicação real-time**: WebSocket (gorilla/websocket) para eventos de processamento
- **Armazenamento de arquivos**: Local (`/tmp/aulaflash-uploads` por padrão), deletado após processamento

---

## 4. Arquitetura do Sistema

### 4.1 Visão Geral

```
┌─────────────────┐      Áudio (WAV)      ┌──────────────────┐
│  Chrome Ext.    │ ──────────────────────►│   Backend Go     │
│  (tabCapture)   │   WebSocket / HTTP     │   (:8081)        │
└─────────────────┘                       └────────┬─────────┘
                                                   │
                                    ┌──────────────┼──────────────┐
                                    │              │              │
                                    ▼              ▼              ▼
                           ┌─────────────┐ ┌─────────────┐ ┌──────────┐
                           │ Groq Whisper│ │  LLM (HF)   │ │PostgreSQL│
                           │   (STT)     │ │  ou Ollama  │ │  (Neon)  │
                           └─────────────┘ └─────────────┘ └──────────┘
                                                   │
                                    ┌──────────────┴──────────────┐
                                    ▼                             ▼
                           ┌──────────────┐              ┌──────────────┐
                           │  Frontend    │              │  WebSocket   │
                           │  Vue 3 (:5173)│              │  Hub         │
                           └──────────────┘              └──────────────┘
```

### 4.2 Pipeline de Processamento

O fluxo de processamento de áudio segue estas etapas:

1. **Captura**: Chrome Extension captura áudio da aba via `chrome.tabCapture`
2. **Upload**: Áudio enviado ao backend via WebSocket stream ou HTTP multipart
3. **Validação**: Backend valida que o arquivo é áudio válido
4. **Conversão**: Áudio convertido para WAV (formato exigido pelo Groq Whisper)
5. **Transcrição**: Groq Whisper converte WAV em texto (PT-BR)
6. **Resumo**: LLM gera resumo estruturado em JSON a partir da transcrição
7. **Flashcards**: LLM gera 10-15 flashcards em JSON a partir da transcrição
8. **Persistência**: Dados salvos no PostgreSQL
9. **Limpeza**: Arquivo de áudio original é deletado (privacidade)
10. **Notificação**: Status atualizado via WebSocket para o frontend

### 4.3 Estrutura de Diretórios

```
clarity-flash/
├── frontend/                 # Vue 3 SPA
│   ├── src/
│   │   ├── components/       # Componentes reutilizáveis
│   │   │   ├── bento/        # Bento Grid cards para listagem
│   │   │   ├── flashcard/    # FlipCard.vue, DeckList.vue
│   │   │   ├── layout/       # Header, Footer, Sidebar
│   │   │   ├── quiz/         # QuizSession.vue
│   │   │   └── ui/           # Componentes UI genéricos
│   │   ├── views/            # Páginas da aplicação
│   │   │   ├── HomeView.vue       # Lista de sessões
│   │   │   ├── SessionDetailView.vue  # Detalhe + flashcards
│   │   │   └── QuizView.vue       # Modo quiz
│   │   ├── store/            # Pinia stores (estado global)
│   │   ├── services/         # API clients (api.ts, auth.ts)
│   │   ├── composables/      # Lógica reutilizável
│   │   ├── router/           # Vue Router config
│   │   ├── directives/       # Custom Vue directives
│   │   ├── utils/            # Funções utilitárias
│   │   ├── styles/           # Global CSS + Tailwind
│   │   ├── assets/           # Imagens, ícones
│   │   ├── App.vue
│   │   └── main.ts
│   ├── index.html
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── package.json
├── backend/
│   ├── cmd/                  # Entry points
│   │   ├── server/main.go    # Servidor HTTP principal
│   │   ├── migrate/main.go   # Executor de migrações
│   │   └── worker/main.go    # Worker em background
│   ├── internal/             # Código privado da aplicação
│   │   ├── api/routes.go     # Definição de todas as rotas
│   │   ├── auth/             # Serviço de tokens JWT
│   │   ├── config/config.go  # Configuração via .env
│   │   ├── domain/           # Camada de domínio
│   │   │   ├── model/        # Entidades (Summary, Flashcard)
│   │   │   ├── repository/   # Interfaces de repositório
│   │   │   └── service/      # Interfaces de serviço
│   │   ├── handler/          # HTTP handlers
│   │   │   ├── session.go    # Handlers de sessão (upload, list, get, delete)
│   │   │   ├── auth.go       # Handlers de auth (register, login)
│   │   │   ├── export.go     # Handlers de export (CSV, TXT)
│   │   │   └── health.go     # Health check
│   │   ├── middleware/       # Middleware HTTP
│   │   │   └── auth.go       # JWT + fallback X-User-ID
│   │   ├── service/          # Lógica de negócio
│   │   │   ├── processor.go  # Pipeline: upload→transcrição→resumo→flashcards
│   │   │   └── auth.go       # AuthService (register, login)
│   │   ├── repository/       # Implementações PostgreSQL
│   │   ├── cache/            # Cache em memória
│   │   ├── websocket/        # WebSocket Hub para eventos real-time
│   │   ├── worker/           # Background worker
│   │   └── logger/           # Logging
│   ├── pkg/                  # Pacotes reutilizáveis
│   │   ├── audio/            # Processamento de áudio (WAV, validação)
│   │   ├── stt/groq.go       # Cliente Groq Whisper API
│   │   ├── llm/huggingface.go # Cliente Hugging Face + Ollama
│   │   └── storage/          # Armazenamento local de arquivos
│   ├── migrations/           # Migrações SQL
│   │   ├── 001_initial.sql   # Tabelas: users, sessions, flashcards
│   │   └── 002_add_auth.sql  # Colunas de autenticação
│   ├── tests/                # Testes
│   ├── scripts/              # Scripts auxiliares
│   ├── Makefile
│   ├── go.mod
│   └── go.sum
├── extension/                # Chrome Extension (Manifest V3)
│   ├── manifest.json
│   ├── assets/
│   └── src/
│       ├── background/       # Service worker (captura de áudio)
│       ├── popup/            # UI do popup
│       ├── content/          # Content script
│       └── icons/            # Ícones da extensão (16, 48, 128)
├── docs/
│   └── index.md              # Esta documentação
├── .env.example              # Template de variáveis de ambiente
├── .gitignore
└── README.md
```

---

## 5. Banco de Dados

### 5.1 Schema

#### Tabela: `users`

| Coluna | Tipo | Restrições | Descrição |
|---|---|---|---|
| `id` | VARCHAR(64) | PRIMARY KEY | ID único do usuário (32 chars hex) |
| `name` | VARCHAR(255) | | Nome do usuário |
| `email` | VARCHAR(255) | UNIQUE | Email (usado para login) |
| `password_hash` | TEXT | NULLABLE | Hash da senha (bcrypt) |
| `mode` | VARCHAR(20) | DEFAULT 'student' | Modo: 'student' ou 'professional' |
| `created_at` | TIMESTAMP | DEFAULT NOW() | Data de criação |
| `updated_at` | TIMESTAMP | DEFAULT NOW() | Data de atualização |

#### Tabela: `sessions`

| Coluna | Tipo | Restrições | Descrição |
|---|---|---|---|
| `id` | VARCHAR(64) | PRIMARY KEY | ID único da sessão |
| `user_id` | VARCHAR(64) | FK → users(id) | Dono da sessão |
| `title` | VARCHAR(500) | | Título da sessão |
| `description` | TEXT | | Descrição |
| `duration` | INTEGER | DEFAULT 0 | Duração em segundos |
| `status` | VARCHAR(20) | DEFAULT 'processing' | Status: processing, completed, failed |
| `mode` | VARCHAR(20) | DEFAULT 'student' | Modo da sessão |
| `transcript` | TEXT | | Transcrição completa do áudio |
| `audio_path` | VARCHAR(1000) | | Caminho do arquivo (temporário) |
| `summary_data` | JSONB | | Resumo estruturado (JSON) |
| `created_at` | TIMESTAMP | DEFAULT NOW() | Data de criação |
| `updated_at` | TIMESTAMP | DEFAULT NOW() | Data de atualização |

#### Tabela: `flashcards`

| Coluna | Tipo | Restrições | Descrição |
|---|---|---|---|
| `id` | VARCHAR(64) | PRIMARY KEY | ID único do flashcard |
| `session_id` | VARCHAR(64) | FK → sessions(id) ON DELETE CASCADE | Sessão origem |
| `front` | TEXT | NOT NULL | Pergunta (frente do card) |
| `back` | TEXT | NOT NULL | Resposta (verso do card) |
| `difficulty` | INTEGER | DEFAULT 2 | 1=fácil, 2=médio, 3=difícil |
| `known` | BOOLEAN | DEFAULT FALSE | Se o usuário já sabe |
| `created_at` | TIMESTAMP | DEFAULT NOW() | Data de criação |

### 5.2 Índices

```sql
idx_sessions_user_id    ON sessions(user_id)
idx_sessions_status     ON sessions(status)
idx_sessions_created    ON sessions(created_at DESC)
idx_flashcards_session  ON flashcards(session_id)
idx_users_email         ON users(email)
```

---

## 6. API Endpoints

### 6.1 Saúde

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/health` | Health check — retorna status do servidor |

### 6.2 WebSocket

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/ws` | WebSocket para eventos em tempo real (identificado por `user_id` query param) |
| `GET` | `/ws-stream` | WebSocket stream de áudio — recebe frames binários da extensão |

### 6.3 Autenticação

| Método | Rota | Request Body | Response |
|---|---|---|---|
| `POST` | `/api/auth/register` | `{name, email, password, mode?}` | `{token, user: {id, name, email, mode}}` |
| `POST` | `/api/auth/login` | `{email, password}` | `{token, user: {id, name, email, mode}}` |

### 6.4 Sessões

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| `POST` | `/api/sessions/upload` | Sim | Upload de arquivo de áudio completo (multipart/form-data, campo `audio`, max 50MB) |
| `POST` | `/api/sessions/stream-init` | Não | Cria sessão vazia para streaming — retorna `session_id` |
| `PATCH` | `/api/sessions/{id}/audio-chunk` | Não | Recebe chunk de áudio em streaming (HTTP fallback) |
| `POST` | `/api/sessions/{id}/audio-complete` | Não | Sinaliza fim do streaming e inicia processamento |
| `POST` | `/api/sessions/{id}/upload-complete` | Não | Upload completo de áudio (fallback offline da extensão) |
| `GET` | `/api/sessions/{id}` | Sim | Retorna detalhes de uma sessão (com summary_data e flashcards) |
| `GET` | `/api/sessions` | Sim | Lista sessões do usuário (query: `user_id`, limite 50) |
| `DELETE` | `/api/sessions/{id}` | Sim | Deleta uma sessão e seus flashcards |

### 6.5 Exportação

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| `GET` | `/api/export/{id}/csv` | Sim | Exporta flashcards em CSV (colunas: Front, Back, Difficulty) — compatível com Anki |
| `GET` | `/api/export/{id}/txt` | Sim | Exporta flashcards em texto simples (formato: Card N / Q: / A:) |

> **Autenticação**: Rotas protegidas aceitam JWT no header `Authorization: Bearer <token>` ou fallback via header `X-User-ID` (para a extensão).

### 6.6 Exemplos de Request

**Upload de áudio:**
```bash
curl -X POST http://localhost:8081/api/sessions/upload \
  -H "X-User-ID: anonymous" \
  -F "audio=@recording.wav" \
  -F "user_id=user123" \
  -F "mode=student"
```

**Listar sessões:**
```bash
curl http://localhost:8081/api/sessions?user_id=user123 \
  -H "Authorization: Bearer <jwt_token>"
```

**Exportar flashcards CSV:**
```bash
curl http://localhost:8081/api/export/session123/csv \
  -H "Authorization: Bearer <jwt_token>" \
  -o flashcards.csv
```

---

## 7. Modelos de IA e Prompt Engineering

### 7.1 Transcrição (STT)

**Modelo**: Groq Whisper `whisper-large-v3`  
**Endpoint**: `https://api.groq.com/openai/v1/audio/transcriptions`  
**Parâmetros**:
- `model`: `whisper-large-v3`
- `language`: `pt`
- `response_format`: `json`

**Implementação**: `backend/pkg/stt/groq.go`

### 7.2 Geração de Resumo

**Modelo**: `meta-llama/llama-3.1-8b-instruct`  
**Endpoint**: `https://api-inference.huggingface.co/models/{model}`  
**Parâmetros**:
- `max_new_tokens`: 2048
- `temperature`: 0.3
- `return_full_text`: false

**Prompt** (definido em `backend/internal/domain/model/summary.go`):
```
Você é um assistente especializado em analisar transcrições de aulas e reuniões em português brasileiro.
Analise a transcrição abaixo e gere um resumo profissional seguindo EXATAMENTE o formato JSON.

Regras:
- Seja conciso e claro
- Destaque apenas o que é relevante
- Action items devem ser práticos e com verbo no infinitivo
- Se não houver decisões ou action items, use arrays vazios

Retorne APENAS o JSON válido, sem markdown ou explicações extras.
Formato esperado:
{"title":"título curto","description":"resumo em 2-3 frases","highlights":["destaque1"],"decisions":["decisão1"],"action_items":["ação1"],"key_concepts":["conceito1"]}
```

**Estrutura de saída** (`Summary`):
```json
{
  "title": "string",
  "description": "string",
  "highlights": ["string"],
  "decisions": ["string"],
  "action_items": ["string"],
  "key_concepts": ["string"]
}
```

### 7.3 Geração de Flashcards

**Modelo**: `meta-llama/llama-3.1-8b-instruct` (mesmo cliente do resumo)

**Prompt** (definido em `backend/internal/domain/model/flashcard.go`):
```
Você é um especialista em criação de flashcards para estudo.
Com base na transcrição abaixo, crie exatamente 10-15 flashcards em português brasileiro.

Regras:
- Cada flashcard deve ter uma pergunta clara na frente e resposta objetiva no verso
- Cubra os conceitos mais importantes do conteúdo
- Evite perguntas triviais ou de "sim/não"
- Varie a dificuldade entre fácil, médio e difícil

Retorne APENAS um JSON válido no formato abaixo, sem markdown adicional:
{"cards":[{"front":"pergunta","back":"resposta","difficulty":1},...]}
```

**Estrutura de saída** (`FlashcardDeck`):
```json
{
  "session_id": "string",
  "title": "string",
  "description": "string",
  "cards": [
    {
      "front": "string",
      "back": "string",
      "difficulty": 1
    }
  ]
}
```

O campo `difficulty` usa a escala: **1** = fácil, **2** = médio, **3** = difícil.

### 7.4 Fallback Ollama

Para usar Ollama local como alternativa ao Hugging Face:

```env
USE_OLLAMA=true
OLLAMA_URL=http://localhost:11434
LLM_MODEL=llama3.1:8b
```

O cliente Ollama (`backend/pkg/llm/huggingface.go:86-142`) usa o endpoint `/api/generate` com `stream: false`.

---

## 8. Chrome Extension

### 8.1 Configuração

**Manifest**: V3  
**Nome**: AulaFlash - Gravar e Resumir  
**Versão**: 0.1.0

**Permissões**:
- `tabCapture` — captura de áudio da aba ativa
- `storage` — armazenamento local de configurações
- `activeTab` — acesso à aba ativa
- `scripting` — injeção de scripts

**Host Permissions**:
- `http://localhost:8081/*` — backend local

### 8.2 Componentes

| Arquivo | Função |
|---|---|
| `src/background/` | Service worker — gerencia captura de áudio e envio |
| `src/popup/` | Popup UI — controles de gravação |
| `src/content/` | Content script — interação com a página |

### 8.3 Fluxo de Captura

1. Usuário clica no ícone da extensão
2. Popup exibe controles (gravar/parar)
3. Ao iniciar, `chrome.tabCapture.getMediaStream` captura áudio da aba
4. Áudio é enviado ao backend via:
   - **WebSocket stream** (`/ws-stream`) — envio em tempo real de frames binários
   - **HTTP chunks** (`PATCH /api/sessions/{id}/audio-chunk`) — fallback
5. Ao parar, sinaliza conclusão (`POST /api/sessions/{id}/audio-complete`)
6. Backend inicia pipeline de processamento automaticamente

### 8.4 Instalação

1. Abra `chrome://extensions/` ou `brave://extensions/`
2. Ative **Modo desenvolvedor** (toggle no canto superior direito)
3. Clique em **Carregar sem compactação**
4. Selecione a pasta `extension/` do projeto
5. O ícone do ClarityFlash aparece na barra de extensões

---

## 9. Frontend — Componentes e Views

### 9.1 Views

| Componente | Rota | Descrição |
|---|---|---|
| `HomeView.vue` | `/` | Lista de sessões em Bento Grid |
| `SessionDetailView.vue` | `/session/:id` | Resumo + flashcards de uma sessão |
| `QuizView.vue` | `/quiz/:id` | Modo quiz interativo |

### 9.2 Componentes Principais

| Componente | Caminho | Descrição |
|---|---|---|
| `FlipCard.vue` | `components/flashcard/` | Card com animação de flip (frente/verso) |
| `DeckList.vue` | `components/flashcard/` | Lista de flashcards de uma sessão |
| `QuizSession.vue` | `components/quiz/` | Interface de quiz com pontuação |
| Bento Grid | `components/bento/` | Cards modulares para listagem de sessões |

### 9.3 Serviços

| Módulo | Caminho | Descrição |
|---|---|---|
| `api.ts` | `services/api.ts` | Axios client com interceptor `X-User-ID` |
| `auth.ts` | `services/auth.ts` | Funções de autenticação |

### 9.4 Autenticação no Frontend

O frontend utiliza o header `X-User-ID` injetado via interceptor do Axios, lido do `localStorage`:

```typescript
api.interceptors.request.use((config) => {
  const userId = localStorage.getItem('user_id')
  if (userId) {
    config.headers['X-User-ID'] = userId
  }
  return config
})
```

---

## 10. Funcionalidades Não-Funcionais

- **Discrição total**: Nenhum bot ou indicador visível na reunião
- **Performance**: Resumo + flashcards entregues em até 7 minutos após gravação
- **Usabilidade**: Interface limpa, intuitiva e minimalista (inspirado em Notion + Linear + Anki)
- **Acessibilidade**: Suporte completo a dark mode e boa legibilidade
- **Escalabilidade**: Suporte a cache para reutilização de conteúdo similar (planejado)
- **Privacidade**: Anonimização de dados sensíveis e opção fácil de exclusão; áudio original deletado após processamento
- **Resiliência**: Gravação local + upload em background (para internet ruim)

---

## 11. Design e Layout

- **Estilo geral**: Clean, minimalista e moderno (inspirado em Notion + Linear + Anki)
- **Cores principais**: Roxo/azul suave como cor principal + neutros (cinza claro/escuro)
- **Dark mode**: Padrão do app
- **Layout**: Baseado em Bento Grid (cards modulares e scannáveis)
- **Animações**: Suaves (flip dos cards, loading states, micro-interações)
- **Responsividade**: Mobile-first

---

## 12. Status do Desenvolvimento

### MVP: ~80% Completo

- ✅ Chrome Extension com captura de áudio (tabCapture)
- ✅ Backend API em Go com Clean Architecture
- ✅ Upload de áudio (arquivo completo e streaming via WebSocket/HTTP)
- ✅ Transcrição com Groq Whisper Large V3
- ✅ Geração de resumo com Llama 3.1 8B (Hugging Face Inference API)
- ✅ Geração de flashcards (10-15 cards por sessão)
- ✅ PostgreSQL com migrações (users, sessions, flashcards + auth)
- ✅ WebSocket para atualizações em tempo real
- ✅ Frontend Vue 3 com Tailwind CSS
- ✅ Sistema básico de autenticação (JWT + fallback X-User-ID)
- ✅ Exportação (CSV para Anki, TXT)
- ⚠️ Componentes de Flashcard (FlipCard, DeckList) — necessita testes
- ⚠️ Modo Quiz — necessita testes
- ❌ Assistente IA em tempo real (planejado para v1.0)
- ❌ Sistema de cache inteligente (planejado para v1.0)
- ❌ Modos Estudante/Profissional (planejado para v1.0)
- ❌ Geração de mapa mental (planejado para v2.0)
- ❌ OCR de slides (planejado para v2.0)

---

## 13. Roadmap

### MVP (Atual — ~80%)
- Gravação via Chrome Extension
- Transcrição + Resumo + Flashcards
- Interface Vue básica com flashcards e quiz
- Exportação CSV/TXT

### Versão 1.0
- Assistente Inteligente leve em tempo real (dicas, action items, sugestões)
- Cache básico de resumos e flashcards (reaproveitamento por hash)
- Modo Estudante vs Profissional
- Dark mode completo
- Exportações aprimoradas (WhatsApp, Email)

### Versão 2.0
- Mapa Mental automático gerado a partir dos flashcards
- OCR de slides (visão computacional em tela compartilhada)
- Assistente mais avançado com pesquisa contextual

---

## 14. Como Executar

### Pré-requisitos

- Node.js 20+
- Go 1.21+
- PostgreSQL 15+ (ou conta Neon DB)
- Chave de API Groq (https://console.groq.com)
- Chave de API Hugging Face (ou Ollama local)
- Chrome/Brave browser

### Passo a Passo

**1. Banco de Dados**
```bash
createdb aulaflash
psql aulaflash < backend/migrations/001_initial.sql
psql aulaflash < backend/migrations/002_add_auth.sql
# Ou: make -C backend migrate
```

**2. Backend**
```bash
cd backend
cp ../.env.example ../.env
# Edite .env com suas chaves de API
make run
# Servidor em http://localhost:8081
```

**3. Frontend**
```bash
cd frontend
npm install
npm run dev
# App em http://localhost:5173
```

**4. Chrome Extension**
1. Abra `chrome://extensions/`
2. Ative **Modo desenvolvedor**
3. Clique em **Carregar sem compactação**
4. Selecione a pasta `extension/`

### Comandos do Makefile

| Comando | Descrição |
|---|---|
| `make run` | Executa servidor em modo desenvolvimento |
| `make build` | Compila binário em `bin/aulaflash` |
| `make migrate` | Executa migrações do banco |
| `make test` | Roda testes (`go test -v ./...`) |
| `make clean` | Remove binários e uploads temporários |

---

## 15. Configuração do .env

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

---

## 16. Habilidades Necessárias para Desenvolvimento

- **Go** — desenvolvimento backend, Clean Architecture, HTTP handlers
- **Vue 3 + TypeScript** — desenvolvimento frontend, Composition API, Pinia
- **PostgreSQL** — gerenciamento de banco de dados, migrações SQL, JSONB
- **WebSockets** — comunicação em tempo real (gorilla/websocket)
- **Chrome Extension (Manifest V3)** — captura de áudio com tabCapture API, service workers
- **Integração com IA/LLM** — Groq API, Hugging Face Inference, Ollama, prompt engineering
- **Processamento de áudio** — conversão WAV, streaming, validação de formato

---

## 17. Regras de Desenvolvimento

- Priorizar simplicidade e velocidade de entrega
- Manter o código modular, limpo e bem documentado
- Toda chamada de LLM deve ser otimizada (cache quando possível)
- Feedback visual claro durante todo o processamento (progress bars)
- Testar com áudios reais em português brasileiro desde o início
- Áudio original nunca deve persistir após o processamento

---

## 18. Fluxo de Uso Completo

1. **Iniciar gravação**: Abra uma reunião/aula no navegador, clique no ícone do ClarityFlash e inicie a gravação
2. **Processamento automático**: Ao terminar, o áudio é enviado ao backend que executa: transcrição → resumo → flashcards
3. **Visualizar resultados**: Acesse `http://localhost:5173` para ver a lista de sessões
4. **Detalhes**: Clique em uma sessão para ver o resumo e os flashcards gerados
5. **Estudar com flashcards**: Use os cards com animação de flip para estudar
6. **Modo Quiz**: Acesse o quiz para testar seu conhecimento com os cards da sessão
7. **Exportar**: Exporte flashcards em CSV (para importar no Anki) ou em texto simples

---

## 19. Licença

MIT
