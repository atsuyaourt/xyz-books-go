<template>
  <v-text-field
    v-model="searchString"
    density="compact"
    variant="outlined"
    :label="label"
    append-inner-icon="i-mdi:magnify"
    clear-icon="i-mdi:close-circle"
    @click:append-inner="handleSearch"
    @keydown.enter.prevent="handleSearch"
    @update:focused="handleSearch"
    @click:clear="handleSearch"
    class="un-w-full sm:un-w-96"
    clearable
  ></v-text-field>
</template>

<script setup lang="ts">
  defineProps<{
    label: string
  }>()
  const emit = defineEmits<{
    search: [searchString: string]
  }>()

  const searchString = ref('')

  const handleSearch = () => {
    emit('search', searchString.value)
  }
</script>
