import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

// --- Type definitions ---
export interface Flashcard {
  id: string
  session_id: string
  front: string
  back: string
  difficulty: 1 | 2 | 3  // 1=facil, 2=medio, 3=dificil
  mastered: boolean
}

export interface SummaryData {
  title: string
  description: string
  highlights: string[]
  decisions: string[]
  action_items: string[]
  key_concepts: string[]
}

export interface Session {
  id: string
  user_id: string
  title: string
  description?: string
  created_at: string
  duration: number
  status: 'processing' | 'completed' | 'failed'
  mode: 'student' | 'professional'
  flashcards: Flashcard[]
  summary?: SummaryData
}

export interface QuizState {
  currentIndex: number
  answers: Record<string, 'correct' | 'incorrect'>
  showResults: boolean
}

// --- Store ---
export const useSessionStore = defineStore('session', () => {
  const sessions = ref<Session[]>([])
  const activeSession = ref<Session | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const totalFlashcards = computed(() =>
    sessions.value.reduce((sum, s) => sum + (s.flashcards?.length ?? 0), 0)
  )

  const totalRecordings = computed(() => sessions.value.length)

  const readySessions = computed(() =>
    sessions.value.filter((s) => s.status === 'completed')
  )

  const sessionFlashcards = computed(() => activeSession.value?.flashcards ?? [])

  // Actions
  async function fetchSessions(): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/sessions')
      sessions.value = response.data
    } catch (e: any) {
      error.value = 'Erro ao carregar gravacoes'
      console.error('fetchSessions:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchSession(id: string): Promise<Session | null> {
    loading.value = true
    error.value = null
    try {
      const response = await api.get<Session>(`/sessions/${id}`)
      activeSession.value = response.data
      return response.data
    } catch (e: any) {
      error.value = e.response?.data ?? 'Erro ao carregar sessao'
      console.error('fetchSession:', e)
      activeSession.value = null
      return null
    } finally {
      loading.value = false
    }
  }

  async function deleteSession(id: string): Promise<void> {
    try {
      await api.delete(`/sessions/${id}`)
      sessions.value = sessions.value.filter((s) => s.id !== id)
      if (activeSession.value?.id === id) {
        activeSession.value = null
      }
    } catch (e: any) {
      error.value = 'Erro ao deletar sessao'
      console.error('deleteSession:', e)
    }
  }

  function setActiveSession(session: Session | null): void {
    activeSession.value = session
  }

  function resetError(): void {
    error.value = null
  }

  return {
    sessions,
    activeSession,
    loading,
    error,
    totalFlashcards,
    totalRecordings,
    readySessions,
    sessionFlashcards,
    fetchSessions,
    fetchSession,
    deleteSession,
    setActiveSession,
    resetError,
  }
})
