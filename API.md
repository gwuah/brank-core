# Rest API docs for brank

### Staging URL - `https://brank-core.herokuapp.com`

## Endpoints

- [Link](#linkservice)

  - [Exchange](#exchange-contract-code)
  
### Link

#### Exchange Contract Code

##### Headers 
```json
{
  "brank-access-token": "beu3b4y34vewvd39e3esudbsdusdibsbdbd",
}
```

##### Params

```json
{
  "code": "230284",
  "public_key": "in3beu3breibruowgh4hqfshaihnhfhxal"
}
```

| Key      | Type  | Remark   |
| -------- | ----- | -------- |
| `code` | `string` | Required |
| `public_key`  | `string` | Required |

##### Result

```json
{
    "data": {
        "link": {
            "id": 2,
            "created_at": "2020-12-25T08:08:54.880073Z",
            "updated_at": "2020-12-25T08:08:54.880073Z",
            "raw": null,
            "code": "24423027",
            "bank_id": 1,
            "app_id": 1,
            "username": "banku",
            "password": "stew"
        }
    },
    "message": "exchange was successful"
}
```

