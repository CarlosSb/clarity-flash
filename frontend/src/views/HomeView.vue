<template>
  <div class="min-h-screen bg-bg">
    <section class="px-4 py-6 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="mx-auto mb-8 max-w-6xl">
        <div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h1 class="text-3xl font-bold tracking-tight sm:text-4xl">
              <span class="text-gradient">AulaFlash</span>
            </h1>
            <p class="mt-1 text-text-secondary">
              Transform your lecture recordings into study tools
            </p>
          </div>
          <button class="btn-primary mt-4 sm:mt-0" @click="$emit('upload')">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            New Session
          </button>
        </div>
      </div>

      <!-- Stats Grid -->
      <StatsGrid class="mx-auto mb-8 max-w-6xl" />

      <!-- Content -->
      <div class="mx-auto max-w-6xl">
        <!-- Loading State -->
        <div v-if="store.loading" class="space-y-4">
          <div v-for="i in 3" :key="i" class="card h-32 animate-pulse">
            <div class="flex items-start gap-4">
              <div class="h-3 w-3 shrink-0 rounded-full bg-card-border"></div>
              <div class="flex-1 space-y-3">
                <div class="h-5 w-2/3 rounded bg-card-border"></div>
                <div class="h-4 w-1/2 rounded bg-card-border"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="store.error" class="card border-danger/30 bg-danger/5">
          <div class="flex items-center gap-3">
            <svg class="h-5 w-5 shrink-0 text-danger" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M5.07 19H19a2 2 0 001.75-2.97l-6.93-12a2 2 0 00-3.5 0l-6.93 12A2 2 0 005.07 19z" />
            </svg>
            <div>
              <p class="text-sm font-medium text-danger">{{ store.error }}</p>
              <button class="mt-2 btn-secondary text-xs" @click="store.fetchSessions()">
                Try Again
              </button>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="store.sessions.length === 0" class="text-center py-16">
          <div class="mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-card">
            <svg class="h-10 w-10 text-text-dim" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <h3 class="text-xl font-semibold text-text">No sessions yet</h3>
          <p class="mt-2 text-text-secondary">
            Start by uploading your first lecture recording to generate flashcards.
          </p>
          <button class="btn-primary mt-6" @click="$emit('upload')">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Create Your First Session
          </button>
        </div>

        <!-- Bento Grid of Sessions -->
        <div v-else class="bento-grid">
          <SessionCard
            v-for="session in store.sessions"
            :key="session.id"
            :session="session"
          />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSessionStore } from '@/store'
import SessionCard from '@/components/bento/SessionCard.vue'
import StatsGrid from '@/components/bento/StatsGrid.vue'

defineEmits<{
  upload: []
}>()

const store = useSessionStore()

onMounted(() => {
  store.fetchSessions()
})
</script>
