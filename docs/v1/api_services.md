# Service API Calls

## [POST] Manage Services - Create New Service

This request creates a new service.

#### Request

`POST /v1/services`

### Example request
```
curl -X POST -H "Content-Type: application/json"
  "https://{URL}/v1/services?key={key_in_the_config}"
```


##### Post Body

```json
{
 	"name": "string",
 	"hosts": ["host1", "host2"],
 	"auth_types": ["x509", "oidc"],
 	"auth_method": "api-key",
 	"retrieval_field": "token"
 }
 ```
 
 ### Response
 
 If the request is successful, the response contains the newly created service.
 
 Success Response
 
 `201 CREATED`
 
 ```json
 {
  	"name": "string",
  	"hosts": ["host1", "host2"],
  	"auth_types": ["x509", "oidc"],
  	"auth_method": "api-key",
  	"retrieval_field": "token",
  	"created_on": "2018-05-05T18:04:05Z" 
  }
  ```
  
  ### Errors

  Please refer to section [Errors](api_errors.md) to see all possible Errors
  
  ## [GET] Manage Services - ListAllServices
  
  ### Request
  
  `GET/v1/services`
  
   ### Response
   
   If the request is successful, the response contains the newly created service.
   
   Success Response
   
   `200 OK`
   
   ```json
{
    "services": [
        {
            "name": "s1",
            "hosts": [
                "example.gr",
                "127.0.0.1"
            ],
            "auth_types": [
                "x509"
            ],
            "auth_method": "api-key",
            "retrieval_field": "token",
            "created_on": ""
        },
        {
            "name": "s2",
            "hosts": [
                "127.0.0.1",
                "example2.gr"
            ],
            "auth_types": [
                "x509",
                "oidc"
            ],
            "auth_method": "api-key",
            "retrieval_field": "token",
            "created_on": "2018-05-13T21:52:58Z"
        }
    ]
}
```

  ## [GET] Manage Services - ListOneService
  
  ### Request
  
  `GET/v1/services/{NAME}`
  
   If the request is successful, the response contains information for the requested service.
   
   Success Response
   
   `200 OK`
   
```json
   {
    	"name": "string",
    	"hosts": ["host1", "host2"],
    	"auth_types": ["x509", "oidc"],
    	"auth_method": "api-key",
    	"retrieval_field": "token",
    	"created_on": "2018-05-05T18:04:05Z" 
    }
```
  Please refer to section [Errors](api_errors.md) to see all possible Errors


