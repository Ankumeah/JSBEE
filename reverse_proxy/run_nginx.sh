#! /bin/sh

envsubst '$HOSTNAME $BACKEND_PORT $API_VERSION $REVERSE_PROXY_PORT' < /opt/nginx/nginx.conf.template > /opt/nginx/nginx.conf

echo "Nginx running on port $REVERSE_PROXY_PORT"

exec nginx -g "daemon off;" -c /opt/nginx/nginx.conf
