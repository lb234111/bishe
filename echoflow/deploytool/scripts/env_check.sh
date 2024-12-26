#! /bin/bash

set -euo pipefail

workspace=${PWD}/`basename $(dirname $0)`
source ${workspace}/log.sh

command_check(){
    command -v $1 > /dev/null 2>&1 && info "check $1 installed." || { error "command $1 not exist";exit 1; }
}

command_check docker
command_check go
command_check python3
command_check docker-compose
command_check gcc
command_check 7z