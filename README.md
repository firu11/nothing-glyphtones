# Glyphtones
[Nothing](https://nothing.tech/) is making phones with programmable LED lights on the back.
They call it "Glyph Interface" and there is a bunch of ringtones preinstalled. Nothing also made an app called Glyph Composer,
which allows users to create new ringtones, but the options are quite limited. Users have figured out a way to create custom
ringtones and started making popular songs with matching lights.

Glyphtones is a platform, where people can either share their custom compositions, or find those they like.

![Screenshots](https://s3-nothing-prod.s3.eu-central-1.amazonaws.com/2025-01-04/1735987786-859251-render.png)

## Tech stack
The app uses [Go](https://go.dev/) + [echo](https://echo.labstack.com/) + [templ](https://github.com/a-h/templ) to render HTML pages for the client (and a little bit of [htmx](https://htmx.org/)).
This approach is called "server-side rendering". Data are stored in a [PostgreSQL](https://www.postgresql.org/) database.

## Production
The website is running in Germany, Falkenstein on [Hetzner](https://www.hetzner.com/cloud/) VPS.

---

## How to run (for developers)
1. Install [Go](https://go.dev/doc/install) compiler and [PostgreSQL](https://www.postgresql.org/download/) server
2. Install [Templ](https://templ.guide/quick-start/installation) via `go install`
3. Create a new database in psql
4. Clone this repository
5. Run the _init.sql_ file to setup the database
6. Rename _.env.example_ to _.env_ and configure your enviroment variables
7. Run the project (`templ generate && go run .`)

### MacOS example:
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
git clone https://github.com/firu11/nothing_glyphtones.git  # clone the repository

# CONFIGURATION
cd nothing_glyphtones       # go into the project
psql glyphtones < init.sql  # load the init.sql file into the database
cp .env.example .env        # duplicate .env.example -> .env
open -e .env                # open .env file
# (instead of `open -e` you can use `code`, `zed`, `vim`, `nano`)
# edit the configuration:
#   when following this tutorial, only 3 variables need to be changed
#   DB_NAME=glyphtones
#   DB_USER=chris
#   DB_NAME=password
# save the file

# RUN
templ generate  # generate html templates
go run .        # run the code
# go to: http://localhost:1323 and voilà
```
