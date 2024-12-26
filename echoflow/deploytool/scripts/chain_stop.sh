#! /bin/bash
set -euo pipefail

CHAINMAKER_GO_PATH=`dirname $(dirname ${PWD})`/chainmaker-csv-adapter

cd ${CHAINMAKER_GO_PATH}/scripts && ./cluster_quick_stop.sh
# echo "y" | docker container prune
