# TSBS Supplemental Guide: QuestDB

QuestDB is a high-performance open-source time series database with SQL as a
query language with time-oriented extensions. QuestDB implements PostgreSQL wire
protocol, a REST API, and supports ingestion using InfluxDB line protocol.

This guide explains how the data for TSBS is generated along with additional
flags available when using the data importer (`tsbs_load_questdb`).
**This should be read _after_ the main README.**

## Data format

Data generated by `tsbs_generate_data` is in InfluxDB line protocol format where each
reading is composed of the following:

- the table name followed by a comma
- several comma-separated items of tags in the format `<label>=<value>` followed
  by a space
- several comma-separated items of fields in the format `<label>=<value>`
  followed by a space
- a timestamp for the record
- a newline character `\n`

An example reading from the `iot` use case looks like the following:

```text
diagnostics,name=truck_3985,fleet=West,driver=Seth,model=H-2,device_version=v1.5 load_capacity=1500,fuel_capacity=150,nominal_fuel_consumption=12,fuel_state=0.8,current_load=482,status=4i 1451609990000000000
```

## `tsbs_load_questdb` additional flags

**`--ilp-bind-to`** (type: `string`, default `127.0.0.1:9009`)

QuestDB InfluxDB line protocol TCP port in the format `<ip>:<port>`

**`--url`** (type: `string`, default: `http://localhost:9000/`)

QuestDB REST end point.

**`-help`**

Prints available flags and their defaults:

```bash
~/tmp/go/bin/tsbs_load_questdb -help
```

## How to run the test (FreeBSD example)

Firstly, install and build the benchmark suite

### Set up TSBS

Create a temporary directory for the Go binaries

```bash
mkdir -p ~/tmp/go/src/github.com/timescale/
cd ~/tmp/go/src/github.com/timescale/
```

Clone the TSBS repository, build test and install Go binaries:

```bash
git clone git@github.com:questdb/tsbs.git
cd ~/tmp/go/src/github.com/timescale/tsbs/ && git checkout questdb-tsbs-load
GOPATH=~/tmp/go go build -v ./...
GOPATH=~/tmp/go go test -v github.com/timescale/tsbs/cmd/tsbs_load_questdb
GOPATH=~/tmp/go go install -v ./...
```

### Generating data

Data is generated using the `influx` format. To generate a small dataset for
quick benchmarks:

```bash
~/tmp/go/bin/tsbs_generate_data \
--use-case="iot" --seed=123 --scale=4000 \
--timestamp-start="2016-01-01T00:00:00Z" --timestamp-end="2016-01-01T01:00:00Z" \
--log-interval="10s" --format="influx" > /tmp/data
```

To generate a full data set for more intensive benchmarks:

```bash
~/tmp/go/bin/tsbs_generate_data \
--use-case="iot" --seed=123 --scale=4000 \
--timestamp-start="2016-01-01T00:00:00Z" --timestamp-end="2016-01-04T00:00:00Z" \
--log-interval="10s" --format="influx" > /tmp/data
```

### Running the benchmark tool

Generated data can be loaded directly using the tool:

```bash
cat /tmp/data | ~/tmp/go/bin/tsbs_load_questdb
```

Alternatively, shell scripts are provided which can be used to generate and load
data:

```bash
cd ~/tmp/go/src/github.com/timescale/tsbs

# generates data file /tmp/bulk_data/influx-data.gz
PATH=${PATH}:~/tmp/go/bin FORMATS=influx TS_END=2016-01-01T02:00:00Z bash ./scripts/generate_data.sh
# load data into QuestDB
PATH=${PATH}:~/tmp/go/bin NUM_WORKERS=1 ./scripts/load/load_questdb.sh
```

### Query benchmarks for iot data set (single-groupby-5-8-1 type)

Queries are generated using the `questdb` format.

**single-groupby-5-8-1:**

The dataset used to run the queries is created with the following commands for
`single-groupby-5-8-1`:

```bash
cd ~/tmp/go/src/github.com/timescale/

~/tmp/go/bin/tsbs_generate_queries \
--use-case="cpu-only" --seed=123 --scale=4000 \
--timestamp-start="2016-01-01T00:00:00Z" --timestamp-end="2016-01-02T00:00:01Z" \
--queries=1000 --query-type="single-groupby-5-8-1" \
--format="questdb" > /tmp/queries_questdb

~/tmp/go/bin/tsbs_run_queries_questdb --file /tmp/queries_questdb --print-interval 500
```

### Query benchmarks for iot data set (high-cpu-1 use case)

Queries are generated using the `questdb` format.

**high-cpu-1:**

The dataset used to run the queries is created with the following commands for
`high-cpu-1`:

```bash
cd ~/tmp/go/src/github.com/timescale/

~/tmp/go/bin/tsbs_generate_queries \
--use-case="cpu-only" --seed=123 --scale=4000 \
--timestamp-start="2016-01-01T00:00:00Z" --timestamp-end="2016-01-02T00:00:01Z" \
--queries=1000 --query-type="high-cpu-1" --format="questdb" > /tmp/queries_questdb

~/tmp/go/bin/tsbs_run_queries_questdb --file /tmp/queries_questdb --print-interval 500
```

### Query benchmark shell scripts

Additionally, shell scripts are provided which can be used to generate and run
the queries:

```bash
cd ~/tmp/go/src/github.com/timescale/
PATH=${PATH}:~/tmp/go/bin FORMATS=questdb TS_END=2016-01-02T00:00:00Z bash ./scripts/generate_queries.sh
PATH=${PATH}:~/tmp/go/bin ./scripts/run_queries/run_queries_questdb.sh
```