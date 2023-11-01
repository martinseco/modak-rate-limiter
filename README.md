# modak-rate-limiter

This is the solution for the take-home test required as a part of the interview process for Modak.

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

### Building and Running the Project

Once you check out the code from GitHub, the project can be run using from the project source:
```
go mod tidy

go build -o notifications src/api/main.go

./notifications
```
