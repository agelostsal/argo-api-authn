  # Auth method API Calls

## [POST] Manage Auth Methods - Create New Auth Method

This request creates a new service.

#### Request

`POST /v1/authM`

### Example request
```
curl -X POST -H "Content-Type: application/json"
  "https://{URL}/v1/authM?key={key_in_the_config}"
```


##### Post Body

```json
        {
            "access_key": "key1",
            "host": "127.0.0.1",
            "path": "/path/{{identifier}}?key={{access_key}}",
            "port": 9000,
            "service": "s1",
            "type": "api-key"
        }
 ```
 
  ### Response
  
  If the request is successful, the response contains the newly created auth method.
  
  Success Response
  
  `201 CREATED`
  
  ```json
        {
            "access_key": "key1",
            "host": "127.0.0.1",
            "path": "/path/{{identifier}}?key={{access_key}}",
            "port": 9000,
            "service": "s1",
            "type": "api-key"
        }
   ```
 
  ## [GET] Manage Auth Methods - ListOneAuthMethod
  
  ### Request
  
  `GET /v1/services/{service}/hosts/{host}/authM`
  
  ### Example request
  ```
  curl -X GET -H "Content-Type: application/json"
  "https://{URL}/v1/services/{service}/hosts/{host}/authM?key={key_in_the_config}"
  ```
  
   If the request is successful, the response contains information for the requested auth method.
   
   Success Response
   
   `200 OK`
   
```json
{
    "access_key": "b328c3861f061f87cbd34cf34f36ba2ae20883a5",
    "host": "127.0.0.1",
    "path": "/v1/users:byUUID/{{identifier}}?key={{access_key}}",
    "port": 8081,
    "service": "ams",
    "type": "api-key"
}
```


  ## [GET] Manage Auth Methods - ListAllAuthMethods
  
  ### Request
  
  `GET /v1/authM`
  
  ### Example request
  ```
  curl -X GET -H "Content-Type: application/json"
  "https://{URL}/v1/authM?key={key_in_the_config}"
  ```
  
   If the request is successful, the response contains information for all the auth methods.
   
   Success Response
   
   `200 OK`
   
```json
{
  "auth_methods": [
    {
      "access_key": "key",
      "host": "host2",
      "path": "path",
      "port": 9000,
      "service": "s1",
      "type": "api-key"
    },
    {
      "access_key": "key",
      "host": "host1",
      "path": "path",
      "port": 9000,
      "service": "s2",
      "type": "api-key"
    }
  ]
}
```
  Please refer to section [Errors](api_errors.md) to see all possible Errors