# Crawler

## Clone repo:

```
git clone https://github.com/ohal/go-second.git
```

## Build:

NOTE: Go 1.9 have to be preinstalled before build.

To build:

```
$ cd ./go-second/BE/go-crawler
$ make
```

## Start:

Before start make sure if app can bind to port `:3000` and MongoDB is accessible without authentication on `127.0.0.1:27017`.

To start app:

```
$ ./apps/go-crawler/go-crawler
```

It will scrape reviews from `https://apps.shopify.com/omnisend#reviews-heading` and save all scraped reviews to DB.

## Requests:

To get all reviews stored in DB:

```
$ curl -s -X GET http://localhost:3000/api/v1/reviews
```