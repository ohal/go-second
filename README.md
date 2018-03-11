# Get list of five most common three-word phrases from Omnisend reviews

## Clone repo:

Use git for cloning repository:

```
git clone https://github.com/ohal/go-second.git
```

### BE build:

NOTE: Go 1.9 have to be preinstalled before build.

To build:

```
$ cd ./go-second/BE/go-crawler
$ make
```

Or:

```
$ cd ./go-second/BE/go-crawler
$ export GOPATH=$PWD && pushd ./src/go-crawler/ && go build -ldflags="-s -X main.ApplicationVersion=0.0.1" -a -o ../../apps/go-crawler/go-crawler && popd
```

### BE start:

Before start make sure if app can bind to port `:3000` and MongoDB is accessible without authentication on `127.0.0.1:27017`.

To start app:

```
$ ./apps/go-crawler/go-crawler
```

It will scrape reviews from `https://apps.shopify.com/omnisend#reviews-heading` and save all scraped reviews to DB.

Progress of reviews scraping will be logged to screen.

### BE REST API endpoints:

To get all reviews stored in DB:

```
$ curl -s -X GET http://localhost:3000/api/v1/reviews
```

To get reviews in date range:

```
$ curl -s -X POST -H "Content-Type: application/json" -d '{"beginDate":{"day":9,"month":10,"year":2010},"endDate":{"day":19,"month":10,"year":2018}}' http://localhost:3000/api/v1/range
```

To get list of  five most common three-word phrases in date range:

```
$ curl -s -X POST -H "Content-Type: application/json" -d '{"beginDate":{"day":9,"month":10,"year":2010},"endDate":{"day":19,"month":10,"year":2018}}' http://localhost:3000/api/v1/shingle
```

### FE install:

Before start make sure Node.js®, npm and Angular CLI are installed on your machine.

npm version details:

```
$ npm version
{ 'o-reviews': '0.0.0',
  npm: '5.6.0',
  ares: '1.10.1-DEV',
  cldr: '32.0',
  http_parser: '2.7.0',
  icu: '60.1',
  modules: '57',
  nghttp2: '1.25.0',
  node: '8.10.0',
  openssl: '1.0.2n',
  tz: '2017c',
  unicode: '10.0',
  uv: '1.19.1',
  v8: '6.2.414.50',
  zlib: '1.2.11' }
$ cd ./go-second/FE/o-reviews
$ npm install -g @angular/cli
```

Install app dependencies:

```
$ cd ./go-second/FE/o-reviews
$ npm install
```

### FE start:

Before start make sure if app can bind to port `:4200` and BE REST API is accessible as `http://localhost:3000/api/v1/`.

To start application:

```
$ ng serve --open
```

It will automatically open your browser on `http://localhost:4200/`

Now you should be able to access page and get list of five most common three-word phrases from Omnisend reviews using date range picker to choose date range.
