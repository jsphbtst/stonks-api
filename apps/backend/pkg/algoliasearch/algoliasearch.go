package algoliasearch

import "github.com/algolia/algoliasearch-client-go/v3/algolia/search"

type AlgoliaSearch struct {
	Client *search.Client
	Index  *search.Index
}

func Init(appId string, apiKey string, indexName string) (*search.Client, *search.Index) {
	client := search.NewClient(appId, apiKey)
	index := client.InitIndex(indexName)

	return client, index
}
