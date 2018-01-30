# DB Perf Comp

## Sample Run

```bash
$ go run main.go --db-path /tmp/perfcomp-sqlite.db --db-type sqlite --mode reinit --init-size 1000
... (starts sqlite, initializes DB)
$ go run main.go --db-path /tmp/perfcomp-sqlite.db --db-type sqlite --mode select
... (starts sqlite, performs query)
$ go run main.go --db-type bolt --db-path /tmp/perfcomp-bolt.db --mode reinit --init-size 1000
... (starts boltdb, initializes DB)
$ go run main.go --db-type bolt --db-path /tmp/perfcomp-bolt.db --mode select
... (starts boltdb, performs query)
$ go run main.go --db-path /tmp/perfcomp-sqlite.db --db-type bolt --mode select --ot 00000009
... (starts boltdb, performs query with offset token)
```
