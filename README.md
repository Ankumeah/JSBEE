# JSBEE
This is a website to share and upload research papers to

## Stack

| Component        | Stack                                     |
|------------------|-------------------------------------------|
| Primary Language | [Golang](https://go.dev/)                 |
| Web server       | [Gin](https://gin-gonic.com/)             |
| Database         | [Postgresql](https://www.postgresql.org/) |
| Cache            | [Redis](https://redis.io/)                |
| Auth             | [Firebase](https://firebase.google.com/)  |

## How to run
> [!NOTE]
> This repo requires database (postgresql), cache (redis) and auth (firebase) credentials.
These are expected to be provided by the user in .env

- Clone the repo
- Copy env.exmaple to .env and fill/replace the values to your liking
- `docker compose up`!
