## How to connect to db

```
make docker-up
connect to db container
psql -U housi -d housi_db
```

#### Alternative way to connect to db (DSN)
```
psql $DB_DSN
```

## Install citext extension
This adds a case-insensitive character string type to PostgreSQL, which we will use later in the book to store user email addresses.
```
CREATE EXTENSION IF NOT EXISTS citext;
```

#### Generate suggested values based on your available system hardware.
https://pgtune.leopard.in.ua/