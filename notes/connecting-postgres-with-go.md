## Import driver

> Import a DB driver. We will use github.com/jackc/pgx/v4 as Postgres Driver.

## Connect to DB

> sql.Open() function is used to connect to a SQL-like DB.

```go
    // Parameters are set according to docker-compose.yaml
    db, err := sql.Open("pgx", "host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable")

    // Check if the DB is up and running by pinging it.
    err = db.Ping()
    if err != nil {
        panic(err)
    }
```

> If you change docker-compose.yaml file, restart images by 'docker compose up & down'
