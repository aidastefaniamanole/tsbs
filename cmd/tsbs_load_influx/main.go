// bulk_load_influx loads an InfluxDB daemon with data from stdin.
//
// The caller is responsible for assuring that the database is empty before
// bulk load.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/timescale/tsbs/pkg/targets"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/timescale/tsbs/internal/utils"
	"github.com/timescale/tsbs/load"
)

// Program option vars:
var (
	daemonURLs        []string
	replicationFactor int
	backoff           time.Duration
	useGzip           bool
	doAbortOnExist    bool
	consistency       string
)

// Global vars
var (
	loader  *load.BenchmarkRunner
	bufPool sync.Pool
)

var consistencyChoices = map[string]struct{}{
	"any":    struct{}{},
	"one":    struct{}{},
	"quorum": struct{}{},
	"all":    struct{}{},
}

// allows for testing
var fatal = log.Fatalf

// Parse args:
func init() {
	var config load.BenchmarkRunnerConfig
	config.AddToFlagSet(pflag.CommandLine)
	var csvDaemonURLs string

	pflag.String("urls", "http://localhost:8086", "InfluxDB URLs, comma-separated. Will be used in a round-robin fashion.")
	pflag.Int("replication-factor", 1, "Cluster replication factor (only applies to clustered databases).")
	pflag.String("consistency", "all", "Write consistency. Must be one of: any, one, quorum, all.")
	pflag.Duration("backoff", time.Second, "Time to sleep between requests when server indicates backpressure is needed.")
	pflag.Bool("gzip", true, "Whether to gzip encode requests (default true).")

	pflag.Parse()

	err := utils.SetupConfigFile()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode config: %s", err))
	}

	csvDaemonURLs = viper.GetString("urls")
	replicationFactor = viper.GetInt("replication-factor")
	consistency = viper.GetString("consistency")
	backoff = viper.GetDuration("backoff")
	useGzip = viper.GetBool("gzip")

	if _, ok := consistencyChoices[consistency]; !ok {
		log.Fatalf("invalid consistency settings")
	}

	daemonURLs = strings.Split(csvDaemonURLs, ",")
	if len(daemonURLs) == 0 {
		log.Fatal("missing 'urls' flag")
	}
	config.HashWorkers = false
	loader = load.GetBenchmarkRunner(config)
}

type benchmark struct{}

func (b *benchmark) GetDataSource() targets.DataSource {
	return &fileDataSource{scanner: bufio.NewScanner(load.GetBufferedReader(loader.FileName))}
}

func (b *benchmark) GetBatchFactory() targets.BatchFactory {
	return &factory{}
}

func (b *benchmark) GetPointIndexer(_ uint) targets.PointIndexer {
	return &targets.ConstantIndexer{}
}

func (b *benchmark) GetProcessor() targets.Processor {
	return &processor{}
}

func (b *benchmark) GetDBCreator() targets.DBCreator {
	return &dbCreator{}
}

func main() {
	bufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 4*1024*1024))
		},
	}

	loader.RunBenchmark(&benchmark{})
}
