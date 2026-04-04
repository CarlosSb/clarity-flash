<template>
  <div class="mx-auto max-w-3xl">
    <!-- Progress Bar -->
    <div class="mb-6">
      <div class="mb-2 flex items-center justify-between text-sm text-text-secondary">
        <span>Card {{ currentIndex + 1 }} de {{ props.cards.length }}</span>
        <span>{{ progressPercent }}%</span>
      </div>
      <div class="h-2 overflow-hidden rounded-full bg-card-border">
        <div
          class="h-full rounded-full bg-primary transition-all duration-500 ease-out"
          :style="{ width: `${progressPercent}%` }"
        ></div>
      </div>
    </div>

    <!-- Results Screen -->
    <div v-if="quizComplete" class="card overflow-hidden">
      <div class="border-b border-card-border bg-primary-muted p-6 text-center">
        <h2 class="text-2xl font-bold text-primary-light">Quiz Finalizado!</h2>
        <p class="mt-1 text-sm text-text-secondary">Confira seus resultados</p>
      </div>
      <div class="p-6">
        <div class="grid grid-cols-3 gap-4 text-center">
          <div>
            <p class="text-3xl font-bold text-text">{{ props.cards.length }}</p>
            <p class="text-xs text-text-secondary">Total</p>
          </div>
          <div>
            <p class="text-3xl font-bold text-success">{{ stats.correct }}</p>
            <p class="text-xs text-text-secondary">Acertos</p>
          </div>
          <div>
            <p class="text-3xl font-bold text-danger">{{ stats.incorrect }}</p>
            <p class="text-xs text-text-secondary">Erros</p>
          </div>
        </div>
        <p class="mt-4 text-center text-sm text-text-secondary">
          {{ stats.correct >= props.cards.length * 0.7 ? 'Otimo desempenho!' : 'Continue praticando!' }}
          - {{ Math.round((stats.correct / props.cards.length) * 100) }}%
        </p>
        <div class="mt-6 flex justify-center gap-3">
          <button class="btn-secondary" @click="restart">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Tentar Novamente
          </button>
        </div>
      </div>
    </div>

    <!-- Quiz Card with Flip Animation -->
    <div v-else-if="currentCard" class="flip-card-perspective mb-6 h-72 sm:h-80 cursor-pointer" @click="handleCardClick">
      <div class="flip-card-inner" :class="{ 'is-flipped': isRevealed }">
        <!-- Front - Question -->
        <div class="flip-card-face rounded-xl border border-card-border bg-card p-6">
          <div class="flex h-full flex-col">
            <div class="mb-4 flex items-center gap-2">
              <div class="flex h-6 w-6 items-center justify-center rounded-full bg-primary-muted">
                <svg class="h-3.5 w-3.5 text-primary-light" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <span class="text-xs font-medium uppercase tracking-wider text-text-secondary">Pergunta</span>
            </div>
            <div class="flex flex-1 items-center">
              <p class="text-lg font-medium leading-relaxed text-text">{{ currentCard.front }}</p>
            </div>
            <div class="mt-4 flex items-center justify-between">
              <DifficultyBadge :difficulty="currentCard.difficulty" />
              <span class="text-xs text-text-secondary">Clique para ver a resposta</span>
            </div>
          </div>
        </div>

        <!-- Back - Answer -->
        <div class="flip-card-face flip-card-back rounded-xl border border-primary-muted/30 bg-card p-6">
          <div class="flex h-full flex-col">
            <div class="mb-4 flex items-center gap-2">
              <div class="flex h-6 w-6 items-center justify-center rounded-full bg-primary-muted">
                <svg class="h-3.5 w-3.5 text-primary-light" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                </svg>
              </div>
              <span class="text-xs font-medium uppercase tracking-wider text-text-secondary">Resposta</span>
            </div>
            <div class="flex flex-1 items-center">
              <p class="text-base leading-relaxed text-text">{{ currentCard.back }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Answer Buttons (shown after reveal) -->
    <div v-if="isRevealed && !quizComplete" class="flex items-center justify-center gap-3">
      <button class="btn-primary bg-danger hover:bg-danger/90" @click.stop="rateAndNext('incorrect')">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
        Errei
      </button>
      <button class="btn-primary bg-success hover:bg-success/90" @click.stop="rateAndNext('correct')">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        Acertei
      </button>
    </div>

    <!-- Reveal Button -->
    <div v-if="!isRevealed && !quizComplete" class="flex justify-center">
      <button class="btn-secondary w-full max-w-xs" @click="handleCardClick">
        Mostrar Resposta
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface Flashcard {
  id: string
  front: string
  back: string
  difficulty: 1 | 2 | 3
}

const props = defineProps<{
  cards: Flashcard[]
  session?: { title?: string }
}>()

const emit = defineEmits<{
  finish: [results: { correct: number; incorrect: number; total: number }]
  'reset-quiz': []
}>()

const currentIndex = ref(0)
const isRevealed = ref(false)
const quizComplete = ref(false)
const stats = ref({ correct: 0, incorrect: 0 })

const currentCard = computed(() => props.cards[currentIndex.value] ?? null)
const progressPercent = computed(() => {
  if (props.cards.length === 0) return 0
  return Math.round(((stats.value.correct + stats.value.incorrect) / props.cards.length) * 100)
})

function handleCardClick(): void {
  if (!isRevealed.value) {
    isRevealed.value = true
  }
}

function rateAndNext(rating: 'correct' | 'incorrect'): void {
  if (rating === 'correct') {
    stats.value.correct++
  } else {
    stats.value.incorrect++
  }

  isRevealed.value = false

  if (currentIndex.value < props.cards.length - 1) {
    currentIndex.value++
  } else {
    quizComplete.value = true
    emit('finish', {
      correct: stats.value.correct,
      incorrect: stats.value.incorrect,
      total: props.cards.length,
    })
  }
}

function restart(): void {
  currentIndex.value = 0
  isRevealed.value = false
  quizComplete.value = false
  stats.value = { correct: 0, incorrect: 0 }
  emit('reset-quiz')
}

// Difficulty Badge sub-component
const DifficultyBadge = {
  props: {
    difficulty: { type: Number as () => number, required: true },
  },
  template: `
    <span
      class="badge"
      :class="{
        'badge-success': difficulty === 1,
        'badge-warning': difficulty === 2,
        'badge-danger': difficulty === 3,
      }"
    >
      {{ difficulty === 1 ? 'Facil' : difficulty === 2 ? 'Medio' : 'Dificil' }}
    </span>
  `,
}
</script>
