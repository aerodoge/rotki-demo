<template>
  <div class="rpc-nodes-settings">
    <!-- Header Section -->
    <div class="settings-header">
      <div>
        <h2>RPC node setting</h2>
        <p class="subtitle">Manage and view your RPC node.</p>
      </div>
      <button class="btn-add-node" @click="openAddNodeModal(activeChain)">+ Add Node</button>
    </div>

    <!-- Chain Tabs -->
    <div class="chain-tabs-container">
      <div class="chain-tabs">
        <button
          v-for="chain in supportedChains"
          :key="chain.id"
          :class="['chain-tab', { active: activeChain === chain.id }]"
          @click="activeChain = chain.id"
        >
          <img v-if="chain.logo_url" :src="chain.logo_url" :alt="chain.name" class="chain-icon" />
          <span>{{ chain.name }}</span>
        </button>
      </div>
    </div>

    <!-- RPC Nodes Table -->
    <div class="nodes-table-container">
      <div class="nodes-table-header">
        <div class="col-node">Node</div>
        <div class="col-weight">Node Weight</div>
        <div class="col-connectivity">Connectivity</div>
        <div class="col-actions"></div>
      </div>

      <div v-if="loading" class="loading-state">Loading nodes...</div>

      <div v-else-if="currentChainNodes.length === 0" class="empty-state">
        <p>No RPC nodes configured for this chain.</p>
        <button class="btn-secondary" @click="openAddNodeModal(activeChain)">Add your first node</button>
      </div>

      <div v-else class="nodes-list">
        <div v-for="node in currentChainNodes" :key="node.id" class="node-row">
          <div class="col-node">
            <div class="node-icon">🌐</div>
            <div class="node-info">
              <div class="node-name">{{ node.name }}</div>
              <div class="node-url">{{ node.url }}</div>
            </div>
          </div>

          <div class="col-weight">
            <div class="weight-display">{{ node.weight }}%</div>
          </div>

          <div class="col-connectivity">
            <span :class="['connectivity-badge', node.is_connected ? 'connected' : 'disconnected']">
              {{ node.is_connected ? 'CONNECTED' : 'DISCONNECTED' }}
            </span>
          </div>

          <div class="col-actions">
            <label class="toggle-switch">
              <input type="checkbox" :checked="node.is_enabled" @change="toggleNodeEnabled(node)" />
              <span class="toggle-slider"></span>
            </label>
            <button class="btn-icon" @click="editNode(node)" title="Edit">✏️</button>
            <button class="btn-icon" @click="deleteNode(node)" title="Delete">🗑️</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Node Modal -->
    <div
      v-if="showAddNodeModal || showEditNodeModal"
      class="modal-overlay"
      @click.self="closeModals"
    >
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ showEditNodeModal ? 'Edit Node' : 'Add Node' }}</h3>
          <button class="btn-close" @click="closeModals">×</button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>Chain *</label>
            <select v-model="nodeForm.chain_id" :disabled="showEditNodeModal">
              <option value="">Select chain</option>
              <option v-for="chain in supportedChains" :key="chain.id" :value="chain.id">
                {{ chain.name }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>Node Name *</label>
            <input v-model="nodeForm.name" type="text" placeholder="e.g., 0xRPC, PublicNode" />
          </div>

          <div class="form-group">
            <label>RPC URL *</label>
            <input v-model="nodeForm.url" type="text" placeholder="https://..." />
          </div>

          <div class="form-group">
            <label>Weight (0-100) *</label>
            <input
              v-model.number="nodeForm.weight"
              type="number"
              min="0"
              max="100"
              placeholder="100"
            />
            <small>Higher weight means more requests will be routed to this node</small>
          </div>

          <div class="form-group">
            <label>Priority</label>
            <input v-model.number="nodeForm.priority" type="number" min="0" placeholder="0" />
            <small>Higher priority nodes are preferred when weights are equal</small>
          </div>

          <div class="form-group">
            <label>Timeout (seconds)</label>
            <input v-model.number="nodeForm.timeout" type="number" min="1" placeholder="30" />
          </div>

          <div class="form-group checkbox-group">
            <label>
              <input type="checkbox" v-model="nodeForm.is_enabled" />
              Enable this node
            </label>
          </div>

          <div v-if="nodeFormError" class="error-message">
            {{ nodeFormError }}
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-secondary" @click="closeModals">Cancel</button>
          <button class="btn-primary" @click="saveNode" :disabled="!isFormValid">
            {{ showEditNodeModal ? 'Update' : 'Add' }} Node
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { rpcNodesAPI, chainsAPI } from '../api/client'

// State
const loading = ref(true)
const rpcNodesByChain = ref({})
const supportedChains = ref([])
const activeChain = ref('eth')
const showAddNodeModal = ref(false)
const showEditNodeModal = ref(false)
const nodeFormError = ref('')

// Node form
const nodeForm = ref({
  id: null,
  chain_id: '',
  name: '',
  url: '',
  weight: 100,
  priority: 0,
  timeout: 30,
  is_enabled: true
})

// Computed
const currentChainNodes = computed(() => {
  return rpcNodesByChain.value[activeChain.value] || []
})

const isFormValid = computed(() => {
  return (
    nodeForm.value.chain_id &&
    nodeForm.value.name &&
    nodeForm.value.url &&
    nodeForm.value.weight >= 0 &&
    nodeForm.value.weight <= 100
  )
})

// Methods
const loadChains = async () => {
  try {
    const response = await chainsAPI.list()
    supportedChains.value = response.data || []
    // Set default active chain to first chain if available
    if (supportedChains.value.length > 0 && !activeChain.value) {
      activeChain.value = supportedChains.value[0].id
    }
  } catch (error) {
    console.error('Failed to load chains:', error)
  }
}

const loadRPCNodes = async () => {
  loading.value = true
  try {
    const response = await rpcNodesAPI.grouped()
    rpcNodesByChain.value = response.data || {}
  } catch (error) {
    console.error('Failed to load RPC nodes:', error)
  } finally {
    loading.value = false
  }
}

const resetForm = (preselectedChainId = null) => {
  nodeForm.value = {
    id: null,
    chain_id: preselectedChainId || activeChain.value,
    name: '',
    url: '',
    weight: 100,
    priority: 0,
    timeout: 30,
    is_enabled: true
  }
  nodeFormError.value = ''
}

const openAddNodeModal = (chainId = null) => {
  resetForm(chainId)
  showAddNodeModal.value = true
}

const closeModals = () => {
  showAddNodeModal.value = false
  showEditNodeModal.value = false
  resetForm()
}

const editNode = (node) => {
  nodeForm.value = {
    id: node.id,
    chain_id: node.chain_id,
    name: node.name,
    url: node.url,
    weight: node.weight,
    priority: node.priority,
    timeout: node.timeout,
    is_enabled: node.is_enabled
  }
  showEditNodeModal.value = true
}

const saveNode = async () => {
  nodeFormError.value = ''

  try {
    if (showEditNodeModal.value) {
      // Update existing node
      await rpcNodesAPI.update(nodeForm.value.id, nodeForm.value)
    } else {
      // Create new node
      await rpcNodesAPI.create(nodeForm.value)
    }

    await loadRPCNodes()
    closeModals()
  } catch (error) {
    console.error('Failed to save node:', error)
    nodeFormError.value = error.response?.data?.error || 'Failed to save node'
  }
}

const deleteNode = async (node) => {
  if (!confirm(`Are you sure you want to delete "${node.name}"?`)) {
    return
  }

  try {
    await rpcNodesAPI.delete(node.id)
    await loadRPCNodes()
  } catch (error) {
    console.error('Failed to delete node:', error)
    alert('Failed to delete node')
  }
}

const toggleNodeEnabled = async (node) => {
  try {
    await rpcNodesAPI.update(node.id, {
      ...node,
      is_enabled: !node.is_enabled
    })
    await loadRPCNodes()
  } catch (error) {
    console.error('Failed to toggle node:', error)
    alert('Failed to update node status')
  }
}

// Lifecycle
onMounted(async () => {
  await loadChains()
  await loadRPCNodes()
})
</script>

<style scoped>
.rpc-nodes-settings {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.settings-header h2 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 4px 0;
  color: #1a1a1a;
}

.subtitle {
  color: #666;
  margin: 0;
  font-size: 14px;
}

.btn-add-node {
  background: #5865f2;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-add-node:hover {
  background: #4752c4;
}

/* Chain Tabs */
.chain-tabs-container {
  margin-bottom: 24px;
  border-bottom: 1px solid #e0e0e0;
}

.chain-tabs {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 2px;
}

.chain-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: all 0.2s;
  white-space: nowrap;
}

.chain-tab:hover {
  color: #1a1a1a;
}

.chain-tab.active {
  color: #5865f2;
  border-bottom-color: #5865f2;
  font-weight: 500;
}

.chain-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
}

/* Nodes Table */
.nodes-table-container {
  background: white;
  border-radius: 8px;
  border: 1px solid #e0e0e0;
  overflow: hidden;
}

.nodes-table-header {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 140px;
  gap: 16px;
  padding: 16px 20px;
  background: #f8f8f8;
  font-size: 13px;
  font-weight: 600;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.nodes-list {
  display: flex;
  flex-direction: column;
}

.node-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 140px;
  gap: 16px;
  padding: 20px;
  border-top: 1px solid #e0e0e0;
  align-items: center;
}

.node-row:first-child {
  border-top: none;
}

.col-node {
  display: flex;
  align-items: center;
  gap: 12px;
}

.node-icon {
  font-size: 24px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f0f0;
  border-radius: 8px;
}

.node-info {
  flex: 1;
  min-width: 0;
}

.node-name {
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.node-url {
  font-size: 13px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.weight-display {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
}

.connectivity-badge {
  display: inline-block;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.connectivity-badge.connected {
  background: #d4f4dd;
  color: #1e7e34;
}

.connectivity-badge.disconnected {
  background: #ffe0e0;
  color: #c00;
}

.col-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: flex-end;
}

/* Toggle Switch */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
  cursor: pointer;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  border-radius: 24px;
  transition: 0.3s;
}

.toggle-slider:before {
  position: absolute;
  content: '';
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  border-radius: 50%;
  transition: 0.3s;
}

input:checked + .toggle-slider {
  background-color: #5865f2;
}

input:checked + .toggle-slider:before {
  transform: translateX(20px);
}

.btn-icon {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  padding: 6px;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.btn-icon:hover {
  opacity: 1;
}

/* Empty/Loading States */
.loading-state,
.empty-state {
  padding: 60px 20px;
  text-align: center;
  color: #666;
}

.empty-state p {
  margin-bottom: 16px;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e0e0e0;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
}

.btn-close {
  background: none;
  border: none;
  font-size: 28px;
  line-height: 1;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 32px;
  height: 32px;
}

.modal-body {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
  font-size: 14px;
  color: #1a1a1a;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d0d0d0;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #5865f2;
}

.form-group small {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: #666;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-weight: normal;
}

.checkbox-group input[type='checkbox'] {
  width: auto;
  cursor: pointer;
}

.error-message {
  padding: 12px;
  background: #ffe0e0;
  color: #c00;
  border-radius: 6px;
  font-size: 14px;
  margin-top: 16px;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #e0e0e0;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-primary,
.btn-secondary {
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.btn-primary {
  background: #5865f2;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #4752c4;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f0f0f0;
  color: #1a1a1a;
}

.btn-secondary:hover {
  background: #e0e0e0;
}
</style>
