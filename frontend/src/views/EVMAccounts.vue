<template>
  <div class="container mx-auto p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <h1 class="text-3xl font-bold tracking-tight">EVM</h1>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input
            type="checkbox"
            v-model="hideSmallBalances"
            class="w-4 h-4 rounded border-gray-300 cursor-pointer"
          />
          <span class="text-muted-foreground">隐藏小额 (&lt;10U)</span>
        </label>
      </div>
      <div class="flex items-center gap-3">
        <CurrencySelector />
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex gap-3 justify-end">
      <Button
        @click="showAddWalletModal = true"
        class="bg-gradient-primary hover:bg-gradient-primary-hover shadow-lg shadow-primary/25"
      >
        <span class="mr-2">+</span>
        Add Wallet
      </Button>
      <Button
        @click="showAddAddressModal = true"
        variant="outline"
        class="border-primary/20 hover:bg-gradient-accent hover:border-primary/40"
      >
        <span class="mr-2">+</span>
        Add Address
      </Button>
      <Button
        @click="refreshAllWallets"
        variant="outline"
        class="border-primary/20 hover:bg-accent"
        :disabled="loading"
      >
        <svg
          class="w-4 h-4 mr-2"
          :class="{ 'animate-spin': loading }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
        {{ loading ? 'Refreshing...' : 'Refresh All' }}
      </Button>
    </div>

    <!-- Accounts Card -->
    <Card>
      <CardContent class="p-0">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="border-b bg-muted/50">
              <tr>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                >
                  Account
                </th>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                >
                  Chains
                </th>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                >
                  Tags
                </th>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                >
                  Assets
                </th>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                >
                  {{ selectedCurrency }} value
                </th>
                <th
                  class="px-4 py-3 text-right text-xs font-medium text-muted-foreground uppercase tracking-wider"
                >
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              <template v-for="wallet in filteredWallets" :key="wallet.id">
                <!-- Wallet Row -->
                <tr
                  class="border-b hover:bg-accent/50 cursor-pointer transition-all group"
                  @click="toggleWallet(wallet.id)"
                >
                  <td class="px-4 py-4">
                    <div class="flex items-center gap-3">
                      <span
                        class="text-xs text-muted-foreground transform transition-transform group-hover:scale-110"
                        :class="{ 'rotate-90': expandedWallets[wallet.id] }"
                      >
                        ▶
                      </span>
                      <img
                        :src="getWalletAvatar(wallet.name, wallet.id)"
                        :alt="wallet.name"
                        class="w-8 h-8 rounded-lg shadow-sm group-hover:shadow-md transition-all"
                      />
                      <span class="font-semibold group-hover:text-primary transition-colors">{{
                        wallet.name
                      }}</span>
                    </div>
                  </td>
                  <td class="px-4 py-4">
                    <div class="flex items-center gap-1">
                      <template
                        v-for="(chainId, index) in getWalletUniqueChains(wallet.id).slice(0, 3)"
                        :key="chainId"
                      >
                        <img
                          v-if="getChainLogo(chainId) && !failedChainImages[chainId]"
                          :src="getChainLogo(chainId)"
                          :alt="chainId"
                          :title="getChainName(chainId)"
                          class="w-6 h-6 rounded-full object-cover border border-border"
                          @error="
                            () => {
                              console.log(
                                'Image load failed for chain:',
                                chainId,
                                getChainLogo(chainId)
                              )
                              failedChainImages[chainId] = true
                            }
                          "
                        />
                        <div
                          v-else
                          class="w-6 h-6 rounded-full bg-gradient-to-br from-purple-500 to-purple-700 text-white flex items-center justify-center text-[10px] font-semibold border border-purple-400"
                          :title="getChainName(chainId)"
                        >
                          {{ getChainName(chainId).charAt(0).toUpperCase() }}
                        </div>
                      </template>
                      <div
                        v-if="getWalletChainCount(wallet.id) > 3"
                        class="w-6 h-6 rounded-full bg-muted hover:bg-accent flex items-center justify-center text-[10px] font-semibold text-muted-foreground border border-border cursor-pointer transition-colors"
                        :title="`${getWalletChainCount(wallet.id) - 3} more chains`"
                      >
                        +{{ getWalletChainCount(wallet.id) - 3 }}
                      </div>
                    </div>
                  </td>
                  <td class="px-4 py-4">
                    <div class="flex gap-2 flex-wrap">
                      <Badge
                        v-for="tag in wallet.tags || []"
                        :key="tag"
                        variant="outline"
                        class="border-primary/30 text-primary/80"
                      >
                        {{ tag }}
                      </Badge>
                    </div>
                  </td>
                  <td class="px-4 py-4">
                    <div class="flex items-center gap-2">
                      <span class="text-sm text-muted-foreground"
                        >{{ getWalletAssetCount(wallet.id) }} assets</span
                      >
                    </div>
                  </td>
                  <td class="px-4 py-4">
                    <div class="font-mono font-semibold text-primary">
                      {{ currencySymbols[selectedCurrency]
                      }}{{ formatValue(getTotalValueByWallet(wallet.id)) }}
                    </div>
                  </td>
                  <td class="px-4 py-4">
                    <div class="flex items-center justify-end gap-1">
                      <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8"
                        @click.stop="refreshWallet(wallet.id)"
                        title="Refresh"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        >
                          <polyline points="23 4 23 10 17 10" />
                          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10" />
                        </svg>
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8"
                        @click.stop="editWallet(wallet)"
                        title="Edit"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        >
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                        </svg>
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8 hover:text-destructive"
                        @click.stop="deleteWallet(wallet.id)"
                        title="Delete"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        >
                          <polyline points="3 6 5 6 21 6" />
                          <path
                            d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                          />
                        </svg>
                      </Button>
                    </div>
                  </td>
                </tr>

                <!-- Expanded Wallet: Addresses -->
                <template v-if="expandedWallets[wallet.id]">
                  <!-- Address Header -->
                  <tr
                    v-if="getAddressesByWallet(wallet.id).length > 0"
                    class="bg-muted/30 border-b"
                  >
                    <th
                      class="px-4 py-2 text-left text-xs font-medium text-muted-foreground uppercase"
                    >
                      Account
                    </th>
                    <th
                      class="px-4 py-2 text-left text-xs font-medium text-muted-foreground uppercase"
                    >
                      Chains
                    </th>
                    <th
                      class="px-4 py-2 text-left text-xs font-medium text-muted-foreground uppercase"
                    >
                      Tags
                    </th>
                    <th
                      class="px-4 py-2 text-left text-xs font-medium text-muted-foreground uppercase"
                    >
                      Assets
                    </th>
                    <th
                      class="px-4 py-2 text-left text-xs font-medium text-muted-foreground uppercase"
                    >
                      {{ selectedCurrency }} value
                    </th>
                    <th
                      class="px-4 py-2 text-right text-xs font-medium text-muted-foreground uppercase"
                    >
                      Actions
                    </th>
                  </tr>

                  <!-- Address Rows -->
                  <template v-for="address in getAddressesByWallet(wallet.id)" :key="address.id">
                    <tr
                      class="bg-muted/30 border-b hover:bg-muted/50 cursor-pointer transition-all group"
                      @click="toggleAddress(address.id)"
                    >
                      <td class="px-4 py-3 pl-8">
                        <div class="flex items-center gap-2">
                          <span
                            class="text-xs text-muted-foreground transform transition-transform group-hover:scale-110"
                            :class="{ 'rotate-90': expandedAddresses[address.id] }"
                          >
                            ▶
                          </span>
                          <div
                            class="w-6 h-6 rounded-md bg-primary/10 flex items-center justify-center ring-1 ring-primary/20 group-hover:ring-primary/40 transition-all"
                          >
                            <span class="text-xs text-primary font-semibold">
                              {{ address.label?.[0]?.toUpperCase() || 'A' }}
                            </span>
                          </div>
                          <div class="flex flex-col">
                            <span
                              class="text-sm font-medium group-hover:text-primary transition-colors"
                              >{{ address.label || 'Unlabeled' }}</span
                            >
                            <span class="text-xs text-muted-foreground font-mono">{{
                              formatAddress(address.address)
                            }}</span>
                          </div>
                        </div>
                      </td>
                      <td class="px-4 py-3">
                        <div class="flex items-center gap-1">
                          <template
                            v-for="chainId in getAddressChains(address).slice(0, 3)"
                            :key="chainId"
                          >
                            <img
                              v-if="getChainLogo(chainId) && !failedChainImages[chainId]"
                              :src="getChainLogo(chainId)"
                              :alt="chainId"
                              :title="getChainName(chainId)"
                              class="w-5 h-5 rounded-full object-cover border border-border"
                              @error="failedChainImages[chainId] = true"
                            />
                            <div
                              v-else
                              class="w-5 h-5 rounded-full bg-gradient-to-br from-purple-500 to-purple-700 text-white flex items-center justify-center text-[8px] font-semibold border border-purple-400"
                              :title="getChainName(chainId)"
                            >
                              {{ getChainName(chainId).charAt(0).toUpperCase() }}
                            </div>
                          </template>
                          <div
                            v-if="getAddressChains(address).length > 3"
                            class="w-5 h-5 rounded-full bg-muted hover:bg-accent flex items-center justify-center text-[9px] font-semibold text-muted-foreground border border-border cursor-pointer transition-colors"
                            :title="`${getAddressChains(address).length - 3} more chains`"
                          >
                            +{{ getAddressChains(address).length - 3 }}
                          </div>
                        </div>
                      </td>
                      <td class="px-4 py-3">
                        <div class="flex gap-1 flex-wrap">
                          <Badge
                            v-for="tag in address.tags || []"
                            :key="tag"
                            variant="secondary"
                            class="text-xs bg-primary/5 text-primary/70 border-primary/20"
                          >
                            {{ tag }}
                          </Badge>
                        </div>
                      </td>
                      <td class="px-4 py-3">
                        <div class="flex items-center gap-1">
                          <div
                            v-for="token in (address.tokens || []).slice(0, 3)"
                            :key="token.id"
                            class="w-6 h-6"
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
                              class="w-6 h-6 rounded-full border object-cover"
                              @error="handleImageError($event, token.id)"
                            />
                            <div
                              v-else
                              class="w-6 h-6 rounded-full bg-primary/10 flex items-center justify-center text-[9px] font-bold border"
                              :title="token.symbol"
                            >
                              {{ token.symbol.substring(0, 2) }}
                            </div>
                          </div>
                          <span
                            v-if="(address.tokens || []).length > 3"
                            class="text-xs text-muted-foreground ml-1"
                          >
                            +{{ (address.tokens || []).length - 3 }}
                          </span>
                        </div>
                      </td>
                      <td class="px-4 py-3">
                        <div class="font-mono text-sm">
                          {{ currencySymbols[selectedCurrency]
                          }}{{ formatValue(getAddressValue(address)) }}
                        </div>
                      </td>
                      <td class="px-4 py-3">
                        <div class="flex items-center justify-end gap-1">
                          <Button
                            variant="ghost"
                            size="icon"
                            class="h-8 w-8"
                            @click.stop="refreshAddress(address.id)"
                            title="Refresh"
                          >
                            <svg
                              xmlns="http://www.w3.org/2000/svg"
                              class="h-4 w-4"
                              viewBox="0 0 24 24"
                              fill="none"
                              stroke="currentColor"
                              stroke-width="2"
                              stroke-linecap="round"
                              stroke-linejoin="round"
                            >
                              <polyline points="23 4 23 10 17 10" />
                              <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10" />
                            </svg>
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            class="h-8 w-8"
                            @click.stop="editAddress(address)"
                            title="Edit"
                          >
                            <svg
                              xmlns="http://www.w3.org/2000/svg"
                              class="h-4 w-4"
                              viewBox="0 0 24 24"
                              fill="none"
                              stroke="currentColor"
                              stroke-width="2"
                              stroke-linecap="round"
                              stroke-linejoin="round"
                            >
                              <path
                                d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
                              />
                              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                            </svg>
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            class="h-8 w-8 hover:text-destructive"
                            @click.stop="deleteAddress(address.id)"
                            title="Delete"
                          >
                            <svg
                              xmlns="http://www.w3.org/2000/svg"
                              class="h-4 w-4"
                              viewBox="0 0 24 24"
                              fill="none"
                              stroke="currentColor"
                              stroke-width="2"
                              stroke-linecap="round"
                              stroke-linejoin="round"
                            >
                              <polyline points="3 6 5 6 21 6" />
                              <path
                                d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                              />
                            </svg>
                          </Button>
                        </div>
                      </td>
                    </tr>

                    <!-- Expanded Address: Tokens -->
                    <template v-if="expandedAddresses[address.id]">
                      <!-- No tokens message -->
                      <tr v-if="!address.tokens || address.tokens.length === 0" class="bg-muted/10">
                        <td colspan="6" class="px-4 py-8">
                          <div class="text-center text-sm text-muted-foreground italic">
                            No tokens found for this address. Click the refresh button to sync data.
                          </div>
                        </td>
                      </tr>

                      <!-- View Mode Tabs -->
                      <tr v-if="address.tokens && address.tokens.length > 0" class="bg-muted/10">
                        <td colspan="6" class="p-0">
                          <Tabs
                            :model-value="addressViewMode[address.id] || 'aggregated'"
                            @update:model-value="(val) => setAddressViewMode(address.id, val)"
                            class="w-full"
                          >
                            <TabsList
                              class="w-full rounded-none border-b bg-transparent p-0 h-auto"
                            >
                              <TabsTrigger
                                value="aggregated"
                                class="rounded-none border-b-2 border-transparent data-[state=active]:border-primary"
                              >
                                Aggregated assets
                              </TabsTrigger>
                              <TabsTrigger
                                value="perChain"
                                class="rounded-none border-b-2 border-transparent data-[state=active]:border-primary"
                              >
                                Per chain
                              </TabsTrigger>
                            </TabsList>
                          </Tabs>
                        </td>
                      </tr>

                      <!-- Aggregated View -->
                      <template v-if="addressViewMode[address.id] === 'aggregated'">
                        <!-- Token Header -->
                        <tr
                          v-if="address.tokens && address.tokens.length > 0"
                          class="bg-muted/30 border-b"
                        >
                          <td colspan="6" class="p-0">
                            <div
                              class="grid grid-cols-[50px_1fr_150px_150px_150px_150px] gap-4 pl-[52px] pr-4 py-3 text-xs font-medium text-muted-foreground uppercase"
                            >
                              <div>Asset</div>
                              <div></div>
                              <div>Location</div>
                              <div class="text-right">Price in {{ selectedCurrency }}</div>
                              <div class="text-right">Amount</div>
                              <div class="text-right">{{ selectedCurrency }} Value</div>
                            </div>
                          </td>
                        </tr>
                        <TokenDetailRow
                          v-for="token in getPaginatedTokens(address.id)"
                          :key="token.id"
                          :token="token"
                        />

                        <!-- Pagination -->
                        <tr
                          v-if="address.tokens && address.tokens.length > 0"
                          class="bg-muted/5 border-t"
                        >
                          <td colspan="6" class="px-4 py-3">
                            <div class="flex items-center justify-between">
                              <div class="flex items-center gap-4 text-sm text-muted-foreground">
                                <span class="font-medium text-foreground">Rows per page:</span>
                                <Select
                                  :model-value="String(tokenPagination[address.id]?.pageSize || 10)"
                                  @update:model-value="
                                    (val) => {
                                      tokenPagination[address.id].pageSize = Number(val)
                                      onPageSizeChange(address.id)
                                    }
                                  "
                                >
                                  <SelectTrigger class="w-20 h-8">
                                    <SelectValue />
                                  </SelectTrigger>
                                  <SelectContent>
                                    <SelectItem value="10">10</SelectItem>
                                    <SelectItem value="20">20</SelectItem>
                                    <SelectItem value="50">50</SelectItem>
                                    <SelectItem value="100">100</SelectItem>
                                  </SelectContent>
                                </Select>
                                <span>
                                  Items {{ getStartIndex(address.id) }}-{{
                                    getEndIndex(address.id)
                                  }}
                                  of {{ address.tokens.length }}
                                </span>
                                <Badge variant="secondary">
                                  Page {{ tokenPagination[address.id]?.currentPage || 1 }} of
                                  {{ getTotalPages(address.id) }}
                                </Badge>
                              </div>
                              <div class="flex gap-1">
                                <Button
                                  variant="outline"
                                  size="icon"
                                  class="h-8 w-8"
                                  @click="goToFirstPage(address.id)"
                                  :disabled="isFirstPage(address.id)"
                                  title="First page"
                                >
                                  ⟨⟨
                                </Button>
                                <Button
                                  variant="outline"
                                  size="icon"
                                  class="h-8 w-8"
                                  @click="goToPreviousPage(address.id)"
                                  :disabled="isFirstPage(address.id)"
                                  title="Previous page"
                                >
                                  ⟨
                                </Button>
                                <Button
                                  variant="outline"
                                  size="icon"
                                  class="h-8 w-8"
                                  @click="goToNextPage(address.id)"
                                  :disabled="isLastPage(address.id)"
                                  title="Next page"
                                >
                                  ⟩
                                </Button>
                                <Button
                                  variant="outline"
                                  size="icon"
                                  class="h-8 w-8"
                                  @click="goToLastPage(address.id)"
                                  :disabled="isLastPage(address.id)"
                                  title="Last page"
                                >
                                  ⟩⟩
                                </Button>
                              </div>
                            </div>
                          </td>
                        </tr>
                      </template>

                      <!-- Per Chain View -->
                      <template v-if="addressViewMode[address.id] === 'perChain'">
                        <template
                          v-for="chain in getAddressChainGroups(address.id)"
                          :key="chain.chainId"
                        >
                          <!-- Chain Group Row -->
                          <tr
                            class="bg-muted/20 hover:bg-muted/30 cursor-pointer transition-colors"
                            @click="toggleAddressChain(address.id, chain.chainId)"
                          >
                            <td class="pl-12 pr-4 py-3">
                              <div class="flex items-center gap-3">
                                <span class="text-xs text-muted-foreground">
                                  {{
                                    expandedAddressChains[`${address.id}_${chain.chainId}`]
                                      ? '▼'
                                      : '▶'
                                  }}
                                </span>
                                <img
                                  v-if="chain.logoUrl && !failedChainImages[chain.chainId]"
                                  :src="chain.logoUrl"
                                  :alt="chain.name"
                                  class="w-7 h-7 rounded-full object-cover"
                                  @error="failedChainImages[chain.chainId] = true"
                                />
                                <div
                                  v-else
                                  class="w-7 h-7 rounded-full bg-gradient-to-br from-purple-500 to-purple-700 text-white flex items-center justify-center text-xs font-semibold"
                                >
                                  {{ (chain.name || chain.chainId || '?').charAt(0).toUpperCase() }}
                                </div>
                                <span class="font-medium">{{ chain.name || chain.chainId }}</span>
                              </div>
                            </td>
                            <td class="px-4 py-3"></td>
                            <td class="px-4 py-3"></td>
                            <td class="px-4 py-3">
                              <div class="flex items-center gap-1">
                                <div
                                  v-for="token in chain.tokens.slice(0, 2)"
                                  :key="token.id"
                                  class="w-6 h-6"
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
                                    class="w-6 h-6 rounded-full border object-cover"
                                    @error="handleImageError($event, token.id)"
                                  />
                                  <div
                                    v-else
                                    class="w-6 h-6 rounded-full bg-primary/10 flex items-center justify-center text-[9px] font-bold border"
                                    :title="token.symbol"
                                  >
                                    {{ token.symbol.substring(0, 2) }}
                                  </div>
                                </div>
                                <span
                                  v-if="chain.tokens.length > 2"
                                  class="text-xs text-muted-foreground"
                                >
                                  +{{ chain.tokens.length - 2 }}
                                </span>
                              </div>
                            </td>
                            <td class="px-4 py-3 text-right">
                              <span class="font-mono font-semibold">
                                {{ currencySymbols[selectedCurrency]
                                }}{{ formatValue(chain.totalValue) }}
                              </span>
                            </td>
                            <td class="px-4 py-3"></td>
                          </tr>

                          <!-- Expanded Chain Tokens -->
                          <template v-if="expandedAddressChains[`${address.id}_${chain.chainId}`]">
                            <!-- Chain Token Header -->
                            <tr class="bg-muted/30 border-b">
                              <td colspan="6" class="p-0">
                                <div
                                  class="grid grid-cols-[50px_1fr_150px_150px_150px_150px] gap-4 pl-[52px] pr-4 py-3 text-xs font-medium text-muted-foreground uppercase"
                                >
                                  <div>Asset</div>
                                  <div></div>
                                  <div>Location</div>
                                  <div class="text-right">Price in {{ selectedCurrency }}</div>
                                  <div class="text-right">Amount</div>
                                  <div class="text-right">{{ selectedCurrency }} Value</div>
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

                            <!-- Chain Pagination -->
                            <tr class="bg-muted/5 border-t">
                              <td colspan="6" class="px-4 py-3">
                                <div class="flex items-center justify-between">
                                  <div
                                    class="flex items-center gap-4 text-sm text-muted-foreground"
                                  >
                                    <span class="font-medium text-foreground">Rows per page:</span>
                                    <Select
                                      :model-value="
                                        String(
                                          addressChainPagination[`${address.id}_${chain.chainId}`]
                                            ?.pageSize || 10
                                        )
                                      "
                                      @update:model-value="
                                        (val) => {
                                          addressChainPagination[
                                            `${address.id}_${chain.chainId}`
                                          ].pageSize = Number(val)
                                          onAddressChainPageSizeChange(address.id, chain.chainId)
                                        }
                                      "
                                    >
                                      <SelectTrigger class="w-20 h-8">
                                        <SelectValue />
                                      </SelectTrigger>
                                      <SelectContent>
                                        <SelectItem value="10">10</SelectItem>
                                        <SelectItem value="20">20</SelectItem>
                                        <SelectItem value="50">50</SelectItem>
                                        <SelectItem value="100">100</SelectItem>
                                      </SelectContent>
                                    </Select>
                                    <span>
                                      Items
                                      {{ getAddressChainStartIndex(address.id, chain.chainId) }}-{{
                                        getAddressChainEndIndex(address.id, chain.chainId)
                                      }}
                                      of {{ chain.tokens.length }}
                                    </span>
                                    <Badge variant="secondary">
                                      Page
                                      {{
                                        addressChainPagination[`${address.id}_${chain.chainId}`]
                                          ?.currentPage || 1
                                      }}
                                      of {{ getAddressChainTotalPages(address.id, chain.chainId) }}
                                    </Badge>
                                  </div>
                                  <div class="flex gap-1">
                                    <Button
                                      variant="outline"
                                      size="icon"
                                      class="h-8 w-8"
                                      @click="goToAddressChainFirstPage(address.id, chain.chainId)"
                                      :disabled="isAddressChainFirstPage(address.id, chain.chainId)"
                                      title="First page"
                                    >
                                      ⟨⟨
                                    </Button>
                                    <Button
                                      variant="outline"
                                      size="icon"
                                      class="h-8 w-8"
                                      @click="
                                        goToAddressChainPreviousPage(address.id, chain.chainId)
                                      "
                                      :disabled="isAddressChainFirstPage(address.id, chain.chainId)"
                                      title="Previous page"
                                    >
                                      ⟨
                                    </Button>
                                    <Button
                                      variant="outline"
                                      size="icon"
                                      class="h-8 w-8"
                                      @click="goToAddressChainNextPage(address.id, chain.chainId)"
                                      :disabled="isAddressChainLastPage(address.id, chain.chainId)"
                                      title="Next page"
                                    >
                                      ⟩
                                    </Button>
                                    <Button
                                      variant="outline"
                                      size="icon"
                                      class="h-8 w-8"
                                      @click="goToAddressChainLastPage(address.id, chain.chainId)"
                                      :disabled="isAddressChainLastPage(address.id, chain.chainId)"
                                      title="Last page"
                                    >
                                      ⟩⟩
                                    </Button>
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

              <!-- Total Row -->
              <tr class="bg-muted font-semibold border-t-2">
                <td class="px-4 py-4">Total</td>
                <td></td>
                <td></td>
                <td></td>
                <td class="px-4 py-4 font-mono">
                  {{ currencySymbols[selectedCurrency] }}{{ formatValue(getTotalValue()) }}
                </td>
                <td></td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>

    <!-- Add Wallet Dialog -->
    <Dialog :open="showAddWalletModal" @update:open="showAddWalletModal = $event">
      <DialogContent class="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Add Wallet</DialogTitle>
          <DialogDescription>Create a new wallet to organize your addresses</DialogDescription>
        </DialogHeader>
        <form @submit.prevent="handleAddWallet" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="wallet-name">Wallet Name *</Label>
            <Input id="wallet-name" v-model="newWallet.name" required placeholder="My Wallet" />
          </div>
          <div class="space-y-2">
            <Label for="wallet-description">Description (optional)</Label>
            <Input
              id="wallet-description"
              v-model="newWallet.description"
              placeholder="Wallet description"
            />
          </div>
          <div class="space-y-2">
            <Label>Tags (optional)</Label>
            <div class="border rounded-lg p-3 space-y-2">
              <div class="flex flex-wrap gap-2 min-h-[24px]">
                <Badge
                  v-for="(tag, index) in newWallet.tags"
                  :key="index"
                  variant="secondary"
                  class="gap-1"
                >
                  {{ tag }}
                  <button
                    type="button"
                    @click="removeNewWalletTag(index)"
                    class="ml-1 hover:text-destructive"
                  >
                    ×
                  </button>
                </Badge>
              </div>
              <div class="flex gap-2">
                <Input
                  v-model="newWalletTagInput"
                  placeholder="Add a tag"
                  @keyup.enter="addNewWalletTag"
                  class="flex-1"
                />
                <Button type="button" @click="addNewWalletTag" size="sm">Add</Button>
              </div>
            </div>
          </div>
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <Label>Enabled Chains (optional)</Label>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                class="text-xs h-7"
                @click="toggleAllNewWalletChains"
              >
                {{ isAllNewWalletChainsSelected ? 'Deselect All' : 'Select All' }}
              </Button>
            </div>
            <div class="grid grid-cols-2 gap-3 max-h-[300px] overflow-y-auto p-1">
              <div
                v-for="chain in availableChains"
                :key="chain.id"
                class="flex items-center gap-2 p-3 border rounded-lg hover:bg-accent cursor-pointer transition-colors"
                :class="{
                  'bg-primary/5 border-primary': newWallet.enabled_chains.includes(chain.id)
                }"
              >
                <input
                  type="checkbox"
                  :id="`chain-${chain.id}`"
                  :value="chain.id"
                  v-model="newWallet.enabled_chains"
                  class="cursor-pointer"
                />
                <Label
                  :for="`chain-${chain.id}`"
                  class="flex items-center gap-2 cursor-pointer flex-1"
                >
                  <img
                    v-if="!failedChainImages[chain.id]"
                    :src="`/images/chains/${chain.id}.png`"
                    :alt="chain.name"
                    class="w-6 h-6 rounded-full object-cover"
                    @error="failedChainImages[chain.id] = true"
                  />
                  <div
                    v-else
                    class="w-6 h-6 rounded-full bg-gradient-to-br from-purple-500 to-purple-700 text-white flex items-center justify-center text-xs font-semibold"
                  >
                    {{ (chain.name || chain.id || '?').charAt(0).toUpperCase() }}
                  </div>
                  <span class="text-sm">{{ chain.name || chain.id || 'Unknown' }}</span>
                </Label>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" @click="showAddWalletModal = false"
              >Cancel</Button
            >
            <Button type="submit">Add Wallet</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- Edit Wallet Dialog -->
    <Dialog :open="showEditWalletModal" @update:open="showEditWalletModal = $event">
      <DialogContent class="sm:max-w-[500px]" v-if="editingWallet">
        <DialogHeader>
          <DialogTitle>Edit Wallet</DialogTitle>
          <DialogDescription>Update wallet information</DialogDescription>
        </DialogHeader>
        <form @submit.prevent="handleEditWallet" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="edit-wallet-name">Wallet Name *</Label>
            <Input
              id="edit-wallet-name"
              v-model="editingWallet.name"
              required
              placeholder="My Wallet"
            />
          </div>
          <div class="space-y-2">
            <Label for="edit-wallet-description">Description (optional)</Label>
            <Input
              id="edit-wallet-description"
              v-model="editingWallet.description"
              placeholder="Wallet description"
            />
          </div>
          <div class="space-y-2">
            <Label>Tags (optional)</Label>
            <div class="border rounded-lg p-3 space-y-2">
              <div class="flex flex-wrap gap-2 min-h-[24px]">
                <Badge
                  v-for="(tag, index) in editingWallet.tags"
                  :key="index"
                  variant="secondary"
                  class="gap-1"
                >
                  {{ tag }}
                  <button
                    type="button"
                    @click="removeEditWalletTag(index)"
                    class="ml-1 hover:text-destructive"
                  >
                    ×
                  </button>
                </Badge>
              </div>
              <div class="flex gap-2">
                <Input
                  v-model="editWalletTagInput"
                  placeholder="Add a tag"
                  @keyup.enter="addEditWalletTag"
                  class="flex-1"
                />
                <Button type="button" @click="addEditWalletTag" size="sm">Add</Button>
              </div>
            </div>
          </div>
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <Label>Enabled Chains (optional)</Label>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                class="text-xs h-7"
                @click="toggleAllEditWalletChains"
              >
                {{ isAllEditWalletChainsSelected ? 'Deselect All' : 'Select All' }}
              </Button>
            </div>
            <div class="grid grid-cols-2 gap-3 max-h-[300px] overflow-y-auto p-1">
              <div
                v-for="chain in availableChains"
                :key="chain.id"
                class="flex items-center gap-2 p-3 border rounded-lg hover:bg-accent cursor-pointer transition-colors"
                :class="{
                  'bg-primary/5 border-primary': editingWallet.enabled_chains.includes(chain.id)
                }"
              >
                <input
                  type="checkbox"
                  :id="`edit-chain-${chain.id}`"
                  :value="chain.id"
                  v-model="editingWallet.enabled_chains"
                  class="cursor-pointer"
                />
                <Label
                  :for="`edit-chain-${chain.id}`"
                  class="flex items-center gap-2 cursor-pointer flex-1"
                >
                  <img
                    v-if="!failedChainImages[chain.id]"
                    :src="`/images/chains/${chain.id}.png`"
                    :alt="chain.name"
                    class="w-6 h-6 rounded-full object-cover"
                    @error="failedChainImages[chain.id] = true"
                  />
                  <div
                    v-else
                    class="w-6 h-6 rounded-full bg-gradient-to-br from-purple-500 to-purple-700 text-white flex items-center justify-center text-xs font-semibold"
                  >
                    {{ (chain.name || chain.id || '?').charAt(0).toUpperCase() }}
                  </div>
                  <span class="text-sm">{{ chain.name || chain.id || 'Unknown' }}</span>
                </Label>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" @click="showEditWalletModal = false"
              >Cancel</Button
            >
            <Button type="submit">Save Changes</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- Add Address Dialog -->
    <Dialog :open="showAddAddressModal" @update:open="showAddAddressModal = $event">
      <DialogContent class="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Add Address</DialogTitle>
          <DialogDescription>Add a new blockchain address to a wallet</DialogDescription>
        </DialogHeader>
        <form @submit.prevent="handleAddAddress" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="address-wallet">Wallet *</Label>
            <Select v-model="newAddress.wallet_id" required>
              <SelectTrigger id="address-wallet">
                <SelectValue placeholder="Select a wallet" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="wallet in wallets" :key="wallet.id" :value="String(wallet.id)">
                  {{ wallet.name }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-2">
            <Label for="address-value">Address *</Label>
            <Input
              id="address-value"
              v-model="newAddress.address"
              required
              placeholder="0x..."
              class="font-mono"
            />
          </div>
          <div class="space-y-2">
            <Label for="address-label">Label (optional)</Label>
            <Input id="address-label" v-model="newAddress.label" placeholder="Main Address" />
          </div>
          <div class="space-y-2">
            <Label>Tags (optional)</Label>
            <div class="border rounded-lg p-3 space-y-2">
              <div class="flex flex-wrap gap-2 min-h-[24px]">
                <Badge
                  v-for="(tag, index) in newAddress.tags"
                  :key="index"
                  variant="secondary"
                  class="gap-1"
                >
                  {{ tag }}
                  <button
                    type="button"
                    @click="removeNewAddressTag(index)"
                    class="ml-1 hover:text-destructive"
                  >
                    ×
                  </button>
                </Badge>
              </div>
              <div class="flex gap-2">
                <Input
                  v-model="newAddressTagInput"
                  placeholder="Add a tag"
                  @keyup.enter="addNewAddressTag"
                  class="flex-1"
                />
                <Button type="button" @click="addNewAddressTag" size="sm">Add</Button>
              </div>
            </div>
          </div>
          <div class="space-y-2">
            <Label for="chain-type">Chain Type</Label>
            <Select v-model="newAddress.chain_type">
              <SelectTrigger id="chain-type">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="EVM">EVM</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" @click="showAddAddressModal = false"
              >Cancel</Button
            >
            <Button type="submit">Add Address</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- Edit Address Dialog -->
    <Dialog :open="showEditAddressModal" @update:open="showEditAddressModal = $event">
      <DialogContent class="sm:max-w-[500px]" v-if="editingAddress">
        <DialogHeader>
          <DialogTitle>Edit Address</DialogTitle>
          <DialogDescription>Update address information</DialogDescription>
        </DialogHeader>
        <form @submit.prevent="handleEditAddress" class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="edit-address-value">Address</Label>
            <Input
              id="edit-address-value"
              v-model="editingAddress.address"
              disabled
              placeholder="0x..."
              class="font-mono"
            />
          </div>
          <div class="space-y-2">
            <Label for="edit-address-label">Label</Label>
            <Input
              id="edit-address-label"
              v-model="editingAddress.label"
              placeholder="Address Label"
            />
          </div>
          <div class="space-y-2">
            <Label>Tags (optional)</Label>
            <div class="border rounded-lg p-3 space-y-2">
              <div class="flex flex-wrap gap-2 min-h-[24px]">
                <Badge
                  v-for="(tag, index) in editingAddress.tags"
                  :key="index"
                  variant="secondary"
                  class="gap-1"
                >
                  {{ tag }}
                  <button
                    type="button"
                    @click="removeEditAddressTag(index)"
                    class="ml-1 hover:text-destructive"
                  >
                    ×
                  </button>
                </Badge>
              </div>
              <div class="flex gap-2">
                <Input
                  v-model="editAddressTagInput"
                  placeholder="Add a tag"
                  @keyup.enter="addEditAddressTag"
                  class="flex-1"
                />
                <Button type="button" @click="addEditAddressTag" size="sm">Add</Button>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" @click="showEditAddressModal = false"
              >Cancel</Button
            >
            <Button type="submit">Save Changes</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useWalletStore } from '../stores/wallet'
import { storeToRefs } from 'pinia'
import { chainsAPI } from '../api/client'
import { useCurrency } from '../composables/useCurrency'
import type { Wallet, Address, Chain } from '../types'
import CurrencySelector from '../components/CurrencySelector.vue'
import TokenDetailRow from '../components/TokenDetailRow.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
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
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'

const walletStore = useWalletStore()
const { wallets, addresses, loading } = storeToRefs(walletStore)

const availableChains = ref<Chain[]>([])

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

const expandedWallets = reactive<Record<number, boolean>>({})
const expandedAddresses = reactive<Record<number, boolean>>({})
const addressViewMode = reactive<Record<number, string>>({})
const expandedAddressChains = reactive<Record<string, boolean>>({})
const failedImages = reactive<Record<number, boolean>>({})
const failedChainImages = reactive<Record<string, boolean>>({})
const tokenPagination = reactive<Record<number, { currentPage: number; pageSize: number }>>({})
const addressChainPagination = reactive<Record<string, { currentPage: number; pageSize: number }>>(
  {}
)

// 隐藏小额资产
const hideSmallBalances = ref(false)

// 过滤小额钱包
const filteredWallets = computed(() => {
  if (!hideSmallBalances.value) {
    return wallets.value
  }
  return wallets.value.filter(wallet => {
    // 计算钱包的USD总值
    const addrs = getAddressesByWallet(wallet.id)
    const totalUsdValue = addrs.reduce((sum, addr) => sum + getAddressValue(addr), 0)
    return totalUsdValue >= 10
  })
})

const showAddWalletModal = ref(false)
const showEditWalletModal = ref(false)
const showAddAddressModal = ref(false)
const showEditAddressModal = ref(false)
const editingWallet = ref<any>(null)
const editingAddress = ref<any>(null)

const newWallet = ref({
  name: '',
  description: '',
  tags: [] as string[],
  enabled_chains: [] as string[]
})

const newAddress = ref({
  wallet_id: '',
  address: '',
  label: '',
  tags: [] as string[],
  chain_type: 'EVM'
})

const newWalletTagInput = ref('')
const editWalletTagInput = ref('')
const newAddressTagInput = ref('')
const editAddressTagInput = ref('')

onMounted(async () => {
  restoreSavedCurrency()
  await walletStore.fetchWallets()
  await walletStore.fetchAddresses()

  try {
    const response = await chainsAPI.list()
    availableChains.value = response.data
  } catch (error) {
    console.error('Failed to fetch chains:', error)
  }

  await updateExchangeRates(addresses.value)
})

const toggleWallet = (walletId: number) => {
  expandedWallets[walletId] = !expandedWallets[walletId]
}

const toggleAddress = (addressId: number) => {
  expandedAddresses[addressId] = !expandedAddresses[addressId]
  if (!addressViewMode[addressId]) {
    addressViewMode[addressId] = 'aggregated'
  }
  if (!tokenPagination[addressId]) {
    tokenPagination[addressId] = { currentPage: 1, pageSize: 10 }
  }
}

const setAddressViewMode = (addressId: number, mode: string) => {
  addressViewMode[addressId] = mode
}

const initPagination = (addressId: number) => {
  if (!tokenPagination[addressId]) {
    tokenPagination[addressId] = { currentPage: 1, pageSize: 10 }
  }
}

const getPaginatedTokens = (addressId: number) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens) return []

  initPagination(addressId)
  const pagination = tokenPagination[addressId]

  const sortedTokens = [...address.tokens].sort((a, b) => {
    const usdValueA = a.usd_value || 0
    const usdValueB = b.usd_value || 0
    const rate = exchangeRates.value[selectedCurrency.value] || 1
    const valueA = rate === 0 ? usdValueA : usdValueA / rate
    const valueB = rate === 0 ? usdValueB : usdValueB / rate
    return valueB - valueA
  })

  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return sortedTokens.slice(start, end)
}

const getTotalPages = (addressId: number) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens) return 1
  initPagination(addressId)
  const pagination = tokenPagination[addressId]
  return Math.ceil(address.tokens.length / pagination.pageSize)
}

const getStartIndex = (addressId: number) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens || address.tokens.length === 0) return 0
  initPagination(addressId)
  const pagination = tokenPagination[addressId]
  return (pagination.currentPage - 1) * pagination.pageSize + 1
}

const getEndIndex = (addressId: number) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens) return 0
  initPagination(addressId)
  const pagination = tokenPagination[addressId]
  const end = pagination.currentPage * pagination.pageSize
  return Math.min(end, address.tokens.length)
}

const isFirstPage = (addressId: number) => {
  initPagination(addressId)
  return tokenPagination[addressId].currentPage === 1
}

const isLastPage = (addressId: number) => {
  initPagination(addressId)
  return tokenPagination[addressId].currentPage >= getTotalPages(addressId)
}

const goToFirstPage = (addressId: number) => {
  initPagination(addressId)
  tokenPagination[addressId].currentPage = 1
}

const goToPreviousPage = (addressId: number) => {
  initPagination(addressId)
  if (tokenPagination[addressId].currentPage > 1) {
    tokenPagination[addressId].currentPage--
  }
}

const goToNextPage = (addressId: number) => {
  initPagination(addressId)
  if (tokenPagination[addressId].currentPage < getTotalPages(addressId)) {
    tokenPagination[addressId].currentPage++
  }
}

const goToLastPage = (addressId: number) => {
  initPagination(addressId)
  tokenPagination[addressId].currentPage = getTotalPages(addressId)
}

const onPageSizeChange = (addressId: number) => {
  initPagination(addressId)
  tokenPagination[addressId].currentPage = 1
}

const getAddressChainGroups = (addressId: number) => {
  const address = addresses.value.find((a) => a.id === addressId)
  if (!address || !address.tokens) return []

  const chainMap = new Map()

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

  const chains = Array.from(chainMap.values()).map((chain: any) => {
    const sortedTokens = [...chain.tokens].sort((a: any, b: any) => {
      const usdValueA = a.usd_value || 0
      const usdValueB = b.usd_value || 0
      const rate = exchangeRates.value[selectedCurrency.value] || 1
      const valueA = rate === 0 ? usdValueA : usdValueA / rate
      const valueB = rate === 0 ? usdValueB : usdValueB / rate
      return valueB - valueA
    })

    chain.tokens = sortedTokens
    chain.totalValue = chain.tokens.reduce(
      (sum: number, token: any) => sum + (token.usd_value || 0),
      0
    )
    return chain
  })

  return chains.sort((a, b) => b.totalValue - a.totalValue)
}

const toggleAddressChain = (addressId: number, chainId: string) => {
  const key = `${addressId}_${chainId}`
  expandedAddressChains[key] = !expandedAddressChains[key]
  if (!addressChainPagination[key]) {
    addressChainPagination[key] = { currentPage: 1, pageSize: 10 }
  }
}

const initAddressChainPagination = (addressId: number, chainId: string) => {
  const key = `${addressId}_${chainId}`
  if (!addressChainPagination[key]) {
    addressChainPagination[key] = { currentPage: 1, pageSize: 10 }
  }
}

const getPaginatedAddressChainTokens = (addressId: number, chainId: string) => {
  const chains = getAddressChainGroups(addressId)
  const chain = chains.find((c: any) => c.chainId === chainId)
  if (!chain || !chain.tokens) return []

  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  const pagination = addressChainPagination[key]
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return chain.tokens.slice(start, end)
}

const getAddressChainTotalPages = (addressId: number, chainId: string) => {
  const chains = getAddressChainGroups(addressId)
  const chain = chains.find((c: any) => c.chainId === chainId)
  if (!chain || !chain.tokens) return 1
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  const pagination = addressChainPagination[key]
  return Math.ceil(chain.tokens.length / pagination.pageSize)
}

const getAddressChainStartIndex = (addressId: number, chainId: string) => {
  const chains = getAddressChainGroups(addressId)
  const chain = chains.find((c: any) => c.chainId === chainId)
  if (!chain || !chain.tokens || chain.tokens.length === 0) return 0
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  const pagination = addressChainPagination[key]
  return (pagination.currentPage - 1) * pagination.pageSize + 1
}

const getAddressChainEndIndex = (addressId: number, chainId: string) => {
  const chains = getAddressChainGroups(addressId)
  const chain = chains.find((c: any) => c.chainId === chainId)
  if (!chain || !chain.tokens) return 0
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  const pagination = addressChainPagination[key]
  const end = pagination.currentPage * pagination.pageSize
  return Math.min(end, chain.tokens.length)
}

const isAddressChainFirstPage = (addressId: number, chainId: string) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  return addressChainPagination[key].currentPage === 1
}

const isAddressChainLastPage = (addressId: number, chainId: string) => {
  initAddressChainPagination(addressId, chainId)
  return (
    addressChainPagination[`${addressId}_${chainId}`].currentPage >=
    getAddressChainTotalPages(addressId, chainId)
  )
}

const goToAddressChainFirstPage = (addressId: number, chainId: string) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  addressChainPagination[key].currentPage = 1
}

const goToAddressChainPreviousPage = (addressId: number, chainId: string) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  if (addressChainPagination[key].currentPage > 1) {
    addressChainPagination[key].currentPage--
  }
}

const goToAddressChainNextPage = (addressId: number, chainId: string) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  if (addressChainPagination[key].currentPage < getAddressChainTotalPages(addressId, chainId)) {
    addressChainPagination[key].currentPage++
  }
}

const goToAddressChainLastPage = (addressId: number, chainId: string) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  addressChainPagination[key].currentPage = getAddressChainTotalPages(addressId, chainId)
}

const onAddressChainPageSizeChange = (addressId: number, chainId: string) => {
  initAddressChainPagination(addressId, chainId)
  const key = `${addressId}_${chainId}`
  addressChainPagination[key].currentPage = 1
}

const getAddressesByWallet = (walletId: number) => {
  return addresses.value.filter((addr) => addr.wallet_id === walletId)
}

const getTotalValueByWallet = (walletId: number) => {
  return walletStore.getTotalValueByWallet(walletId)
}

const getTotalValue = () => {
  return walletStore.getTotalValue
}

const getAddressValue = (address: Address) => {
  // 只计算钱包代币（不属于任何协议的代币）的价值
  const walletTokenValue =
    address.tokens?.reduce((sum, token) => {
      // 只累加不属于协议的代币
      if (!token.protocol_id) {
        return sum + (token.usd_value || 0)
      }
      return sum
    }, 0) || 0

  // 协议净值已经包含了协议代币的价值
  const protocolValue =
    address.protocols?.reduce((sum, protocol) => sum + (protocol.net_usd_value || 0), 0) || 0

  return walletTokenValue + protocolValue
}

const getWalletChainCount = (walletId: number) => {
  const wallet = wallets.value.find((w) => w.id === walletId)

  // If wallet has enabled_chains, use those
  if (wallet?.enabled_chains && wallet.enabled_chains.length > 0) {
    return wallet.enabled_chains.length
  }

  // Otherwise fall back to chains with tokens
  const addrs = getAddressesByWallet(walletId)
  const chains = new Set()
  addrs.forEach((addr) => {
    addr.tokens?.forEach((token) => chains.add(token.chain_id))
  })
  return chains.size
}

const getWalletUniqueChains = (walletId: number) => {
  const wallet = wallets.value.find((w) => w.id === walletId)

  // If wallet has enabled_chains, use those
  if (wallet?.enabled_chains && wallet.enabled_chains.length > 0) {
    return wallet.enabled_chains
  }

  // Otherwise fall back to chains with tokens
  const addrs = getAddressesByWallet(walletId)
  const chains = new Set<string>()
  addrs.forEach((addr) => {
    addr.tokens?.forEach((token) => chains.add(token.chain_id))
  })
  return Array.from(chains)
}

const getWalletAssetCount = (walletId: number) => {
  const addrs = getAddressesByWallet(walletId)
  return addrs.reduce((sum, addr) => sum + (addr.tokens?.length || 0), 0)
}

const getChainLogo = (chainId: string) => {
  // Always use local logo files
  const localLogoPath = `/images/chains/${chainId}.png`
  return localLogoPath
}

const getChainName = (chainId: string) => {
  const chain = availableChains.value.find((c) => c.id === chainId)
  return chain?.name || chainId
}

const getAddressChains = (address: Address) => {
  const chains = new Set()
  address.tokens?.forEach((token) => chains.add(token.chain_id))
  return Array.from(chains)
}

// Generate avatar for wallet using DiceBear API (GitHub identicon style)
const getWalletAvatar = (name: string, id: number) => {
  // Use identicon style for GitHub-like geometric pattern avatars
  const seed = `${name}-${id}`
  return `https://api.dicebear.com/7.x/identicon/svg?seed=${encodeURIComponent(seed)}&size=32`
}

const getChainIcon = (chainId: string) => {
  const icons: Record<string, string> = {
    eth: '⟠',
    bsc: 'B',
    polygon: 'P',
    arbitrum: 'A',
    optimism: 'O'
  }
  return icons[chainId] || '⛓️'
}

const formatAddress = (address: string) => {
  return `${address.slice(0, 6)}...${address.slice(-4)}`
}

const handleAddWallet = async () => {
  try {
    await walletStore.createWallet(newWallet.value)
    showAddWalletModal.value = false
    newWallet.value = { name: '', description: '', tags: [], enabled_chains: [] }
    newWalletTagInput.value = ''
  } catch (error: any) {
    alert('Failed to add wallet: ' + error.message)
  }
}

const editWallet = (wallet: Wallet) => {
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

    // Reload chains to get any newly enabled chains
    try {
      const response = await chainsAPI.list()
      availableChains.value = response.data
    } catch (error) {
      console.error('Failed to refresh chains:', error)
    }

    await walletStore.fetchWallets()
    await walletStore.fetchAddresses()
  } catch (error: any) {
    alert('Failed to update wallet: ' + error.message)
  }
}

const handleAddAddress = async () => {
  try {
    await walletStore.createAddress(newAddress.value)
    showAddAddressModal.value = false
    newAddress.value = { wallet_id: '', address: '', label: '', tags: [], chain_type: 'EVM' }
    newAddressTagInput.value = ''
    alert('Address added! Syncing data in background...')
    setTimeout(async () => {
      await walletStore.fetchAddresses()
    }, 2000)
  } catch (error: any) {
    alert('Failed to add address: ' + error.message)
  }
}

const handleImageError = (event: Event, tokenId: number) => {
  failedImages[tokenId] = true
}

const refreshWallet = async (walletId: number) => {
  try {
    await walletStore.refreshWallet(walletId)
  } catch (error: any) {
    alert('Failed to refresh wallet: ' + error.message)
  }
}

const refreshAllWallets = async () => {
  if (!confirm('确定要刷新所有钱包吗？这可能需要一些时间。')) {
    return
  }

  try {
    const totalWallets = wallets.value.length
    let completedWallets = 0

    // 串行刷新每个钱包，避免同时发送太多请求
    for (const wallet of wallets.value) {
      try {
        await walletStore.refreshWallet(wallet.id)
        completedWallets++
        console.log(`已刷新 ${completedWallets}/${totalWallets} 个钱包`)
      } catch (error: any) {
        console.error(`刷新钱包 ${wallet.name} 失败:`, error)
        // 继续刷新下一个钱包
      }
    }

    alert(`刷新完成！成功刷新 ${completedWallets}/${totalWallets} 个钱包`)
  } catch (error: any) {
    alert('刷新失败: ' + error.message)
  }
}

const refreshAddress = async (addressId: number) => {
  try {
    await walletStore.refreshAddress(addressId)
  } catch (error: any) {
    alert('Failed to refresh address: ' + error.message)
  }
}

const deleteWallet = async (walletId: number) => {
  if (confirm('Are you sure you want to delete this wallet?')) {
    try {
      await walletStore.deleteWallet(walletId)
    } catch (error: any) {
      alert('Failed to delete wallet: ' + error.message)
    }
  }
}

const editAddress = (address: Address) => {
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
  } catch (error: any) {
    alert('Failed to update address: ' + error.message)
  }
}

const deleteAddress = async (addressId: number) => {
  if (confirm('Are you sure you want to delete this address?')) {
    try {
      await walletStore.deleteAddress(addressId)
    } catch (error: any) {
      alert('Failed to delete address: ' + error.message)
    }
  }
}

const addNewWalletTag = () => {
  const tag = newWalletTagInput.value.trim()
  if (tag && !newWallet.value.tags.includes(tag)) {
    newWallet.value.tags.push(tag)
    newWalletTagInput.value = ''
  }
}

const removeNewWalletTag = (index: number) => {
  newWallet.value.tags.splice(index, 1)
}

const addEditWalletTag = () => {
  const tag = editWalletTagInput.value.trim()
  if (tag && !editingWallet.value.tags.includes(tag)) {
    editingWallet.value.tags.push(tag)
    editWalletTagInput.value = ''
  }
}

const removeEditWalletTag = (index: number) => {
  editingWallet.value.tags.splice(index, 1)
}

const addNewAddressTag = () => {
  const tag = newAddressTagInput.value.trim()
  if (tag && !newAddress.value.tags.includes(tag)) {
    newAddress.value.tags.push(tag)
    newAddressTagInput.value = ''
  }
}

const removeNewAddressTag = (index: number) => {
  newAddress.value.tags.splice(index, 1)
}

const addEditAddressTag = () => {
  const tag = editAddressTagInput.value.trim()
  if (tag && !editingAddress.value.tags.includes(tag)) {
    editingAddress.value.tags.push(tag)
    editAddressTagInput.value = ''
  }
}

const removeEditAddressTag = (index: number) => {
  editingAddress.value.tags.splice(index, 1)
}

// Computed properties for Select All functionality
const isAllNewWalletChainsSelected = computed(() => {
  return (
    availableChains.value.length > 0 &&
    newWallet.value.enabled_chains.length === availableChains.value.length
  )
})

const isAllEditWalletChainsSelected = computed(() => {
  return (
    availableChains.value.length > 0 &&
    editingWallet.value &&
    editingWallet.value.enabled_chains.length === availableChains.value.length
  )
})

// Methods to toggle all chains
const toggleAllNewWalletChains = () => {
  if (isAllNewWalletChainsSelected.value) {
    // Deselect all
    newWallet.value.enabled_chains = []
  } else {
    // Select all
    newWallet.value.enabled_chains = availableChains.value.map((chain) => chain.id)
  }
}

const toggleAllEditWalletChains = () => {
  if (isAllEditWalletChainsSelected.value) {
    // Deselect all
    editingWallet.value.enabled_chains = []
  } else {
    // Select all
    editingWallet.value.enabled_chains = availableChains.value.map((chain) => chain.id)
  }
}
</script>
