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

QUERY="single-groupby-5-8-1"
QUERY_TYPES=${QUERY:-${QUERY_TYPES_ALL}}

# Number of hosts to generate data about
SCALE=${SCALE:-"8"}

# Number of queries to generate
QUERIES=${QUERIES:-"1000"}

TIMESTAMP_START="2022-01-18T19:48:41Z"
TIMESTAMP_END="2022-01-20T19:48:41Z"
DIFF=$(( $(date +%s -d ${TIMESTAMP_END})-$(date +%s -d ${TIMESTAMP_START}) ))
EXPERIMENT_DAYS=$(( $DIFF / (60 * 60 * 24) )) # to get the number of days the data spans over

echo "Started query generation"
for QUERY_TYPE in ${QUERY_TYPES}; do
    echo -e "Number of hosts to emit queries about=${SCALE}
Number of queries=${QUERIES}
Timestamp start=${TIMESTAMP_START}
Timestamp end=${TIMESTAMP_END}
Queries span in days=${EXPERIMENT_DAYS}" >> "${OUTPUT_DIR}/${QUERY_TYPE}_${SCALE}_hosts.csv"
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