const startTime = new Date().getTime()

import { createClient } from '@libsql/client'

if (!process.env.OPENAI_API_KEY || !process.env.TURSO_URL || !process.env.TURSO_AUTH_TOKEN) {
  console.error('Missing env vars')
  process.exit(1)
}

const client = createClient({
  url: process.env.TURSO_URL,
  authToken: process.env.TURSO_AUTH_TOKEN
})

console.log('✅ Successfully connected to Turso', client)

type OpenAiResponse = {
  object: string
  model: string
  data: {
    object: string
    index: number
    embedding: number[]
  }[]
  usage: {
    prompt_tokens: number
    total_tokens: number
  }
}

type VectorDbRow = { id: number; text: string; embedding: number[] }

async function main() {
  // TODO: if success json file exists, just parse that as
  // that's already pre-populated with the embeddings data...
  // at least I hope it does, past J - future J

  const dbResponse = await client.execute({
    sql: 'SELECT * FROM Companies',
    args: []
  })
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const dbResponseJson = dbResponse.toJSON() as { columns: string[]; rows: any[] }

  /*
  "0": "id",
  "1": "symbol",
  "2": "addressId",
  "3": "postalAddressId",
  "4": "sector",
  "5": "industry",
  "6": "phone",
  "7": "website",
  "8": "about",
  "9": "mission",
  "10": "vision",
  "11": "name"
  */
  const vectorDbData: VectorDbRow[] = []
  for (let idx = 0; idx < dbResponseJson.rows.length; idx++) {
    const currentRow = dbResponseJson.rows[idx]
    const id = currentRow[0]
    const symbol = currentRow[1]
    const name = currentRow[11]
    const sector = currentRow[4]
    const industry = currentRow[5]
    const about = currentRow[8]

    const text = `name: ${name} (${symbol}). sector: ${sector}. industry: ${industry}. description: ${about}`
    vectorDbData.push({ id, text, embedding: [] })
  }

  // Attach OpenAI Generated Embedding to Each Data
  const openAiApiKey = process.env.OPENAI_API_KEY as string
  const url = 'https://api.openai.com/v1/embeddings'

  const successVectorDbData: VectorDbRow[] = []
  const failedVectorDbData: VectorDbRow[] = []
  for (let idx = 0; idx < vectorDbData.length; idx++) {
    const row = vectorDbData[idx]
    console.log('Generating embedding for: ', row.id)

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${openAiApiKey}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        input: 'Hello, world!',
        model: 'text-embedding-ada-002'
      })
    })

    if (!response?.ok) {
      failedVectorDbData.push(row)
      console.log('❌ Failed to generate embedding for: ', row.id)
      continue
    }

    const responseJson: Awaited<OpenAiResponse> = await response.json()
    if (!responseJson?.data?.length) {
      failedVectorDbData.push(row)
      console.log('❌ Failed to generate embedding for: ', row.id)
      continue
    }

    const openAiEmbedding = responseJson.data[0].embedding
    row.embedding = openAiEmbedding
    successVectorDbData.push(row)

    console.log('Finished saving embedding for: ', row.id)
  }

  console.log('Saving successful generated embeddings of length: ', successVectorDbData.length)
  await Bun.write(
    'tmp-success-generate-embeddings.json',
    JSON.stringify({ data: successVectorDbData }, null, 2)
  )

  console.log('Saving failed to generate embeddings of length: ', failedVectorDbData.length)
  await Bun.write(
    'tmp-failed-generate-embeddings.json',
    JSON.stringify({ data: failedVectorDbData }, null, 2)
  )
}

main()
  .catch(err => {
    console.error(err)
    process.exit(1)
  })
  .finally(() => {
    client.close()
    console.log('✅ Closed turso connection.')

    const endTime = new Date().getTime()
    const duration = (endTime - startTime) / 1_000
    console.log(`Program ran for ${duration} s`)
  })
