# Authenticate Using an x509 certificate

## [GET] Authenticate Via x509

This request will use the provided x509 certificate in order to retrieve
a token from the given service type.

### Example Request

```
curl -X GET -H "Content-Type: application/json"
  "https://{URL}/v1/service-types/{Name}/hosts/{host}:authX509" 
   --cert /path/to/a/cert/file --key /path/to/the/respective/key -k
```

 ### Response
 
 If the request is successful, the response contains the token that is associated with the provided certificate.
 
 Success Response
 
 `200 CREATED`
 
 ```json
 {
    "token": "some-service-type-token"
 }
  ```