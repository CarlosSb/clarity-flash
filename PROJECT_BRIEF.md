# ClarityFlash - Documentação Completa do Projeto

## 1. Visão Geral do Produto

**Nome:** ClarityFlash  
**Tagline:** Grava sua aula ou reunião e transforma em resumo claro + flashcards inteligentes de forma discreta e automática.

**Objetivo principal:**  
Criar um assistente inteligente que escuta aulas e reuniões (Zoom, Google Meet, Teams) e entrega automaticamente:
- Resumo profissional claro e acionável
- Flashcards inteligentes para revisão rápida

O produto foi projetado para servir tanto **estudantes** quanto **profissionais em home office**, com foco em simplicidade, discrição e qualidade em português brasileiro.

## 2. Regras de Negócio

- A gravação deve ser **100% discreta** (sem bot visível na reunião para outros participantes).
- O app deve ter **dupla utilidade**: acadêmica (estudantes e professores) e profissional (reuniões de trabalho).
- Todo processamento começa gratuito (utilizando free tiers sempre que possível).
- O usuário deve ter total controle sobre o que é gravado e armazenado.
- Conteúdo similar de uma mesma aula/reunião pode ser reutilizado por múltiplos usuários (cache inteligente).
- Privacidade: nunca armazenar áudio original após o processamento; apenas texto anonimizado quando necessário.

## 3. Funcionalidades Funcionais

### MVP (Versão Inicial)

- Gravação de áudio via Chrome Extension (captura discreta da aba)
- Transcrição automática com boa qualidade em PT-BR
- Geração de resumo profissional (com action items, decisões e destaques)
- Geração automática de 10-15 flashcards
- Interface com animação de flip card
- Modo Quiz básico
- Exportação (CSV para Anki, texto simples, WhatsApp/Email)

### Versão 1.0 (Próxima)

- Assistente Inteligente leve em tempo real (dicas, action items, sugestões de respostas)
- Dois modos claros: “Estudante” e “Profissional”
- Cache básico de resumos e flashcards (reaproveitamento por hash)
- Dark mode completo

### Funcionalidades Futuras

- Mapa Mental gerado automaticamente a partir dos flashcards
- OCR de slides (visão computacional em tela compartilhada)
- Assistente mais avançado com pesquisa contextual

## 4. Funcionalidades Não-Funcionais

- **Discrição total**: Nenhum bot ou indicador visível na reunião
- **Performance**: Resumo + flashcards entregues em até 7 minutos após gravação
- **Usabilidade**: Interface limpa, intuitiva e minimalista
- **Acessibilidade**: Suporte completo a dark mode e boa legibilidade
- **Escalabilidade**: Suporte a cache para reutilização de conteúdo similar
- **Privacidade**: Anonimização de dados sensíveis e opção fácil de exclusão
- **Resiliência**: Gravação local + upload em background (para internet ruim)

## 5. Tecnologias e Stack

**Frontend:**
- Vue 3 + Vite + Tailwind CSS + Pinia

**Backend:**
- Go (Gin ou Fiber)

**Banco de dados:**
- PostgreSQL

**Captura de áudio:**
- Chrome Extension (Manifest V3) com `chrome.tabCapture`

**IA:**
- STT: Groq Whisper (prioridade) ou fallback para Whisper local
- LLM: Hugging Face Inference (Llama 3.1 8B ou Qwen2.5 7B) ou Ollama

**Real-time:**
- WebSocket (para assistente futuro)

**Hospedagem inicial:**
- Hetzner ou VPS equivalente (baixo custo)

## 6. Descrição do Layout e Design

- **Estilo geral**: Clean, minimalista e moderno (inspirado em Notion + Linear + Anki)
- **Cores principais**: Roxo/azul suave como cor principal + neutros (cinza claro/escuro)
- **Dark mode**: Padrão do app
- **Layout**: Baseado em Bento Grid (cards modulares e scannáveis)
- **Animações**: Suaves (flip dos cards, loading states, micro-interações)
- **Responsividade**: Mobile-first

**Telas principais:**
- Home: Lista de gravações em Bento Grid
- Detalhe da reunião: Resumo à esquerda + Flashcards à direita
- Flashcard: Animação de flip com frente/verso
- Quiz: Tela simples e focada

## 7. Roadmap de Funcionalidades

**MVP (6-7 semanas)**
- Gravação via Chrome Extension
- Transcrição + Resumo + Flashcards
- Interface Vue básica

**Versão 1.0 (9-11 semanas)**
- Assistente Inteligente leve em tempo real
- Cache básico
- Modo Estudante vs Profissional
- Exportações aprimoradas

**Versão 2.0 (futuro)**
- Mapa Mental automático
- OCR de slides
- Assistente mais avançado

## 8. Regras Importantes para o Desenvolvimento

- Priorizar simplicidade e velocidade de entrega.
- Manter o código modular, limpo e bem documentado.
- Toda chamada de LLM deve ser otimizada (cache quando possível).
- Feedback visual claro durante todo o processamento (progress bars).
- Testar com áudios reais em português brasileiro desde o início.

---

Este documento está pronto para ser usado como briefing completo para qualquer IA de código (Cursor, Claude, Qwen, Grok, etc.).

Quer que eu ajuste algo antes de prosseguir?  
Exemplos:
- Deixar mais curto?
- Adicionar seção de arquitetura técnica?
- Incluir exemplos de prompts para o LLM?

Me diga o que quer modificar ou se já podemos usar essa versão!
