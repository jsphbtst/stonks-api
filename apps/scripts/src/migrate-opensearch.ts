import { createClient } from '@libsql/client'
import { CompaniesSchema } from './schema/turso.schema'

const startTime = new Date().getTime()

if (!process.env.TURSO_URL || !process.env.TURSO_AUTH_TOKEN || !process.env.FINNWORLDS_API_KEY) {
  console.error('Missing env vars')
  process.exit(1)
}

const client = createClient({
  url: process.env.TURSO_URL,
  authToken: process.env.TURSO_AUTH_TOKEN
})

console.log('✅ Successfully connected to Turso', client)

async function main() {
  const response = await client.execute({
    sql: 'select * from Companies',
    args: {}
  })

  const rows = response.rows
  console.log('📝 Total number of stonks: ', rows.length)

  const TOTAL = rows.length
  const BUFFER_LENGTH = 40
  for (let idx = 0; idx < TOTAL; idx++) {
    const array: string[] = []

    let jdx = idx
    for (; jdx < idx + BUFFER_LENGTH; jdx++) {
      if (jdx > TOTAL - 1) {
        break
      }

      const row = CompaniesSchema.parse(rows[jdx])
      const ticker = row.symbol
      array.push(ticker)
    }

    idx = jdx - 1
    // TODO: bulk commit to OpenSearch here
    console.log(array)
  }

  client.close()
  console.log('✅ Closed turso connection.')

  const endTime = new Date().getTime()
  const totalTimeInS = (endTime - startTime) / 1_000
  console.log(`Program ran for ${totalTimeInS} seconds.`)
}

main()
