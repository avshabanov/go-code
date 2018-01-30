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

### BoltDB

Parallel access tests:

```bash
$ go run main.go --db-path /tmp/perfcomp-bolt-100k.db --db-type bolt --mode reinit --init-size 100000
... (init db)
$ go run main.go --db-path /tmp/perfcomp-bolt-100k.db --db-type bolt --mode parallel-select
... (run tests)
```

Sample run:

```raw
# MacBook Pro 15'', Mid 2014
job 6 done, totalUsersFetched=17230, timeSpent=1.682705173s
job 1 done, totalUsersFetched=18040, timeSpent=1.706430234s
job 8 done, totalUsersFetched=17910, timeSpent=1.721068692s
job 0 done, totalUsersFetched=17520, timeSpent=1.764967563s
job 7 done, totalUsersFetched=18640, timeSpent=1.805118965s
job 2 done, totalUsersFetched=18080, timeSpent=1.805226579s
job 4 done, totalUsersFetched=19200, timeSpent=1.809998407s
job 9 done, totalUsersFetched=19160, timeSpent=1.850357077s
job 3 done, totalUsersFetched=20280, timeSpent=1.883559468s
job 5 done, totalUsersFetched=29010, timeSpent=2.187639415s
```

### Sqlite

```bash
$ go run main.go --db-path /tmp/perfcomp-sqlite-100k.db --db-type sqlite --mode reinit --init-size 100000
... (init db)
$ go run main.go --db-path /tmp/perfcomp-sqlite-100k.db --db-type sqlite --mode parallel-select
... (run tests)
```

Sample run:

```raw
# MacBook Pro 15'', Mid 2014
job 6 done, totalUsersFetched=18120, timeSpent=2.285631353s
job 8 done, totalUsersFetched=18800, timeSpent=2.404640239s
job 4 done, totalUsersFetched=20120, timeSpent=2.668034817s
job 2 done, totalUsersFetched=21560, timeSpent=2.911040734s
job 1 done, totalUsersFetched=21520, timeSpent=2.961946309s
job 0 done, totalUsersFetched=21000, timeSpent=3.004260151s
job 7 done, totalUsersFetched=22320, timeSpent=3.040617887s
job 9 done, totalUsersFetched=22840, timeSpent=3.078309619s
job 3 done, totalUsersFetched=23960, timeSpent=3.114378004s
job 5 done, totalUsersFetched=31740, timeSpent=3.543083141s
```
