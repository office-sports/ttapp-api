# https://github.com/office-sports/ttapp-api
### API for TTAPP (table tennis scoring app)

The TTAPP web application consists of two parts:
- API, current repository, backend for web frontend (Golang + MySQL DB)
- frontend (Vue) https://github.com/office-sports/ttapp-frontend

### Prerequisites:

1. Go language 

Follow these instructions to install and run go:
https://go.dev/doc/tutorial/getting-started

2. MySQL database instance with database schema.

empty database schema can be found in `config/empty_schema.sql`, 
minimal schema for example tournament with single match in `config/minimal_schema.sql` 

## Running API

Having Go installed and database schema imported, minimal setup requires these values to be set in `config/config.yaml`:

```
db:
  address: MYSQL_DB_HOST
  port: MYSQL_DB_PORT
  username: MYSQL_USER
  password: MYSQL_PASS
  schema: MYSQL_SCHEMA_NAME
api:
  port: API_PORT
```

Any other config values are bound to Slack integration with the TTAPP.

Once you have go running, you can run the backend API with `go run api.go`. 
The API will be exposed to `ttapp-frontend` on port defined in config file's `API_PORT`  

