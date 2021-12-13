## Expose an unreliable API as reliable one

Requirements:
- Docker
- Make

First, create `.env` file (`.env.example` can be used as a guide).

To start services, run:
```make prod && docker-compose up reliable-api```

This will expose an API on port 8080 with two routes
- `/crypto/sign` where you should pass message as URL parameter
- `/job/{jobID}` where you will be able to check your job status

If call to an HTTP server will fail - we return response with job ID, which will be used to check job status in the future.

Full Job json response looks like:
```
{
    "id":"61b63b80703f8d8157aae2c0",
    "message":"fake message",
    "successful":true,
    "lastTry":"2021-12-12T18:13:14.219Z",
    "result":"result of api call"}
```

If call has succeded - successful will be marked as true and result will be the result of an API call.

Note:
At the time of writing this Readme, the provided API that is supposed to test out behaviour was not working, so we have implemented a custom unreliable API that fails in 50% of the cases.
To run it, run `docker-compose up unreliable-api` and set correct values inside .env file (use `unreliable-api` as API_URL).

###Arhitecture
TODO



### Notes:
- tests are not done in DRY way (read: copy-pasting), also some e2e tests are missing as we have only covered unit/integration part