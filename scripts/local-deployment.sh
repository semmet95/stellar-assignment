#!/bin/bash
set -ex

cleanup() {
    make cleanup-containers
}
trap cleanup ERR SIGINT SIGTERM

make start-containers