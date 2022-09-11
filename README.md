# Efishery-Task #
This repo is my submission for efishery backend engineer skill test. Fingers crossed!
# Table of Contents
* [Prequisites](#prequisites)
* [Installation](#installation)
* [Auth-App](#auth-app)
	* [Environment Variables](#auth-app-environment-variables)
  
## Prequisites
* Docker
* go 1.18 (optional, used to run directly in local machine)


## Installation
To set up the application, you need to have `docker` installed in your machine or `go1.18` if you want to run the app directly in your local machine. Ensure you have nothing on port **5000** (auth-app) 
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
