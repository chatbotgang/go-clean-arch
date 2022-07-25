# Development Guidelines

## Unit Test

```shell
make test
```

## Run Crescendo Barter on Local

### 1. Start Dependent Services

```shell
docker compose up -d
```

### 2. Init Local DB

```shell
make migrate-db-up
make init-local-db
```

### 3. Setup Application Configs

Some application configs have already configured in `Makefile` for running Crescendo Barter on local successfully.
These pre-configured configs provide basic functions, such as accessing DB.

For testing advanced features, we need assign related configs by ourselves, and the application configs can be assigned through:

1. Environment Variables
```shell
ENV=staging ./bin/applicatopn 
```

2. Command-Line Flags
```shell
./bin/applicatopn --env="staging"
```

Here lists the configurable application configs:
<details>
<summary> Common configs </summary>

| Env Var / Flag Var              | Description                                                             | Type    | Required | Default |
|---------------------------------|-------------------------------------------------------------------------|---------|----------|---------|
| `CB_ENV` <br> `env`             | The running environment.                                                | string  |          | staging |
| `CB_LOG_LEVEL` <br> `log_level` | Log filtering level.<br>Support error, warn, info, debug, and disabled. | string  |          | info    |
| `CB_PORT` <br> `port`           | The HTTP server port.                                                   | integer |          | 9000    |

</details>

<details>
<summary> Data systems </summary>

| Env Var / Flag Var                    | Description                                                 | Type    | Required | Default |
|---------------------------------------|-------------------------------------------------------------|---------|----------|---------|
| `CB_DATABASE_DSN` <br> `database_dsn` | The used Postgres DSN.                                      | string  | v        |         |                                            | string  |          |         |

</details>

<details>
<summary> Application Features </summary>

| Env Var / Flag Var                                                | Description                                               | Type    | Required | Default           |
|-------------------------------------------------------------------|-----------------------------------------------------------|---------|----------|-------------------|
| `CB_TOKEN_SIGNING_KEY` <br> `token_signing_key`                   | JWT Token signing key.                                    | string  |          | cb-signing-key    |
| `CB_TOKEN_ISSUER` <br> `token_issuer`                             | JWT Token issuer.                                         | string  |          | crescendo-barter  |
| `CB_TOKEN_EXPIRY_DURATION_HOUR` <br> `token_expiry_duration_hour` | JWT Token expiry hours used for customer-facing APIs.     | integer |          | 8 (8h)            |

</details>

### 4. Start Cantata on Local

```shell
make run
```

## Migrate DB Schema

1. Add migration script in `migrations/`
    * `{no.}_{description}.up.sql`: migration script
    * `{no.}_{description}.down.sql`: rollback script
2. Run `make migrate-db-up`
3. If the migration result is unexpected, run `make migrate-db-down`