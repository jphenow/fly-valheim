# Valheim on Fly.io

This is an experimental stab at quickly booting a [Valheim](https://www.valheimgame.com/) server on [fly.io](https://fly.io).

## Getting Started

```shell
$> git clone https://github.com/jphenow/fly-valheim.git
$> cd fly-valheim
$> curl -L https://fly.io/install.sh | sh # if you don't have flyctl yet
$> flyctl auth login
$> flyctl launch
# Choose an app name
# Choose a region near most of your Valheim players
# "No" for postgres db
# "No" for deploy now

# see source for details
# Turns the VM into a dedicated CPU with more RAM
# Sets up volume for storing world/config data
$> ./configure-fly.sh
# When prompted, choose the same region as for your app above

$> flyctl deploy

# Now you can watch logs; it may take a few minutes for Valheim to be ready
# Watch for the message "Waiting for server to listen on UDP port 3456" to stop
$> flyctl logs

# Once its ready, you can connect to valheim using the "Hostname" from:
$> flyctl info
```
