#!/usr/bin/env bash

set -euo pipefail

mkdir -p /volume/data
mkdir -p /volume/config

realDataLocation=/opt/valheim
rmdir $realDataLocation
ln -s /volume/data/ $realDataLocation

realConfigLocation=/home/valheim/.config/unity3d/IronGate/Valheim
metaConfigLocation=/config
rm -rf $realConfigLocation
[ -d "$metaConfigLocation" ] && rm -rf $metaConfigLocation
ln -s /volume/config/ $realConfigLocation
ln -s /volume/config/ $metaConfigLocation

udp-proxy -p 2456 -P 3456 &

/usr/local/sbin/bootstrap
