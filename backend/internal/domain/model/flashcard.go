package model

import "fmt"

// Flashcard representa um cartao de estudo com frente e verso
type Flashcard struct {
	Front    string `json:"front"`
	Back     string `json:"back"`
	Difficulty int  `json:"difficulty"` // 1=facil, 2=medio, 3=dficil
}

// FlashcardDeck representa um conjunto de flashcards
type FlashcardDeck struct {
	SessionID   string      `json:"session_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Cards       []Flashcard `json:"cards"`
}

// FlashcardPrompt gera o prompt para gerar flashcards em PT-BR
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
