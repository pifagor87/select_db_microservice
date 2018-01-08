# A simple select Postgresql DB implemented using pgx.
For DB works, we use a lightweight high-performance DB driver pgx.
To start using "select_db_microservice" you need:
* Install and configure the latest version of Golang
* Be sure to modify the data access data in the connect_db.json to more complex.
* Move the connect_db.json file to a secure location on the server. Provide the necessary access rights.

    Add const AccessDbPatch string = "./connect_db.json". It is access connect DB.
   In file main.go, specify the correct path to the file connect_db.json. Replace "./connect_db.json" with the desired path to the file.
* go get github.com/pifagor87/select_db_microservice
* To use "select_db_microservice", add in import - "github.com/pifagor87/select_db_microservice".

    When constructing your own microservice, use them in the following way, for example:

    Use select_db_microservice.SelectDb(AccessDbPatch)

    Structure POST data:

    key = "data"

    Value, for example:
    {"tables":{"origin":{"table":"your_table","alias":"p"},"join":[{"table":"your_table_join","name":"LEFT JOIN","alias":"ps","left":"p.id","right":"ps.id"}]},"filters":{"and":[{"column":"p.id","val":["your_n"],"operator":">"},{"column":"p.id","val":["your_n1"],"operator":"<"}],"or":[{"column":"p.url","val":["your_val"],"operator":"ILIKE"}]},"fields":["p.id", "p.url", "p.status", "p.created"],"params":{"order":{"fields":["p.id"],"sort":["DESC"]},"group":["p.id"],"limit":"10"}}

## Dependencies
* github.com/jackc/pgx
* github.com/pifagor87/conect_db_microservice