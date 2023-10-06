<template>
  <div class="un-flex un-shadow-md">
    <img v-if="book.imageUrl" :src="book.imageUrl" class="un-w-56 un-h-80 un-object-cover" />
    <div v-else class="un-flex un-flex-col un-w-56 un-h-80 un-bg-blue-300 un-p-4">
      <div class="un-border un-border-gray-700 un-w-11/12 un-mx-auto un-text-center">
        <div class="un-bg-gray-700 un-text-gray-200 un-m-1 un-py-1.5 un-px-1">
          <div class="un-text-lg un-font-bold">{{ book.title }}</div>
          <div class="un-text-xs un-font-semibold un-mt-2">{{ book.authors.join(', ') }}</div>
        </div>
      </div>
      <div class="un-mx-auto un-text-center un-text-xs un-font-light un-text-gray-700">{{ book.edition }}</div>
      <div class="un-flex-grow"></div>
      <div class="un-border un-border-gray-700 un-h-1 un-w-4/5 un-mx-auto"></div>
      <div class="un-w-4/5 un-mx-auto un-mt-3 un-text-center un-text-gray-800">{{ book.publisher }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import type { Book } from '@/types/book'
  defineProps<{ book: Book }>()
</script>
