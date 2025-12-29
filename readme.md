# Expense Management Systems

## Setup and Run
Dont forget to Clone this repo first

### Dependencies
This project uses:
- Go version 1.24 [[install instruction](https://go.dev/doc/install)]
- docker and docker compose for containerization [[install instruction](https://go.dev/doc/install)]
- golang-migrate for database migration [[install instruction](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)]
- swag/swaggo for API docs [[just do instruction no 2 on the getting started section](https://github.com/swaggo/swag?tab=readme-ov-file#getting-started)]

so you need to install these cli tools and language in order to run the project.


### Database Initialization & Migration
after cloning the repo and installing all dependencies, go inside the project then:
```bash
cd platform/postgres
```
this will move you to platform/postgres folder, and there will be an `example.env` file, make a `.env` file from it and fill the empty env file with credentials that you want

after that from the platform/postgres folder, do:
```bash
docker compose up -d
```

now that postgres is set, time to migrate the database, go back the project root using
```
cd ../..
```

and run 
```
migrate -database 'postgres://<POSTGRES_USER>:<POSTGRES_PASSWORD>M@localhost:5432/<POSTGRES_DB>?sslmode=disable' -path migration up
```
this will initialized all tables on the database 


### Running the API
from the project root, copy the `example.env` to `.env` and fill the env with value.
Note [1]: for DB_MASTER you can use `postgres://<POSTGRES_USER>:<POSTGRES_PASSWORD>@ems-db:5432/<POSTGRES_DB>?sslmode=disable`, same as the migration one but the host is `ems-db` instead of `localhost`
Note [2]: for PAYMENT_GATEWAY_URL use the base url of the Payment processor mock from the test **without the trailing /**

then you can run the program using docker compose
```bash
docker compose up --build -d
```


## API Docs
once the API run successfully, the API docs can access using in http://localhost:20000/docs/index.html this will open the swagger docs for the API.
