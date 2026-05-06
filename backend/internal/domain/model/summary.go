package model

import "fmt"

// Summary - Estrutura do resumo inteligente gerado pela IA
// Esta struct representa o resultado final da análise de uma aula/reunião
type Summary struct {
	Title       string   `json:"title"`        // Título gerado automaticamente
	Description string   `json:"description"`  // Resumo executivo (2-3 frases)
	Highlights  []string `json:"highlights"`   // Pontos principais (bullet points)
	Decisions   []string `json:"decisions"`    // Decisões tomadas (para reuniões)
	ActionItems []string `json:"action_items"` // Tarefas a fazer (follow-ups)
	KeyConcepts []string `json:"key_concepts"` // Conceitos importantes aprendidos
}

// SummaryPrompt - Cria o "comando" enviado para a IA gerar resumos
// Esta função cria um prompt muito específico que "ensina" a IA como
// analisar transcrições em português brasileiro e gerar resumos estruturados
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
