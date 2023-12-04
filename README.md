# jsontosql

to do:
- add dataschema to postgres
- add user input
- add main loop func asking what user wants to do, 1 adds users, 2 checks current db state...
- handle passwords https://x-team.com/blog/storing-secure-passwords-with-postgresql/


issues:
- // panic: Get "https://dummyjson.com/users": tls: failed to verify certificate: x509: certificate signed by unknown authority
    - overcame by making a client with InsecureSkipVerify
    - exploring more appropriate solutions...
- docker-compose up just hangs when i type to cli, but docker-compose run gocli works fine


-----------

sudo docker exec -it json_to_sql_db_1 bash -U
