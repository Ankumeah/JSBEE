#! /bin/sh

set -eu

echo "STARTING INIT"
/opt/nginx/init
echo "COMPLETED INIT"

echo "STARTING API SERVER"
/opt/nginx/api-server &
sleep 3

echo "STARTING NGINX"
/opt/nginx/run_nginx.sh
