# jsontosql

to do:
- handle passwords https://x-team.com/blog/storing-secure-passwords-with-postgresql/
- host on free oracle!
- allow for user to create users...

issues:

-----------

sudo docker exec -it jsontosql_db_1 bash
sudo docker compose run gocli

---- 

Final performance:
| Max Speed | 10 tokens/second | else randomuser.me/api throws 429 |
| :----------- | :------: | ------------: |
| 0363534365 ns/op | 10223976 B/op | 50450 allocs/op |
