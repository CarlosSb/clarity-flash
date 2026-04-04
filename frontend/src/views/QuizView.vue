<template>
  <div class="min-h-screen bg-bg">
    <div class="px-4 py-6 sm:px-6 lg:px-8">
      <!-- Back Button -->
      <div class="mx-auto mb-6 max-w-3xl">
        <RouterLink
          :to="`/session/${sessionId}`"
          class="inline-flex items-center gap-1.5 text-sm text-text-secondary transition-colors hover:text-primary"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Voltar para sessao
        </RouterLink>
      </div>

      <!-- Loading State -->
      <div v-if="store.loading" class="mx-auto max-w-3xl">
        <div class="card animate-pulse py-12 text-center">
          <div class="mx-auto mb-4 h-8 w-32 rounded bg-card-border"></div>
          <div class="mx-auto h-4 w-48 rounded bg-card-border"></div>
        </div>
      </div>

      <!-- No Flashcards -->
      <div v-else-if="flashcards.length === 0" class="mx-auto max-w-3xl">
        <div class="card py-12 text-center">
          <svg class="mx-auto mb-4 h-12 w-12 text-text-dim" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          <h2 class="text-lg font-semibold text-text">Sem flashcards ainda</h2>
          <p class="mt-2 text-text-secondary">
            Esta sessao nao tem flashcards para o quiz.
          </p>
          <RouterLink :to="`/session/${sessionId}`" class="btn-primary mt-4">
            Voltar
          </RouterLink>
        </div>
      </div>

      <!-- Results -->
      <div v-else-if="showResults" class="mx-auto max-w-3xl">
        <div class="card overflow-hidden">
          <div class="border-b border-card-border bg-primary-muted p-6 text-center">
            <h2 class="text-2xl font-bold text-primary-light">Quiz Finalizado!</h2>
            <p class="mt-1 text-sm text-text-secondary">Confira seus resultados</p>
          </div>
          <div class="p-6">
            <div class="grid grid-cols-3 gap-4 text-center">
              <div>
                <p class="text-3xl font-bold text-text">{{ quizStats.total }}</p>
                <p class="text-xs text-text-secondary">Total</p>
              </div>
              <div>
                <p class="text-3xl font-bold text-success">{{ quizStats.correct }}</p>
                <p class="text-xs text-text-secondary">Acertos</p>
              </div>
              <div>
                <p class="text-3xl font-bold text-danger">{{ quizStats.incorrect }}</p>
                <p class="text-xs text-text-secondary">Erros</p>
              </div>
            </div>
            <div class="mt-6 flex justify-center gap-3">
              <button class="btn-secondary" @click="resetQuiz()">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Tentar Novamente
              </button>
              <RouterLink :to="`/session/${sessionId}`" class="btn-primary">
                Voltar para Sessao
              </RouterLink>
            </div>
          </div>
        </div>

        <!-- Answer Review -->
        <div class="mt-6 space-y-3">
          <div
            v-for="(card, index) in flashcards"
            :key="card.id"
            class="card border-2"
            :class="{
              'border-success/30 bg-success/5': quizAnswers[card.id] === 'correct',
              'border-danger/30 bg-danger/5': quizAnswers[card.id] === 'incorrect',
              'border-card-border': !quizAnswers[card.id],
            }"
          >
            <div class="flex items-center gap-3">
              <span class="flex h-6 w-6 items-center justify-center rounded-full bg-card-border text-xs font-medium">
                {{ index + 1 }}
              </span>
              <div class="flex-1">
                <p class="text-sm font-medium text-text">{{ card.front }}</p>
                <p class="mt-1 text-xs text-text-secondary">{{ card.back }}</p>
              </div>
              <span
                class="badge"
                :class="{
                  'badge-success': quizAnswers[card.id] === 'correct',
                  'badge-danger': quizAnswers[card.id] === 'incorrect',
                }"
              >
                {{ quizAnswers[card.id] === 'correct' ? 'Acertou' : quizAnswers[card.id] === 'incorrect' ? 'Errou' : 'Nao respondido' }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Quiz Session -->
      <QuizSession v-else :cards="flashcards" @reset-quiz="resetQuiz" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, watch, computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useSessionStore } from '@/store'
import { useFlashcards } from '@/composables/useFlashcards'
import QuizSession from '@/components/quiz/QuizSession.vue'

const route = useRoute()
const store = useSessionStore()

const props = defineProps<{
  id: string
}>()

const sessionId = computed(() => props.id ?? (route.params.id as string))

const {
  flashcards,
  showResults,
  quizAnswers,
  quizStats,
  resetQuiz: resetQuizFn,
} = useFlashcards(sessionId.value)

function resetQuiz(): void {
  resetQuizFn()
}

async function loadSession(): Promise<void> {
  await store.fetchSession(sessionId.value)
}

onMounted(loadSession)
watch(() => route.params.id, loadSession)
</script>
