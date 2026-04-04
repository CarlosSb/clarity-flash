<template>
  <div class="grid grid-cols-2 gap-4 sm:grid-cols-4">
    <!-- Total Sessions -->
    <div class="card text-center">
      <div class="mx-auto mb-2 flex h-10 w-10 items-center justify-center rounded-full bg-primary-muted/40">
        <svg class="h-5 w-5 text-primary-light" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
      </div>
      <p class="text-2xl font-bold text-text">{{ store.sessions.length }}</p>
      <p class="text-xs text-text-secondary">Sessions</p>
    </div>

    <!-- Total Recordings -->
    <div class="card text-center">
      <div class="mx-auto mb-2 flex h-10 w-10 items-center justify-center rounded-full bg-success/10">
        <svg class="h-5 w-5 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
        </svg>
      </div>
      <p class="text-2xl font-bold text-text">{{ store.totalRecordings }}</p>
      <p class="text-xs text-text-secondary">Recordings</p>
    </div>

    <!-- Flashcards -->
    <div class="card text-center">
      <div class="mx-auto mb-2 flex h-10 w-10 items-center justify-center rounded-full bg-warning/10">
        <svg class="h-5 w-5 text-warning" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
      </div>
      <p class="text-2xl font-bold text-text">{{ store.totalFlashcards }}</p>
      <p class="text-xs text-text-secondary">Flashcards</p>
    </div>

    <!-- Ready Sessions -->
    <div class="card text-center">
      <div class="mx-auto mb-2 flex h-10 w-10 items-center justify-center rounded-full bg-primary-muted/40">
        <svg class="h-5 w-5 text-primary-light" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <p class="text-2xl font-bold text-text">{{ store.readySessions.length }}</p>
      <p class="text-xs text-text-secondary">Ready</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSessionStore } from '@/store'

const store = useSessionStore()
</script>
