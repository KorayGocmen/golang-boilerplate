# API

# Architecture

Implements a modified version of the "Onion Architecture".
Layers are as follows:

1. Transport (http)
2. Service (application logic)
3. Repo (database interface)
4. DB (sql)

Each layer requires the layer belows it and written in a way that it's dependencies can be mocked for unit testing.

### File Structure

API follows the Golang Standards application layout described here: https://github.com/golang-standards/project-layout

### Config

API only uses env vars for application configuration to enforce best practices in order to not commit application secrets or have them in plain text.

---

# Database

## Local

### Setting up database

In psql:

Create the database

```
create database boilerplate;
```

Create the admin user

```
create user boilerplate_admin with encrypted password 'password';
grant all privileges on database boilerplate to boilerplate_admin;
alter database boilerplate owner to boilerplate_admin;
```

Set up your .env file.

In terminal:

```
make build db_reset db_up db_seed
```
