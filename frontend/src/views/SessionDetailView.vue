<template>
  <div class="min-h-screen bg-bg">
    <div class="px-4 py-6 sm:px-6 lg:px-8">
      <!-- Back Navigation -->
      <div class="mx-auto mb-6 max-w-6xl">
        <RouterLink to="/" class="inline-flex items-center gap-1.5 text-sm text-text-secondary transition-colors hover:text-primary">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Voltar para gravacoes
        </RouterLink>
      </div>

      <!-- Loading State -->
      <div v-if="store.loading" class="mx-auto max-w-6xl">
        <div class="animate-pulse space-y-6">
          <div class="card h-48"></div>
          <div class="card h-96"></div>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="store.error" class="mx-auto max-w-6xl">
        <div class="card border-danger/30 bg-danger/5">
          <p class="text-sm font-medium text-danger">{{ store.error }}</p>
          <button class="mt-3 btn-secondary text-xs" @click="loadSession">
            Tentar Novamente
          </button>
        </div>
      </div>

      <!-- Session Detail -->
      <div v-else-if="session" class="mx-auto max-w-6xl">
        <!-- Title Header -->
        <div class="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <div class="flex items-center gap-2">
              <div
                class="h-2.5 w-2.5 rounded-full"
                :class="{
                  'bg-success': session.status === 'completed',
                  'bg-warning': session.status === 'processing',
                  'bg-danger': session.status === 'failed',
                }"
              ></div>
              <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">{{ session.title }}</h1>
            </div>
            <p class="mt-1 text-text-secondary">
              {{ session.duration }}min duracao &middot;
              {{ formatDate(session.created_at) }} &middot;
              {{ session.flashcards?.length || 0 }} flashcards
            </p>
          </div>
          <RouterLink
            :to="`/quiz/${session.id}`"
            class="btn-primary"
            :class="session.flashcards?.length === 0 ? 'opacity-50 cursor-not-allowed pointer-events-none' : ''"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
            </svg>
            Iniciar Quiz
          </RouterLink>
        </div>

        <div class="grid grid-cols-1 gap-6 lg:grid-cols-5">
          <!-- Left: Summary Panel -->
          <div class="lg:col-span-2">
            <div v-if="session.summary" class="card">
              <h2 class="mb-2 text-lg font-semibold">{{ session.summary.title || 'Resumo' }}</h2>
              <p class="mb-4 text-sm leading-relaxed text-text-secondary">
                {{ session.summary.description || 'Sem resumo disponivel.' }}
              </p>

              <!-- Highlights -->
              <div v-if="session.summary.highlights?.length" class="mb-4">
                <h3 class="mb-2 flex items-center gap-2 text-sm font-semibold text-text">
                  <svg class="h-4 w-4 text-primary-light" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
                  </svg>
                  Destaques
                </h3>
                <ul class="space-y-1.5">
                  <li
                    v-for="(highlight, i) in session.summary.highlights"
                    :key="i"
                    class="flex items-start gap-2 text-sm text-text-secondary"
                  >
                    <span class="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary-light"></span>
                    {{ highlight }}
                  </li>
                </ul>
              </div>

              <!-- Decisions -->
              <div v-if="session.summary.decisions?.length" class="mb-4">
                <h3 class="mb-2 flex items-center gap-2 text-sm font-semibold text-text">
                  <svg class="h-4 w-4 text-primary-light" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Decisoes
                </h3>
                <ul class="space-y-1.5">
                  <li
                    v-for="(decision, i) in session.summary.decisions"
                    :key="i"
                    class="flex items-start gap-2 text-sm text-text-secondary"
                  >
                    <svg class="mt-0.5 h-4 w-4 shrink-0 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    {{ decision }}
                  </li>
                </ul>
              </div>

              <!-- Action Items -->
              <div v-if="session.summary.action_items?.length">
                <h3 class="mb-2 flex items-center gap-2 text-sm font-semibold text-text">
                  <svg class="h-4 w-4 text-primary-light" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                  </svg>
                  Action Items
                </h3>
                <ul class="space-y-1.5">
                  <li
                    v-for="(item, i) in session.summary.action_items"
                    :key="i"
                    class="flex items-start gap-2 text-sm text-text-secondary"
                  >
                    <svg class="mt-0.5 h-4 w-4 shrink-0 text-warning" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {{ item }}
                  </li>
                </ul>
              </div>
            </div>
            <div v-else class="card">
              <p class="text-sm text-text-secondary">Resumo indisponivel.</p>
            </div>
          </div>

          <!-- Right: Flashcards Panel -->
          <div class="lg:col-span-3">
            <DeckList />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, watch, computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useSessionStore } from '@/store'
import DeckList from '@/components/flashcard/DeckList.vue'

const props = defineProps<{
  id: string
}>()

const route = useRoute()
const store = useSessionStore()

const session = computed(() => store.activeSession)

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('pt-BR', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

async function loadSession(): Promise<void> {
  await store.fetchSession(props.id ?? (route.params.id as string))
}

onMounted(loadSession)

watch(() => props.id, loadSession)
</script>
