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
once the API is running, the API docs can access using in http://localhost:20000/docs/index.html this will open the swagger docs for the API.

## Architecture decisions and trade-offs 
Pretty much follow the data models that the test gave me:
```
Users (id, email, name, role, created_at)
Expenses (id, user_id, amount_idr, description, receipt_url, status, submitted_at, processed_at)
Approvals (id, expense_id, approver_id, status, notes, created_at)
```
While I pretty much follow this data model to the tee. I do have some opinion on how this can be improved (maybe) which i will cover in the **Assumptions & what to improve** sections


## Business rules implementation explanation
The business rules I do is pretty much the Access Control and the auto approve Threshold

### Access Control
For the access control, there are checks on routes layer and services layer.
on the routes layer I created a RBAC middleware to block user that is not fulfill the required role.

The other checks is in service layer. This one i use simple checks if the user is manager then, show all expense if not show only user's expenses.
**example line in service/expensesv/expense.go line 58**
```go
userData := ctxutils.GetUserDataFromCtx(ctx)
if !userData.IsManager() {
	userId = &userData.UserId
}
```

### Auto-approve Threshold
This is the same as access control on service level, its a simple checks if this is an auto approved case or not.
**example line in service/expensesv/expense.go line 79**
```go
status := "pending"
now := timeNow()
processedAt := time.Time{} // zero value

isAutoApproved := req.IsAutoApproved()
if isAutoApproved {
    status = "auto-approved"
    isAutoApproved = true
    processedAt = now
}
```

## Assumptions & what to improve
there are some assumptions I made along the way, and I had it written down in my notes, So I'm just going to copy that here and add some more:
- there is only 2 role: user and manager
- expenses status : rejected, pending, approved, auto-approved
- completed will be on separate column
- login only by email
- completed status is a bit confusing (how are we going to filter by approved and auto-approved if in the end the status will be completed) and since I don't think there any use of it on the FE, I'm going ignore it for now
- approvals status : rejected, approved, auto-approved

as you can see the completed status is a bit confusing to me, so i ignore it. 
there is some solutions that I've thought of but there also some trade offs to it, for example:

1. We can store status on both expenses and approvals instead just expenses and update only expense status to "complete" once the payment is done. But in turn **for each list request with (or without if not handled) filter you would need to join expenses and approvals** to find if its on pending, rejected, approved or auto-approved

2. I think its better to just create a new column called **is_payment_complete as a bool or a vartext(50) or something like that**. Since I don't think the value have any use for the test in the frontend side (cmiiw), but I do understand the need to at least store it. The trade-off of this probably more data to store. Imho, I would prefer this approach.

as to what other things I would do if i had more time, probably refactoring the router layer, I feel it can be split up a bit more, and also I haven't add indexing yet to the tabels since I didn't know what indexes I need when starting the project.

Other things that I can do is maybe add more unit tests, on the services level especially, and create mocks for the repos.
