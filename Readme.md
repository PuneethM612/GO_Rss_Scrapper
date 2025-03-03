## MOTIVE OF THE PROJECT

- This Go program is designed to scrape RSS feeds periodically and store the fetched posts in a PostgreSQL database.
- Fetching feeds to scrape from the database.
- Using goroutines for concurrent processing of multiple feeds.
- Parsing RSS feed data and storing it in the database.
- Handling errors gracefully to ensure continuous execution.

Make use of chi-router for the server part

- cors are used for sending a bunch of HTTP headers
- cors also allow you to send all the requests and headers and there are ways to also put constraints
- sqlc and goose were isnyalled as they work on raw SQL

## MAIN.GO

- This is where the entire application starts and
- It initializes the database, middleware, routes, and HTTP Handlers
- Calls the Scrapper.go file in case the server is running in the background.

## SCRAPPER.GO

- The file is responsible for periodically fetching the RSS feeds and storing the posts in the database.
- It also makes use of GoRoutines for handling the concurrency.
- Once the scrapping starts in main.go, the function runs in the background continously and fetches the RSS feed regularly.
- It also makes use of the database to store the posts.

## GO ROUTINES

- Go routines are used to run the scrapping function in the background, using concurrency without interepting the actual flow of the main application.

# OVERVIEW OF THE RSS SCRAPPER.GO FILE

- Fetches the feed from the database and then it is updated at the database to overcome the duplicay
- Scrappes each feed maing use of concurrency and go routines.
- Parsing the feed data to `URLtoFeed` in rss.go
- storing the new feed in the database.

## SCRAPPER. GO File Explanation

- This function, startScrapping, is responsible for periodically fetching RSS feeds from the database and processing them concurrently using goroutines
- We make use of the tickers so that the process of fetching the feed is done on a regular basis.
- And we make use of the context package in go to ensured for controlling the deadlines, timeouts, cancellation and other operations.
- So basically, `GetNextFeedsToFetch` is used to get the feeds from the Database and we make use of the contecxt package in go to maintain the fetching process with the help of Tickers.

## Understanding SyncWaitGroup

- It is responsible for ensuring that multiple web scrapping tasks are organsied concurrently and it also ensures that one batch of scrapping tasks are finished before it moving on to the next batch.
- And we basically run a for loop to hover through the feeds and keep the tasks of fetching and processed being gone.
- counter keeps track of number of go routines currently and `.Add` functions increases them accordingly, if the counter keeps track and if its 5 then it is increases by 5.

## channels and Items

- we know the feeds are in the format of XML and whe we need ot parse them we make use of the channel and items specifically to ensure that, channel basically stores the Meta Data of the feeds and items are the actual feeds.

## OVERVIEW OF THE RSS.GO FILE

- We basicalyy make an HTTP request to fetch an RSS feed from the given url and timeout if more than 10 seconds
- Read the response which is in XMl Format
- Parsing the XMl data to the Go Structs that is Channel and Items.
- Return the Parsed feed so that it can be processed further.

## RSS.GO URLTOFEED FUNCTION

- This Function is basically used for converting the URL into a particular feed using the Parsing method.\

## Overview of Handler_user.go

1. Handler Create user : It basically creates a user by taking the name in the JSON.

- It generates the unique ID using the uuid package and it then stores in the database.
- Rest of the code can be explaine using the Code

2. Handler Get User:

- The handler get user function basically gets the data about the user that is being stored in the database.

3. Handler Get Posts for User:

- Basically this code of method is used for getting the posts for the user that have been stored in the database and have been atgged along.

## OVERVIEW OF HANDLER_FEED. GO

1. HandlerCreateFeed - It basically takes the name and feed url as the input and stores in the database with a unique id
2. HandlerGetFeed - It basically gets the feed from the database and returns it.

## OVERVIEW OF HANDLER FEED FOLLOWS.GO

1. Handler Create Feed Follows - It basically lets the user to follow a feed by taking the json request.

- Also creates a new entry in the database with the specified data parameters.

2. Handler Get Feed Follows - It basically returns the data that is being stored in the data base after being created.

3. Handler Delete Feed Follows - It basically deletes the feed follows from the database by taking the feedID

## OVERVIEW FOR HANDLER READINESS AND ERR.GO

1. Handler Readiness - ensures that the server is ready to take on requests or no.
2. Handler err - It is used to handle the errors using the internal error message.

## OVERVIEW ON JSON.GO FILE

- Respond with error: Meaning if the error is more than 500 then it asks for a debug error.
- Respomf with JSON:

## OVERVIEW OF MIDDLEWARE AUTH.GO

- Basically this is used to authenticate the user before accessing any of the routes.
- It checks with the APIKEY and the data base to authenticate the user.

{
"name": "Hacker News",
"url": "https://news.ycombinator.com/rss"
}

{
"name": "BBC Technology",
"url": "http://feeds.bbci.co.uk/news/technology/rss.xml"
}

{
"name": "Ars Technica",
"url": "http://feeds.arstechnica.com/arstechnica/index"
}

{
"name": "New York Times - World News",
"url": "https://rss.nytimes.com/services/xml/rss/nyt/World.xml"
}

{
"name": "CNN Top Stories",
"url": "http://rss.cnn.com/rss/edition.rss"
}
