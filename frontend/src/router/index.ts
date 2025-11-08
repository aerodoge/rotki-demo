import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import EVMAccounts from '../views/EVMAccounts.vue'
import RPCNodesSettings from '../views/RPCNodesSettings.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard
  },
  {
    path: '/accounts/evm',
    name: 'EVMAccounts',
    component: EVMAccounts
  },
  {
    path: '/settings/rpc-nodes',
    name: 'RPCNodesSettings',
    component: RPCNodesSettings
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
