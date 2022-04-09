#!/usr/bin/env bash

set -exuo pipefail

REGION="${1:-ord}"

flyctl scale vm dedicated-cpu-1x
flyctl scale memory 4096

flyctl volumes create flyvalheim -s 10

