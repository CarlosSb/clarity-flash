# ClarityFlash 🎓

> Transforme aulas e reuniões em conhecimento estruturado com IA

O ClarityFlash é um assistente inteligente que **grava áudio automaticamente** e converte em:
- 📝 **Resumos profissionais** com destaques, decisões e conceitos-chave
- 🃏 **Flashcards de estudo** (10-15 cartões gerados por IA)
- 🎯 **Modo quiz** para revisão interativa
- 📊 **Exportação** para Anki e outros formatos

## 🚀 Demonstração Rápida

```bash
# 1. Instalar dependências
npm install && go mod tidy

# 2. Configurar banco e APIs
createdb clarityflash
cp .env.example .env  # Editar com suas chaves

# 3. Executar tudo
npm run dev &  # Frontend na porta 5173
go run backend/cmd/server/main.go  # Backend na porta 8081
```

Acesse `http://localhost:5173` e faça upload de um arquivo de áudio!

## 🏗️ Como Funciona (Arquitetura)

```
🎤 Áudio → Extensão Chrome → Backend Go → IA Whisper → IA Llama → Resumo + Flashcards
```

### Fluxo Detalhado:
1. **Upload**: Usuário envia áudio via extensão Chrome
2. **Transcrição**: IA Whisper converte fala → texto
3. **Análise**: IA Llama gera resumo estruturado
4. **Flashcards**: IA cria 10-15 cartões de estudo
5. **Resultado**: Frontend exibe tudo organizado

## 🛠️ Tecnologias

| Componente | Tecnologia | Propósito |
|------------|------------|-----------|
| **Frontend** | Vue 3 + TypeScript | Interface web moderna |
| **Backend** | Go + Fiber | API REST performante |
| **Banco** | PostgreSQL | Dados estruturados |
| **IA Fala** | Groq Whisper Large V3 | Transcrição automática |
| **IA Texto** | Llama 3.1 8B | Geração de conteúdo |
| **Extensão** | Chrome Manifest V3 | Captura de áudio |

## 📦 Instalação Completa

### Pré-requisitos
- **Node.js 20+** (frontend)
- **Go 1.21+** (backend)
- **PostgreSQL 15+** (banco)
- **Chrome** (extensão)
- **Chaves API** do [Groq](https://console.groq.com)

### Passo a Passo

```bash
# 1. Clonar repositório
git clone https://github.com/seu-usuario/clarityflash.git
cd clarityflash

# 2. Configurar banco
createdb clarityflash
psql clarityflash < backend/migrations/001_initial.sql

# 3. Configurar ambiente
cp .env.example .env
# Editar .env com GROQ_API_KEY e outras configurações

# 4. Instalar dependências
cd frontend && npm install
cd ../backend && go mod tidy

# 5. Executar
cd frontend && npm run dev     # http://localhost:5173
cd ../backend && go run cmd/server/main.go  # http://localhost:8081
```

## 🎯 Como Usar

### 1. Fazer Login/Cadastro
- Acesse o frontend
- Crie conta ou faça login

### 2. Gravar Áudio
- Instale a extensão Chrome
- Clique no botão da extensão durante aula/reunião
- Ou faça upload direto de arquivo MP3/WAV

### 3. Ver Resultado
- Aguarde processamento (1-2 minutos)
- Visualize resumo estruturado
- Estude com flashcards gerados
- Exporte para Anki

## 🔧 Desenvolvimento

### Scripts Úteis
```bash
# Backend
make -C backend run     # Servidor dev
make -C backend test    # Executar testes
make -C backend build   # Build produção

# Frontend
npm run dev            # Servidor dev
npm run build          # Build produção
npm run preview        # Preview build
```

### Estrutura do Código
```
backend/
├── cmd/server/         # Ponto de entrada
├── internal/
│   ├── api/           # Configuração rotas
│   ├── auth/          # Autenticação JWT
│   ├── handler/       # Controllers HTTP
│   ├── service/       # Lógica negócio
│   ├── repository/    # Acesso banco
│   └── middleware/    # Middlewares
└── pkg/               # Utilitários

frontend/
├── src/
│   ├── components/    # Componentes Vue
│   ├── views/         # Páginas
│   ├── stores/        # Estado (Pinia)
│   └── router/        # Rotas
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie branch: `git checkout -b feature/nova-funcionalidade`
3. Commit: `git commit -m "feat: descrição clara"`
4. Push: `git push origin feature/nova-funcionalidade`
5. Abra Pull Request

### Áreas para Contribuir
- ✨ Melhorar UI/UX
- 🤖 Otimizar prompts de IA
- 📱 Suporte mobile
- 🌐 Internacionalização
- 📊 Analytics e métricas

## 📄 Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

## 🙏 Agradecimentos

- **Groq** por APIs de IA acessíveis
- **Vue.js** por framework frontend incrível
- **Go** por linguagem backend performante
- **Comunidade open source** por todas as ferramentas

---

**Feito com ❤️ para estudantes e profissionais que querem aprender melhor**</content>
<parameter name="filePath">README.md