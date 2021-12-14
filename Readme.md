# Make me reliable

Expose an unreliable API as a reliable one

## Getting started

### Prerequisites

- Docker - [Get Docker](https://docs.docker.com/get-docker/)
- Make

### Setup

- Create a `.env` file ([`.env.example`](.env.example) can be used as a guide)

### Usage

```sh
# Build docker image
$ make prod

# Start server
$ docker-compose up reliable-api
```

The server will run on the `SERVER_PORT` specified in the `.env` file

### Tests

```sh
# Run tests and linter
$ make ci
```

## API Documentation

### Create job

```http
GET /crypto/sign?message=...
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `message` | `string` | **Required** Message query parameter |

### Get job status

```http
GET /job/{jobID}
```

### Responses

The response for `/crypto/sign` and `/job/{jobID}` are the same

```javascript
{
    "id"         : string,   // Job ID (Mongo ObjectId)
    "message"    : string,   // The message requested
    "successful" : boolean,  // true or false
    "lastTry"    : string,   // eg. "2021-12-12T18:13:14.219Z"
    "result"     : string    // Result of the API call
}
```

If the call to the unreliable API fails - returns a response with the job ID, which will be used to check the job status in the future.

If the call succedes - `successful` will be `true` and `result` will be the result of the API call.

## Description

### Arhitecture

We used [Go](https://go.dev/) as the programming language, which allows easy concurrency and [MongoDB](https://www.mongodb.com/) for persistent storage.

We didn't use any MQ, as requirements for RPS are low and a regular database can easily handle this. We used MongoDB and a background worker for polling the database every few seconds to check if there is any job to try again.

For the limitation of 10 RPS we used Token buckets, where every 1 minute we get 10 tokens. After we spend 10 tokens, requests in the remaining time will automatically create a job for the client which will be done in the backgorund.

### Notes

- At the time of writing this Readme, the provided API that is supposed to test out behaviour was not working, so we have implemented a custom unreliable API that fails in 50% of the time.
  To run it, run `docker-compose up unreliable-api` and set correct values inside the `.env` file (use `unreliable-api` as `API_URL`).

- Tests are not done in the DRY way (read: copy-pasting), also some e2e tests are missing as we have only covered unit/integration part

- Repository Job would in production code be completly database agnostic (no bson mapping would be inside, ObjectId would be string, etc.). Some corners were cut :frowning_face:.
