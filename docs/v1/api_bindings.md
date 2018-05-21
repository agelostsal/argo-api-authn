# Binding API Calls

## [POST] Manage Bindings - Create New Binding

This request creates a new binding.

#### Request

`POST /v1/bindings`

### Example request
```
curl -X POST -H "Content-Type: application/json"
  "https://{URL}/v1/bindings?key={key_in_the_config}"
```

##### Post Body

```json
{
	"name": "b1",
	"service_uuid": "b61030d9-bef3-4768-9a03-7b1ff36e8af4cc",
	"host": "host1",
	"dn":"dn",
	"oidc_token": "token",
	"unique_key": "key"
}
 ```
 
 ### Response
 
 If the request is successful, the response contains the newly created binding.
 
 Success Response
 
 `201 CREATED`
 
 ```json
 {
     "name": "b1",
     "service_uuid": "b61030d9-bef3-4768-9a03-7b1ff36e8af4cc",
     "host": "host1",
     "dn": "host1",
     "oidc_token": "token",
     "unique_key": "key",
     "created_on": "2018-05-24T09:58:17Z"
 }
  ```
  
  ### Errors

  Please refer to section [Errors](api_errors.md) to see all possible Errors