# Go Api
This is an API built with Go <code>v1.20</code>

### Setup
```
cp docker/local/.env.example docker/local/.env
add UID and GID in .env
make docker/build
make docker/up
make db/migrate/up
```
### Database - PostgreSQL
In this application we make use of PostgreSQL v.13.<br/>
We are using [migrate](https://github.com/golang-migrate/migrate) for Database migrations.<br/>
About our database setup you can read more here:
[Read More](./docs/postgreSQL.md)

### Emails
In order to set up  emails you need to mailgun on any other provider for development.
Configure the following environmental variables:
```
SMTP_HOST=""
SMTP_PORT=""
SMTP_USERNAME=""
SMTP_PASSWORD=""
SMTP_SENDER=""
```
