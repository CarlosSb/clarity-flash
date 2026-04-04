package model

import "fmt"

// Summary representa o resumo gerado por IA de uma sessao
type Summary struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Highlights    []string `json:"highlights"`
	Decisions     []string `json:"decisions"`
	ActionItems   []string `json:"action_items"`
	KeyConcepts   []string `json:"key_concepts"`
}

// SummaryPrompt gera o prompt de sistema para geracao de resumo em PT-BR
func SummaryPrompt(transcript string) string {
	return fmt.Sprintf(`Voce e um assistente especializado em analisar transcricoes de aulas e reunioes em portugues brasileiro.
Analise a transcricao abaixo e gere um resumo profissional seguindo EXATAMENTE o formato JSON.

Regras:
- Seja conciso e claro
- Destaque apenas o que e relevante
- Action items devem ser práticos e com verbo no infinitivo
- Se nao houver decisoes ou action items, use arrays vazios

Transcricao:
%s

Retorne APENAS o JSON valido, sem markdown ou explicacoes extras.
Formato esperado:
{"title":"titulo curto","description":"resumo em 2-3 frases","highlights":["destaque1"],"decisions":["decisao1"],"action_items":["acao1"],"key_concepts":["conceito1"]}`, transcript)
}
