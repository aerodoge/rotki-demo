<template>
  <aside class="w-64 border-r bg-card h-screen flex flex-col shadow-sm">
    <!-- Logo -->
    <div class="p-5 border-b">
      <h2 class="text-2xl font-bold text-foreground">rotki</h2>
    </div>

    <!-- Search -->
    <div class="p-4 border-b">
      <Input type="text" placeholder="Search..." class="w-full" />
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto py-2">
      <router-link
        to="/dashboard"
        class="flex items-center px-5 py-3 mx-2 my-0.5 text-sm text-muted-foreground hover:text-foreground hover:bg-gradient-accent rounded-lg transition-all cursor-pointer group"
        active-class="bg-gradient-primary text-white font-medium shadow-md"
      >
        <span class="mr-3">📊</span>
        <span>Dashboard</span>
      </router-link>

      <!-- Accounts Section -->
      <div class="my-1">
        <div
          class="flex items-center justify-between px-5 py-3 mx-2 my-0.5 text-sm text-muted-foreground hover:text-foreground hover:bg-gradient-accent rounded-lg transition-all cursor-pointer group"
          @click="toggleAccounts"
        >
          <div class="flex items-center">
            <span class="mr-3">👤</span>
            <span>Accounts</span>
          </div>
          <span
            class="text-xs transform transition-transform"
            :class="{ 'rotate-90': accountsExpanded }"
            >▶</span
          >
        </div>

        <div v-if="accountsExpanded" class="ml-4 space-y-0.5 mt-1">
          <router-link
            to="/accounts/evm"
            class="flex items-center px-5 py-2.5 mx-2 text-sm text-muted-foreground hover:text-foreground hover:bg-gradient-accent rounded-lg transition-all"
            active-class="bg-gradient-primary text-white font-medium shadow-sm"
          >
            <span class="mr-3">⟠</span>
            <span>EVM Accounts</span>
          </router-link>
          <div
            class="flex items-center px-5 py-2.5 mx-2 text-sm text-muted-foreground/50 cursor-not-allowed rounded-lg"
          >
            <span class="mr-3">₿</span>
            <span>Bitcoin Accounts</span>
          </div>
          <div
            class="flex items-center px-5 py-2.5 mx-2 text-sm text-muted-foreground/50 cursor-not-allowed rounded-lg"
          >
            <span class="mr-3">◎</span>
            <span>Solana Accounts</span>
          </div>
          <div
            class="flex items-center px-5 py-2.5 mx-2 text-sm text-muted-foreground/50 cursor-not-allowed rounded-lg"
          >
            <span class="mr-3">●</span>
            <span>Substrate Accounts</span>
          </div>
        </div>
      </div>

      <!-- Settings Section -->
      <div class="my-1">
        <div
          class="flex items-center justify-between px-5 py-3 mx-2 my-0.5 text-sm text-muted-foreground hover:text-foreground hover:bg-gradient-accent rounded-lg transition-all cursor-pointer group"
          @click="toggleSettings"
        >
          <div class="flex items-center">
            <span class="mr-3">⚙️</span>
            <span>Settings</span>
          </div>
          <span
            class="text-xs transform transition-transform"
            :class="{ 'rotate-90': settingsExpanded }"
            >▶</span
          >
        </div>

        <div v-if="settingsExpanded" class="ml-4 space-y-0.5 mt-1">
          <router-link
            to="/settings/rpc-nodes"
            class="flex items-center px-5 py-2.5 mx-2 text-sm text-muted-foreground hover:text-foreground hover:bg-gradient-accent rounded-lg transition-all"
            active-class="bg-gradient-primary text-white font-medium shadow-sm"
          >
            <span class="mr-3">🔗</span>
            <span>RPC Nodes</span>
          </router-link>
        </div>
      </div>
    </nav>

    <!-- Version -->
    <div class="p-4 border-t text-xs text-muted-foreground flex items-center justify-between">
      <span>v1.0.0</span>
      <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Input } from '@/components/ui/input'

const accountsExpanded = ref(true)
const settingsExpanded = ref(false)

const toggleAccounts = () => {
  accountsExpanded.value = !accountsExpanded.value
}

const toggleSettings = () => {
  settingsExpanded.value = !settingsExpanded.value
}
</script>
