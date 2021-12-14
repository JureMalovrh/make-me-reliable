## Expose an unreliable API as reliable one

Requirements:
- Docker
- Make

### How to run:
1. create `.env` file (`.env.example` can be used as a guide)
2. run `make prod` 
3. run `docker-compose up reliable-api`

Server will run on port specified inside `.env` file

###Tests:
1. run `make ci` to run tests and linter


###API doc
- `/crypto/sign` where you should pass message as URL parameter the same as in specification
- `/job/{jobID}` where you will be able to check your job status if your response was unsuccesfull

If call to an HTTP server will fail - we return response with job ID, which will be used to check job status in the future.

Full Job json response looks like:
```
{
    "id":ID,
    "message":message requested,
    "successful":false/true,
    "lastTry":"2021-12-12T18:13:14.219Z",
    "result":"result of api call"
}
```

If call has succeded - successful will be marked as true and result will be the result of an API call.


Note:
At the time of writing this Readme, the provided API that is supposed to test out behaviour was not working, so we have implemented a custom unreliable API that fails in 50% of the cases.
To run it, run `docker-compose up unreliable-api` and set correct values inside .env file (use `unreliable-api` as API_URL).

###Arhitecture
We used Go as programming language, which allows easy concurrent model. For database, we used MongoDB for persistent storage.

We didn't use any MQ as requirements for RPS are really low and regular database can easily handle this. We used Mongo and background worker will pool database every few seconds and check if there is anything to try again.

For limitation of 10 RPS - we used Token buckets, where every 1 minute we get 10 tokens. After we spend more than 10 tokens, for remaining of the minute we will automatically create a job for client which will be done in the backgorund.  

### Notes:
- At the time of writing this Readme, the provided API that is supposed to test out behaviour was not working, so we have implemented a custom unreliable API that fails in 50% of the cases.
  To run it, run `docker-compose up unreliable-api` and set correct values inside .env file (use `unreliable-api` as API_URL).
- tests are not done in DRY way (read: copy-pasting), also some e2e tests are missing as we have only covered unit/integration part
- Repository Job would in production code be completly database agnostic (no bson mapping would be inside, ObjectId would be string, etc.). Some corners were cut :(.