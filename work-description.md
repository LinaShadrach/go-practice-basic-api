# Netlify Technical Challenge

_Below is a description of the work I did_

## Part 1: Getting the server up and running:

#### Get postgres running

_The process outlined below is specific to running postgres on linux_

* used `service postgresql start` to start db server
* used `service postgresql status` to check status, saw it was 'active'
* used `lsof -nP | grep LISTEN` to check which port it was listening at, saw no entry for 'postgres'
* configured postgres to listen to the correct address by editting _/etc/postgresql/9.5/main/postgres.conf_
    * set `listen_addresses` to my instance of the server (`142.93.30.134`)
* restarted postgres after changing configuration file, confirmed that it was listening at '142.93.30.134:5432'
* check db contents
    * ran `psql server server_read` as postgres user
    * in psql, ran `SELECT * from viewers`, saw expected data
* added temporary hard-coded db connection string to _server.go_ to test whether connection could be established
* ran _server.go_ received error `Failed to ping the DB: pq: no pg_hba.conf entry for host "142.93.30.134", user "server_read", database "server", SSL on`
* added entry to _/etc/postgresql/9.5/main/pg\_hba.conf_ to allow all users with encrypted passwords to use all databases
    * entry: 
        ```
        # TYPE  DATABASE        USER            ADDRESS                 METHOD
        host    all             all             0.0.0.0/0               md5
        ```
* restarted postgres and ran _server.go_
* used curl to send GET request to `/count`, got expected value

#### Add configuration to server

* created _config.json_, added the port and db connection string
    * replaced temporary hard-coded data (see above) with input from config file
    * switched config.Port to config.PORT for consistency
    * switched config.DBURL to config.DBCONNSTR to be more general
* ran _server.go_ and passed in _config.json_, queried `/count`, got expected value

#### Fix Makefile

* opened Makefile, noticed commands were not in sequence
* switched order of commands (moved `deploy`  above `build`)
* changed `go build` output to _src/server_ in build task, change source for _server_ in deploy task
* added command for copying config file
* run `> make`

## Part 2: Bug Fix

* updated SQL query to use `id` parameter instead of `name`
* ran _server.go_, queried `/count?id=3`, got expected value

## Part 3: Service Improvements

* add function `handleNameRequest()` to process requests to `/count?name=`
* add conditional to `main()` to invoke `handleNameRequest()` when name parameter supplied to `/count` endpoint
* ran _server.go_, queried `/count?name=sundance`, got expected value

_NOTE: I've started a branch called refactor where I've played around with making the code more resusable. I've never programmed in Go, so I took the opportunity to get to know the language_