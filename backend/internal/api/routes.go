package api

// Routes documenta todas as rotas disponiveis no servidor Fiber
var Routes = []Route{
	{Method: "GET", Path: "/health", Desc: "Health check"},
	{Method: "POST", Path: "/api/auth/register", Desc: "Registro de usuario"},
	{Method: "POST", Path: "/api/auth/login", Desc: "Login de usuario"},
	{Method: "POST", Path: "/api/sessions/upload", Desc: "Upload de audio para processamento", Auth: true},
	{Method: "GET", Path: "/api/sessions/:id", Desc: "Detalhes de uma sessao", Auth: true},
	{Method: "GET", Path: "/api/sessions", Desc: "Lista sessoes do usuario", Auth: true},
	{Method: "DELETE", Path: "/api/sessions/:id", Desc: "Deleta uma sessao", Auth: true},
	{Method: "GET", Path: "/api/export/:id/csv", Desc: "Exporta flashcards em CSV (Anki)", Auth: true},
	{Method: "GET", Path: "/api/export/:id/txt", Desc: "Exporta flashcards em texto", Auth: true},
}

type Route struct {
	Method string
	Path   string
	Desc   string
	Auth   bool // requires JWT auth
}

// This file serves as documentation for the API routes.
// The actual routing is handled by SetupFiberApp in fiber_app.go
