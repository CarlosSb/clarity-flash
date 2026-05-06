# 💡 Comentários Didáticos para Desenvolvedores

Este arquivo contém explicações detalhadas das partes mais complexas do código ClarityFlash, para ajudar novos desenvolvedores a entenderem rapidamente.

## 🎯 Middleware de Autenticação (`middleware/jwt.go`)

### Por que é complexo?
O middleware precisa aceitar **múltiplas formas de autenticação** porque:
- Web apps usam JWT tokens
- Extensões Chrome não conseguem armazenar tokens facilmente
- Uploads precisam funcionar mesmo sem login completo

### Fluxo de Tentativas:
```go
// 🥇 JWT Bearer Token (web app normal)
// Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

// 🥈 Header customizado (fallback para extensão)
// X-User-ID: user123

// 🥉 Query parameter (fallback adicional)
// ?user_id=user123

// 🔍 Campo no formulário (para uploads)
// FormData: { audio: file, user_id: "user123" }
```

### Dica para Debug:
Se autenticação falhar, verifique os logs nesta ordem:
1. Token JWT presente e válido?
2. Header X-User-ID existe?
3. Query param user_id existe?
4. Formulário multipart tem user_id?

## 🎼 Pipeline de IA (`service/processor.go`)

### O "Maestro" que coordena tudo:
```go
Process()     // ← Ponto de entrada (síncrono)
├── runPipeline() // ← Executado em goroutine (assíncrono)
    ├── ValidateAudio()
    ├── ConvertToWAV()
    ├── Transcribe()      // ← IA fala→texto
    ├── Generate Summary  // ← IA texto→resumo
    ├── Generate Flashcards // ← IA texto→cartões
    └── Cleanup()
```

### Por que assíncrono?
- Upload retorna **IMEDIATAMENTE** (usuário não espera)
- Processamento pode levar **1-2 minutos**
- Evita timeout de HTTP (geralmente 30s)
- Usuário vê progresso em tempo real

### Possíveis problemas:
- **Goroutine leak**: Se panic acontecer, ninguém limpa
- **Arquivo temporário**: Sempre usar `defer cleanup()`
- **Estado da sessão**: Atualizar status corretamente

## 📊 ExtractJSON (`processor.go`)

### Problema que resolve:
IAs de linguagem nem sempre retornam JSON "limpo":
```
Resposta da IA: "Aqui está seu JSON: {"title": "Aula"} Espero que ajude!"
```

### Como funciona:
```go
extractJSON("texto {json aqui} mais texto")
// Retorna: {json aqui}
```

### Algoritmo (parsing manual):
- Percorre string caractere por caractere
- Conta `{` e `}` para encontrar JSON válido
- Retorna primeiro bloco JSON completo encontrado

## 🔄 Fiber vs HTTP Padrão

### Por que migramos para Fiber?
- **Performance**: 2-3x mais rápido que net/http
- **API moderna**: Inspirado em Express.js
- **Middleware**: Mais fácil de usar
- **Context**: Melhor gerenciamento de requisição

### Diferenças principais:
```go
// HTTP padrão
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

// Fiber
func handler(c fiber.Ctx) error {
    return c.JSON(data)  // Tudo automático!
}
```

## 🗄️ Repositório Pattern

### Por que usar repositórios?
- **Separação**: Lógica negócio ≠ acesso dados
- **Testável**: Mock fácil para testes
- **Reutilizável**: Mesmo repo para web, CLI, API

### Estrutura típica:
```go
type SessionRepository interface {
    Create(ctx context.Context, session *Session) error
    GetByID(ctx context.Context, id string) (*Session, error)
    UpdateStatus(ctx context.Context, id, status string) error
}

type postgresSessionRepo struct {
    db *sql.DB
}

func (r *postgresSessionRepo) Create(ctx context.Context, session *Session) error {
    // SQL aqui
}
```

## 🎤 Processamento de Áudio

### Pipeline completo:
```mermaid
graph TD
    A[Upload MP3/WAV] --> B[Validar arquivo]
    B --> C[Converter para WAV]
    C --> D[Enviar para Whisper API]
    D --> E[Receber transcrição]
    E --> F[Gerar prompt resumo]
    F --> G[Enviar para Llama API]
    G --> H[Receber resumo JSON]
    H --> I[Gerar prompt flashcards]
    I --> J[Enviar para Llama API]
    J --> K[Receber flashcards JSON]
    K --> L[Salvar no banco]
    L --> M[Marcar como completo]
```

### Pontos críticos:
- **Formato WAV**: Whisper funciona melhor com WAV
- **Limpeza**: Deletar arquivos temporários sempre
- **Privacidade**: Remover áudio original após processamento

## 🔐 Autenticação JWT

### Como funciona:
```go
// 1. Login: usuário envia email/senha
user := authenticate(email, password)

// 2. Gerar token com informações do usuário
token := jwt.Sign(user.ID, user.Email)

// 3. Cliente guarda token
// Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

// 4. Servidor valida token em cada requisição
claims := jwt.Verify(token)
userID := claims.UserID
```

### Segurança importante:
- **Segredo forte**: Nunca commit segredo no código
- **Expiração**: Tokens devem expirar (normalmente 24h)
- **Refresh tokens**: Para renovação automática

## 🌐 CORS (Cross-Origin Resource Sharing)

### Por que precisamos?
Frontend roda em `localhost:5173`
Backend roda em `localhost:8081`
Navegador bloqueia por segurança!

### Solução Fiber:
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:5173",
    AllowMethods: "GET,POST,PUT,DELETE",
    AllowHeaders: "Content-Type,Authorization",
}))
```

## 📋 Checklist para Novos Recursos

Antes de implementar algo novo:

### 1. Interface (se aplicável)
- [ ] Componente Vue criado?
- [ ] Store Pinia atualizado?
- [ ] Rota adicionada?

### 2. Backend
- [ ] Handler criado?
- [ ] Service implementado?
- [ ] Repository atualizado?
- [ ] Rota registrada no Fiber?

### 3. Banco de dados
- [ ] Migration criada?
- [ ] Model atualizado?
- [ ] Repository testado?

### 4. Testes
- [ ] Unit tests escritos?
- [ ] Integration tests passando?
- [ ] E2E testado manualmente?

### 5. Documentação
- [ ] README atualizado?
- [ ] Comentários no código?
- [ ] API documentada?

## 🐛 Debug Comum

### Problema: "erro ao abrir arquivo"
```go
// Verificar se arquivo foi enviado
file, err := c.FormFile("audio")
if err != nil {
    return c.Status(400).JSON(fiber.Map{"error": "arquivo não enviado"})
}
```

### Problema: Autenticação falhando
```go
// Verificar se middleware está rodando
// Verificar se user_id está no contexto
userID := c.Locals("user_id")
log.Printf("UserID: %v", userID)
```

### Problema: JSON da IA inválido
```go
// Log da resposta bruta antes do parse
log.Printf("Resposta IA: %s", rawResponse)
summaryJSON := extractJSON(rawResponse)
```

## 🎯 Dicas Gerais

1. **Leia os comentários**: Adicionei explicações detalhadas no código
2. **Use o guia didático**: `docs/guia-didatico.md` explica o projeto
3. **Teste pequeno**: Faça mudanças pequenas e teste frequentemente
4. **Logs ajudam**: Adicione `log.Printf()` temporários para debug
5. **Pergunte**: Código complexo? Pergunte nos comentários!

---

**Lembre-se**: Todo código era complexo até alguém explicar! 🚀</content>
<parameter name="filePath">docs/comentarios-desenvolvedor.md