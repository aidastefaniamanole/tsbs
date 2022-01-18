package victoriametrics

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/timescale/tsbs/cmd/tsbs_generate_queries/uses/devops"
	"github.com/timescale/tsbs/cmd/tsbs_generate_queries/utils"
	iutils "github.com/timescale/tsbs/internal/utils"
	"github.com/timescale/tsbs/pkg/query"
)

type BaseGenerator struct{}

// GenerateEmptyQuery returns an empty query.HTTP.
func (g *BaseGenerator) GenerateEmptyQuery() query.Query {
	return query.NewHTTP()
}

// NewDevops creates a new devops use case query generator.
func (g *BaseGenerator) NewDevops(start, end time.Time, scale int) (utils.QueryGenerator, error) {
	core, err := devops.NewCore(start, end, scale)
	if err != nil {
		return nil, err
	}
	return &Devops{
		BaseGenerator: g,
		Core:          core,
		Metrics:       []string{"disk", "diskio", "kernel", "mem", "net", "nginx", "postgresql", "redis"},
	}, nil
}

type queryInfo struct {
	// prometheus query
	query string
	// label to describe type of query
	label string
	// desc to describe type of query
	desc string
	// time range for query executing
	interval *iutils.TimeInterval
	// time period to group by in seconds
	step string
}

// fill Query fills the query struct with data
func (g *BaseGenerator) fillInQuery(qq query.Query, qi *query.QueryInfo) {
	q := qq.(*query.HTTP)
	q.HumanLabel = []byte(qi.Label)
	if qi.Interval != nil {
		q.HumanDescription = []byte(fmt.Sprintf("%s: %s", qi.Label, qi.Interval.StartString()))
	}
	q.Method = []byte("GET")

	v := url.Values{}
	v.Set("query", qi.Query)
	v.Set("start", strconv.FormatInt(qi.Interval.StartUnixNano()/1e9, 10))
	v.Set("end", strconv.FormatInt(qi.Interval.EndUnixNano()/1e9, 10))
	v.Set("step", qi.Step)
	q.Path = []byte(fmt.Sprintf("/api/v1/query_range?%s", v.Encode()))
	q.Body = nil

	q.SetQueryInfo(qi)
}
