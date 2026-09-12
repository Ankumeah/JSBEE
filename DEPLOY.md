# How to Deploy

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
| `PORT`                   | Port which backend runs on                       |
| `HOSTNAME`               | URL to the backend                               |
| `REVERSE_PROXY_PORT`     | Port which reverse proxy runs on                 |
| `CACHE_URL`              | URL to cache. Currently unused. Any value works  |
| `DB_URL`                 | URL of the db, also supports some config options |
| `FIREBASE_CREDENTIALS`   | Firebase auth credentials                        |
| `FIREBASE_CLIENT_CONFIG` | Firebase OAuth 2.0 Client ID                     |
| `OBJECT_STORE_CONFIG`    | look at `env.example` for more info              |
| `FRONTEND_SAVE_DIR`      | Dir to save static files to be served            |

> Look at `env.example` for more info

## 2. Deploy

```sh
docker compose up --build
```

Open `http://<host>:8000`.

## 3. First run setup

- New sign ups default to the `viewer` role. Promote yourself directly in SQLite:
  ```sh
  sqlite3 data/jsbee.db "UPDATE users SET role = 'owner' WHERE email = 'your email';"
  ```
- Edit the About page (Admin → About), the seed is placeholder text.
