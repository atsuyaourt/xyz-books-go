import { z } from 'zod'

import { camelize } from '@/schemas/common'

export const createPaginatedList = <ItemType extends z.ZodTypeAny>(itemSchema: ItemType) => {
  return z
    .object({
      page: z.number(),
      per_page: z.number(),
      total_pages: z.number().default(0),
      next_page: z.number(),
      prev_page: z.number(),
      count: z.number(),
      items: z.array(itemSchema),
    })
    .transform(camelize)
}
