import { z } from 'zod'

import { camelize } from '@/schemas/common'

export const Book = z
  .object({
    title: z.string(),
    isbn13: z.string().optional(),
    isbn10: z.string().optional(),
    publication_year: z.number(),
    price: z.number(),
    image_url: z.string().optional(),
    edition: z.string().optional(),
    publisher: z.string(),
    authors: z.string().array(),
  })
  .transform(camelize)
