package model

import "fmt"

// Flashcard - O cartão de estudo individual
// Cada flashcard é como um mini-teste: pergunta na frente, resposta atrás
type Flashcard struct {
	Front      string `json:"front"`      // Pergunta (ex: "O que é polimorfismo?")
	Back       string `json:"back"`       // Resposta (ex: "Capacidade de um objeto assumir várias formas")
	Difficulty int    `json:"difficulty"` // 1=fácil, 2=médio, 3=difícil (para algoritmos de repetição)
}

// FlashcardDeck - Um baralho completo de flashcards
// Agrupa vários cartões relacionados a uma mesma aula/sessão
type FlashcardDeck struct {
	SessionID   string      `json:"session_id"`   // ID da sessão que gerou estes cartões
	Title       string      `json:"title"`        // Título do deck (opcional)
	Description string      `json:"description"`  // Descrição do conteúdo (opcional)
	Cards       []Flashcard `json:"cards"`        // Lista de 10-15 cartões gerados
}

// FlashcardPrompt - Cria instruções para a IA gerar flashcards de qualidade
// Esta função "ensina" a IA como criar boas flashcards de estudo
func FlashcardPrompt(transcript string) string {
	return fmt.Sprintf(`Voce e um especialista em criacao de flashcards para estudo.
Com base na transcricao abaixo, crie exatamente 10-15 flashcards em portugues brasileiro.

Regras:
- Cada flashcard deve ter uma pergunta clara na frente e resposta objetiva no verso
- Cubra os conceitos mais importantes do conteudo
- Evite questions triviais ou de "sim/nao"
- Varie a dificuldade entre facil, medio e dificil

Transcricao:
%s

Retorne APENAS um JSON valido no formato abaixo, sem markdown adicional:
{"cards":[{"front":"pergunta","back":"resposta","difficulty":1},...]}`, transcript)
}
