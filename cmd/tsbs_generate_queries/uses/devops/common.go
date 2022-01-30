package devops

import (
	"fmt"
	"time"

	"github.com/timescale/tsbs/cmd/tsbs_generate_queries/uses/common"
	"github.com/timescale/tsbs/pkg/query"
)

const (
	allHosts                = "all hosts"
	errNHostsCannotNegative = "nHosts cannot be negative"
	errNoMetrics            = "cannot get 0 metrics"
	errTooManyMetrics       = "too many metrics asked for"

	// TableName is the name of the table where the time series data is stored for devops use case.
	TableName = "cpu"

	// DoubleGroupByDuration is the how big the time range for DoubleGroupBy query is
	DoubleGroupByDuration = 12 * time.Hour
	// HighCPUDuration is the how big the time range for HighCPU query is
	HighCPUDuration = 12 * time.Hour
	// MaxAllDuration is the how big the time range for MaxAll query is
	MaxAllDuration = 8 * time.Hour

	// LabelSingleGroupby is the label prefix for queries of the single groupby variety
	LabelSingleGroupby = "single-groupby"
	// LabelDoubleGroupby is the label prefix for queries of the double groupby variety
	LabelDoubleGroupby = "double-groupby"
	// LabelLastpoint is the label for the lastpoint query
	LabelLastpoint = "lastpoint"
	// LabelMaxAll is the label prefix for queries of the max all variety
	LabelMaxAll = "cpu-max-all"
	// LabelGroupbyOrderbyLimit is the label for groupby-orderby-limit query
	LabelGroupbyOrderbyLimit = "groupby-orderby-limit"
	// LabelHighCPU is the prefix for queries of the high-CPU variety
	LabelHighCPU = "high-cpu"
)

// Core is the common component of all generators for all systems
type Core struct {
	*common.Core
}

// NewCore returns a new Core for the given time range and cardinality
func NewCore(start, end time.Time, scale int) (*Core, error) {
	c, err := common.NewCore(start, end, scale)
	return &Core{Core: c}, err

}

// GetRandomHosts returns a random set of nHosts from a given Core
func (d *Core) GetRandomHosts(nHosts int) ([]string, error) {
	return getRandomHosts(nHosts, d.Scale)
}

// cpuMetrics is the list of metric names for CPU
var cpuMetrics = []string{
	"usage_user",
	"usage_system",
	"usage_idle",
	"usage_nice",
	"usage_iowait",
	"usage_irq",
	"usage_softirq",
	"usage_steal",
	"usage_guest",
	"usage_guest_nice",
}

// GetCPUMetricsSlice returns a subset of metrics for the CPU
func GetCPUMetricsSlice(numMetrics int) ([]string, error) {
	if numMetrics <= 0 {
		return nil, fmt.Errorf(errNoMetrics)
	}
	if numMetrics > len(cpuMetrics) {
		return nil, fmt.Errorf(errTooManyMetrics)
	}
	return cpuMetrics[:numMetrics], nil
}

func GetMetricsSlice(numMetrics int, metrics []string) ([]string, error) {
	if numMetrics <= 0 {
		return nil, fmt.Errorf(errNoMetrics)
	}
	if numMetrics > len(metrics) {
		return nil, fmt.Errorf(errTooManyMetrics)
	}
	return metrics[:numMetrics], nil
}

// GetAllCPUMetrics returns all the metrics for CPU
func GetAllCPUMetrics() []string {
	return cpuMetrics
}

// GetCPUMetricsLen returns the number of metrics in CPU
func GetCPUMetricsLen() int {
	return len(cpuMetrics)
}

var diskMetrics = []string{
	"total",
	"free",
	"used",
	"used_percent",
	"inodes_total",
	"inodes_free",
	"inodes_used",
}

func GetAllDiskMetrics() []string {
	return diskMetrics
}

func GetDiskMetricsLen() int {
	return len(diskMetrics)
}

var diskioMetrics = []string{
	"reads",
	"writes",
	"read_bytes",
	"write_bytes",
	"read_time",
	"write_time",
	"io_time",
}

func GetAllDiskioMetrics() []string {
	return diskioMetrics
}

func GetDiskioMetricsLen() int {
	return len(diskioMetrics)
}

var kernelMetrics = []string{
	"interrupts",
	"context_switches",
	"processes_forked",
	"disk_pages_in",
	"disk_pages_out",
}

func GetAllKernelMetrics() []string {
	return kernelMetrics
}

func GetKernelMetricsLen() int {
	return len(kernelMetrics)
}

var memMetrics = []string{
	"total",
	"available",
	"used",
	"free",
	"cached",
	"buffered",
	"used_percent",
	"available_percent",
	"buffered_percent",
}

func GetAllMemMetrics() []string {
	return memMetrics
}

func GetMemMetricsLen() int {
	return len(memMetrics)
}

var netMetrics = []string{
	"bytes_sent",
	"bytes_recv",
	"packets_sent",
	"packets_recv",
	"err_in",
	"err_out",
	"drop_in",
	"drop_out",
}

func GetAllNetMetrics() []string {
	return netMetrics
}

func GetNetMetricsLen() int {
	return len(netMetrics)
}

var nginxMetrics = []string{
	"accepts",
	"active",
	"handled",
	"reading",
	"requests",
	"waiting",
	"writing",
}

func GetAllNginxMetrics() []string {
	return nginxMetrics
}

func GetNginxMetricsLen() int {
	return len(nginxMetrics)
}

var postgresqlMetrics = []string{
	"numbackends",
	"xact_commit",
	"xact_rollback",
	"blks_read",
	"blks_hit",
	"tup_returned",
	"tup_fetched",
	"tup_inserted",
	"tup_updated",
	"tup_deleted",
	"conflicts",
	"temp_files",
	"temp_bytes",
	"deadlocks",
	"blk_read_time",
	"blk_write_time",
}

func GetAllPostgresqlMetrics() []string {
	return postgresqlMetrics
}

func GetPostgresqlMetricsLen() int {
	return len(postgresqlMetrics)
}

var redisMetrics = []string{
	"total_connections_received",
	"expired_keys",
	"evicted_keys",
	"keyspace_hits",
	"keyspace_misses",
	"instantaneous_ops_per_sec",
	"instantaneous_input_kbps",
	"instantaneous_output_kbps",
	"connected_clients",
	"used_memory",
	"used_memory_rss",
	"used_memory_peak",
	"used_memory_lua",
	"rdb_changes_since_last_save",
	"sync_full",
	"sync_partial_ok",
	"sync_partial_err",
	"pubsub_channels",
	"pubsub_patterns",
	"latest_fork_usec",
	"connected_slaves",
	"master_repl_offset",
	"repl_backlog_active",
	"repl_backlog_size",
	"repl_backlog_histlen",
	"mem_fragmentation_ratio",
	"used_cpu_sys",
	"used_cpu_user",
	"used_cpu_sys_children",
	"used_cpu_user_children",
}

func GetAllRedisMetrics() []string {
	return redisMetrics
}

func GetRedisMetricsLen() int {
	return len(redisMetrics)
}

func GetMetrics(metric string) []string {
	switch metric {
	case "disk":
		return diskMetrics
	case "diskio":
		return diskioMetrics
	case "kernel":
		return kernelMetrics
	case "mem":
		return memMetrics
	case "net":
		return netMetrics
	case "nginx":
		return nginxMetrics
	case "postgresql":
		return postgresqlMetrics
	case "redis":
		return redisMetrics
	default:
		fmt.Printf("Return the CPU related metrics")
	}
	return cpuMetrics
}

// SingleGroupbyFiller is a type that can fill in a single groupby query
type SingleGroupbyFiller interface {
	GroupByTime(query.Query, int, int, time.Duration)
}

// DoubleGroupbyFiller is a type that can fill in a double groupby query
type DoubleGroupbyFiller interface {
	GroupByTimeAndPrimaryTag(query.Query, int)
}

// LastPointFiller is a type that can fill in a last point query
type LastPointFiller interface {
	LastPointPerHost(query.Query)
}

// MaxAllFiller is a type that can fill in a max all CPU metrics query
type MaxAllFiller interface {
	MaxAllCPU(query.Query, int)
}

// GroupbyOrderbyLimitFiller is a type that can fill in a groupby-orderby-limit query
type GroupbyOrderbyLimitFiller interface {
	GroupByOrderByLimit(query.Query)
}

// HighCPUFiller is a type that can fill in a high-cpu query
type HighCPUFiller interface {
	HighCPUForHosts(query.Query, int)
}

// GetDoubleGroupByLabel returns the Query human-readable label for DoubleGroupBy queries
func GetDoubleGroupByLabel(dbName string, numMetrics int) string {
	return fmt.Sprintf("%s mean of %d metrics, all hosts, random %s by 1h", dbName, numMetrics, DoubleGroupByDuration)
}

// GetHighCPULabel returns the Query human-readable label for HighCPU queries
func GetHighCPULabel(dbName string, nHosts int) (string, error) {
	label := dbName + " CPU over threshold, "
	if nHosts > 0 {
		label += fmt.Sprintf("%d host(s)", nHosts)
	} else if nHosts == 0 {
		label += allHosts
	} else {
		return "", fmt.Errorf("nHosts cannot be negative")
	}
	return label, nil
}

// GetMaxAllLabel returns the Query human-readable label for MaxAllCPU queries
func GetMaxAllLabel(dbName string, nHosts int) string {
	return fmt.Sprintf("%s max of all CPU metrics, random %4d hosts, random %s by 1h", dbName, nHosts, MaxAllDuration)
}

// getRandomHosts returns a subset of numHosts hostnames of a permutation of hostnames,
// numbered from 0 to totalHosts.
// Ex.: host_12, host_7, host_25 for numHosts=3 and totalHosts=30 (3 out of 30)
func getRandomHosts(numHosts int, totalHosts int) ([]string, error) {
	if numHosts < 1 {
		return nil, fmt.Errorf("number of hosts cannot be < 1; got %d", numHosts)
	}
	if numHosts > totalHosts {
		return nil, fmt.Errorf("number of hosts (%d) larger than total hosts. See --scale (%d)", numHosts, totalHosts)
	}

	randomNumbers, err := common.GetRandomSubsetPerm(numHosts, totalHosts)
	if err != nil {
		return nil, err
	}

	hostnames := []string{}
	for _, n := range randomNumbers {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	return hostnames, nil
}
