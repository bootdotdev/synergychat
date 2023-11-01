# SynergyChat

SynergyChat is not only the best workforce chatting tool, but the best analytics engine. SynergyChat is powered by several microservies:

* "web" - The web frontend. This small micro-service serves static HTML, CSS, and JavaScript files that form the shell of the application.
* "api" - The web API. This micro-service is exposed as our public API, and powers the data for the web frontend.
* "crawler" - The analytics crawler. This micro-service scrapes a public repository of books for various keywords and stores them in a database. It is not exposed as a public API, but the "api" micro-service can access its data through internal-to-k8s HTTP requests.

## API

### Environment Variables

| Name             | Description                                                                                        | Required | Example                            |
| ---------------- | -------------------------------------------------------------------------------------------------- | -------- | ---------------------------------- |
| API_PORT         | The port the server will listen on                                                                 | True     | 8080                               |
| API_DB_FILEPATH  | The file path where the database be created and stored. If omitted, ephemeral memory will be used. | False    | `/var/lib/synergychat/api/db.json` |
| CRAWLER_BASE_URL | The base URL of the crawler service. If not provided slash commands won't work.                    | False    | `http://localhost:8081`            |

### HTTP Endpoints

#### `GET /healthz`

Returns a 200 OK if the service is healthy.

#### `POST /messages`

Creates a new message. Request body example:

```json
{
  "AuthorUsername": "john_doe",
  "Text": "Hello, world!"
}
```

Profane words like "heck", "darn", and "fetch" might cause... problems.

The `/stats` slash command can be used at the beginning of the `Text` field to get a response from the crawler bot. Here are some examples:

* `/stats`: Returns a summary of all keywords crawled in all books.
* `/stats keywords=love`: Returns a summary of the keyword "love" in all books.
* `/stats keywords=love,hate`: Returns a summary of the keywords "love" and "hate" in all books.
* `/stats title=Frankenstein`: Returns a summary of all keywords in the book "Frankenstein".
* `/stats keywords=love,hate title=Frankenstein`: Returns a summary of the keywords "love" and "hate" in the book "Frankenstein".

The keywords need to actually be crawled by the crawler before they can be queried, so make sure the crawler is configured and has been running for a while before querying for keywords.

The array of previously created messages is returned in the response body.

#### `GET /messages`

An array of previously created messages is returned in the response body:

```json
[
    {
        "AuthorUsername": "john_doe",
        "Text": "Hello, world!"
    },
    {
        "AuthorUsername": "jane_sue",
        "Text": "Hello to you ;)"
    }
]
```

## Crawler

### Environment Variables

| Name             | Description                                                                                                    | Required | Example                           |
| ---------------- | -------------------------------------------------------------------------------------------------------------- | -------- | --------------------------------- |
| CRAWLER_PORT     | The port the server will listen on                                                                             | True     | 8081                              |
| TO_CRAWL_URL     | The base URL of the website to crawl                                                                           | True     | https://www.gutenberg.org/books   |
| CRAWLER_KEYWORDS | The keywords to search for. Only included keywords will be counted                                             | True     | love,hate                         |
| CRAWLER_DB_PATH  | The directory path where the database files be created and stored.  If omitted, ephemeral memory will be used. | False    | `/var/lib/synergychat/crawler/db` |

### HTTP Endpoints

#### `GET /healthz`

Returns a 200 OK if the service is healthy.

#### `GET /stats`

Optional query parameters:

* `keywords`: A comma-separated list of keywords to filter the results by. If omitted, all keywords will be returned.
* `title`: A book title to filter the results by. If omitted, all books will be returned.

Returns an array of JSON objects containing the counts of the keywords in the database. For example:

```json
[
  {
    "Keyword": "hate",
    "BookTitle": "The Strange Case of Dr. Jekyll and Mr. Hyde",
    "Count": 10
  },
  {
    "Keyword": "love",
    "BookTitle": "The Strange Case of Dr. Jekyll and Mr. Hyde",
    "Count": 12
  }
  {
    "Keyword": "love",
    "BookTitle": "Frankenstein; Or, The Modern Prometheus",
    "Count": 17
  }
]
```
