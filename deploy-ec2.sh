#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <ec2-host> [ssh-key]" >&2
  exit 1
fi

HOST="$1"
KEY="${2:-}"

SSH_OPTS="-o StrictHostKeyChecking=no"
if [ -n "$KEY" ]; then
  SSH_OPTS="$SSH_OPTS -i $KEY"
fi

# Build the server image
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)

docker build -t squirrel-server -f "$SCRIPT_DIR/infra/Dockerfile" "$SCRIPT_DIR"
docker pull postgres:15

# Save images to tar files
server_tar=$(mktemp /tmp/server.XXXXXX.tar)
postgres_tar=$(mktemp /tmp/postgres.XXXXXX.tar)

docker save -o "$server_tar" squirrel-server
docker save -o "$postgres_tar" postgres:15

# Copy images to the remote host
scp $SSH_OPTS "$server_tar" "$postgres_tar" "$HOST:/tmp/"

# Clean up local temporary files
rm -f "$server_tar" "$postgres_tar"

# Remote setup and container run
ssh $SSH_OPTS "$HOST" <<'EOSSH'
set -e
if ! command -v docker >/dev/null 2>&1; then
  sudo yum update -y
  sudo amazon-linux-extras install docker -y
  sudo service docker start
fi

sudo docker load -i /tmp/postgres*.tar
sudo docker load -i /tmp/server*.tar

sudo docker network create squirrel_net >/dev/null 2>&1 || true
sudo docker rm -f squirrel_server squirrel_db >/dev/null 2>&1 || true
sudo docker run -d --name squirrel_db --network squirrel_net \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=squirrel \
  -p 5432:5432 postgres:15
sudo docker run -d --name squirrel_server --network squirrel_net \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://postgres:postgres@squirrel_db:5432/squirrel?sslmode=disable \
  squirrel-server
EOSSH

echo "Deployment complete"
