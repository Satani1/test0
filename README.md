# Wildberries Task Level 0


### API

Requests:
1. `/` GET - return index page with html form to search order by ID;
2. `/order/{id}` POST - return page with order data by ID;
3. `/create` POST - create an order via a POST request using json data;

Generate Data:

For generate some data and publish that into 'orders' channel in nats:
```
go run ./cmd/publisher
```

### Launch service

`docker compose up -d --build`

To access: `localhost:8080`

