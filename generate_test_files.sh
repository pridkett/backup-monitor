#!/usr/bin/env bash

TARGET_DIR=$1
FILE_PREFIX=${2:-dummy}
NUM_FILES=${3:-10}

if [ -z "$TARGET_DIR" ]; then
    echo "Usage: $0 <target_dir> [prefix] [num_files]"
    exit 1
fi

if [ ! -d "$TARGET_DIR" ]; then
    echo "Error: $TARGET_DIR is not a directory"
    exit 1
fi

. shell_helpers.sh

for i in $(seq 1 ${NUM_FILES}); do
    filename="${TARGET_DIR}/${FILE_PREFIX}-$i.txt"
    echo "This is file number $i" > "$filename"
    _log "Created $filename"

    days_ago=$((RANDOM % 10 + 1))  # Random number between 1 and 10
    gtouch -d "$days_ago days ago" "$filename"
    _log "Modified $filename with a timestamp $days_ago days ago"

done

