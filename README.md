# Netlify Technical Interview Project

##### _Candidate: Lina Shadrach_

This work was done for the Netlify interview process. The goal was to get running, fix, and improve a simple server written in Go that depends on a Postgres database. The [original program](https://gitlab.com/linashadrach/server/tree/7e1c1c69ea9ad335b1c77ee9871a476630c298a0) was provided by Netlify. 
* Click [here](https://gitlab.com/linashadrach/server/blob/master/technical-inteview-instructions.md) to view the full instructions.
* Click [here](https://gitlab.com/linashadrach/server/blob/master/work-description.md) to see a description of the work completed

### What it does

* The program can be used to query a dataset, where each entry in the dataset has `id`, `name`, and `count` properties. 
* Three different queries are handled:
    * `/count` returns the sum of all of the entries `count` property
    * `/count?id=INT` returns the sum of all of the entries' `count` property where the entries' `id` property is a provided integer
    * `/count?name=STR` returns the sum of all of the entries' `count` property where the entries' `name` property is a provided string

## Installation and Configuration

* Make sure you have the following tools installed properly:
    * Go [Download here](https://golang.org/dl/). _This project was built with version go1.11.1 linux/amd64_
    * Postgres [Download here](https://www.postgresql.org/download/) _This project was built with version 9.5_
    
* [Download this project's source code](https://gitlab.com/linashadrach/server)

#### Database Setup

* Start postgres.
* Use the following psql commands for recreating the database schema
    ```
    CREATE DATABASE counter;
    \c counter;
    create table viewers (id serial primary key,name varchar(256) not null,count bigint default 0);
    ```

    * optional: use the following to add sample data:
        ```
        insert into viewers (id, name, count) values (2, 'ryan', 10),(3, 'aaron', 5),(4, 'sue', 200),(5, 'amanda', 50),(6, 'emily', 1),(7, 'sundance', 123);
        ```

#### _config.json_ Setup

* Create a file in top level of project named _config.json_. Supply the port number where the server should run and the database connection string in the following format: 
    ```
    {
      "port" : [CHOSEN PORT NUMBER],
      "db_conn_str" : [DATABASE CONNECTION STRING]
    }
    ```
Check out this documentation on [the Go package pq](https://godoc.org/github.com/lib/pq) for the correct formatting and examples of connection strings. 

#### Makefile

* Run `> make` in top level of project
    * Program will be installed in `/usr/local/bin/`


### Run the Program

* In the directory containing sever.go, use `> go run server.go config.json` to run the project
* After running `> make`, the project can be run from installed location:
    * Use `> /usr/local/bin/server config.json`
* Make HTTP requests to `/count` using the port number supplied in _config.json_

### Things to add

* Test coverage.
* Logging file
* Custom Error handler

# License

[MIT License](LICENSE).
[License Info](https://writing.kemitchell.com/2016/09/21/MIT-License-Line-by-Line.html)


# Thank you!

##### Contact
_The author can be reached through GH at github.com/**linashadrach**_
