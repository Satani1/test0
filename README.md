# Wildberries Task Level 0


### API

Requests:
1. `/` GET - return index page with html form to search order by ID;
2. `/order/{id}` POST - return page with order data by ID;
3. `/create` POST - create order with random uID and test data;


### Launch service

`docker compose up -d --build`

To access: `localhost:8080`
