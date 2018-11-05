Context

There is a simple webserver running on the provided instance. It should be responding to web requests but it is not. Once we get it fixed, we will need to address some features and bugs. 

Part 1: Get it running

I'd like you to login, investigate and fix the service. Please document all the things that you change, what you find, and your general process. 

Some information about the instance:

- The user-facing service should be listening on port 80
- The web server
  - is in go
  - the source code lives in `/opt/app``s``/sillyserver/code/`
  - is installed as `/usr/local/bin/server`
  - GET requests to `/count` will return the total number of views
  - GET requests to `/count?id=<id>` return the total number of views for that one user
  - there is a Makefile to help with the build and deploy capabilities
- the go tool chain is appropriately installed on the machine 
- postgres is installed on the machine
  - the database that had the data is `server` 
  - the username is `server_read` and the password is `super_secret_value`
  - it has this data:
    > select * from viewers;
    +------+----------+---------+
    | id   | name     | count   |
    |------+----------+---------|
    | 2    | ryan     | 10      |
    | 3    | aaron    | 5       |
    | 4    | sue      | 200     |
    | 5    | amanda   | 50      |
    | 6    | emily    | 1       |
    | 7    | sundance | 123     |
    +------+----------+---------+
- there are no external blocks to accessing the machine

Evaluation Criteria
Part of the problem is if you can get the service running, the other is the documentation about the things that you fixed. The goal is to really see if you can debug a not running service and talk about what was amiss. 

If you get stuck, please reach out in the channel and we can discuss what is going on! 


Part 2: bug fix

Now that the service is running, there seems to be a bug. When I curl `/count` I get the right amount back (389), but if I query `/count?id=3` I get back a 404. I'd expect to get back a 200 and the value 5 (the value for the user `aaron`).

- the code lives in `/opt/apps/sillyserver/code/`
- you can replace the binary at `/usr/local/bin/server`

Please investigate, write up the bug, and fix and deploy the service. 


Part 3: service improvements

Now that the service can be queried by an `id` field, people want to be able to search by `name`. Please extend the server to detect and respond to the query parameter `name`. For instance, when you visit `/count?name=sundance` you should get back 123. 