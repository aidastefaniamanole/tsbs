package query

import (
	"fmt"

	iutils "github.com/timescale/tsbs/internal/utils"
)

// TODO: remove ugly hack
type QueryInfo struct {
	// prometheus array of queries when we want to use more than CPU related metrics
	Queries []string
	// prometheus query
	Query string
	// label to describe type of query
	Label string
	// desc to describe type of query
	Desc string
	// time range for query executing
	Interval *iutils.TimeInterval
	// time period to group by in seconds
	Step string
	// metric to be used in the query
	Metric string
}

// Query is an interface used for encoding a benchmark query for different databases
type Query interface {
	Release()
	HumanLabelName() []byte
	HumanDescriptionName() []byte
	GetID() uint64
	SetID(uint64)
	SetQueryInfo(qi *QueryInfo)
	GetQueryInfo() *QueryInfo
	fmt.Stringer
}
