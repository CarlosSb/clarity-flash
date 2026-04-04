package stt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// GroqClient client para a API Groq Whisper
type GroqClient struct {
	APIKey  string
	Model   string
	BaseURL string
}

func NewGroqClient(apiKey, model string) *GroqClient {
	return &GroqClient{
		APIKey:  apiKey,
		Model:   model,
		BaseURL: "https://api.groq.com/openai/v1",
	}
}

// TranscribeResponse estrutura da resposta da Groq
type TranscribeResponse struct {
	Text string `json:"text"`
}

// Transcribe envia um audio para transcricao via Groq Whisper
func (c *GroqClient) Transcribe(audioPath string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Anexa arquivo de audio
	fw, err := writer.CreateFormFile("file", audioPath)
	if err != nil {
		return "", fmt.Errorf("criar form file: %w", err)
	}

	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("abrir audio: %w", err)
	}
	defer file.Close()

	if _, err = io.Copy(fw, file); err != nil {
		return "", fmt.Errorf("copiar audio: %w", err)
	}

	// Parametros da transcricao
	_ = writer.WriteField("model", c.Model)
	_ = writer.WriteField("language", "pt")
	_ = writer.WriteField("response_format", "json")
	writer.Close()

	req, err := http.NewRequest("POST", c.BaseURL+"/audio/transcriptions", body)
	if err != nil {
		return "", fmt.Errorf("criar request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request groq: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("groq error %d: %s", resp.StatusCode, string(respBody))
	}

	var result TranscribeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if result.Text == "" {
		return "", fmt.Errorf("transcricao vazia")
	}

	return result.Text, nil
}
