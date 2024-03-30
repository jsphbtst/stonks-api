import { createClient } from '@libsql/client'
import algoliasearch from 'algoliasearch'
import { CompaniesSchema } from './schema/turso.schema'

const startTime = new Date().getTime()

if (
  !process.env.ALGOLIA_APP_ID ||
  !process.env.ALGOLIA_API_KEY ||
  !process.env.ALGOLIA_INDEX_NAME ||
  !process.env.TURSO_URL ||
  !process.env.TURSO_AUTH_TOKEN
) {
  console.error('Missing ENV keys')
  process.exit(1)
}

const algoliaClient = algoliasearch(process.env.ALGOLIA_APP_ID, process.env.ALGOLIA_API_KEY)
const index = algoliaClient.initIndex(process.env.ALGOLIA_INDEX_NAME)

console.log('✅ Successfully connected to Algolia')

const tursoClient = createClient({
  url: process.env.TURSO_URL,
  authToken: process.env.TURSO_AUTH_TOKEN
})

console.log('✅ Successfully connected to Turso', tursoClient)

async function main() {
  const response = await tursoClient.execute({
    sql: 'select * from Companies',
    args: {}
  })

  const rows = response.rows
  console.log('📝 Total number of stonks: ', rows.length)

  const algoliaObjects = []
  for (let idx = 0; idx < rows.length; idx++) {
    const row = CompaniesSchema.parse(rows[idx])
    algoliaObjects.push({
      objectID: rows[idx].id,
      ...row
    })
  }

  await index.saveObjects(algoliaObjects)
  console.log('✅ Successfully indexed Algolia')

  tursoClient.close()
  console.log('✅ Closed turso connection.')

  const endTime = new Date().getTime()
  const totalTimeInS = (endTime - startTime) / 1_000
  console.log(`Program ran for ${totalTimeInS} seconds.`)
}

main()
