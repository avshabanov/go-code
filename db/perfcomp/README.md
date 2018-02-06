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
$ go run main.go --db-path /tmp/perfcomp-bolt-100k.db --db-type bolt --mode random-get
```

Sample run:

```raw
# MacBookPro11,3
job 6 done, totalUsersFetched=17230, timeSpent=1.772132812s
job 0 done, totalUsersFetched=17520, timeSpent=1.796367396s
job 8 done, totalUsersFetched=17910, timeSpent=1.823644427s
job 1 done, totalUsersFetched=18040, timeSpent=1.836107206s
job 2 done, totalUsersFetched=18080, timeSpent=1.837655139s
job 7 done, totalUsersFetched=18640, timeSpent=1.874763292s
job 9 done, totalUsersFetched=19160, timeSpent=1.908456461s
job 4 done, totalUsersFetched=19200, timeSpent=1.914502361s
job 3 done, totalUsersFetched=20280, timeSpent=1.98823104s
job 5 done, totalUsersFetched=29010, timeSpent=2.556816733s
```

### Sqlite

```bash
$ go run main.go --db-path /tmp/perfcomp-sqlite-100k.db --db-type sqlite --mode reinit --init-size 100000
... (init db)
$ go run main.go --db-path /tmp/perfcomp-sqlite-100k.db --db-type sqlite --mode parallel-select
... (run tests, do random scans in parallel)
$ go run main.go --db-path /tmp/perfcomp-sqlite-100k.db --db-type sqlite --mode random-get
... (run tests, get random users in parallel)
```

Sample run:

```raw
# (Optimized) MacBookPro11,3
job 6 done, totalUsersFetched=17230, timeSpent=1.316952223s
job 8 done, totalUsersFetched=17910, timeSpent=1.36384122s
job 4 done, totalUsersFetched=19200, timeSpent=1.445517728s
job 0 done, totalUsersFetched=17520, timeSpent=1.64475786s
job 2 done, totalUsersFetched=18080, timeSpent=1.693972656s
job 7 done, totalUsersFetched=18640, timeSpent=1.703829041s
job 1 done, totalUsersFetched=18040, timeSpent=1.704083697s
job 9 done, totalUsersFetched=19160, timeSpent=1.719454138s
job 3 done, totalUsersFetched=20280, timeSpent=1.763559057s
job 5 done, totalUsersFetched=29010, timeSpent=1.945161213s

# (Pre-optimized) MacBookPro11,3
job 6 done, totalUsersFetched=17230, timeSpent=2.350053822s
job 8 done, totalUsersFetched=17910, timeSpent=2.479878825s
job 4 done, totalUsersFetched=19200, timeSpent=2.572471238s
job 0 done, totalUsersFetched=17520, timeSpent=2.633651209s
job 2 done, totalUsersFetched=18080, timeSpent=2.690347777s
job 1 done, totalUsersFetched=18040, timeSpent=2.726054241s
job 7 done, totalUsersFetched=18640, timeSpent=2.744656555s
job 9 done, totalUsersFetched=19160, timeSpent=2.75299892s
job 3 done, totalUsersFetched=20280, timeSpent=2.824421861s
job 5 done, totalUsersFetched=29010, timeSpent=3.239841158s
```

### KV Sqlite

```bash
$ go run main.go --db-path /tmp/perfcomp-kvsqlite-100k.db --db-type kvsqlite --mode reinit --init-size 100000
...
$ go run main.go --db-path /tmp/perfcomp-kvsqlite-100k.db --db-type kvsqlite --mode parallel-select
...
$ go run main.go --db-path /tmp/perfcomp-kvsqlite-100k.db --db-type kvsqlite --mode random-get
...
```

Sample run:

```raw
# MacBookPro11,3
job 6 done, totalUsersFetched=17230, timeSpent=1.802435994s
job 8 done, totalUsersFetched=17910, timeSpent=1.877926901s
job 4 done, totalUsersFetched=19200, timeSpent=2.045990441s
job 2 done, totalUsersFetched=18080, timeSpent=2.308357453s
job 1 done, totalUsersFetched=18040, timeSpent=2.308585171s
job 0 done, totalUsersFetched=17520, timeSpent=2.327896002s
job 7 done, totalUsersFetched=18640, timeSpent=2.35772531s
job 9 done, totalUsersFetched=19160, timeSpent=2.405297984s
job 3 done, totalUsersFetched=20280, timeSpent=2.424307993s
job 5 done, totalUsersFetched=29010, timeSpent=2.690742239s
```
