import { z } from 'zod'

export const CompaniesSchema = z.object({
  symbol: z.string(),
  name: z.string(),
  about: z.string(),
  sector: z.string(),
  industry: z.string(),
  mission: z.string(),
  vision: z.string(),
  phone: z.string(),
  website: z.string()
})
