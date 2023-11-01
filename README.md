# modak-rate-limiter

This is the solution for the take-home test required as a part of the interview process for Modak.

### Building and Running the Project

Once you check out the code from GitHub, the project can be run using from the project source:
```
go mod tidy

go build -o notifications src/api/main.go

./notifications
```

### Example Request/Responses

#### Send notification for user
Request:
```
curl --location --request POST 'localhost:8080/notifications' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id": "user_123",
    "type": "news",
    "message": "message to send"
}'
```
##### Successful Response :
Response with HTTP Status: **202**

##### Error Responses:
```
{
    "error": "Forbidden",
    "message": "news's notification type for user_123 user reached it's limit!",
    "status": 403
}
```

### Sample Data
The project builds with the following sample rate limiting rules in the DB:
* **Status**: not more than 1 per 10 seconds for each recipient 
* **News**: not more than 3 per 30 seconds for each recipient 
* **Marketing**: not more than 4 per minute for each recipient

These can be edited in the applications/services.go file.

NOTE: _The endpoints to manage rules were left out of the project scope._
