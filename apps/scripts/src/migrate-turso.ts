import { createClient } from '@libsql/client'
import { getCompanyInformation } from './utils/finnworld'

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
  const file = Bun.file('apps/scripts/src/assets/snp500.csv')
  const text = await file.text()
  const textSplit = text.split('\n')

  const regex = /(".*?"|[^,]+)/g
  const tickers: string[] = []
  for (let idx = 1; idx < textSplit.length; idx++) {
    const row = textSplit[idx].match(regex)
    if (!row) {
      console.error('Could not split: ', textSplit[idx])
      continue
    }

    const symbol = row[0]
    tickers.push(symbol)
  }

  const sortedTickers: string[] = tickers.sort()
  console.log('📝 Total tickers parsed: ', sortedTickers.length)

  const failedTickers: string[] = []
  for (let idx = 0; idx < sortedTickers.length; idx++) {
    const symbol = sortedTickers[idx]
    console.log(`[${idx}]: Processing ${symbol}...`)

    try {
      const response = await client.execute({
        sql: 'SELECT * FROM Companies where symbol = ?',
        args: [symbol]
      })
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const responseJson = response.toJSON() as { rows: any[] }
      const isInDb = responseJson?.rows?.length > 0
      if (isInDb) {
        console.log(`${symbol} is already in Turso. Skipping...`)
        continue
      }

      const company = await getCompanyInformation(symbol)

      const [address, postalAddress] = await Promise.all([
        client.execute({
          sql: `INSERT INTO Addresses (street1, street2, city, state, postalCode) VALUES (:street1, :street2, :city, :state, :postalCode)`,
          args: {
            street1: company?.result?.output?.address?.street1 ?? '',
            street2: company?.result?.output?.address?.street2 ?? '',
            city: company?.result?.output?.address?.city ?? '',
            state: company?.result?.output?.address?.stateOrCountry ?? '',
            postalCode: company?.result?.output?.address?.postal_code ?? ''
          }
        }),
        client.execute({
          sql: `INSERT INTO Addresses (street1, street2, city, state, postalCode) VALUES (:street1, :street2, :city, :state, :postalCode)`,
          args: {
            street1: company?.result?.output?.post_address?.street1 ?? '',
            street2: company?.result?.output?.post_address?.street2 ?? '',
            city: company?.result?.output?.post_address?.city ?? '',
            state: company?.result?.output?.post_address?.stateOrCountry ?? '',
            postalCode: company?.result?.output?.post_address?.postal_code ?? ''
          }
        })
      ])

      const addressJson = address.toJSON() as { lastInsertRowid: string }
      const postalAddressJson = postalAddress.toJSON() as { lastInsertRowid: string }

      await client.execute({
        sql: 'INSERT INTO Companies (symbol, name, addressId, postalAddressId, sector, industry, phone, website, about, mission, vision) VALUES (:symbol, :name, :addressId, :postalAddressId, :sector, :industry, :phone, :website, :about, :mission, :vision)',
        args: {
          symbol: company.result.basics.ticker,
          name: company.result.basics.company,
          addressId: addressJson.lastInsertRowid,
          postalAddressId: postalAddressJson.lastInsertRowid,
          sector: company?.result?.basics?.sector ?? '',
          industry: company?.result?.basics?.industry ?? '',
          phone: company?.result?.output?.phone ?? '',
          website: company?.result?.output?.website ?? '',
          about: company?.result?.output?.about ?? '',
          mission: company?.result?.output?.mission ?? '',
          vision: company?.result?.output?.vision ?? ''
        }
      })

      console.log(`Finished processing ${symbol}`)
    } catch (e) {
      console.error(`Failed to process ${symbol}: `, e)
      failedTickers.push(symbol)
    }
  }

  console.log('✅ Finished processing stonks.')

  client.close()
  console.log('✅ Closed turso connection.')

  console.log('Here are the failed tickers: ', failedTickers)

  const endTime = new Date().getTime()
  const totalTimeInS = (endTime - startTime) / 1_000
  console.log(`Program ran for ${totalTimeInS} seconds.`)
}

main()
