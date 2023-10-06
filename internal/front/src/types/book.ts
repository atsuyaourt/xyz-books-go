import { z } from 'zod'

import { Book as BooksSchema } from '@/schemas/book'

export type Book = z.infer<typeof BooksSchema>
