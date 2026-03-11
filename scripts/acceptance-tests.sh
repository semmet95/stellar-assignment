#!/bin/bash
set -ex

cleanup() {
    make cleanup-containers
}
# trap cleanup EXIT ERR SIGINT SIGTERM

make start-containers