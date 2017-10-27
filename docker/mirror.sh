#!/bin/sh

# This is required to get java to handle signals correctly (and stop immediately
# instead of waiting around for a KILL).
trap 'kill -TERM $PID' TERM INT

java -Xmx2G -XX:+HeapDumpOnOutOfMemoryError \
  -cp /mirror/mirror-all.jar mirror.Mirror "$@" &

PID=$!
wait $PID
trap - TERM INT
wait $PID

exit $?
