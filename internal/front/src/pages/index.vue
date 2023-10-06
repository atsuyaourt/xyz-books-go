<template>
  <div class="un-flex un-flex-col un-justify-center un-items-center un-gap-4 un-w-full">
    <div
      class="un-flex un-flex-col md:un-flex md:un-flex-row un-space-y-1 md:un-space-x-4 md:un-space-y-0 un-w-full md:un-w-5/6"
    >
      <SearchBar label="Search books" @search="titleFilter = $event" class="un-w-full md:un-w-1/3"></SearchBar>
      <SearchBar label="Search by author" @search="authorFilter = $event" class="un-w-full md:un-w-1/3"></SearchBar>
      <SearchBar label="Search by publisher" @search="publisherFilter = $event"></SearchBar>
    </div>

    <div class="un-p-4 un-grid un-grid-cols-3 md:un-grid-cols-5 un-gap-4" v-if="isSuccess">
      <router-link v-for="book in data?.items" :to="generateUrl(book)" class="un-no-underline hover:un-shadow-md">
        <Book :book="book" />
      </router-link>
    </div>

    <v-pagination
      v-if="isSuccess"
      v-model="curPage"
      :length="data?.totalPages"
      :total-visible="data?.totalPages ?? 0 > 6 ? 6 : 3"
      prev-icon="i-mdi:chevron-left"
      next-icon="i-mdi:chevron-right"
    ></v-pagination>
  </div>
</template>

<script setup lang="ts">
  import { Book as BookSchema } from '@/schemas/book'
  import { createPaginatedList } from '@/schemas/page'

  import type { Book } from '@/types/book'

  const titleFilter = ref('')
  const authorFilter = ref('')
  const publisherFilter = ref('')
  const curPage = ref(1)

  const { data, isSuccess } = useQuery({
    queryKey: ['books', curPage, { title: titleFilter, author: authorFilter, publisher: publisherFilter }],
    queryFn: async ({ queryKey }) => {
      const [_key, page, filterParams] = queryKey as [string, number, Record<string, string | null>]
      const filterStr = Object.keys(filterParams)
        .filter((key) => filterParams[key])
        .map((key) => key + '=' + filterParams[key])
        .join('&')
      const { data } = await axios.get(`/api/v1/books?page=${page}&${filterStr}`)

      return createPaginatedList(BookSchema).parse(data)
    },
  })

  const generateUrl = (book: Book) => {
    const isbn13 = book.isbn13 ? book.isbn13 : isbn10toisbn13(book.isbn10)
    return `/${isbn13}`
  }

  const isbn10toisbn13 = (isbn10?: string) => {
    if (isbn10 === undefined || !/^\d{9}[\dX]$/.test(isbn10)) {
      return
    }

    const partial = '978' + isbn10.slice(0, 9)

    const digits = partial.split('').map(Number)
    const checkDigit = (10 - (digits.reduce((acc, val, index) => acc + (index % 2 ? val * 3 : val), 0) % 10)) % 10

    return partial + checkDigit
  }
</script>
