# How to Deploy

> [!NOTE]
> Make sure to read docker-compose.yaml
> to make sure you understand what is going on

## What you need

- A Firebase project with **Google sign in enabled
- An S3 compatible object store
- A domain
- A VPS / server
- A `.env` file (see below)

## Firebase project
- Create a new firebase project and app
- Copy your firebase config and condense it to one line
- Set this as your `FIREBASE_CLIENT_CONFIG` in `.env`
- Add a new service account
- Generate a private key
- Copy the entire key and condense it to one line
- Set this as your `FIREBASE_CREDENTIALS` in `.env`
- Go to Authentication and enable Google sign in
- Open Google Cloud Console
- Go to credentials
- Click on the OAuth 2.0 Client ID
- Add your domain to the Authorised domains and redirect URI

## 1. `.env`

```sh
cp env.example .env
```

Fill in `.env`

| Env Var                  | Notes                                            |
|--------------------------|--------------------------------------------------|
| `API_VERSION`            | Api version, try to not to change this           |
| `BACKEND_PORT`                   | Port which backend runs on                       |
| `HOSTNAME`               | URL to the backend                               |
| `REVERSE_PROXY_PORT`     | Port which reverse proxy runs on                 |
| `DB_URL`                 | URL of the db, also supports some config options |
| `FIREBASE_CREDENTIALS`   | Firebase auth credentials                        |
| `FIREBASE_CLIENT_CONFIG` | Firebase OAuth 2.0 Client ID                     |
| `OBJECT_STORE_CONFIG`    | Look at `env.example` for more info              |
| `FRONTEND_SAVE_DIR`      | Dir to save static files to be served            |

> Look at `env.example` for more info

## 2. `backend/internal/provider/values.go`
This file conatins some provider specific values that
aren't exactly secrets. Make sure to set them to your own
values

## 3. Deploy

We have 2 options for deployment

- ### 1. Docker compose
  ```sh
  docker compose up --build
  ```

- ### 2. Single docker container
  Either use our provided docker image
  ```sh
  docker pull ankumeah/jsbee
  ```
  Or build it yourself
  ```sh
  docker build -t ankumeah/jsbee:latest .
  ```

  And then run with
  ```sh
  docker run --env-file .env -v ./data:/data/ -p 8000:8000 ankumeah/jsbee:latest
  ```

Open `http://<host>:8000`.

## 4. First run setup

- New sign ups default to the `viewer` role. Promote yourself directly in SQLite:
  ```sh
  sqlite3 data/jsbee.db "UPDATE users SET role = 'owner' WHERE email = 'your email';"
  ```
