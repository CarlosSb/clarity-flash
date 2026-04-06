package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// LLMClient interface comum para provedores de LLM
type LLMClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// GroqLLMClient client para Groq LLM API (Llama, etc)
type GroqLLMClient struct {
	APIKey string
	Model  string
}

func NewGroqLLMClient(apiKey, model string) *GroqLLMClient {
	return &GroqLLMClient{
		APIKey: apiKey,
		Model:  model,
	}
}

// Generate envia prompt ao modelo Groq LLM
func (c *GroqLLMClient) Generate(ctx context.Context, prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
		"max_tokens":  2048,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	url := "https://api.groq.com/openai/v1/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return "", fmt.Errorf("criar request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request groq llm: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq llm error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("resposta vazia do Groq LLM")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// HuggingFaceClient client para Hugging Face Inference API
type HuggingFaceClient struct {
	Token string
	Model string
}

func NewHuggingFaceClient(token, model string) *HuggingFaceClient {
	return &HuggingFaceClient{
		Token: token,
		Model: model,
	}
}

// Generate envia prompt ao modelo via HF Inference API
func (c *HuggingFaceClient) Generate(ctx context.Context, prompt string) (string, error) {
	payload := map[string]interface{}{
		"inputs": prompt,
		"parameters": map[string]interface{}{
			"max_new_tokens":   2048,
			"temperature":      0.3,
			"return_full_text": false,
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", c.Model)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return "", fmt.Errorf("criar request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request hf: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("hf error %d: %s", resp.StatusCode, string(respBody))
	}

	var results []struct {
		GeneratedText string `json:"generated_text"`
	}
	if err := json.Unmarshal(respBody, &results); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(results) == 0 {
		return "", fmt.Errorf("resposta vazia do LLM")
	}

	return strings.TrimSpace(results[0].GeneratedText), nil
}

// OllamaClient client para Ollama local
type OllamaClient struct {
	URL   string
	Model string
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		URL:   baseURL,
		Model: model,
	}
}

// Generate envia prompt ao Ollama local
func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	payload := map[string]interface{}{
		"model":  c.Model,
		"prompt": prompt,
		"stream": false,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", c.URL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return "", fmt.Errorf("criar request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request ollama: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return strings.TrimSpace(result.Response), nil
}
