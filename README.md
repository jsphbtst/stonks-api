# Stonks API

Current 3rd-Party Dependencies:

- Turso for DB;
- Aiven for Redis (set to use LRU for its Redis data eviction policy); and
- Algolia for Search;
- Finnworlds (paid API) for reliable S&P 500 company data

I already had a prepared CSV for the S&P 500 stocks from an old CLI project. I used the tickers there to grab the data I needed from Finnworlds, committed select parts of those data to Turso, and then indexed Algolia with those rows saved to Turso. This allowed me to create two API endpoints:

- `GET /api/v1/companies/:symbol`; and
- `GET /api/v1/companies/search?query=`

The first endpoint is a simple GET endpoint that returns the ticker data. I implemented cache-aside to make response times faster. For cache misses, I concurrenly do a `SetEx` to redis while responding to the user, which should reduce response times.

The second endpoint makes use of Algolia to handle smart searching for me. I also implemented cache-aside here where the `query` is the key. The same strategy of concurrently peforming `SetEx` while responding to the user for cache misses is implemented here.

Lastly, I also implemented a custom rate limiter middleware that makes use of Redis—as opposed to making it in-memory. The rate limiter middleware uses a custom fixed-window strategy, where you can only make at most 12 API calls per minute.

## Future Plans

- [ ] Implement a custom RAG solution for my search (hosted in `GET /api/v2/search?query=`)
- [ ] Build out the frontend to quickly demo blazingly fast search
