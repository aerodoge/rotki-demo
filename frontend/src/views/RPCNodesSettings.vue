<template>
  <div class="container mx-auto p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">RPC Node Settings</h1>
        <p class="text-muted-foreground">Manage and view your RPC nodes</p>
      </div>
      <Button @click="openAddNodeModal(activeChain)">
        <span class="mr-2">+</span>
        Add Node
      </Button>
    </div>

    <!-- Chain Tabs -->
    <Tabs v-model="activeChain" class="w-full">
      <TabsList>
        <TabsTrigger v-for="chain in supportedChains" :key="chain.id" :value="chain.id">
          <img
            v-if="chain.logo_url"
            :src="chain.logo_url"
            :alt="chain.name"
            class="w-4 h-4 mr-2 rounded-full"
          />
          {{ chain.name }}
        </TabsTrigger>
      </TabsList>
    </Tabs>

    <!-- RPC Nodes Table -->
    <Card>
      <CardHeader>
        <CardTitle>Nodes</CardTitle>
        <CardDescription>View and manage RPC nodes for {{ activeChainName }}</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="text-muted-foreground">Loading nodes...</div>
        </div>

        <div
          v-else-if="currentChainNodes.length === 0"
          class="flex flex-col items-center justify-center py-12 space-y-4"
        >
          <div class="text-muted-foreground">No RPC nodes configured for this chain.</div>
          <Button variant="outline" @click="openAddNodeModal(activeChain)">
            Add your first node
          </Button>
        </div>

        <div v-else class="space-y-4">
          <!-- Table Header -->
          <div
            class="grid grid-cols-12 gap-4 px-4 py-2 border-b text-sm font-medium text-muted-foreground"
          >
            <div class="col-span-4">Node</div>
            <div class="col-span-2">Weight</div>
            <div class="col-span-3">Connectivity</div>
            <div class="col-span-3 text-right">Actions</div>
          </div>

          <!-- Table Rows -->
          <div
            v-for="node in currentChainNodes"
            :key="node.id"
            class="grid grid-cols-12 gap-4 px-4 py-4 items-center border rounded-lg hover:bg-accent/50 transition-colors"
          >
            <!-- Node Info -->
            <div class="col-span-4 flex items-center gap-3">
              <div
                class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center text-lg"
              >
                🌐
              </div>
              <div class="min-w-0">
                <div class="font-medium">{{ node.name }}</div>
                <div class="text-sm text-muted-foreground truncate">{{ node.url }}</div>
              </div>
            </div>

            <!-- Weight -->
            <div class="col-span-2">
              <div class="text-xl font-semibold">{{ node.weight }}%</div>
            </div>

            <!-- Connectivity -->
            <div class="col-span-3">
              <Badge :variant="node.is_connected ? 'default' : 'destructive'">
                {{ node.is_connected ? 'CONNECTED' : 'DISCONNECTED' }}
              </Badge>
            </div>

            <!-- Actions -->
            <div class="col-span-3 flex items-center justify-end gap-2">
              <Switch :checked="node.is_enabled" @update:checked="() => toggleNodeEnabled(node)" />
              <Button variant="ghost" size="icon" @click="editNode(node)" title="Edit"> ✏️ </Button>
              <Button variant="ghost" size="icon" @click="deleteNode(node)" title="Delete">
                🗑️
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Add/Edit Node Dialog -->
    <Dialog :open="showAddNodeModal || showEditNodeModal" @update:open="closeModals">
      <DialogContent class="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{{ showEditNodeModal ? 'Edit Node' : 'Add Node' }}</DialogTitle>
          <DialogDescription>
            {{ showEditNodeModal ? 'Update' : 'Add' }} RPC node configuration
          </DialogDescription>
        </DialogHeader>

        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="chain">Chain *</Label>
            <Select v-model="nodeForm.chain_id" :disabled="showEditNodeModal">
              <SelectTrigger id="chain">
                <SelectValue placeholder="Select chain" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="chain in supportedChains" :key="chain.id" :value="chain.id">
                  {{ chain.name }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="space-y-2">
            <Label for="name">Node Name *</Label>
            <Input id="name" v-model="nodeForm.name" placeholder="e.g., 0xRPC, PublicNode" />
          </div>

          <div class="space-y-2">
            <Label for="url">RPC URL *</Label>
            <Input id="url" v-model="nodeForm.url" type="url" placeholder="https://..." />
          </div>

          <div class="space-y-2">
            <Label for="weight">Weight (0-100) *</Label>
            <Input
              id="weight"
              v-model.number="nodeForm.weight"
              type="number"
              min="0"
              max="100"
              placeholder="100"
            />
            <p class="text-xs text-muted-foreground">
              Higher weight means more requests will be routed to this node
            </p>
          </div>

          <div class="space-y-2">
            <Label for="priority">Priority</Label>
            <Input
              id="priority"
              v-model.number="nodeForm.priority"
              type="number"
              min="0"
              placeholder="0"
            />
            <p class="text-xs text-muted-foreground">
              Higher priority nodes are preferred when weights are equal
            </p>
          </div>

          <div class="space-y-2">
            <Label for="timeout">Timeout (seconds)</Label>
            <Input
              id="timeout"
              v-model.number="nodeForm.timeout"
              type="number"
              min="1"
              placeholder="30"
            />
          </div>

          <div class="flex items-center space-x-2">
            <Switch id="enabled" v-model:checked="nodeForm.is_enabled" />
            <Label for="enabled" class="cursor-pointer">Enable this node</Label>
          </div>

          <Alert v-if="nodeFormError" variant="destructive">
            <AlertDescription>{{ nodeFormError }}</AlertDescription>
          </Alert>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="closeModals">Cancel</Button>
          <Button @click="saveNode" :disabled="!isFormValid">
            {{ showEditNodeModal ? 'Update' : 'Add' }} Node
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { rpcNodesAPI, chainsAPI } from '@/api/client'
import type { Chain, RPCNode } from '@/types'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Switch } from '@/components/ui/switch'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Alert, AlertDescription } from '@/components/ui/alert'

// State
const loading = ref(true)
const rpcNodesByChain = ref<Record<string, RPCNode[]>>({})
const supportedChains = ref<Chain[]>([])
const activeChain = ref('eth')
const showAddNodeModal = ref(false)
const showEditNodeModal = ref(false)
const nodeFormError = ref('')

// Node form
const nodeForm = ref({
  id: null as number | null,
  chain_id: '',
  name: '',
  url: '',
  weight: 100,
  priority: 0,
  timeout: 30,
  is_enabled: true
})

// Computed
const activeChainName = computed(() => {
  const chain = supportedChains.value.find((c) => c.id === activeChain.value)
  return chain?.name || activeChain.value
})

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

const resetForm = (preselectedChainId: string | null = null) => {
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

const openAddNodeModal = (chainId: string | null = null) => {
  resetForm(chainId)
  showAddNodeModal.value = true
}

const closeModals = () => {
  showAddNodeModal.value = false
  showEditNodeModal.value = false
  resetForm()
}

const editNode = (node: RPCNode) => {
  nodeForm.value = {
    id: node.id,
    chain_id: node.chain_id,
    name: node.name,
    url: node.url,
    weight: node.weight || 100,
    priority: node.priority || 0,
    timeout: 30,
    is_enabled: node.is_active
  }
  showEditNodeModal.value = true
}

const saveNode = async () => {
  nodeFormError.value = ''

  try {
    if (showEditNodeModal.value && nodeForm.value.id) {
      await rpcNodesAPI.update(nodeForm.value.id, nodeForm.value)
    } else {
      await rpcNodesAPI.create(nodeForm.value)
    }

    await loadRPCNodes()
    closeModals()
  } catch (error: any) {
    console.error('Failed to save node:', error)
    nodeFormError.value = error.response?.data?.error || 'Failed to save node'
  }
}

const deleteNode = async (node: RPCNode) => {
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

const toggleNodeEnabled = async (node: RPCNode) => {
  try {
    await rpcNodesAPI.update(node.id, {
      ...node,
      is_active: !node.is_active
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
