import { z } from 'zod'

export const CompanyInformationSchema = z.object({
  result: z.object({
    basics: z.object({
      company: z.string(),
      ticker: z.string(),
      sector: z.string(),
      industry: z.string()
    }),
    output: z.object({
      name: z.string(),
      ticker: z.string(),
      address: z.object({
        street1: z.string(),
        street2: z.string(),
        city: z.string(),
        stateOrCountry: z.string(),
        postal_code: z.string()
      }),
      post_address: z.object({
        street1: z.string(),
        street2: z.string(),
        city: z.string(),
        stateOrCountry: z.string(),
        postal_code: z.string()
      }),
      phone: z.string(),
      website: z.string(),
      about: z.string(),
      mission: z.string(),
      vision: z.string()
    })
  })
})
