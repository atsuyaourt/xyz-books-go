<template>
  <div v-if="isSuccess" class="un-flex un-space-x-4">
    <div class="un-flex un-w-1/3 un-justify-center un-items-center">
      <Book :book="book" />
    </div>
    <div class="un-flex un-flex-col un-w-2/3">
      <div class="un-space-x-2 un-align-middle">
        <span class="un-text-3xl un-font-bold">{{ book?.title }}</span>
        <span class="un-text-3xl un-font-bold">({{ book?.edition }})</span>
        <span class="un-text-2xl un-font-semibold un-text-gray-500">- {{ book?.publicationYear }}</span>
      </div>
      <div class="un-border-b-2 un-border-black un-w-full">
        by <span>{{ book?.authors.join(', ') }}</span>
      </div>
      <div class="un-border-b-2 un-border-black un-w-full">
        {{ bookPrice }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { Book as BookSchema } from '@/schemas/book'

  const route = useRoute()

  const isbn13 = computed(() => {
    const { isbn } = <{ isbn: string }>route.params
    return isbn
  })

  const { data: book, isSuccess } = useQuery({
    queryKey: ['book', isbn13],
    queryFn: async ({ queryKey }) => {
      const [_key, isbn] = queryKey
      const { data } = await axios.get(`/api/v1/books/${isbn}`)
      return BookSchema.parse(data)
    },
  })

  const bookPrice = computed(() => {
    const { price } = book.value ?? { price: undefined }
    if (price) {
      return `$ ${(Math.round(price * 100) / 100).toFixed(2)}`
    }

    return null
  })
</script>
