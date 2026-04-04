<template>
  <div class="card">
    <!-- Header -->
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-lg font-semibold">Flashcards</h2>
      <span class="badge badge-primary">{{ sessionFlashcards.length }} cards</span>
    </div>

    <!-- Empty State -->
    <div v-if="sessionFlashcards.length === 0" class="py-8 text-center">
      <svg class="mx-auto mb-3 h-10 w-10 text-text-dim" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
      <p class="text-sm text-text-secondary">Nenhum flashcard gerado ainda</p>
    </div>

    <!-- Deck List -->
    <div v-else>
      <!-- Selected Card Preview -->
      <div class="mb-4 overflow-hidden rounded-lg">
        <div class="flip-card-perspective h-64">
          <FlipCard
            v-if="selectedCard"
            :card="selectedCard"
            @flip="onFlip"
          />
        </div>
      </div>

      <!-- Navigation -->
      <div class="mb-4 flex items-center justify-between text-sm text-text-secondary">
        <button
          class="rounded-lg px-3 py-1.5 transition-colors hover:bg-card-border hover:text-text disabled:opacity-30"
          :disabled="currentIndex === 0"
          @click="prev()"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <span>{{ currentIndex + 1 }} / {{ sessionFlashcards.length }}</span>
        <button
          class="rounded-lg px-3 py-1.5 transition-colors hover:bg-card-border hover:text-text disabled:opacity-30"
          :disabled="currentIndex === sessionFlashcards.length - 1"
          @click="next()"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>

      <!-- Card List -->
      <div class="space-y-2">
        <button
          v-for="(card, index) in sessionFlashcards"
          :key="card.id"
          class="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-left text-sm transition-colors hover:bg-card-hover"
          :class="index === currentIndex ? 'bg-card-border/50' : ''"
          @click="currentIndex = index"
        >
          <span
            class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-medium"
            :class="index === currentIndex ? 'bg-primary-muted text-primary-light' : 'bg-card-border text-text-secondary'"
          >
            {{ index + 1 }}
          </span>
          <span class="line-clamp-1 flex-1 text-text-secondary">{{ card.front }}</span>
          <span
            class="badge"
            :class="{
              'badge-success': card.difficulty === 1,
              'badge-warning': card.difficulty === 2,
              'badge-danger': card.difficulty === 3,
            }"
          >
            {{ difficultyLabel(card.difficulty) }}
          </span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useSessionStore } from '@/store'
import FlipCard from './FlipCard.vue'

const store = useSessionStore()
const currentIndex = ref(0)

function difficultyLabel(d: number): string {
  switch (d) {
    case 1: return 'Facil'
    case 2: return 'Medio'
    case 3: return 'Dificil'
    default: return 'Medio'
  }
}

const sessionFlashcards = computed(() => store.sessionFlashcards)
const selectedCard = computed(() => sessionFlashcards.value[currentIndex.value] ?? null)

function next(): void {
  if (currentIndex.value < sessionFlashcards.value.length - 1) {
    currentIndex.value++
  }
}

function prev(): void {
  if (currentIndex.value > 0) {
    currentIndex.value--
  }
}

function onFlip(): void {
  // Could log or track flip events
}
</script>
