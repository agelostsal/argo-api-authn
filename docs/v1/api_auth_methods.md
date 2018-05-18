  # Auth method API Calls

  
  ## [GET] Manage Services - ListOneAuthMethod
  
  ### Request
  
  `GET/v1/services/{service}/hosts/{host}/authM`
  
   If the request is successful, the response contains information for the requested auth method.
   
   Success Response
   
   `200 OK`
   
```json
{
    "_id": "5af36ae263927f878860c3f0",
    "access_key": "b328c3861f061f87cbd34cf34f36ba2ae20883a5",
    "host": "127.0.0.1",
    "path": "/v1/users:byUUID/{{identifier}}?key={{access_key}}",
    "port": 8081,
    "service": "ams",
    "type": "api-key"
}
```


  ## [GET] Manage Services - ListAllAuthMethods
  
  ### Request
  
  `GET/v1/services/{service}/hosts/{host}/authM`
  
   If the request is successful, the response contains information for the requested auth method.
   
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