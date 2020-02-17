package main

import (
	"bufio"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/prometheus/prompb"
	"github.com/timescale/tsbs/cmd/tsbs_generate_data/serialize"
	"github.com/timescale/tsbs/internal/utils"
	"github.com/timescale/tsbs/load"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// runs the benchmark
var loader *load.BenchmarkRunner
var promBatchPool sync.Pool = sync.Pool{New: func() interface{} { return &PrometheusBatch{} }}

// PrometheusBatch implements load.Batch interface
type PrometheusBatch struct {
	series []prompb.TimeSeries
}

func (pb *PrometheusBatch) Len() int {
	return len(pb.series)
}

func (pb *PrometheusBatch) Append(item *load.Point) {
	pb.series = append(pb.series, item.Data.(prompb.TimeSeries))
}

// PrometheusDecoder implements scan.PointDecoder interface
type PrometheusDecoder struct {
	iterator *serialize.PrometheusIterator
}

func (pd *PrometheusDecoder) Decode(reader *bufio.Reader) *load.Point {
	if pd.iterator.HasNext() {
		ts, err := pd.iterator.Next()
		if err != nil {
			panic(err)
		}
		return load.NewPoint(*ts)
	}
	return nil
}

// PrometheusProcessor implements load.Processor interface
type PrometheusProcessor struct {
	client *Client
}

func (pp *PrometheusProcessor) Init(_ int, _ bool) {}

// ProcessBatch ..
func (pp *PrometheusProcessor) ProcessBatch(b load.Batch, doLoad bool) (uint64, uint64) {
	promBatch := b.(*PrometheusBatch)
	nrSamples := uint64(promBatch.Len())
	if doLoad {
		err := pp.client.Post(promBatch.series)
		if err != nil {
			panic(err)
		}
	}
	// reset batch
	promBatch.series = promBatch.series[:0]
	promBatchPool.Put(promBatch)
	return nrSamples, nrSamples
}

// PrometheusBatchFactory implements Factory interface
type PrometheusBatchFactory struct{}

func (pbf *PrometheusBatchFactory) New() load.Batch {
	return promBatchPool.Get().(*PrometheusBatch)
}

// PrometheusBenchmark implements Benchmark interface
type PrometheusBenchmark struct{}

func (pm *PrometheusBenchmark) GetPointDecoder(br *bufio.Reader) load.PointDecoder {
	promIter, err := serialize.NewPrometheusIterator(br)
	if err != nil {
		panic(err)
	}
	return &PrometheusDecoder{iterator: promIter}
}

func (pm *PrometheusBenchmark) GetBatchFactory() load.BatchFactory {
	return &PrometheusBatchFactory{}
}

func (pm *PrometheusBenchmark) GetPointIndexer(_ uint) load.PointIndexer {
	// We always have one shared queue and 1+ workers
	return &load.ConstantIndexer{}
}

func (pm *PrometheusBenchmark) GetProcessor() load.Processor {
	client, err := NewClient(adapterWriteUrl, time.Second*30)
	if err != nil {
		panic(err)
	}
	return &PrometheusProcessor{client: client}
}

func (pm *PrometheusBenchmark) GetDBCreator() load.DBCreator {
	return nil
}

var adapterWriteUrl string

func init() {
	var config load.BenchmarkRunnerConfig
	config.AddToFlagSet(pflag.CommandLine)

	pflag.StringVar(&adapterWriteUrl, "adapter-write-url", "", "Prometheus adapter url to send data to")
	pflag.Parse()

	err := utils.SetupConfigFile()
	if err != nil {
		panic(fmt.Errorf("error setting up a config file: %s", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode config: %s", err))
	}

	loader = load.GetBenchmarkRunner(config)
}

func main() {
	loader.RunBenchmark(&PrometheusBenchmark{}, load.SingleQueue)
}
