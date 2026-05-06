package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "JSON puro",
			input:    `{"nome": "João"}`,
			expected: `{"nome": "João"}`,
		},
		{
			name:     "JSON com texto antes",
			input:    `Aqui está: {"nome": "João"} e depois`,
			expected: `{"nome": "João"}`,
		},
		{
			name:     "JSON com texto depois",
			input:    `{"nome": "João"} Espero que ajude!`,
			expected: `{"nome": "João"}`,
		},
		{
			name:     "JSON aninhado",
			input:    `Resposta: {"user": {"name": "João", "age": 30}} ok`,
			expected: `{"user": {"name": "João", "age": 30}}`,
		},
		{
			name:     "Sem JSON",
			input:    "Texto sem JSON",
			expected: "Texto sem JSON",
		},
		{
			name:     "JSON incompleto",
			input:    `{"nome": "João"`,
			expected: `{"nome": "João"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractJSON(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateID(t *testing.T) {
	id1, err1 := generateID()
	id2, err2 := generateID()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)  // IDs devem ser únicos
	assert.Equal(t, 32, len(id1)) // 16 bytes * 2 chars hex = 32
	assert.Equal(t, 32, len(id2))
}
