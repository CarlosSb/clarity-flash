package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Processor lida com conversao e preprocessamento de audio
type Processor struct {
	TempDir string
}

func NewProcessor(tempDir string) (*Processor, error) {
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("criar temp dir: %w", err)
	}
	return &Processor{TempDir: tempDir}, nil
}

// ValidateAudio verifica se o arquivo e valido
func (p *Processor) ValidateAudio(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("arquivo nao encontrado: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("arquivo vazio")
	}
	// Limite de 200MB
	if info.Size() > 200*1024*1024 {
		return fmt.Errorf("arquivo muito grande (max 200MB)")
	}
	return nil
}

// ConvertToWAV converte audio para WAV 16kHz mono (ideal para Whisper)
func (p *Processor) ConvertToWAV(inputPath string) (string, error) {
	baseName := filepath.Base(inputPath)
	outputPath := filepath.Join(p.TempDir, "converted_"+baseName+".wav")

	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath,
		"-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le",
		outputPath)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %w", err)
	}

	return outputPath, nil
}

// Cleanup remove arquivos temporarios apos processamento
func (p *Processor) Cleanup(paths ...string) {
	for _, path := range paths {
		_ = os.Remove(path)
	}
}
