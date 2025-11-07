<template>
  <div class="evm-accounts">
    <div class="header">
      <h1>EVM accounts</h1>
      <div class="header-actions">
        <CurrencySelector />
      </div>
    </div>

    <div class="actions-bar">
      <div class="actions-left">
        <button class="btn-secondary" @click="showAddWalletModal = true">Add Wallet</button>
        <button class="btn-secondary" @click="showAddAddressModal = true">Add Address</button>
      </div>
    </div>

    <div class="content">
      <div class="table-container">
        <table class="accounts-table">
          <thead>
            <tr>
              <th>Account</th>
              <th>Chains</th>
              <th>Tags</th>
              <th>Assets</th>
              <th>{{ selectedCurrency }} value</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="wallet in wallets" :key="wallet.id">
              <tr class="wallet-row" @click="toggleWallet(wallet.id)">
                <td>
                  <div class="account-cell">
                    <span class="expand-icon">{{ expandedWallets[wallet.id] ? '▼' : '▶' }}</span>
                    <div class="account-avatar">{{ wallet.name[0].toUpperCase() }}</div>
                    <span class="account-name">{{ wallet.name }}</span>
                  </div>
                </td>
                <td>
                  <div class="chains">
                    <span class="chain-badge">{{ getWalletChainCount(wallet.id) }}</span>
                  </div>
                </td>
                <td>
                  <div class="tags">
                    <span v-for="tag in wallet.tags || []" :key="tag" class="tag">{{ tag }}</span>
                  </div>
                </td>
                <td>
                  <div class="assets">
                    <div class="asset-icons">
                      <span class="asset-icon">ETH</span>
                      <span class="asset-count">{{ getWalletAssetCount(wallet.id) }}+</span>
                    </div>
                  </div>
                </td>
                <td>
                  <div class="eth-value">
                    {{ currencySymbols[selectedCurrency]
                    }}{{ formatValue(getTotalValueByWallet(wallet.id)) }}
                  </div>
                </td>
                <td>
                  <div class="actions">
                    <button class="icon-btn" @click.stop="refreshWallet(wallet.id)" title="Refresh">
                      🔄
                    </button>
                    <button class="icon-btn" @click.stop="editWallet(wallet)" title="Edit">
                      ✏️
                    </button>
                    <button class="icon-btn" @click.stop="deleteWallet(wallet.id)" title="Delete">
                      🗑️
                    </button>
                  </div>
                </td>
              </tr>

              <template v-if="expandedWallets[wallet.id]">
                <!-- Address list header -->
                <tr v-if="getAddressesByWallet(wallet.id).length > 0" class="address-header-row">
                  <th class="address-header-cell">Account</th>
                  <th class="address-header-cell">Chains</th>
                  <th class="address-header-cell">Tags</th>
                  <th class="address-header-cell">Assets</th>
                  <th class="address-header-cell">{{ selectedCurrency }} value</th>
                  <th class="address-header-cell">Actions</th>
                </tr>
                <template v-for="address in getAddressesByWallet(wallet.id)" :key="address.id">
                  <tr class="address-row" @click="toggleAddress(address.id)">
                    <td>
                      <div class="account-cell address-cell">
                        <span class="expand-icon">{{
                          expandedAddresses[address.id] ? '▼' : '▶'
                        }}</span>
                        <div class="account-avatar small">
                          {{ address.label?.[0]?.toUpperCase() || 'A' }}
                        </div>
                        <div class="address-info">
                          <div class="address-label">{{ address.label || 'Unlabeled' }}</div>
                          <div class="address-hash">{{ formatAddress(address.address) }}</div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div class="chains">
                        <span
                          class="chain-icon"
                          v-for="chain in getAddressChains(address)"
                          :key="chain"
                        >
                          {{ getChainIcon(chain) }}
                        </span>
                      </div>
                    </td>
                    <td>
                      <div class="tags">
                        <span v-for="tag in address.tags || []" :key="tag" class="tag-small">{{
                          tag
                        }}</span>
                      </div>
                    </td>
                    <td>
                      <div class="assets">
                        <div class="asset-icons">
                          <div
                            v-for="token in (address.tokens || []).slice(0, 3)"
                            :key="token.id"
                            class="token-logo-wrapper"
                          >
                            <img
                              v-if="
                                token.logo_url &&
                                token.logo_url.length > 0 &&
                                !failedImages[token.id]
                              "
                              :src="token.logo_url"
                              :alt="token.symbol"
                              :title="token.symbol"
                              class="asset-token-logo"
                              @error="handleImageError($event, token.id)"
                            />
                            <div v-else class="asset-token-fallback" :title="token.symbol">
                              {{ token.symbol.substring(0, 2) }}
                            </div>
                          </div>
                          <span v-if="(address.tokens || []).length > 3" class="asset-count">
                            +{{ (address.tokens || []).length - 3 }}
                          </span>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div class="eth-value">
                        {{ currencySymbols[selectedCurrency]
                        }}{{ formatValue(getAddressValue(address)) }}
                      </div>
                    </td>
                    <td>
                      <div class="actions">
                        <button
                          class="icon-btn"
                          @click.stop="refreshAddress(address.id)"
                          title="Refresh"
                        >
                          🔄
                        </button>
                        <button class="icon-btn" @click.stop="editAddress(address)" title="Edit">
                          ✏️
                        </button>
                        <button
                          class="icon-btn"
                          @click.stop="deleteAddress(address.id)"
                          title="Delete"
                        >
                          🗑️
                        </button>
                      </div>
                    </td>
                  </tr>

                  <!-- Expanded token details for address -->
                  <template v-if="expandedAddresses[address.id]">
                    <tr v-if="!address.tokens || address.tokens.length === 0" class="no-tokens-row">
                      <td colspan="6">
                        <div class="no-tokens-message">
                          No tokens found for this address. Click the refresh button to sync data.
                        </div>
                      </td>
                    </tr>
                    <!-- View Tabs for address tokens -->
                    <tr v-if="address.tokens && address.tokens.length > 0" class="tabs-row">
                      <td colspan="6">
                        <div class="view-tabs">
                          <button
                            :class="[
                              'tab-button',
                              { active: addressViewMode[address.id] === 'aggregated' }
                            ]"
                            @click="setAddressViewMode(address.id, 'aggregated')"
                          >
                            Aggregated assets
                          </button>
                          <button
                            :class="[
                              'tab-button',
                              { active: addressViewMode[address.id] === 'perChain' }
                            ]"
                            @click="setAddressViewMode(address.id, 'perChain')"
                          >
                            Per chain
                          </button>
                        </div>
                      </td>
                    </tr>

                    <!-- Aggregated view: show all tokens -->
                    <template v-if="addressViewMode[address.id] === 'aggregated'">
                      <!-- Token details header -->
                      <tr
                        v-if="address.tokens && address.tokens.length > 0"
                        class="token-header-row"
                      >
                        <td colspan="6">
                          <div class="token-detail token-header">
                            <div class="token-icon-wrapper">Asset</div>
                            <div class="token-info-detail"></div>
                            <div class="token-location">Location</div>
                            <div class="token-price">Price in {{ selectedCurrency }}</div>
                            <div class="token-balance">Amount</div>
                            <div class="token-value">{{ selectedCurrency }} Value</div>
                          </div>
                        </td>
                      </tr>
                      <TokenDetailRow
                        v-for="token in getPaginatedTokens(address.id)"
                        :key="token.id"
                        :token="token"
                      />
                      <!-- Pagination controls -->
                      <tr v-if="address.tokens && address.tokens.length > 0" class="pagination-row">
                        <td colspan="6">
                          <div class="pagination-controls">
                            <div class="pagination-info">
                              <span class="pagination-label">Rows per page:</span>
                              <select
                                v-model="tokenPagination[address.id].pageSize"
                                @change="onPageSizeChange(address.id)"
                                class="page-size-select"
                              >
                                <option :value="10">10</option>
                                <option :value="20">20</option>
                                <option :value="50">50</option>
                                <option :value="100">100</option>
                              </select>
                              <span class="pagination-range">
                                Items #
                                {{ getStartIndex(address.id) }}-{{ getEndIndex(address.id) }} of
                                {{ address.tokens.length }}
                              </span>
                              <span class="pagination-page">
                                Page {{ tokenPagination[address.id]?.currentPage || 1 }} of
                                {{ getTotalPages(address.id) }}
                              </span>
                            </div>
                            <div class="pagination-buttons">
                              <button
                                @click="goToFirstPage(address.id)"
                                :disabled="isFirstPage(address.id)"
                                class="pagination-btn"
                                title="First page"
                              >
                                ⟨⟨
                              </button>
                              <button
                                @click="goToPreviousPage(address.id)"
                                :disabled="isFirstPage(address.id)"
                                class="pagination-btn"
                                title="Previous page"
                              >
                                ⟨
                              </button>
                              <button
                                @click="goToNextPage(address.id)"
                                :disabled="isLastPage(address.id)"
                                class="pagination-btn"
                                title="Next page"
                              >
                                ⟩
                              </button>
                              <button
                                @click="goToLastPage(address.id)"
                                :disabled="isLastPage(address.id)"
                                class="pagination-btn"
                                title="Last page"
                              >
                                ⟩⟩
                              </button>
                            </div>
                          </div>
                        </td>
                      </tr>
                    </template>

                    <!-- Per Chain view: show tokens grouped by chain -->
                    <template v-if="addressViewMode[address.id] === 'perChain'">
                      <template
                        v-for="chain in getAddressChainGroups(address.id)"
                        :key="chain.chainId"
                      >
                        <tr
                          class="chain-group-row"
                          @click="toggleAddressChain(address.id, chain.chainId)"
                        >
                          <td colspan="6">
                            <div class="chain-group-header">
                              <span class="expand-icon">{{
                                expandedAddressChains[`${address.id}_${chain.chainId}`] ? '▼' : '▶'
                              }}</span>
                              <img
                                v-if="chain.logoUrl && !failedChainImages[chain.chainId]"
                                :src="chain.logoUrl"
                                :alt="chain.name"
                                class="chain-logo-small"
                                @error="failedChainImages[chain.chainId] = true"
                              />
                              <div v-else class="chain-logo-placeholder-tiny">
                                {{ (chain.name || chain.chainId || '?').charAt(0).toUpperCase() }}
                              </div>
                              <span class="chain-name">{{ chain.name || chain.chainId }}</span>
                              <div class="chain-assets">
                                <div class="asset-icons">
                                  <div
                                    v-for="token in chain.tokens.slice(0, 2)"
                                    :key="token.id"
                                    class="token-logo-wrapper-small"
                                  >
                                    <img
                                      v-if="
                                        token.logo_url &&
                                        token.logo_url.length > 0 &&
                                        !failedImages[token.id]
                                      "
                                      :src="token.logo_url"
                                      :alt="token.symbol"
                                      :title="token.symbol"
                                      class="asset-token-logo-small"
                                      @error="handleImageError($event, token.id)"
                                    />
                                    <div
                                      v-else
                                      class="asset-token-fallback-small"
                                      :title="token.symbol"
                                    >
                                      {{ token.symbol.substring(0, 2) }}
                                    </div>
                                  </div>
                                  <span v-if="chain.tokens.length > 2" class="asset-count-small">
                                    +{{ chain.tokens.length - 2 }}
                                  </span>
                                </div>
                              </div>
                              <span class="chain-value">
                                {{ currencySymbols[selectedCurrency]
                                }}{{ formatValue(chain.totalValue) }}
                              </span>
                            </div>
                          </td>
                        </tr>
                        <template v-if="expandedAddressChains[`${address.id}_${chain.chainId}`]">
                          <!-- Chain token header -->
                          <tr class="token-header-row">
                            <td colspan="6">
                              <div class="token-detail token-header">
                                <div class="token-icon-wrapper">Asset</div>
                                <div class="token-info-detail"></div>
                                <div class="token-location">Location</div>
                                <div class="token-price">Price in {{ selectedCurrency }}</div>
                                <div class="token-balance">Amount</div>
                                <div class="token-value">{{ selectedCurrency }} Value</div>
                              </div>
                            </td>
                          </tr>
                          <TokenDetailRow
                            v-for="token in getPaginatedAddressChainTokens(
                              address.id,
                              chain.chainId
                            )"
                            :key="token.id"
                            :token="token"
                          />
                          <!-- Chain pagination -->
                          <tr class="pagination-row">
                            <td colspan="6">
                              <div class="pagination-controls">
                                <div class="pagination-info">
                                  <span class="pagination-label">Rows per page:</span>
                                  <select
                                    v-model="
                                      addressChainPagination[`${address.id}_${chain.chainId}`]
                                        .pageSize
                                    "
                                    @change="
                                      onAddressChainPageSizeChange(address.id, chain.chainId)
                                    "
                                    class="page-size-select"
                                  >
                                    <option :value="10">10</option>
                                    <option :value="20">20</option>
                                    <option :value="50">50</option>
                                    <option :value="100">100</option>
                                  </select>
                                  <span class="pagination-range">
                                    Items #
                                    {{ getAddressChainStartIndex(address.id, chain.chainId) }}-{{
                                      getAddressChainEndIndex(address.id, chain.chainId)
                                    }}
                                    of {{ chain.tokens.length }}
                                  </span>
                                  <span class="pagination-page">
                                    Page
                                    {{
                                      addressChainPagination[`${address.id}_${chain.chainId}`]
                                        ?.currentPage || 1
                                    }}
                                    of
                                    {{ getAddressChainTotalPages(address.id, chain.chainId) }}
                                  </span>
                                </div>
                                <div class="pagination-buttons">
                                  <button
                                    @click="goToAddressChainFirstPage(address.id, chain.chainId)"
                                    :disabled="isAddressChainFirstPage(address.id, chain.chainId)"
                                    class="pagination-btn"
                                    title="First page"
                                  >
                                    ⟨⟨
                                  </button>
                                  <button
                                    @click="goToAddressChainPreviousPage(address.id, chain.chainId)"
                                    :disabled="isAddressChainFirstPage(address.id, chain.chainId)"
                                    class="pagination-btn"
                                    title="Previous page"
                                  >
                                    ⟨
                                  </button>
                                  <button
                                    @click="goToAddressChainNextPage(address.id, chain.chainId)"
                                    :disabled="isAddressChainLastPage(address.id, chain.chainId)"
                                    class="pagination-btn"
                                    title="Next page"
                                  >
                                    ⟩
                                  </button>
                                  <button
                                    @click="goToAddressChainLastPage(address.id, chain.chainId)"
                                    :disabled="isAddressChainLastPage(address.id, chain.chainId)"
                                    class="pagination-btn"
                                    title="Last page"
                                  >
                                    ⟩⟩
                                  </button>
                                </div>
                              </div>
                            </td>
                          </tr>
                        </template>
                      </template>
                    </template>
                  </template>
                </template>
              </template>
            </template>

            <tr class="total-row">
              <td><strong>Total</strong></td>
              <td></td>
              <td></td>
              <td></td>
              <td>
                <strong
                  >{{ currencySymbols[selectedCurrency] }}{{ formatValue(getTotalValue()) }}</strong
                >
              </td>
              <td></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Wallet Modal -->
    <div v-if="showAddWalletModal" class="modal-overlay" @click="showAddWalletModal = false">
      <div class="modal" @click.stop>
        <h2>Add Wallet</h2>
        <form @submit.prevent="handleAddWallet">
          <div class="form-group">
            <label>Wallet Name</label>
            <input v-model="newWallet.name" type="text" required placeholder="My Wallet" />
          </div>
          <div class="form-group">
            <label>Description (optional)</label>
            <textarea v-model="newWallet.description" placeholder="Wallet description"></textarea>
          </div>
          <div class="form-group">
            <label>Tags (optional)</label>
            <div class="tags-input">
              <div class="tags-display">
                <span v-for="(tag, index) in newWallet.tags" :key="index" class="tag-item">
                  {{ tag }}
                  <button type="button" class="tag-remove" @click="removeNewWalletTag(index)">
                    ×
                  </button>
                </span>
              </div>
              <div class="tag-add-section">
                <input
                  v-model="newWalletTagInput"
                  type="text"
                  placeholder="Add a tag"
                  @keyup.enter="addNewWalletTag"
                />
                <button type="button" class="btn-add-tag" @click="addNewWalletTag">Add</button>
              </div>
            </div>
          </div>
          <div class="form-group">
            <label>Enabled Chains (optional - leave empty for all chains)</label>
            <div class="chains-selection">
              <div v-for="chain in availableChains" :key="chain.id" class="chain-item">
                <img
                  v-if="!failedChainImages[chain.id]"
                  :src="`/images/chains/${chain.id}.png`"
                  :alt="chain.name"
                  class="chain-logo-inline"
                  @error="failedChainImages[chain.id] = true"
                />
                <div v-else class="chain-logo-placeholder">
                  {{ (chain.name || chain.id || '?').charAt(0).toUpperCase() }}
                </div>
                <div class="chain-name">{{ chain.name || chain.id || 'Unknown' }}</div>
                <input
                  type="checkbox"
                  :id="`chain-${chain.id}`"
                  :value="chain.id"
                  v-model="newWallet.enabled_chains"
                  class="chain-checkbox"
                />
              </div>
            </div>
          </div>
          <div class="form-actions">
            <button type="button" class="btn-secondary" @click="showAddWalletModal = false">
              Cancel
            </button>
            <button type="submit" class="btn-primary">Add Wallet</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Wallet Modal -->
    <div
      v-if="showEditWalletModal && editingWallet"
      class="modal-overlay"
      @click="showEditWalletModal = false"
    >
      <div class="modal" @click.stop>
        <h2>Edit Wallet</h2>
        <form @submit.prevent="handleEditWallet">
          <div class="form-group">
            <label>Wallet Name</label>
            <input v-model="editingWallet.name" type="text" required placeholder="My Wallet" />
          </div>
          <div class="form-group">
            <label>Description (optional)</label>
            <textarea
              v-model="editingWallet.description"
              placeholder="Wallet description"
            ></textarea>
          </div>
          <div class="form-group">
            <label>Tags (optional)</label>
            <div class="tags-input">
              <div class="tags-display">
                <span v-for="(tag, index) in editingWallet.tags" :key="index" class="tag-item">
                  {{ tag }}
                  <button type="button" class="tag-remove" @click="removeEditWalletTag(index)">
                    ×
                  </button>
                </span>
              </div>
              <div class="tag-add-section">
                <input
                  v-model="editWalletTagInput"
                  type="text"
                  placeholder="Add a tag"
                  @keyup.enter="addEditWalletTag"
                />
                <button type="button" class="btn-add-tag" @click="addEditWalletTag">Add</button>
              </div>
            </div>
          </div>
          <div class="form-group">
            <label>Enabled Chains (optional - leave empty for all chains)</label>
            <div class="chains-selection">
              <div v-for="chain in availableChains" :key="chain.id" class="chain-item">
                <img
                  v-if="!failedChainImages[chain.id]"
                  :src="`/images/chains/${chain.id}.png`"
                  :alt="chain.name"
                  class="chain-logo-inline"
                  @error="failedChainImages[chain.id] = true"
                />
                <div v-else class="chain-logo-placeholder">
                  {{ (chain.name || chain.id || '?').charAt(0).toUpperCase() }}
                </div>
                <div class="chain-name">{{ chain.name || chain.id || 'Unknown' }}</div>
                <input
                  type="checkbox"
                  :id="`edit-chain-${chain.id}`"
                  :value="chain.id"
                  v-model="editingWallet.enabled_chains"
                  class="chain-checkbox"
                />
              </div>
            </div>
          </div>
          <div class="form-actions">
            <button type="button" class="btn-secondary" @click="showEditWalletModal = false">
              Cancel
            </button>
            <button type="submit" class="btn-primary">Save Changes</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Add Address Modal -->
    <div v-if="showAddAddressModal" class="modal-overlay" @click="showAddAddressModal = false">
      <div class="modal" @click.stop>
        <h2>Add Address</h2>
        <form @submit.prevent="handleAddAddress">
          <div class="form-group">
            <label>Wallet</label>
            <select v-model="newAddress.wallet_id" required>
              <option value="">Select a wallet</option>
              <option v-for="wallet in wallets" :key="wallet.id" :value="wallet.id">
                {{ wallet.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>Address</label>
            <input v-model="newAddress.address" type="text" required placeholder="0x..." />
          </div>
          <div class="form-group">
            <label>Label (optional)</label>
            <input v-model="newAddress.label" type="text" placeholder="Main Address" />
          </div>
          <div class="form-group">
            <label>Tags (optional)</label>
            <div class="tags-input">
              <div class="tags-display">
                <span v-for="(tag, index) in newAddress.tags" :key="index" class="tag-item">
                  {{ tag }}
                  <button type="button" class="tag-remove" @click="removeNewAddressTag(index)">
                    ×
                  </button>
                </span>
              </div>
              <div class="tag-add-section">
                <input
                  v-model="newAddressTagInput"
                  type="text"
                  placeholder="Add a tag"
                  @keyup.enter="addNewAddressTag"
                />
                <button type="button" class="btn-add-tag" @click="addNewAddressTag">Add</button>
              </div>
            </div>
          </div>
          <div class="form-group">
            <label>Chain Type</label>
            <select v-model="newAddress.chain_type">
              <option value="EVM">EVM</option>
            </select>
          </div>
          <div class="form-actions">
            <button type="button" class="btn-secondary" @click="showAddAddressModal = false">
              Cancel
            </button>
            <button type="submit" class="btn-primary">Add Address</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Address Modal -->
    <div
      v-if="showEditAddressModal && editingAddress"
      class="modal-overlay"
      @click="showEditAddressModal = false"
    >
      <div class="modal" @click.stop>
        <h2>Edit Address</h2>
        <form @submit.prevent="handleEditAddress">
          <div class="form-group">
            <label>Address</label>
            <input v-model="editingAddress.address" type="text" disabled placeholder="0x..." />
          </div>
          <div class="form-group">
            <label>Label</label>
            <input v-model="editingAddress.label" type="text" placeholder="Address Label" />
          </div>
          <div class="form-group">
            <label>Tags (optional)</label>
            <div class="tags-input">
              <div class="tags-display">
                <span v-for="(tag, index) in editingAddress.tags" :key="index" class="tag-item">
                  {{ tag }}
                  <button type="button" class="tag-remove" @click="removeEditAddressTag(index)">
                    ×
                  </button>
                </span>
              </div>
              <div class="tag-add-section">
                <input
                  v-model="editAddressTagInput"
                  type="text"
                  placeholder="Add a tag"
                  @keyup.enter="addEditAddressTag"
                />
                <button type="button" class="btn-add-tag" @click="addEditAddressTag">Add</button>
              </div>
            </div>
          </div>
          <div class="form-actions">
            <button type="button" class="btn-secondary" @click="showEditAddressModal = false">
              Cancel
            </button>
            <button type="submit" class="btn-primary">Save Changes</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useWalletStore } from '../stores/wallet'
import { storeToRefs } from 'pinia'
import { chainsAPI } from '../api/client'
import { useCurrency } from '../composables/useCurrency'
import CurrencySelector from '../components/CurrencySelector.vue'
import TokenDetailRow from '../components/TokenDetailRow.vue'

const walletStore = useWalletStore()
const { wallets, addresses, loading } = storeToRefs(walletStore)

const availableChains = ref([])

// Use currency composable
const {
  currencySymbols,
  selectedCurrency,
  exchangeRates,
  formatValue,
  formatTokenPrice,
  formatTokenValue,
  updateExchangeRates,
  restoreSavedCurrency
} = useCurrency()

const expandedWallets = reactive({})
const expandedAddresses = reactive({})
const addressViewMode = reactive({}) // Track view mode for each address ('aggregated' or 'perChain')
const expandedAddressChains = reactive({}) // Track expanded chains for each address
const failedImages = reactive({}) // Track failed image loads
const failedChainImages = reactive({}) // Track failed chain logo loads
const tokenPagination = reactive({}) // Track pagination state for each address
const addressChainPagination = reactive({}) // Track pagination state for each address's chain
const showAddWalletModal = ref(false)
const showEditWalletModal = ref(false)
const showAddAddressModal = ref(false)
const showEditAddressModal = ref(false)
const editingWallet = ref(null)
const editingAddress = ref(null)

const newWallet = ref({
  name: '',
  description: '',
  tags: [],
  enabled_chains: []
})

const newAddress = ref({
  wallet_id: '',
  address: '',
  label: '',
  tags: [],
  chain_type: 'EVM'
})

// Tag input state
const newWalletTagInput = ref('')
const editWalletTagInput = ref('')
const newAddressTagInput = ref('')
const editAddressTagInput = ref('')

onMounted(async () => {
  // Restore saved currency preference
  restoreSavedCurrency()

  await walletStore.fetchWallets()
  await walletStore.fetchAddresses()

  // Fetch available chains
  try {
    const response = await chainsAPI.list()
    availableChains.value = response.data
    console.log('Available chains loaded:', availableChains.value)
  } catch (error) {
    console.error('Failed to fetch chains:', error)
  }

  // Update exchange rates from token data
  updateExchangeRates(addresses.value)

  // Debug: Check what logo URLs we have
  console.log('Addresses with tokens:', addresses.value)
  addresses.value.forEach((addr) => {
    if (addr.tokens && addr.tokens.length > 0) {
      console.log(
        `Address ${addr.address} tokens:`,
        addr.tokens.map((t) => ({
          symbol: t.symbol,
          logo_url: t.logo_url,
          has_logo: !!t.logo_url
        }))
      )
    }
  })
})

const toggleWallet = (walletId) => {
  console.log('toggleWallet called with:', walletId)
  console.log('Before:', expandedWallets)
  expandedWallets[walletId] = !expandedWallets[walletId]
  console.log('After:', expandedWallets)
}

const toggleAddress = (addressId) => {
  console.log('toggleAddress called with:', addressId)
  expandedAddresses[addressId] = !expandedAddresses[addressId]
  console.log('After:', expandedAddresses)

  // Initialize view mode for this address (default to aggregated)
  if (!addressViewMode[addressId]) {
    addressViewMode[addressId] = 'aggregated'
  }

  // Initialize pagination for this address if not exists
  if (!tokenPagination[addressId]) {
    tokenPagination[addressId] = {
      currentPage: 1,
      pageSize: 10
    }
  }
}

// Set view mode for a specific address
const setAddressViewMode = (addressId, mode) => {
  addressViewMode[addressId] = mode
}

// Pagination helper functions
const initPagination = (addressId) => {
  if (!tokenPagination[addressId]) {
    tokenPagination[addressId] = {
      currentPage: 1,
      pageSize: 10
    }
  }
}

const getPaginatedTokens = (addressId) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens) return []

  initPagination(addressId)
  const pagination = tokenPagination[addressId]

  // Sort tokens by value in selected currency (descending order)
  const sortedTokens = [...address.tokens].sort((a, b) => {
    // Get USD value first
    const usdValueA = a.usd_value || 0
    const usdValueB = b.usd_value || 0

    // Convert to selected currency using exchange rates
    const rate = exchangeRates.value[selectedCurrency.value] || 1
    const valueA = rate === 0 ? usdValueA : usdValueA / rate
    const valueB = rate === 0 ? usdValueB : usdValueB / rate

    return valueB - valueA
  })

  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize

  return sortedTokens.slice(start, end)
}

const getTotalPages = (addressId) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens) return 1

  initPagination(addressId)
  const pagination = tokenPagination[addressId]
  return Math.ceil(address.tokens.length / pagination.pageSize)
}

const getStartIndex = (addressId) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens || address.tokens.length === 0) return 0

  initPagination(addressId)
  const pagination = tokenPagination[addressId]
  return (pagination.currentPage - 1) * pagination.pageSize + 1
}

const getEndIndex = (addressId) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens) return 0

  initPagination(addressId)
  const pagination = tokenPagination[addressId]
  const end = pagination.currentPage * pagination.pageSize
  return Math.min(end, address.tokens.length)
}

const isFirstPage = (addressId) => {
  initPagination(addressId)
  return tokenPagination[addressId].currentPage === 1
}

const isLastPage = (addressId) => {
  initPagination(addressId)
  return tokenPagination[addressId].currentPage >= getTotalPages(addressId)
}

const goToFirstPage = (addressId) => {
  initPagination(addressId)
  tokenPagination[addressId].currentPage = 1
}

const goToPreviousPage = (addressId) => {
  initPagination(addressId)
  if (tokenPagination[addressId].currentPage > 1) {
    tokenPagination[addressId].currentPage--
  }
}

const goToNextPage = (addressId) => {
  initPagination(addressId)
  if (tokenPagination[addressId].currentPage < getTotalPages(addressId)) {
    tokenPagination[addressId].currentPage++
  }
}

const goToLastPage = (addressId) => {
  initPagination(addressId)
  tokenPagination[addressId].currentPage = getTotalPages(addressId)
}

const onPageSizeChange = (addressId) => {
  initPagination(addressId)
  // Reset to first page when page size changes
  tokenPagination[addressId].currentPage = 1
}

// Address chain grouping functions
const getAddressChainGroups = (addressId) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens) return []

  const chainMap = new Map()

  // Collect tokens by chain for this address
  address.tokens.forEach((token) => {
    if (!chainMap.has(token.chain_id)) {
      const chainInfo = availableChains.value.find((c) => c.id === token.chain_id)
      chainMap.set(token.chain_id, {
        chainId: token.chain_id,
        name: chainInfo?.name || token.chain_id,
        logoUrl: chainInfo?.logo_url || '',
        tokens: [],
        totalValue: 0
      })
    }
    chainMap.get(token.chain_id).tokens.push(token)
  })

  // Calculate total value and sort tokens
  const chains = Array.from(chainMap.values()).map((chain) => {
    const sortedTokens = [...chain.tokens].sort((a, b) => {
      const usdValueA = a.usd_value || 0
      const usdValueB = b.usd_value || 0
      const rate = exchangeRates.value[selectedCurrency.value] || 1
      const valueA = rate === 0 ? usdValueA : usdValueA / rate
      const valueB = rate === 0 ? usdValueB : usdValueB / rate
      return valueB - valueA
    })

    chain.tokens = sortedTokens
    chain.totalValue = chain.tokens.reduce((sum, token) => sum + (token.usd_value || 0), 0)
    return chain
  })

  // Sort chains by total value (descending)
  return chains.sort((a, b) => b.totalValue - a.totalValue)
}

const toggleAddressChain = (addressId, chainId) => {
  const key = `${addressId}_${chainId}`
  expandedAddressChains[key] = !expandedAddressChains[key]

  // Initialize pagination
  if (!addressChainPagination[key]) {
    addressChainPagination[key] = {
      currentPage: 1,
      pageSize: 10
    }
  }
}

// Address chain pagination functions
const initAddressChainPagination = (addressId, chainId) => {
  const key = `${addressId}_${chainId}`
  if (!addressChainPagination[key]) {
    addressChainPagination[key] = {
      currentPage: 1,
      pageSize: 10
    }
  }
}

const getPaginatedAddressChainTokens = (addressId, chainId) => {
  const chains = getAddressChainGroups(addressId)
  const chain = chains.find((c) => c.chainId === chainId)
  if (!chain || !chain.tokens) return []

  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  const pagination = addressChainPagination[key]
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize

  return chain.tokens.slice(start, end)
}

const getAddressChainTotalPages = (addressId, chainId) => {
  const chains = getAddressChainGroups(addressId)
  const chain = chains.find((c) => c.chainId === chainId)
  if (!chain || !chain.tokens) return 1

  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  const pagination = addressChainPagination[key]
  return Math.ceil(chain.tokens.length / pagination.pageSize)
}

const getAddressChainStartIndex = (addressId, chainId) => {
  const chains = getAddressChainGroups(addressId)
  const chain = chains.find((c) => c.chainId === chainId)
  if (!chain || !chain.tokens || chain.tokens.length === 0) return 0

  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  const pagination = addressChainPagination[key]
  return (pagination.currentPage - 1) * pagination.pageSize + 1
}

const getAddressChainEndIndex = (addressId, chainId) => {
  const chains = getAddressChainGroups(addressId)
  const chain = chains.find((c) => c.chainId === chainId)
  if (!chain || !chain.tokens) return 0

  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  const pagination = addressChainPagination[key]
  const end = pagination.currentPage * pagination.pageSize
  return Math.min(end, chain.tokens.length)
}

const isAddressChainFirstPage = (addressId, chainId) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  return addressChainPagination[key].currentPage === 1
}

const isAddressChainLastPage = (addressId, chainId) => {
  initAddressChainPagination(addressId, chainId)
  return (
    addressChainPagination[`${addressId}_${chainId}`].currentPage >=
    getAddressChainTotalPages(addressId, chainId)
  )
}

const goToAddressChainFirstPage = (addressId, chainId) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  addressChainPagination[key].currentPage = 1
}

const goToAddressChainPreviousPage = (addressId, chainId) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  if (addressChainPagination[key].currentPage > 1) {
    addressChainPagination[key].currentPage--
  }
}

const goToAddressChainNextPage = (addressId, chainId) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  if (addressChainPagination[key].currentPage < getAddressChainTotalPages(addressId, chainId)) {
    addressChainPagination[key].currentPage++
  }
}

const goToAddressChainLastPage = (addressId, chainId) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  addressChainPagination[key].currentPage = getAddressChainTotalPages(addressId, chainId)
}

const onAddressChainPageSizeChange = (addressId, chainId) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  addressChainPagination[key].currentPage = 1
}

// Chain grouping functions
const getChainGroups = () => {
  const chainMap = new Map()

  // Collect all tokens from all addresses
  addresses.value.forEach((address) => {
    if (address.tokens && address.tokens.length > 0) {
      address.tokens.forEach((token) => {
        if (!chainMap.has(token.chain_id)) {
          // Find chain info from available chains
          const chainInfo = availableChains.value.find((c) => c.id === token.chain_id)
          chainMap.set(token.chain_id, {
            chainId: token.chain_id,
            name: chainInfo?.name || token.chain_id,
            logoUrl: chainInfo?.logo_url || '',
            tokens: [],
            totalValue: 0
          })
        }
        chainMap.get(token.chain_id).tokens.push(token)
      })
    }
  })

  // Calculate total value for each chain and sort tokens by value
  const chains = Array.from(chainMap.values()).map((chain) => {
    // Sort tokens by value in selected currency (descending)
    const sortedTokens = [...chain.tokens].sort((a, b) => {
      const usdValueA = a.usd_value || 0
      const usdValueB = b.usd_value || 0
      const rate = exchangeRates.value[selectedCurrency.value] || 1
      const valueA = rate === 0 ? usdValueA : usdValueA / rate
      const valueB = rate === 0 ? usdValueB : usdValueB / rate
      return valueB - valueA
    })

    chain.tokens = sortedTokens
    chain.totalValue = chain.tokens.reduce((sum, token) => sum + (token.usd_value || 0), 0)
    return chain
  })

  // Sort chains by total value (descending)
  return chains.sort((a, b) => b.totalValue - a.totalValue)
}

const toggleChain = (chainId) => {
  expandedChains[chainId] = !expandedChains[chainId]

  // Initialize pagination for this chain if not exists
  if (!chainPagination[chainId]) {
    chainPagination[chainId] = {
      currentPage: 1,
      pageSize: 10
    }
  }
}

const refreshChain = async (chainId) => {
  // Refresh all addresses that have tokens on this chain
  const addressesToRefresh = addresses.value.filter(
    (addr) => addr.tokens && addr.tokens.some((token) => token.chain_id === chainId)
  )

  for (const address of addressesToRefresh) {
    await refreshAddress(address.id)
  }
}

// Chain pagination functions
const initChainPagination = (chainId) => {
  if (!chainPagination[chainId]) {
    chainPagination[chainId] = {
      currentPage: 1,
      pageSize: 10
    }
  }
}

const getPaginatedChainTokens = (chainId) => {
  const chains = getChainGroups()
  const chain = chains.find((c) => c.chainId === chainId)
  if (!chain || !chain.tokens) return []

  initChainPagination(chainId)
  const pagination = chainPagination[chainId]
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize

  return chain.tokens.slice(start, end)
}

const getChainTotalPages = (chainId) => {
  const chains = getChainGroups()
  const chain = chains.find((c) => c.chainId === chainId)
  if (!chain || !chain.tokens) return 1

  initChainPagination(chainId)
  const pagination = chainPagination[chainId]
  return Math.ceil(chain.tokens.length / pagination.pageSize)
}

const getChainStartIndex = (chainId) => {
  const chains = getChainGroups()
  const chain = chains.find((c) => c.chainId === chainId)
  if (!chain || !chain.tokens || chain.tokens.length === 0) return 0

  initChainPagination(chainId)
  const pagination = chainPagination[chainId]
  return (pagination.currentPage - 1) * pagination.pageSize + 1
}

const getChainEndIndex = (chainId) => {
  const chains = getChainGroups()
  const chain = chains.find((c) => c.chainId === chainId)
  if (!chain || !chain.tokens) return 0

  initChainPagination(chainId)
  const pagination = chainPagination[chainId]
  const end = pagination.currentPage * pagination.pageSize
  return Math.min(end, chain.tokens.length)
}

const isChainFirstPage = (chainId) => {
  initChainPagination(chainId)
  return chainPagination[chainId].currentPage === 1
}

const isChainLastPage = (chainId) => {
  initChainPagination(chainId)
  return chainPagination[chainId].currentPage >= getChainTotalPages(chainId)
}

const goToChainFirstPage = (chainId) => {
  initChainPagination(chainId)
  chainPagination[chainId].currentPage = 1
}

const goToChainPreviousPage = (chainId) => {
  initChainPagination(chainId)
  if (chainPagination[chainId].currentPage > 1) {
    chainPagination[chainId].currentPage--
  }
}

const goToChainNextPage = (chainId) => {
  initChainPagination(chainId)
  if (chainPagination[chainId].currentPage < getChainTotalPages(chainId)) {
    chainPagination[chainId].currentPage++
  }
}

const goToChainLastPage = (chainId) => {
  initChainPagination(chainId)
  chainPagination[chainId].currentPage = getChainTotalPages(chainId)
}

const onChainPageSizeChange = (chainId) => {
  initChainPagination(chainId)
  // Reset to first page when page size changes
  chainPagination[chainId].currentPage = 1
}

const getAddressesByWallet = (walletId) => {
  return addresses.value.filter((addr) => addr.wallet_id === walletId)
}

const getTotalValueByWallet = (walletId) => {
  return walletStore.getTotalValueByWallet(walletId)
}

const getTotalValue = () => {
  return walletStore.getTotalValue
}

const getAddressValue = (address) => {
  return address.tokens?.reduce((sum, token) => sum + (token.usd_value || 0), 0) || 0
}

// Handle chain logo loading errors
const handleChainImageError = (chainId) => {
  failedChainImages[chainId] = true
}

// Get chain initial for placeholder
const getChainInitial = (chain) => {
  const name = chain.name || chain.id || '?'
  return name.charAt(0).toUpperCase()
}

const getWalletChainCount = (walletId) => {
  const addrs = getAddressesByWallet(walletId)
  const chains = new Set()
  addrs.forEach((addr) => {
    addr.tokens?.forEach((token) => chains.add(token.chain_id))
  })
  return chains.size
}

const getWalletAssetCount = (walletId) => {
  const addrs = getAddressesByWallet(walletId)
  return addrs.reduce((sum, addr) => sum + (addr.tokens?.length || 0), 0)
}

const getAddressChains = (address) => {
  const chains = new Set()
  address.tokens?.forEach((token) => chains.add(token.chain_id))
  return Array.from(chains)
}

const getChainIcon = (chainId) => {
  const icons = {
    eth: '⟠',
    bsc: 'B',
    polygon: 'P',
    arbitrum: 'A',
    optimism: 'O'
  }
  return icons[chainId] || '⛓️'
}

const formatAddress = (address) => {
  return `${address.slice(0, 6)}...${address.slice(-4)}`
}

// Currency functions are now in useCurrency composable

const handleAddWallet = async () => {
  try {
    await walletStore.createWallet(newWallet.value)
    showAddWalletModal.value = false
    newWallet.value = { name: '', description: '', tags: [], enabled_chains: [] }
    newWalletTagInput.value = ''
  } catch (error) {
    alert('Failed to add wallet: ' + error.message)
  }
}

const editWallet = (wallet) => {
  editingWallet.value = {
    id: wallet.id,
    name: wallet.name,
    description: wallet.description || '',
    tags: wallet.tags || [],
    enabled_chains: wallet.enabled_chains || []
  }
  editWalletTagInput.value = ''
  showEditWalletModal.value = true
}

const handleEditWallet = async () => {
  try {
    await walletStore.updateWallet(editingWallet.value.id, {
      name: editingWallet.value.name,
      description: editingWallet.value.description,
      tags: editingWallet.value.tags,
      enabled_chains: editingWallet.value.enabled_chains
    })
    showEditWalletModal.value = false
    editingWallet.value = null
    editWalletTagInput.value = ''
    // Refresh data after update
    await walletStore.fetchWallets()
    await walletStore.fetchAddresses()
  } catch (error) {
    alert('Failed to update wallet: ' + error.message)
  }
}

const handleAddAddress = async () => {
  try {
    const newAddr = await walletStore.createAddress(newAddress.value)
    showAddAddressModal.value = false
    newAddress.value = { wallet_id: '', address: '', label: '', tags: [], chain_type: 'EVM' }
    newAddressTagInput.value = ''

    // Show success message
    alert('Address added! Syncing data in background...')

    // Refresh addresses to get updated data with tokens
    // Wait a bit for background sync to complete
    setTimeout(async () => {
      await walletStore.fetchAddresses()
    }, 2000)
  } catch (error) {
    alert('Failed to add address: ' + error.message)
  }
}

const handleImageError = (event, tokenId) => {
  // Mark this image as failed so Vue will show the fallback
  console.log('Image failed to load:', {
    tokenId,
    src: event.target.src,
    alt: event.target.alt
  })
  failedImages[tokenId] = true
}

const refreshWallet = async (walletId) => {
  try {
    await walletStore.refreshWallet(walletId)
  } catch (error) {
    alert('Failed to refresh wallet: ' + error.message)
  }
}

const refreshAddress = async (addressId) => {
  try {
    await walletStore.refreshAddress(addressId)
  } catch (error) {
    alert('Failed to refresh address: ' + error.message)
  }
}

const deleteWallet = async (walletId) => {
  if (confirm('Are you sure you want to delete this wallet?')) {
    try {
      await walletStore.deleteWallet(walletId)
    } catch (error) {
      alert('Failed to delete wallet: ' + error.message)
    }
  }
}

const editAddress = (address) => {
  editingAddress.value = {
    id: address.id,
    address: address.address,
    label: address.label || '',
    tags: address.tags || []
  }
  editAddressTagInput.value = ''
  showEditAddressModal.value = true
}

const handleEditAddress = async () => {
  try {
    await walletStore.updateAddress(editingAddress.value.id, {
      label: editingAddress.value.label,
      tags: editingAddress.value.tags
    })
    showEditAddressModal.value = false
    editingAddress.value = null
    editAddressTagInput.value = ''
    await walletStore.fetchWallets()
    await walletStore.fetchAddresses()
  } catch (error) {
    alert('Failed to update address: ' + error.message)
  }
}

const deleteAddress = async (addressId) => {
  if (confirm('Are you sure you want to delete this address?')) {
    try {
      await walletStore.deleteAddress(addressId)
    } catch (error) {
      alert('Failed to delete address: ' + error.message)
    }
  }
}

// Tag management functions
const addNewWalletTag = () => {
  const tag = newWalletTagInput.value.trim()
  if (tag && !newWallet.value.tags.includes(tag)) {
    newWallet.value.tags.push(tag)
    newWalletTagInput.value = ''
  }
}

const removeNewWalletTag = (index) => {
  newWallet.value.tags.splice(index, 1)
}

const addEditWalletTag = () => {
  const tag = editWalletTagInput.value.trim()
  if (tag && !editingWallet.value.tags.includes(tag)) {
    editingWallet.value.tags.push(tag)
    editWalletTagInput.value = ''
  }
}

const removeEditWalletTag = (index) => {
  editingWallet.value.tags.splice(index, 1)
}

const addNewAddressTag = () => {
  const tag = newAddressTagInput.value.trim()
  if (tag && !newAddress.value.tags.includes(tag)) {
    newAddress.value.tags.push(tag)
    newAddressTagInput.value = ''
  }
}

const removeNewAddressTag = (index) => {
  newAddress.value.tags.splice(index, 1)
}

const addEditAddressTag = () => {
  const tag = editAddressTagInput.value.trim()
  if (tag && !editingAddress.value.tags.includes(tag)) {
    editingAddress.value.tags.push(tag)
    editAddressTagInput.value = ''
  }
}

const removeEditAddressTag = (index) => {
  editingAddress.value.tags.splice(index, 1)
}
</script>

<style scoped>
.evm-accounts {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #1f2937;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

/* Currency styles moved to CurrencySelector component */

.actions-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.actions-left {
  display: flex;
  gap: 12px;
}

/* View Tabs */
.view-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  border-bottom: 2px solid #e5e7eb;
  padding-bottom: 0;
}

.tab-button {
  padding: 12px 24px;
  background: none;
  border: none;
  border-bottom: 3px solid transparent;
  font-size: 15px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  bottom: -2px;
}

.tab-button:hover {
  color: #1f2937;
  background: #f9fafb;
}

.tab-button.active {
  color: #4f46e5;
  border-bottom-color: #4f46e5;
  font-weight: 600;
}

.tabs-row {
  background: #fafbfc;
}

.tabs-row td {
  padding: 0 !important;
}

.chain-group-row {
  background: #f9fafb;
  cursor: pointer;
  transition: background 0.2s;
}

.chain-group-row:hover {
  background: #f3f4f6;
}

.chain-group-header {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  gap: 12px;
}

.chain-logo-small {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.chain-logo-placeholder-tiny {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.chain-assets {
  flex: 1;
  display: flex;
  justify-content: center;
}

.chain-value {
  font-weight: 600;
  color: #1f2937;
  min-width: 150px;
  text-align: right;
}

.token-logo-wrapper-small {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}

.asset-token-logo-small {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid #e5e7eb;
  background: white;
}

.asset-token-fallback-small {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 9px;
  font-weight: 700;
  color: #4f46e5;
  background: linear-gradient(135deg, #e0e7ff 0%, #c7d2fe 100%);
  border: 1px solid #e5e7eb;
  text-transform: uppercase;
}

.asset-count-small {
  padding: 2px 6px;
  background: #e5e7eb;
  color: #6b7280;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 600;
}

.chain-logo {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 12px;
}

.chain-logo-placeholder-small {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  margin-right: 12px;
  flex-shrink: 0;
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
  background: #4f46e5;
  color: white;
}

.btn-primary:hover {
  background: #4338ca;
}

.btn-secondary {
  background: white;
  color: #4f46e5;
  border: 1px solid #4f46e5;
}

.btn-secondary:hover {
  background: #eef2ff;
}

.table-container {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.accounts-table {
  width: 100%;
  border-collapse: collapse;
}

.accounts-table thead {
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}

.accounts-table th {
  padding: 12px 16px;
  text-align: left;
  font-size: 13px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.accounts-table td {
  padding: 16px;
  border-bottom: 1px solid #f3f4f6;
}

.wallet-row {
  cursor: pointer;
  transition: background 0.2s;
}

.wallet-row:hover {
  background: #f9fafb;
}

.address-header-row {
  background: #f3f4f6;
  border-top: 2px solid #e5e7eb;
  border-bottom: 2px solid #e5e7eb;
}

.address-header-cell {
  padding: 10px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background: #f3f4f6;
}

.address-row {
  background: #fafbfc;
  cursor: pointer;
}

.address-row:hover {
  background: #f3f4f6;
}

/* Token detail styles moved to TokenDetailRow component */

.token-header-row {
  background: #f3f4f6;
  border-top: 1px solid #e5e7eb;
}

.token-header {
  display: grid;
  grid-template-columns: 50px 1fr 150px 150px 150px 150px;
  gap: 16px;
  align-items: center;
  background: #f3f4f6;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 12px 16px;
  border-bottom: 2px solid #e5e7eb;
}

.token-header .token-price,
.token-header .token-balance,
.token-header .token-value {
  text-align: right;
}

.total-row {
  background: #f9fafb;
  font-weight: 600;
}

.account-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.address-cell {
  padding-left: 24px;
}

.expand-icon {
  font-size: 10px;
  color: #9ca3af;
}

.account-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 14px;
}

.account-avatar.small {
  width: 28px;
  height: 28px;
  font-size: 12px;
}

.account-name {
  font-weight: 500;
  color: #1f2937;
}

.address-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.address-label {
  font-weight: 500;
  color: #1f2937;
  font-size: 14px;
}

.address-hash {
  font-size: 12px;
  color: #6b7280;
  font-family: monospace;
}

.chains {
  display: flex;
  gap: 6px;
  align-items: center;
}

.chain-badge {
  padding: 4px 8px;
  background: #e0e7ff;
  color: #4f46e5;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.chain-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #e0e7ff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.tags {
  display: flex;
  gap: 6px;
}

.tag {
  padding: 4px 12px;
  background: #fce7f3;
  color: #be185d;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.tag-small {
  padding: 2px 8px;
  background: #e0e7ff;
  color: #4f46e5;
  border-radius: 10px;
  font-size: 11px;
}

.tags-input {
  border: 1px solid #d1d5db;
  border-radius: 8px;
  padding: 8px;
  background: white;
}

.tags-display {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 8px;
  min-height: 24px;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: #e0e7ff;
  color: #4f46e5;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.tag-remove {
  background: none;
  border: none;
  color: #4f46e5;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  padding: 0;
  margin: 0;
  font-weight: 600;
}

.tag-remove:hover {
  color: #3730a3;
}

.tag-add-section {
  display: flex;
  gap: 8px;
}

.tag-add-section input {
  flex: 1;
  padding: 6px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 13px;
}

.tag-add-section input:focus {
  outline: none;
  border-color: #4f46e5;
}

.btn-add-tag {
  padding: 6px 16px;
  background: #4f46e5;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  font-weight: 500;
}

.btn-add-tag:hover {
  background: #4338ca;
}

.assets {
  display: flex;
  align-items: center;
}

.asset-icons {
  display: flex;
  gap: 6px;
  align-items: center;
}

.token-logo-wrapper {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

.asset-token-logo {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid #e5e7eb;
  background: white;
}

.asset-token-fallback {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  color: #4f46e5;
  background: linear-gradient(135deg, #e0e7ff 0%, #c7d2fe 100%);
  border: 1px solid #e5e7eb;
  text-transform: uppercase;
}

.asset-icon,
.token-symbol {
  padding: 4px 8px;
  background: #dbeafe;
  color: #1e40af;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

.asset-count {
  font-size: 12px;
  color: #6b7280;
}

.eth-value {
  font-weight: 500;
  color: #1f2937;
  font-family: monospace;
}

.actions {
  display: flex;
  gap: 8px;
}

.icon-btn {
  padding: 6px;
  background: transparent;
  border: none;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.2s;
  font-size: 16px;
}

.icon-btn:hover {
  background: #f3f4f6;
}

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

.modal {
  background: white;
  border-radius: 12px;
  padding: 32px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
  font-family: inherit;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  border-color: #4f46e5;
}

.form-group textarea {
  min-height: 80px;
  resize: vertical;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
}

.no-tokens-row {
  background-color: #f9fafb;
}

.no-tokens-message {
  padding: 20px;
  text-align: center;
  color: #6b7280;
  font-size: 14px;
  font-style: italic;
}

.chains-selection {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 12px 0;
  max-height: 400px;
  overflow-y: auto;
}

.chain-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border: 1.5px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background-color: #fff;
  height: 56px;
}

.chain-item:hover {
  background-color: #f9fafb;
  border-color: #cbd5e1;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.06);
}

.chain-item:has(input:checked) {
  background-color: #eff6ff;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.chain-logo-inline {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
  margin-right: 12px;
}

.chain-logo-placeholder {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
  margin-right: 12px;
}

.chain-name {
  font-size: 15px;
  font-weight: 500;
  color: #1f2937;
  flex: 1;
  min-width: 0;
  padding: 0 8px;
}

.chain-checkbox {
  width: 20px;
  height: 20px;
  cursor: pointer;
  flex-shrink: 0;
  accent-color: #3b82f6;
  margin-right: 8px;
}

/* Pagination styles */
.pagination-row {
  background: #fafafa;
  border-top: 1px solid #e5e7eb;
}

.pagination-row td {
  padding: 12px 20px;
}

.pagination-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.pagination-info {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  color: #6b7280;
}

.pagination-label {
  font-weight: 500;
  color: #4b5563;
}

.page-size-select {
  padding: 4px 8px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  cursor: pointer;
  color: #1f2937;
}

.page-size-select:focus {
  outline: none;
  border-color: #6366f1;
}

.pagination-range {
  font-weight: 500;
  color: #1f2937;
}

.pagination-page {
  font-weight: 600;
  color: #4f46e5;
  padding: 4px 12px;
  background: #eef2ff;
  border-radius: 6px;
}

.pagination-buttons {
  display: flex;
  gap: 4px;
}

.pagination-btn {
  padding: 6px 12px;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  color: #6b7280;
  transition: all 0.2s;
  min-width: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pagination-btn:hover:not(:disabled) {
  background: #f3f4f6;
  border-color: #9ca3af;
  color: #1f2937;
}

.pagination-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.pagination-btn:active:not(:disabled) {
  background: #e5e7eb;
}
</style>
