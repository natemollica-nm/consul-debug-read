#!/usr/bin/env bash

SNAPSHOT_FILE="$1"
DECODE_OUTPUT_DIR="$2"
DECODE_OUTPUT_IDENTIFIER="$3"

[ -z "$SNAPSHOT_FILE" ] && echo "Snapshot filename required!" && exit 1
[ -z "$DECODE_OUTPUT_DIR" ] && echo "Snapshot output dir required!" && exit 1
[ -z "$DECODE_OUTPUT_IDENTIFIER" ] && echo "Snapshot output ID required!" && exit 1

test -f "$SNAPSHOT_FILE" || {
    echo "No snapshot file at $SNAPSHOT_FILE!"
    exit 1
}

test -d "$DECODE_OUTPUT_DIR" || {
    mkdir -p "$DECODE_OUTPUT_DIR"
}

echo "snapshot-decode: Running consul snapshot decode $SNAPSHOT_FILE"
TYPES="$(consul snapshot decode "$SNAPSHOT_FILE" | jq -r .Type | uniq)"
mapfile -t SNAPSHOT_TYPES <<<"$TYPES"

for type in "${SNAPSHOT_TYPES[@]}"; do
  consul snapshot decode "$SNAPSHOT_FILE" | jq -r --arg TYPE "$type" 'select(.Type==$TYPE)' >"$DECODE_OUTPUT_DIR"/"$type"_"$DECODE_OUTPUT_IDENTIFIER".json
done
echo "done!"
exit 0