#!/usr/bin/env bash

set -euo pipefail

./clean.sh
flyctl deploy
