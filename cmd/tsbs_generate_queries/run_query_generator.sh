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
SCALE=${SCALE:-"8"}

# Number of queries to generate
QUERIES=${QUERIES:-"100"}

TIMESTAMP_START="2022-01-17T23:46:31+01:00"
TIMESTAMP_END="2022-01-18T23:46:31+01:00"

echo "Started query generation"
for QUERY_TYPE in ${QUERY_TYPES}; do
    echo -e "Use case: ${USE_CASE}
    Number of hosts that emit queries about: ${SCALE}
    Timestamp start: ${TIMESTAMP_START}
    Timestamp end: ${TIMESTAMP_END}
    Queries span over: ${EXPERIMENT_HOURS}h" >> "${OUTPUT_DIR}/${QUERY_TYPE}_${SCALE}_hosts.csv"
    echo "Generating ${QUERIES} of type ${QUERY_TYPE} for ${SCALE} hosts"
    echo 
    go run main.go --format "victoriametrics"  \
        --use-case "cpu-only" \
        --scale ${SCALE} \
        --timestamp-start ${TIMESTAMP_START} \
        --timestamp-end ${TIMESTAMP_END} \
        --queries ${QUERIES} \
        --query-type ${QUERY_TYPE} >> "${OUTPUT_DIR}/${QUERY_TYPE}_${SCALE}_hosts.csv"
done