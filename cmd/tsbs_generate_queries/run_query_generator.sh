#!/bin/bash

OUTPUT_DIR="/home/stefania/IdeaProjects/vm-benchmarking-client/src/main/resources"

# All available for generation query types (sorted alphabetically)
QUERY_TYPES_ALL="\
cpu-max-all-1 \
cpu-max-all-8 \
double-groupby-1 \
double-groupby-5 \
double-groupby-all \
single-groupby-1-1-1 \
single-groupby-1-1-12 \
single-groupby-1-8-1 \
single-groupby-5-1-1 \
single-groupby-5-1-12 \
single-groupby-5-8-1"

QUERY_TYPES=${QUERY:-${QUERY_TYPES_ALL}}

# Number of hosts to generate data about
SCALE=${SCALE:-"100"}

# Number of queries to generate
QUERIES=${QUERIES:-"1000"}

echo -e "Use case: ${USE_CASE}
Number of hosts that emit queries about: ${SCALE}
Timestamp start: ${TIMESTAMP_START}
Timestamp end: ${TIMESTAMP_END}
Queries span over: ${EXPERIMENT_HOURS}h" >> "${OUTPUT_DIR}/${QUERY_TYPE}_${SCALE}_hosts.csv"

echo "Started query generation"
for QUERY_TYPE in ${QUERY_TYPES}; do
    echo "Generating ${QUERIES} of type ${QUERY_TYPE} for ${SCALE} hosts"
    echo 
    go run main.go --format "victoriametrics"  \
        --use-case "cpu-only" \
        --scale ${SCALE} \
        --timestamp-start "2022-01-06T00:00:00Z" \
        --timestamp-end "2022-01-07T00:00:00Z" \
        --queries ${QUERIES} \
        --query-type ${QUERY_TYPE} >> "${OUTPUT_DIR}/${QUERY_TYPE}_${SCALE}_hosts.csv"
done