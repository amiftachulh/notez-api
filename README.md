# Migration

## Install golang-migrate

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

To use the tool, make sure to add the `bin` directory to your `PATH`.

## Create a new migration

Make sure to run the script from the root of the project.

```bash
migrate create -ext sql -dir db/migrations <migration_name>
```

## Run migrations

Make sure to run the script from the root of the project.

### Up

```bash
migrate -path db/migrations -database "<database_url>" -verbose up
```

### Down

```bash
migrate -path db/migrations -database "<database_url>" -verbose down
```

### Force

```bash
migrate -path db/migrations -database "<database_url>" -verbose force <version>
```
