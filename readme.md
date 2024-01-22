## Rate Limiter - Go

### Description

This is a simple rate limiter written in Go and is implemented as a middleware and can be used with any HTTP server.

### How to config

- To config the max number of requests per second for IP, change the value of `RATE_LIMITER_MAX_REQUESTS_PER_SECOND_FOR_IP` in .env file.

- To config the block time for IP, change the value of `RATE_LIMITER_BLOCK_TIME_FOR_IP` in .env file (this value is in seconds).

- To config the block time for token, change the value of `RATE_LIMITER_BLOCK_TIME_FOR_TOKEN` in .env file (this value is in seconds).

- To config the max number of requests per second for tokens, change the value of `RATE_LIMITER_TOKENS` in .env file. The format is: `token1:rate1,token2:rate2,...`. For example: `token1:10,token2:20`.

- To config the persistence mechanism, change the value of `PERSISTENCE_MECHANISM` in .env file. The value can be `redis` or `mysql`.

- OBS: when you change the values in .env you need to restart the server.

### How to run

- Run `docker-compose up -d` to start the server.

### How to test

- You can use the test api file in ./tests/api.http to test the server (you need have installed the REST Client (
humao.rest-client) on VSCode to use this file).
