<template>
  <tr class="bg-muted/5 border-b hover:bg-muted/10 transition-colors">
    <td colspan="6" class="p-0">
      <div
        class="grid grid-cols-[50px_1fr_150px_150px_150px_150px] gap-4 pl-[52px] pr-4 py-3 items-center"
      >
        <!-- Token Icon with Chain Badge -->
        <div class="relative w-10 h-10">
          <div
            class="w-10 h-10 rounded-full overflow-hidden bg-background border flex items-center justify-center"
          >
            <img
              v-if="token.logo_url && token.logo_url.length > 0 && !failedImage"
              :src="token.logo_url"
              :alt="token.symbol"
              class="w-full h-full object-cover"
              @error="handleImageError"
            />
            <span v-else class="text-sm font-bold text-primary uppercase">
              {{ token.symbol.substring(0, 2) }}
            </span>
          </div>
          <div
            class="absolute -bottom-0.5 -right-0.5 w-[18px] h-[18px] rounded-full bg-background border-2 border-background flex items-center justify-center overflow-hidden"
            :title="token.chain?.name || token.chain_id"
          >
            <img
              v-if="token.chain?.logo_url"
              :src="token.chain.logo_url"
              :alt="token.chain_id"
              class="w-full h-full object-cover"
            />
            <span v-else class="text-[8px]">{{ getChainIcon(token.chain_id) }}</span>
          </div>
        </div>

        <!-- Token Info -->
        <div class="flex flex-col gap-0.5">
          <div class="text-sm font-semibold">{{ token.symbol }}</div>
          <div class="text-xs text-muted-foreground">{{ token.name }}</div>
        </div>

        <!-- Location -->
        <div class="flex items-center">
          <Badge variant="secondary" class="text-xs">
            {{ token.chain?.name || token.chain_id }}
          </Badge>
        </div>

        <!-- Price -->
        <div class="text-right">
          <div class="text-sm font-medium">
            {{ currencySymbols[selectedCurrency] }}{{ formatTokenPrice(token.price) }}
          </div>
        </div>

        <!-- Balance -->
        <div class="text-right">
          <div class="text-sm font-medium">{{ parseFloat(token.balance).toFixed(4) }}</div>
        </div>

        <!-- Value -->
        <div class="text-right">
          <div class="text-sm font-medium">
            {{ currencySymbols[selectedCurrency] }}{{ formatTokenValue(token.usd_value) }}
          </div>
        </div>
      </div>
    </td>
  </tr>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useCurrency } from '../composables/useCurrency'
import { Badge } from '@/components/ui/badge'
import type { Token } from '@/types'

interface Props {
  token: Token
}

const props = defineProps<Props>()

const { currencySymbols, selectedCurrency, formatTokenPrice, formatTokenValue } = useCurrency()
const failedImage = ref(false)

const handleImageError = () => {
  failedImage.value = true
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
</script>
