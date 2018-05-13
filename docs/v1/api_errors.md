# Errors

In case of Error during handling user’s request the API responds using the following schema:

```json
{
   "error": {
      "code": 500,
      "message": "Something bad happened",
      "status": "INTERNAL"
   }
}
```
## Error Codes

The following error codes are the possinble errors of all methods

Error | Code | Status | Related Requests
------|------|----------|------------------
Invalid JSON | 400 | BAD REQUEST | Create Service (POST)
Service already exists | 409 | CONFLICT | Create Service (POST)
Service Invalid Argument| 422 | UNPROCCESABLE ENTITY| Create Service (POST)
  