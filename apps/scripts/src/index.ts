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
  console.log('Total tickers parsed: ', sortedTickers.length)
  console.log(sortedTickers)
}

main()
