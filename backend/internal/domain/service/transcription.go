package service

// TranscriptionService contem a logica de transcricao de audio
// Encapsula chamadas ao STT e validacoes de qualidade
type TranscriptionService struct{}

// NewTranscriptionService cria uma nova instancia
func NewTranscriptionService() *TranscriptionService {
	return &TranscriptionService{}
}
