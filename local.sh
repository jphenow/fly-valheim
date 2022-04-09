#!/usr/bin/env bash

set -euo pipefail

docker build -t vheim . && docker run -it --publish-all --env-file local.env --rm --name vheim -v "${PWD}/vol:/volume" vheim
