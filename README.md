# Go/Templ/HTMX Example App

You will need to install the sqlc, templ and tailwindcss command line tools to use this repository.

### Helpful Links

* SQLC docs: https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html
* Templ docs: https://templ.guide/
* Tailwind CSS docs: https://tailwindcss.com/docs/installation
* HTMX docs: https://htmx.org/attributes/hx-trigger/

### Makefile Shortcuts

You use the makefile in this repository to generate code and run the api
```
make run_api
```
First though, you need to start the postgres docker container and create a database in it.
And then run the migrations. Those instructions are below.


### Create Local Database
 
Login to the postgres database running in your local docker container
```
psql -U postgres -h localhost
```

In the postgres terminal enter
```
CREATE database db_local;
```

Connect to the database
```
\c db_local
```

You can now connect to the database like this
```
psql -U postgres -h localhost -d db_local
```

### Run Migrations

To get the current migration version run
```
go run cmd/migrate/main.go version
```

To apply the migrations run
```
go run cmd/migrate/main.go upOne
```

To revert a migration run
```
go run cmd/migrate/main.go downOne
```
You will be prompted if you want to revert the migration to prevent accidental reverts.

To force a migration number run
```
go run cmd/migrate/main.go force
```
You will be prompted to enter a migration number to force to in the database.


### Build Queries

In the `core-service/core` directory run this command to generate the queries
```
sqlc generate
```

This is only required after updating SQL queries.


### Build Templ Files

In the `core-service/core` directory run this command to build the templ files
```
templ generate ./core/ui/
```

### Build Tailwind CSS

In the `core-service/core` directory run this command
```
npx tailwindcss -c ./core/ui/tailwind.config.js -i ./core/ui/main.css -o ./core/ui/static/main.css --minify
```

### Run the Server

In the `core-service` directory run this command to start the server
```
go run cmd/api/main.go
```


