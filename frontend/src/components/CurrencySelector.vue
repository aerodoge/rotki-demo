<template>
  <DropdownMenu v-model:open="showMenu">
    <DropdownMenuTrigger as-child>
      <Button variant="outline" class="gap-2">
        <span class="text-lg font-semibold text-primary">{{
          currencySymbols[selectedCurrency]
        }}</span>
        <span class="font-medium">{{ selectedCurrency }}</span>
        <span class="text-xs text-muted-foreground">▼</span>
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end" class="w-72">
      <DropdownMenuItem
        v-for="currency in currencies"
        :key="currency"
        class="flex items-center gap-4 py-4 cursor-pointer"
        :class="{ 'bg-primary/5 border-l-2 border-l-primary': selectedCurrency === currency }"
        @click="handleSelect(currency)"
      >
        <span class="text-3xl font-bold text-primary w-10 text-center">{{
          currencySymbols[currency]
        }}</span>
        <div class="flex-1">
          <div class="text-sm font-semibold">{{ currency }}</div>
          <div class="text-xs text-muted-foreground">Select as the main currency</div>
        </div>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useCurrency } from '../composables/useCurrency'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'

const { currencies, currencySymbols, selectedCurrency, selectCurrency } = useCurrency()
const showMenu = ref(false)

const handleSelect = (currency: string) => {
  selectCurrency(currency)
  showMenu.value = false
}
</script>
