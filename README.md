# JSBEE
JSBEE (Journal of Sustainable Businesses in Emerging Economies) is a website to share and read research papers at

## Stack

| Component        | Stack                                     |
|------------------|-------------------------------------------|
| Primary Language | [Golang](https://go.dev/)                 |
| Web server       | [Gin](https://gin-gonic.com/)             |
| Database         | [SQLite](https://sqlite.org/)             |
| Reverse Proxy    | [Nginx](https://nginx.org/)               |
| Auth             | [Firebase](https://firebase.google.com/)  |

## How to run
- Clone the repo
- Read `DEPLOY.md`
- Copy env.exmaple to .env and fill/replace the values to your liking
- `docker compose up`!

## Contributing
- Follow commit prefix conventions
- Give commits clear messages, use commit bodies for longer messages
- AI assisted is allowed, Slop written isen't
- Make regular use of comments
- Make sure to either use the githooks provided (`./githooks/`) or manually run them before commting
- Every package within `./backend/internal/` contains a README.md, make sure to read those to get details for every package
