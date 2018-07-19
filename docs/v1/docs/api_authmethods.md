# Auth method API Calls

## [POST] Manage Auth Methods - Create New Auth Method

This request creates a new auth method for the given service type. The type of the auth method
as well as some of its predefined fields will be decided by the service-tye's `auth_method` and `type `fields.
E.g. for a service-type of type `ams` with an auth_method of type `api-key`, it will create an api-key auth method
with predeclared fields for `path` and `retrieval_field` that are common across all type `ams` service-types.
Of course you can always override the default's if you like.

#### Fields

- path: Combined with the `host` and the `port` is represents the URL where the external resource is located. We use it to map the x509 certificate or any other auth mechanism to the needed token
- access_key: In the case of an api-key method, the access key specifies te `key` to use in order to access the external resource
- retrieval_field: The response field from the external service that will contain the attribute we are looking for. e.g. `token`

### Request

```
POST /v1/service-types/{Name}/authm`
```


### Example request
```
curl -X POST -H "Content-Type: application/json"
  "https://{URL}/v1/service-types/{Name}/authm?key={key_in_the_config}"
```

### Post Body

```
        {
            "access_key": "key1",
            "host": "127.0.0.1",
            "port": 9000,
        }
```

### Response

If the request is successful, the response contains the newly created auth method.

Success Response

`201 CREATED`

```
        {
            "access_key": "key1",
            "host": "127.0.0.1",
            "service_uuid": "da22b2d4-ba6c-43ca-b28d-400sd0a5d83e",
            "path": "/path/{{identifier}}?key={{access_key}}",
            "port": 9000,
            "retrieval_field": "token",
            "type": "api-key",
            "uuid": "da22b2d4-8ip0-43ca-b28d-500sd0a5d876e",
            "created_on": "2018-05-05T18:04:05Z"
        }
```

Please refer to section [Errors](api_errors.md) to see all possible Errors

## [GET] Manage Auth Methods - List One Auth Method

### Request

```
GET /v1/services/{service}/hosts/{host}/authm
```

### Example request

```
  curl -X GET -H "Content-Type: application/json"
  "https://{URL}/v1/services/{service}/hosts/{host}/authm?key={key_in_the_config}"
```

If the request is successful, the response contains information for the requested auth method.

#### Success Response

`200 OK`

```
        {
            "access_key": "key1",
            "host": "127.0.0.1",
            "service_uuid": "da22b2d4-ba6c-43ca-b28d-400sd0a5d83e",
            "path": "/path/{{identifier}}?key={{access_key}}",
            "port": 9000,
            "retrieval_field": "token",
            "type": "api-key",
            "uuid": "da22b2d4-8ip0-43ca-b28d-500sd0a5d876e",
            "created_on": "2018-05-05T18:04:05Z"
        }
```
Please refer to section [Errors](api_errors.md) to see all possible Errors

## [GET] Manage Auth Methods - List All Auth Methods

### Request

```
GET /v1/authm`
```

### Example request

```
  curl -X GET -H "Content-Type: application/json"
  "https://{URL}/v1/authm?key={key_in_the_config}"
```

If the request is successful, the response contains information for all the auth methods.

#### Success Response

`200 OK`

```
{
  "auth_methods": [
        {
            "access_key": "key1",
            "host": "127.0.0.1",
            "service_uuid": "da22b2d4-ba6c-43ca-b28d-400sd0a5d83e",
            "path": "/path/{{identifier}}?key={{access_key}}",
            "port": 9000,
            "retrieval_field": "token",
            "type": "api-key",
            "uuid": "da22b2d4-8ip0-43ca-b28d-500sd0a5d876e",
            "created_on": "2018-05-05T18:04:05Z"
        },
        {
            "access_key": "key1",
            "host": "host2",
            "service_uuid": "da22b2d4-ba6c-43ca-b28d-400sd0a5d83e",
            "path": "/path/{{identifier}}?key={{access_key}}",
            "port": 9000,
            "retrieval_field": "token",
            "type": "api-key",
            "uuid": "da22b2d4-9kl2-43ca-b28d-500sd0a5d876e",
            "created_on": "2018-05-05T18:04:05Z"
        }
  ]
}
```

Please refer to section [Errors](api_errors.md) to see all possible Errors

## [PUT] Manage Auth Methods - Update an existing auth method

This request updates the auth method for the given service-type and host.
This request can update one or more fields with one call.

```
PUT /v1/service-types/{service-type}/hosts/{host}/authm
```

### Example request
```
curl -X PUT -H "Content-Type: application/json"
  "https://{URL}/v1/service-types/{Name}/authm?key={key_in_the_config}"
```

### Post Body

```
        {
            "port": 8080,
            "access_key": "key2"
        }
```

### Response

If the request is successful, the response contains the updated auth method.

Success Response

`200 OK`

```
        {
            "access_key": "key2",
            "host": "127.0.0.1",
            "service_uuid": "da22b2d4-ba6c-43ca-b28d-400sd0a5d83e",
            "path": "/path/{{identifier}}?key={{access_key}}",
            "port": 8080,
            "retrieval_field": "token",
            "type": "api-key",
            "uuid": "da22b2d4-8ip0-43ca-b28d-500sd0a5d876e",
            "created_on": "2018-05-05T18:04:05Z"
        }
```

Please refer to section [Errors](api_errors.md) to see all possible Errors

## [DELETE] Manage Auth Methods - Delete an auth method

This request deletes an auth method associated with the provided service-type and host.

### Request

```
DELETE /v1/service-types/{service-type}/hosts/{host}/authm
```

### Response

If the request is successful, the response is empty.

#### Success Response

`204 No Content`

Please refer to section [Errors](api_errors.md) to see all possible Errors
