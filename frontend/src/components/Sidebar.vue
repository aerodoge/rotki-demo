<template>
  <aside class="sidebar">
    <div class="logo">
      <h2>rotki</h2>
    </div>

    <div class="search">
      <input type="text" placeholder="Search" />
    </div>

    <nav class="nav">
      <router-link to="/dashboard" class="nav-item">
        <span>Dashboard</span>
      </router-link>

      <div class="nav-section">
        <div class="nav-item" @click="toggleAccounts">
          <span>Accounts</span>
          <span class="arrow">{{ accountsExpanded ? '▼' : '▶' }}</span>
        </div>

        <div v-if="accountsExpanded" class="nav-submenu">
          <router-link to="/accounts/evm" class="nav-subitem">
            <span>EVM Accounts</span>
          </router-link>
          <div class="nav-subitem disabled">
            <span>Bitcoin Accounts</span>
          </div>
          <div class="nav-subitem disabled">
            <span>Solana Accounts</span>
          </div>
          <div class="nav-subitem disabled">
            <span>Substrate Accounts</span>
          </div>
        </div>
      </div>

      <div class="nav-item">
        <span>Balances</span>
      </div>

      <div class="nav-item">
        <span>History</span>
      </div>

      <div class="nav-item">
        <span>Onchain</span>
      </div>

      <div class="nav-item">
        <span>Staking</span>
      </div>

      <div class="nav-item">
        <span>Statistics</span>
      </div>

      <div class="nav-item">
        <span>Profit and Loss Report</span>
      </div>

      <div class="nav-section">
        <div class="nav-item" @click="toggleSettings">
          <span>⚙️ Settings</span>
          <span class="arrow">{{ settingsExpanded ? '▼' : '▶' }}</span>
        </div>

        <div v-if="settingsExpanded" class="nav-submenu">
          <router-link to="/settings/rpc-nodes" class="nav-subitem">
            <span>RPC Nodes</span>
          </router-link>
        </div>
      </div>
    </nav>

    <div class="version">v1.0.0</div>
  </aside>
</template>

<script setup>
import { ref } from 'vue'

const accountsExpanded = ref(true)
const settingsExpanded = ref(false)

const toggleAccounts = () => {
  accountsExpanded.value = !accountsExpanded.value
}

const toggleSettings = () => {
  settingsExpanded.value = !settingsExpanded.value
}
</script>

<style scoped>
.sidebar {
  width: 260px;
  background: white;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.logo {
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.logo h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
}

.search {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.search input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
}

.search input:focus {
  border-color: #6366f1;
}

.nav {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.nav-item {
  padding: 10px 20px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #6b7280;
  font-size: 14px;
  transition: background 0.2s;
  text-decoration: none;
}

.nav-item:hover {
  background: #f9fafb;
  color: #1f2937;
}

.nav-item.router-link-active {
  background: #eef2ff;
  color: #4f46e5;
  font-weight: 500;
}

.nav-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.nav-submenu {
  background: #f9fafb;
}

.nav-subitem {
  padding: 8px 20px 8px 40px;
  cursor: pointer;
  color: #6b7280;
  font-size: 14px;
  transition: background 0.2s;
  text-decoration: none;
  display: block;
}

.nav-subitem:hover:not(.disabled) {
  background: #f3f4f6;
  color: #1f2937;
}

.nav-subitem.router-link-active {
  background: #eef2ff;
  color: #4f46e5;
  font-weight: 500;
}

.nav-subitem.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.arrow {
  font-size: 10px;
}

.version {
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
  font-size: 12px;
  color: #9ca3af;
}
</style>
