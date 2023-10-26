# synergychat

SynergyChat is not only the best workforce chatting tool, but the best analytics engine. SynergyChat is powered by several microservies:

* "web" - The web frontend. This small micro-service serves static HTML, CSS, and JavaScript files that form the shell of the application.
* "api" - The web API. This micro-service is exposed as our public API, and powers the data for the web frontend.
* "crawler" - The analytics crawler. This micro-service crawls the web scraping analytics about various keywords and stores them in a database. It is not exposed as a public API, but the "api" micro-service can access its data through internal-to-k8s HTTP requests.
