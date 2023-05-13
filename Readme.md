# Core Authentication

## Introduction

The authentication service. It focuses on providing oauth2
authentication and other helper features.

### Features

- Oauth2 authentication
- Email verification

### Tech Stack

- Go 1.18 (Gin + Gorm)
- PostgreSQL
- Redis
- Docker

## Prerequisite  

To use this service, follow the instructions below:

- API responses from your service to this service should follow

  ```json
    {
        "message": "Successful",
        "Data": {
            "id" : "idxxx",
            "email": "xxx@email.com", 
            "profile_id": "xxx",
            "groups": [],
            "permissions": [],
            "tenantIds": [], 
        }
    }
  ```

## Local Setup

- Make a copy of `env.example.com` to `env`
- Install [Go](https://go.dev/doc/install)
- Install [PostgreSQL](https://www.postgresql.org/download/)

## docker Setup

- Make a copy of `env.example.com` to `env`
- Install [Docker](https://docs.docker.com/engine/install/)
- Run `make start`

## Improvement