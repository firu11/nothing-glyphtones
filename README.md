# Glyphtones
[Nothing](https://nothing.tech/) is making phones with programmable LED lights on the back.
They call it "Glyph Interface" and there is a bunch of ringtones preinstalled. Nothing also made an app called Glyph Composer,
which allows users to create new ringtones, but the options are quite limited. Users have figured out a way to create custom
ringtones and started making popular songs with matching lights.

Glyphtones is a platform, where people can either share their custom compositions, or find those they like.

![Screenshots](https://cdn.nothing.community/2025-01-04/1735987786-859251-render.png)

## Tech stack
The app uses [Go](https://go.dev/) + [echo](https://echo.labstack.com/) + [templ](https://github.com/a-h/templ) to render HTML pages for the client (and a little bit of [htmx](https://htmx.org/)).
Data is stored in a [PostgreSQL](https://www.postgresql.org/) database.

## Production
The website is running in Germany, Falkenstein on [Hetzner](https://www.hetzner.com/cloud/) VPS.

## Database Schema
![ER Diagram](database/diagram.svg)

---

## How to run (for developers)
### With Docker/Podman
1. Clone the repository
```sh
git clone https://github.com/firu11/nothing-glyphtones.git
```
2. Setup a [PostgreSQL](https://www.postgresql.org/download/) database
```sh
brew install postgresql@17
brew services start postgresql  # start postgresql
psql postgres                   # open postgresql
# --- inside postgres ---
postgres=# CREATE ROLE chris WITH LOGIN PASSWORD 'password';  # create user with password
postgres=# CREATE DATABASE glyphtones OWNER chris;            # create a database called "glyphtones" with "chris" being the owner
postgres=# \q  # exit
```
3. Build the image
```sh
docker build -t glyphtones .
```
4. Run the container
```sh
docker run --rm \
  -p 8080:8080 \
  -e DB_CONNECTION_STRING="postgresql://chris:password@localhost:5432/glyphtones?sslmode=disable" \
  -e GOOGLE_SECRET="<GOOGLE_CONSOLE_SECRET>" \
  -e GOOGLE_ID="<GOOGLE_CONSOLE_ID> \
  -e GOOGLE_REDIRECT_URL="http://localhost:8080/google-callback" \
  glyphtones
```

### Bare-metal

1. Install [Go](https://go.dev/doc/install) compiler and [PostgreSQL](https://www.postgresql.org/download/) server
2. Install [Templ](https://templ.guide/quick-start/installation) via `go install`
3. Create a new database in psql
4. Clone this repository
5. Run the _init.sql_ file to setup the database
6. Configure your enviroment variables
7. Run the project (`templ generate && go run .`)

#### MacOS example:
```sh
# INSTALLATION
brew install go                                   # go
brew install postgresql@17                        # postgresql
go install github.com/a-h/templ/cmd/templ@latest  # templ

# SETUP
brew services start postgresql  # start postgresql
psql postgres                   # open postgresql
# --- inside postgres ---
postgres=# CREATE ROLE chris WITH LOGIN PASSWORD 'password';  # create user with password
postgres=# CREATE DATABASE glyphtones OWNER chris;            # create a database called "glyphtones" with chris being the owner
postgres=# \q  # exit
# --- back in terminal ---
git clone https://github.com/firu11/nothing-glyphtones.git  # clone the repository

# CONFIGURATION
cd nothing-glyphtones       # go into the project
psql glyphtones < init.sql  # load the init.sql file into the database
# edit the configuration:
#   when following this tutorial, only 3 variables need to be changed
#   DB_NAME=glyphtones
#   DB_USER=chris
#   DB_NAME=password
# save the file

# RUN
export DB_CONNECTION_STRING="postgresql://chris:password@localhost:5432/glyphtones?sslmode=disable"
export GOOGLE_SECRET="<GOOGLE_CONSOLE_SECRET>"
export GOOGLE_ID="<GOOGLE_CONSOLE_ID>"
export GOOGLE_REDIRECT_URL="http://localhost:8080/google-callback"

templ generate  # generate html templates
go run .        # run the code
# go to: http://localhost:1323 and voilà
```
