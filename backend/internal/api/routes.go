package api

// Routes documenta todas as rotas disponiveis
var Routes = []Route{
	{Method: "GET", Path: "/health", Desc: "Health check"},
	{Method: "POST", Path: "/api/sessions/upload", Desc: "Upload de audio para processamento"},
	{Method: "GET", Path: "/api/sessions/{id}", Desc: "Detalhes de uma sessao"},
	{Method: "GET", Path: "/api/sessions", Desc: "Lista sessoes do usuario (query: user_id)"},
	{Method: "DELETE", Path: "/api/sessions/{id}", Desc: "Deleta uma sessao"},
	{Method: "GET", Path: "/api/export/{id}/csv", Desc: "Exporta flashcards em CSV (Anki)"},
	{Method: "GET", Path: "/api/export/{id}/txt", Desc: "Exporta flashcards em texto"},
}

type Route struct {
	Method string
	Path   string
	Desc   string
}
