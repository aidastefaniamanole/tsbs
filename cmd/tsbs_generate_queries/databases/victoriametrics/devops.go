package victoriametrics

import (
	"fmt"
	"strings"
	"time"

	"github.com/timescale/tsbs/cmd/tsbs_generate_queries/uses/devops"
	"github.com/timescale/tsbs/pkg/query"
)

// Devops produces PromQL queries for all the devops query types.
type Devops struct {
	*BaseGenerator
	*devops.Core
	Metrics []string
}

// mustGetRandomHosts is the form of GetRandomHosts that cannot error; if it does error,
// it causes a panic.
func (d *Devops) mustGetRandomHosts(nHosts int) []string {
	hosts, err := d.GetRandomHosts(nHosts)
	if err != nil {
		panic(err.Error())
	}
	return hosts
}

func (d *Devops) GroupByOrderByLimit(qi query.Query) {
	panic("GroupByOrderByLimit not supported in PromQL")
}

func (d *Devops) LastPointPerHost(qq query.Query) {
	panic("LastPointPerHost not supported in PromQL")
}

func (d *Devops) HighCPUForHosts(qi query.Query, nHosts int) {
	panic("HighCPUForHosts not supported in PromQL")
}

// GroupByTime selects the MAX for numMetrics metrics under 'cpu'
// per minute for nhosts hosts,
// e.g. in pseudo-PromQL:
// max(
// 	max_over_time(
// 		{__name__=~"metric1|metric2...|metricN",hostname=~"hostname1|hostname2...|hostnameN"}[1m]
// 	)
// ) by (__name__)
func (d *Devops) GroupByTime(qq query.Query, nHosts, numMetrics int, timeRange time.Duration) {
	hosts := d.mustGetRandomHosts(nHosts)

	metrics := mustGetCPUMetricsSlice(numMetrics)
	selectClause := getSelectClause(metrics, hosts, "cpu")

	qi := &query.QueryInfo{
		Query:    fmt.Sprintf("max(max_over_time(%s[1m])) by (__name__)", selectClause),
		Label:    fmt.Sprintf("VictoriaMetrics %d cpu metric(s), random %4d hosts, random %s by 1m", numMetrics, nHosts, timeRange),
		Interval: d.Interval.MustRandWindow(timeRange),
		Step:     "60",
	}

	queries := make([]string, len(d.Metrics))
	for i := 0; i < len(d.Metrics); i++ {
		metrics = mustGetMetricsSlice(numMetrics, devops.GetMetrics(d.Metrics[i]))
		selectClause = getSelectClause(metrics, hosts, d.Metrics[i])
		queries[i] = fmt.Sprintf("max(max_over_time(%s[1m])) by (__name__)", selectClause)
	}

	qi.Queries = queries
	d.fillInQuery(qq, qi)
}

// GroupByTimeAndPrimaryTag selects the AVG of numMetrics metrics under 'cpu' per device per hour for a day,
// e.g. in pseudo-PromQL:
//
// avg(
// 	avg_over_time(
// 		{__name__=~"metric1|metric2...|metricN"}[1h]
// 	)
// ) by (__name__, hostname)
//
// Resultsets:
// double-groupby-1
// double-groupby-5
// double-groupby-all
func (d *Devops) GroupByTimeAndPrimaryTag(qq query.Query, numMetrics int) {
	metrics := mustGetCPUMetricsSlice(numMetrics)
	selectClause := getSelectClause(metrics, nil, "cpu")
	qi := &query.QueryInfo{
		Query:    fmt.Sprintf("avg(avg_over_time(%s[1h])) by (__name__, hostname)", selectClause),
		Label:    devops.GetDoubleGroupByLabel("VictoriaMetrics", numMetrics),
		Interval: d.Interval.MustRandWindow(devops.DoubleGroupByDuration),
		Step:     "3600",
	}

	queries := make([]string, len(d.Metrics))
	for i := 0; i < len(d.Metrics); i++ {
		if numMetrics == devops.GetCPUMetricsLen() {
			metrics = devops.GetMetrics(d.Metrics[i])
		} else {
			metrics = mustGetMetricsSlice(numMetrics, devops.GetMetrics(d.Metrics[i]))
		}
		selectClause = getSelectClause(metrics, nil, d.Metrics[i])
		queries[i] = fmt.Sprintf("avg(avg_over_time(%s[1h])) by (__name__, hostname)", selectClause)
	}

	qi.Queries = queries
	d.fillInQuery(qq, qi)
}

// MaxAllCPU selects the MAX of all metrics under 'cpu' per hour for nhosts hosts,
// e.g. in pseudo-PromQL:
//
// max(
// 	max_over_time(
// 		{hostname=~"hostname1|hostname2...|hostnameN"}[1h]
// 	)
// ) by (__name__)
func (d *Devops) MaxAllCPU(qq query.Query, nHosts int) {
	hosts := d.mustGetRandomHosts(nHosts)
	selectClause := getSelectClause(devops.GetAllCPUMetrics(), hosts, "cpu")

	qi := &query.QueryInfo{
		Query:    fmt.Sprintf("max(max_over_time(%s[1h])) by (__name__)", selectClause),
		Label:    devops.GetMaxAllLabel("VictoriaMetrics", nHosts),
		Interval: d.Interval.MustRandWindow(devops.MaxAllDuration),
		Step:     "3600",
	}

	queries := make([]string, len(d.Metrics))
	for i := 0; i < len(d.Metrics); i++ {
		selectClause = getSelectClause(devops.GetMetrics(d.Metrics[i]), nil, d.Metrics[i])
		queries[i] = fmt.Sprintf("max(max_over_time(%s[1h])) by (__name__)", selectClause)
	}

	qi.Queries = queries
	d.fillInQuery(qq, qi)
}

func getHostClause(hostnames []string) string {
	if len(hostnames) == 0 {
		return ""
	}
	if len(hostnames) == 1 {
		return fmt.Sprintf("hostname='%s'", hostnames[0])
	}
	return fmt.Sprintf("hostname=~'%s'", strings.Join(hostnames, "|"))
}

func getSelectClause(metrics, hosts []string, metricName string) string {
	if len(metrics) == 0 {
		panic("BUG: must be at least one metric name in clause")
	}

	hostsClause := getHostClause(hosts)
	if len(metrics) == 1 {
		return fmt.Sprintf("%s_%s{%s}", metricName, metrics[0], hostsClause)
	}

	metricsClause := strings.Join(metrics, "|")
	if len(hosts) > 0 {
		return fmt.Sprintf("{__name__=~'%s_(%s)', %s}", metricName, metricsClause, hostsClause)
	}
	return fmt.Sprintf("{__name__=~'%s_(%s)'}", metricName, metricsClause)
}

// mustGetCPUMetricsSlice is the form of GetCPUMetricsSlice that cannot error; if it does error,
// it causes a panic.
func mustGetCPUMetricsSlice(numMetrics int) []string {
	metrics, err := devops.GetCPUMetricsSlice(numMetrics)
	if err != nil {
		panic(err.Error())
	}
	return metrics
}

func mustGetMetricsSlice(numMetrics int, metrics []string) []string {
	metrics, err := devops.GetMetricsSlice(numMetrics, metrics)
	if err != nil {
		panic(err.Error())
	}
	return metrics
}
