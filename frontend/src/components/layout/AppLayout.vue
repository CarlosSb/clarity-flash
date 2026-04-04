<template>
  <div class="flex min-h-screen flex-col">
    <!-- Top Nav -->
    <header class="sticky top-0 z-50 border-b border-card-border/50 bg-bg/80 backdrop-blur-xl">
      <div class="mx-auto flex h-14 max-w-6xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <RouterLink to="/" class="flex items-center gap-2">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
            <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <span class="text-lg font-bold tracking-tight text-text">AulaFlash</span>
        </RouterLink>
        <nav class="flex items-center gap-3">
          <RouterLink
            to="/"
            class="rounded-lg px-3 py-1.5 text-sm font-medium text-text-secondary transition-colors hover:bg-card-border/50 hover:text-text"
            active-class="bg-primary-muted text-primary-light"
          >
            Sessions
          </RouterLink>
        </nav>
      </div>
    </header>

    <!-- Content -->
    <main class="flex-1">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="border-t border-card-border/50 py-6 text-center text-xs text-text-dim">
      AulaFlash &mdash; Turn lectures into flashcards
    </footer>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>
