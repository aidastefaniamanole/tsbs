#!/bin/bash

# configurable variables
OUTPUT_DIR="/home/stefania/IdeaProjects/vm-benchmarking-client/src/main/resources"
USE_CASE="devops"
SCALE=${SCALE:-"8"}

# generate the 14 days old data
TIMESTAMP_START=$(date +"%Y-%m-%dT%H:%M:%S%:z" -d "2 days ago")
TIMESTAMP_END=$(date +"%Y-%m-%dT%H:%M:%S%:z")

TIMESTAMP_START=$(date --utc +%FT%TZ -d "2 days ago")
TIMESTAMP_END=$(date --utc +%FT%TZ)

LOG_INTERVAL=${LOG_INTERVAL:-"10"}
DIFF=$(( $(date +%s -d ${TIMESTAMP_END})-$(date +%s -d ${TIMESTAMP_START}) ))
EXPERIMENT_DAYS=$(( $DIFF / (60 * 60 * 24) )) # to get the number of days the data spans over

echo -e "Use case=${USE_CASE}
Number of hosts that emit telemetry data=${SCALE}
Timestamp start=${TIMESTAMP_START}
Timestamp end=${TIMESTAMP_END}
The interval data is being logged at in seconds=${LOG_INTERVAL}
Data span in days=${EXPERIMENT_DAYS}" | tee -a "${OUTPUT_DIR}/${SCALE}hosts_${EXPERIMENT_DAYS}d"

go run main.go --use-case=${USE_CASE} --scale=${SCALE} \
     --timestamp-start=${TIMESTAMP_START} \
     --timestamp-end=${TIMESTAMP_END} \
     --log-interval=${LOG_INTERVAL}s --format="victoriametrics" >> "${OUTPUT_DIR}/${SCALE}hosts_${EXPERIMENT_DAYS}d"
