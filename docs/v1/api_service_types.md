# Service API Calls

## [POST] Manage Service Types - Create New Service Type

This request creates a new service type.

#### Request

`POST /v1/service-types`

### Example request
```
curl -X POST -H "Content-Type: application/json"
  "https://{URL}/v1/service-types?key={key_in_the_config}"
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
 
 If the request is successful, the response contains the newly created service type.
 
 Success Response
 
 `201 CREATED`
 
 ```json
 {
  	"name": "string",
  	"hosts": ["host1", "host2"],
  	"auth_types": ["x509", "oidc"],
  	"auth_method": "api-key",
  	"uuid": "da22b2d4-ba6c-43ca-b28d-400cd0a5d83e",
  	"retrieval_field": "token",
  	"created_on": "2018-05-05T18:04:05Z" 
  }
  ```
  
  ### Errors

  Please refer to section [Errors](api_errors.md) to see all possible Errors
  
  ## [GET] Manage Service Types - ListAllServiceTypes
  
  ### Request
  
  `GET/v1/service-types`
  
   ### Response
   
   If the request is successful, the response contains all the service types.
   
   Success Response
   
   `200 OK`
   
   ```json
{
    "service_types": [
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
            "uuid": "da22b2d4-ba6c-43ca-b28d-400cd0a5d83e",
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
            "uuid": "da22b2d4-ba6c-43ca-b28d-400sd0a5d83e",
            "retrieval_field": "token",
            "created_on": "2018-05-13T21:52:58Z"
        }
    ]
}
```

  ## [GET] Manage Service Types - ListOneServiceType
  
  ### Request
  
  `GET/v1/service-types/{NAME}`
  
   If the request is successful, the response contains information for the requested service type.
   
   Success Response
   
   `200 OK`
   
```json
   {
    	"name": "string",
    	"hosts": ["host1", "host2"],
    	"auth_types": ["x509", "oidc"],
    	"auth_method": "api-key",
    	"uuid": "da22b2d4-ba6c-43ca-b28d-400cd0a5d83e",
    	"retrieval_field": "token",
    	"created_on": "2018-05-05T18:04:05Z" 
    }
```
  Please refer to section [Errors](api_errors.md) to see all possible Errors


