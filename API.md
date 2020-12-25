# Rest API Documentation

### Staging URL - `https://brank-core.herokuapp.com`

## Endpoints

- [Link](#link)

  - [Exchange](#exchange-contract-code)

- [Client](#client)

  - [Create Client](#create-client)

- [Client Application](#client-application)
  - [Create App](#create-app)

### Link

#### Exchange Contract Code

##### Headers

```json
{
  "brank-access-token": "beu3b4y34vewvd39e3esudbsdusdibsbdbd"
}
```

##### Params

```json
{
  "code": "230284",
  "public_key": "in3beu3breibruowgh4hqfshaihnhfhxal"
}
```

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

### Client

#### Create Client

##### Headers

```json
{
  "brank-access-token": "beu3b4y34vewvd39e3esudbsdusdibsbdbd"
}
```

##### Params

```json
{
  "first_name": "Griffith",
  "last_name": "Awuah",
  "email": "gpa20@gmail.com",
  "password": "gwuah",
  "company_name": "brank"
}
```

##### Result

```json
{
  "data": {
    "client": {
      "id": 3,
      "created_at": "2020-12-25T08:44:52.156626Z",
      "updated_at": "2020-12-25T08:44:52.156626Z",
      "first_name": "Griffith",
      "last_name": "Awuah",
      "email": "gpa220@gmail.com",
      "password": "$2a$10$d3EB/vSWwPIs50WXyGFG8.pIwT/p0IwH1czGMK//8mH19RkWwotKa",
      "company_name": "brank",
      "verified": null
    },
    "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhcHBfaWQiOjAsImNsaWVudF9pZCI6MywiZXhwIjoxNjE2NjYxODkyLCJpYXQiOjE2MDg4ODU4OTJ9.3icXWREwdQnsrwWNE4RAXEdyynvXqv4EURkLrtb3zDCmiUZedUJUI-H4ohZxyqL5bA6wQs78x2_C_AWWeU9nLw"
  },
  "message": "client created successfully"
}
```

### Client Application

#### Create App

##### Headers

```json
{
  "brank-access-token": "beu3b4y34vewvd39e3esudbsdusdibsbdbd"
}
```

##### Params

```json
{
  "name": "Float",
  "logo": "https://avatars3.githubusercontent.com/u/24861123?s=60&u=c4496932b4a839c2452eabf5a001e4954666b410&v=4",
  "callback_url": "https://avatars3.githubusercontent.com",
  "client_id": 2
}
```

##### Result

```json
{
  "data": {
    "app": {
      "id": 7,
      "created_at": "2020-12-25T08:46:34.314716Z",
      "updated_at": "2020-12-25T08:46:34.322571Z",
      "public_key": "FLOAT-913557232050",
      "name": "float",
      "logo": "https://avatars3.githubusercontent.com/u/24861123?s=60&u=c4496932b4a839c2452eabf5a001e4954666b410&v=4",
      "callback_url": "https://avatars3.githubusercontent.com",
      "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhcHBfaWQiOjcsImNsaWVudF9pZCI6MiwiaWF0IjoxNjA4ODg1OTk0fQ.CijN-SE1uBSwKtYZbltcD00T7Fx2CkABLBOh4_4Ty0p6SqOegBaS6iyoejz9_jvCoDVLN-PswZMOxy2gD-DWVA",
      "client": 2
    }
  },
  "message": "client application created successfully"
}
```
