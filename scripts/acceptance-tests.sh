#!/bin/bash
set -ex

export MODBUS_HOST='127.0.0.1'
export MODBUS_PORT='5020'
export INFLUX_HOST='127.0.0.1'
export INFLUX_PORT='8086'

make acceptance-tests