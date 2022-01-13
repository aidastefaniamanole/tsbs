#!/bin/bash

OUTPUT_DIR="/home/stefania/IdeaProjects/vm-benchmarking-client/src/main/resources"

# configurable variables
USE_CASE="devops"
SCALE=${SCALE:-"100"}
TIMESTAMP_START=${TIMESTAMP_START:-"2022-01-07T00:00:00Z"}
TIMESTAMP_END=${TIMESTAMP_END:-"2022-01-08T00:00:00Z"}
LOG_INTERVAL=${LOG_INTERVAL:-"10s"}
DIFF=$(( $(date +%s -d ${TIMESTAMP_END})-$(date +%s -d ${TIMESTAMP_START}) ))
EXPERIMENT_HOURS=$(( $DIFF / (60 * 60) )) # to get the number of hours the data spans over

echo -e "Use case: ${USE_CASE}
Number of hosts that emit telemetry data: ${SCALE}
Timestamp start: ${TIMESTAMP_START}
Timestamp end: ${TIMESTAMP_END}
The data is being logged every: ${LOG_INTERVAL}
Data spans over: ${EXPERIMENT_HOURS}h" >> "${OUTPUT_DIR}/${SCALE}hosts_${EXPERIMENT_HOURS}h" 

go run main.go --use-case=${USE_CASE} --scale=${SCALE} \
     --timestamp-start=${TIMESTAMP_START} \
     --timestamp-end=${TIMESTAMP_END} \
     --log-interval=${LOG_INTERVAL} --format="victoriametrics" >> "${OUTPUT_DIR}/${SCALE}hosts_${EXPERIMENT_HOURS}h"
