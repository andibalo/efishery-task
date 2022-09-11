# Efishery-Task #
This repo is my submission for efishery backend engineer skill test. Fingers crossed!
# Table of Contents
* [Prequisites](#prequisites)
* [Installation](#installation)
* [Postman Collection](#postman-collection)
* [Auth-App](#auth-app)
	* [Environment Variables](#auth-app-environment-variables)
	* [Endpoints](#auth-app-api)
* [Fetch-App](#fetch-app)
	* [Environment Variables](#fetch-app-environment-variables)
	* [Endpoints](#fetch-app-api)
	
## Prequisites
* Docker
* go 1.18 (optional, used to run directly in local machine)

## Postman Collection
The postman collection json for this monorepo can be found at `/postman`. You can download it and import it on your local postman application.

## Installation
To run `auth-app` and `fetch-app`, you need to have `docker` installed in your machine or `go1.18` if you want to run the app directly in your local machine. Ensure you have nothing on port **5000** (auth-app) 
or **5030** (fetch-app) running. Then go to root directory and run the following commands.

```
$ cd infra
$ docker-compose up
```

If you are facing issue, try running docker-compose with sudo permission
```
$ sudo docker-compose up --build
```

Otherwise you have started the `auth-app` and `fetch-app` and can start testing out their APIs

## Auth-App
`auth-app` is responsible for managing user creation, password generation, and JWT generation. It is using a file-based database.

### Auth-App Environment Variables
* `SERVER_PORT` defines which port the application will listen to 
	* The default port is `5000`
* `FILE_PATH` defines the path to the filed used in storing user data
* `JWT_SECRET` defines the secret used to sign the JWT Token

### Auth-App API
  - [Health Check](#get-health---health-check)
  - [Create User](#post-v1user---create-user)
  - [Login User](#post-v1userlogin---login-user)
  - [Get User Token Details](#post-v1userdetails---get-user-details-from-token)
  
### GET /health - Health Check
This endpoint can be used to verify that the app is running

### Response
```
ok
```

### POST /v1/user/ - Create User
This endpoint will insert user data into the database and will return the phone number 4 characters password for that user. 

### Request
```
//HTTP Request Body (JSON)
{
	"name": "andi",
	"phone": "0895132442",
	"role": "admin"
}
```
### Response
```
//HTTP Response (Application/JSON)
{
    "code": "AU0000",
    "status": "SUCCESS",
    "message": "Successfully created user",
    "data": {
        "name": "andi",
        "phone": "0895132442",
        "role": "admin",
        "password": "Hh7.",
        "timestamp": "11 Sep 22 18:11 UTC"
    }
}
```

### POST /v1/user/login - Login User
This endpoint will receive a `phone` and `password` and return a generated JWT that contains `name`, `phone`, `role`, and `timestamp` of the user that has the matching `phone` and `password`.

### Request
```
//HTTP Request Body (JSON)
{
	"phone": "0895132442",
	"password": "<4 characters string>"
}
```

### Response
```
//HTTP Response (Application/JSON)
{
    "code": "AU0000",
    "status": "SUCCESS",
    "message": "Successfully logged in",
    "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYW5kaSIsInBob25lIjoiMDg5MzgzMTExMTgiLCJyb2xlIjoiYWRtaW4iLCJ0aW1lc3RhbXAiOiIxMSBTZXAgMjIgMTE6MzAgV0lCIn0.qutZ-QZ509SCS1eT4UK1olmfYe1tvRODU3uQF5CezjU"
}
```

### POST /v1/user/details - Get User Details From Token
This endpoint will receive a `token` which is a JWT token return the data contained inside the token

### Request
```
//HTTP Request Body (JSON)
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYW5kaSIsInBob25lIjoiMDg5MzgzMTExMTgiLCJyb2xlIjoiYWRtaW4iLCJ0aW1lc3RhbXAiOiIxMSBTZXAgMjIgMTE6MzAgV0lCIn0.qutZ-QZ509SCS1eT4UK1olmfYe1tvRODU3uQF5CezjU"
}
```

### Response
```
//HTTP Response (Application/JSON)
{
    "code": "AU0000",
    "status": "SUCCESS",
    "message": "Successfully get token details",
    "data": {
        "name": "andi",
        "phone": "08938311118",
        "role": "admin",
        "timestamp": "11 Sep 22 11:30 WIB"
    }
}
```

## Fetch-App
`fetch-app` is responsible for resources from efishery endpoint. It also uses an in-memory caching to cache conversion rate from https://www.exchangerate-api.com.

### Fetch-App Environment Variables
* `SERVER_PORT` defines which port the application will listen to 
	* The default port is `5030`
* `CURRENCY_SERVICE_API_KEY` defines the api key to be used when calling the currency conversion api from https://www.exchangerate-api.com
* `JWT_SECRET` defines the secret used to verify/parse the JWT Token

### Fetch-App API
- `Protected` API requires a valid jwt token in the `Authorization` field in the request header
- `Admin-only` API requires a valid jwt token in the `Authorization` field in the request header as well as an admin role

  - [Health Check](#get-health---health-check)
  - [Get Commodities (Protected)](#get-v1commodity---get-commodities-protected)
  
### GET /health - Health Check
This endpoint can be used to verify that the app is running

### Response
```
ok
```

### GET /v1/commodity - Get Commodities (Protected)
This endpoint will fetch commodities data from https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list, and fetch the IDR to USD Conversion rate from https://www.exchangerate-api.com to calculate the original price in USD and add the field `price_usd` with the converted price
```
//HTTP Request Header (Bearer Token)
Authorization: Bearer <JWT>
//HTTP Response (Application/JSON)
{
    "code": "FE0000",
    "status": "SUCCESS",
    "message": "Successfully get commodities",
    "data": [
        {
            "uuid": "",
            "komoditas": "",
            "area_provinsi": "",
            "area_kota": "",
            "size": "",
            "price": "",
            "tgl_parsed": "",
            "timestamp": "",
            "price_usd": ""
        },
        {
            "uuid": "d95f61ee-6f76-408f-88c0-8cd8b95d1daf",
            "komoditas": "PATIN",
            "area_provinsi": "JAWA TIMUR",
            "area_kota": "JEMBER",
            "size": "150",
            "price": "10000",
            "tgl_parsed": "2022-01-06T21:21:39Z",
            "timestamp": "1641504099032",
            "price_usd": "0.676400"
        },
        {
            "uuid": "79d96d67-ebfd-44ac-9a81-efc65016188d",
            "komoditas": "LELE",
            "area_provinsi": "JAWA BARAT",
            "area_kota": "CILILIN",
            "size": "150",
            "price": "64000",
            "tgl_parsed": "2022-01-07T00:20:35Z",
            "timestamp": "1641514835520",
            "price_usd": "4.328960"
        },
	...
}
```
