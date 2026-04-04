<template>
  <div
    class="group card relative cursor-pointer overflow-hidden"
    @click="navigateToSession"
  >
    <!-- Status Indicator Dot -->
    <div
      class="absolute -top-1 -right-1 h-3 w-3 rounded-full ring-2 ring-bg"
      :class="{
        'bg-success': session.status === 'completed',
        'bg-warning': session.status === 'processing',
        'bg-danger': session.status === 'failed',
      }"
    >
      <span
        v-if="session.status === 'processing'"
        class="absolute inset-0 animate-ping rounded-full bg-warning opacity-75"
      ></span>
    </div>

    <!-- Card Content -->
    <div class="flex flex-col gap-3">
      <!-- Icon + Title -->
      <div class="flex items-center gap-3">
        <div
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-primary-muted/40 text-primary-light transition-colors group-hover:bg-primary-muted/60"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>
        <h3 class="line-clamp-1 text-base font-semibold text-text transition-colors group-hover:text-primary-light">
          {{ session.title }}
        </h3>
      </div>

      <!-- Meta Info -->
      <div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-text-secondary">
        <span class="flex items-center gap-1">
          <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {{ session.duration }}min
        </span>
        <span class="flex items-center gap-1">
          <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          {{ session.flashcards?.length || 0 }} cards
        </span>
        <span>{{ formatDate(session.created_at) }}</span>
      </div>

      <!-- Status Badge -->
      <div class="flex items-center justify-between">
        <span
          class="badge"
          :class="{
            'badge-success': session.status === 'completed',
            'badge-warning': session.status === 'processing',
            'badge-danger': session.status === 'failed',
          }"
        >
          {{ session.status === 'completed' ? 'Pronto' : session.status === 'processing' ? 'Processando' : 'Erro' }}
        </span>

        <!-- Hover Action -->
        <button
          class="rounded-md p-1 text-text-dim opacity-0 transition-all hover:bg-card-border hover:text-text group-hover:opacity-100"
          @click.stop="deleteAndConfirm"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useSessionStore, type Session } from '@/store'

const props = defineProps<{
  session: Session
}>()

const router = useRouter()
const store = useSessionStore()

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffHrs = diffMs / (1000 * 60 * 60)

  if (diffHrs < 1) return 'Agora'
  if (diffHrs < 24) return `${Math.floor(diffHrs)}h atras`

  return date.toLocaleDateString('pt-BR', { month: 'short', day: 'numeric' })
}

function navigateToSession(): void {
  if (props.session.status === 'processing') return
  router.push({ name: 'SessionDetail', params: { id: props.session.id } })
}

async function deleteAndConfirm(): Promise<void> {
  if (confirm(`Deletar "${props.session.title}"?`)) {
    await store.deleteSession(props.session.id)
  }
}
</script>
