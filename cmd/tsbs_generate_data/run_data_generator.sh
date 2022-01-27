#!/bin/bash

# configurable variables
OUTPUT_DIR="/home/stefania/IdeaProjects/vm-benchmarking-client/src/main/resources"
USE_CASE="devops"
SCALE=${SCALE:-"8"}
LOG_INTERVAL=${LOG_INTERVAL:-"10"}

TIMESTAMP_START=$(date --utc +%FT%TZ -d "2 days ago")
if [[ $1 == "days" ]]; then
  echo "Generating data with starting timestamp from $2 day(s) ago"
  TIMESTAMP_START=$(date --utc +%FT%TZ -d "$2 days ago")
  TIMESTAMP_END=$(date --utc +%FT%TZ)
  DIFF=$(( $(date +%s -d ${TIMESTAMP_END})-$(date +%s -d ${TIMESTAMP_START}) ))
  EXPERIMENT_SPAN=$(( $DIFF / (60 * 60 * 24) )) # to get the number of days the data spans over
else
  echo "Generating data with starting timestamp from $2 hour(s) ago"
  TIMESTAMP_START=$(date --utc +%FT%TZ -d "$2 hours ago")
  TIMESTAMP_END=$(date --utc +%FT%TZ)
  DIFF=$(( $(date +%s -d ${TIMESTAMP_END})-$(date +%s -d ${TIMESTAMP_START}) ))
  EXPERIMENT_SPAN=$(( $DIFF / (60 * 60) )) # to get the number of hours the data spans over
fi


echo -e "Use case=${USE_CASE}
Number of hosts that emit telemetry data=${SCALE}
Timestamp start=${TIMESTAMP_START}
Timestamp end=${TIMESTAMP_END}
The interval data is being logged at in seconds=${LOG_INTERVAL}
Data span=${EXPERIMENT_SPAN}$1" | tee -a "${OUTPUT_DIR}/${SCALE}hosts"

go run main.go --use-case=${USE_CASE} --scale=${SCALE} \
     --timestamp-start=${TIMESTAMP_START} \
     --timestamp-end=${TIMESTAMP_END} \
     --log-interval=${LOG_INTERVAL}s --format="victoriametrics" >> "${OUTPUT_DIR}/${SCALE}hosts"
