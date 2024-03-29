import { CompanyInformationSchema } from '../schema/finnworld.schema'

export async function getCompanyInformation(symbol: string) {
  const params = new URLSearchParams()
  params.append('key', process.env.FINNWORLDS_API_KEY)
  params.append('ticker', symbol)

  const url = `https://api.finnworlds.com/api/v1/information?${params.toString()}`
  const response = await fetch(url, {
    method: 'GET',
    headers: {
      accept: 'application/json'
    }
  })
  const responseJson = await response.json()
  const parsedResponseJson = CompanyInformationSchema.parse(responseJson)
  return parsedResponseJson
}
