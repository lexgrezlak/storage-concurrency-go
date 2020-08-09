## How to run

```
docker-compose up
```

## Configuration
There are default values provided, but you can change them either in `configs/development.yml` or `docker.env`

## Features
Concurrently reads records from the provided csv file, adds them to Redis, and serves the data using the endpoint below.


### Endpoint
```
curl http://localhost:1321/promotions/promotion-id-from-promotions.csv
```
returns
```
{"id":"172FFC14-D229-4C93-B06B-F48B8C095512", "price":9.68, "expiration_date": "2018-06-04 06:01:20"}
```