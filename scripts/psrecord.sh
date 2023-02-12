#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"/.. || exit 1

if [ $# -lt 3 ]; then
    echo "Usage: $0 <pid> <fx|nofx> <name> [duration]"
    exit 1
fi

duration=$4
if [[ -z "$duration" ]]; then
  duration=$(( 60 * 15 ))
fi

echo "Recording $1 for $duration seconds"
exit

psrecord $1 \
  --include-children \
  --interval 1 \
  --duration ${duration} \
  --plot \
  benchmarks/$2-cpu-memory-plot-${duration}.$3.png
