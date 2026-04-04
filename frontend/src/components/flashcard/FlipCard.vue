<template>
  <div class="flip-card-perspective h-full w-full cursor-pointer" @click="handleClick">
    <div class="flip-card-inner" :class="{ 'is-flipped': isFlipped }">
      <!-- Front -->
      <div class="flip-card-face rounded-xl border border-card-border bg-card p-6">
        <div class="flex h-full flex-col">
          <div class="mb-4 flex items-center gap-2">
            <div class="flex h-5 w-5 items-center justify-center rounded-full bg-primary-muted">
              <span class="text-xs text-primary-light">Q</span>
            </div>
            <span class="text-xs font-medium uppercase tracking-wider text-text-secondary">
              Pergunta
            </span>
          </div>
          <div class="flex-1 flex-col justify-center">
            <p class="text-lg font-medium leading-relaxed text-text">{{ card.front }}</p>
          </div>
          <div class="mt-4 flex items-center justify-between">
            <DifficultyBadge :difficulty="card.difficulty" />
            <span class="text-xs text-text-secondary">Clique para revelar</span>
          </div>
        </div>
      </div>

      <!-- Back -->
      <div class="flip-card-face flip-card-back rounded-xl border border-primary-muted/30 bg-card p-6">
        <div class="flex h-full flex-col">
          <div class="mb-4 flex items-center gap-2">
            <div class="flex h-5 w-5 items-center justify-center rounded-full bg-primary-muted">
              <span class="text-xs text-primary-light">A</span>
            </div>
            <span class="text-xs font-medium uppercase tracking-wider text-text-secondary">
              Resposta
            </span>
          </div>
          <div class="flex-1">
            <p class="text-base leading-relaxed text-text">{{ card.back }}</p>
          </div>
          <div class="mt-4">
            <DifficultyBadge :difficulty="card.difficulty" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface Card {
  id: string
  front: string
  back: string
  difficulty: 1 | 2 | 3
}

defineProps<{
  card: Card
}>()

const emit = defineEmits<{
  flip: []
}>()

const flipped = ref(false)
const isFlipped = computed(() => flipped.value)

function handleClick(): void {
  flipped.value = !flipped.value
  emit('flip')
}
</script>
