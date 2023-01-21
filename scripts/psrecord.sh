#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"/.. || exit 1

if [ $# -ne 2 ]; then
    echo "Usage: $0 <pid> <fx|nofx>"
    exit 1
fi

duration=$(( 60 * 5 ))

psrecord $1 --include-children --interval 1 --duration ${duration} --plot benchmarks/$2-cpu-memory-plot-${duration}.$(date +%Y%m%d).png
