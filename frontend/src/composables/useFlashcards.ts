import { ref, computed, readonly } from 'vue'
import { useSessionStore } from '@/store'

export function useFlashcards(_sessionId?: string) {
  const store = useSessionStore()
  const flippedCards = ref<Set<string>>(new Set())
  const currentCardIndex = ref(0)

  // Quiz state
  const quizAnswers = ref<Record<string, 'correct' | 'incorrect'>>({})
  const showResults = ref(false)

  const flashcards = computed(() => store.sessionFlashcards)
  const currentCard = computed(() => {
    const cards = flashcards.value
    if (cards.length === 0) return null
    return cards[currentCardIndex.value % cards.length]
  })

  const progress = computed(() => {
    const total = flashcards.value.length
    if (total === 0) return 0
    return Math.round((Object.keys(quizAnswers.value).length / total) * 100)
  })

  const quizStats = computed(() => {
    const values = Object.values(quizAnswers.value)
    return {
      total: flashcards.value.length,
      answered: values.length,
      correct: values.filter((v) => v === 'correct').length,
      incorrect: values.filter((v) => v === 'incorrect').length,
    }
  })

  function isFlipped(cardId: string): boolean {
    return flippedCards.value.has(cardId)
  }

  function toggleFlip(cardId: string): void {
    const newSet = new Set(flippedCards.value)
    if (newSet.has(cardId)) {
      newSet.delete(cardId)
    } else {
      newSet.add(cardId)
    }
    flippedCards.value = newSet
  }

  function nextCard(): void {
    const cards = flashcards.value
    if (cards.length === 0) return
    currentCardIndex.value = (currentCardIndex.value + 1) % cards.length
  }

  function prevCard(): void {
    const cards = flashcards.value
    if (cards.length === 0) return
    currentCardIndex.value =
      (currentCardIndex.value - 1 + cards.length) % cards.length
  }

  function goToCard(index: number): void {
    const cards = flashcards.value
    if (index >= 0 && index < cards.length) {
      currentCardIndex.value = index
    }
  }

  function resetQuiz(): void {
    quizAnswers.value = {}
    showResults.value = false
    currentCardIndex.value = 0
    flippedCards.value = new Set()
  }

  function markAnswer(cardId: string, result: 'correct' | 'incorrect'): void {
    quizAnswers.value = { ...quizAnswers.value, [cardId]: result }

    if (Object.keys(quizAnswers.value).length >= flashcards.value.length) {
      showResults.value = true
    }
  }

  return {
    flashcards,
    currentCard,
    currentCardIndex: readonly(currentCardIndex),
    flippedCards: readonly(flippedCards),
    quizAnswers: readonly(quizAnswers),
    showResults: readonly(showResults),
    progress,
    quizStats,
    isFlipped,
    toggleFlip,
    nextCard,
    prevCard,
    goToCard,
    resetQuiz,
    markAnswer,
  }
}
