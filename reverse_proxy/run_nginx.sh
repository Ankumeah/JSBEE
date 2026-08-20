#! /bin/sh

envsubst '$HOSTNAME $PORT' < /opt/nginx/nginx.conf.template > /opt/nginx/nginx.conf

nginx -g "daemon off;" -c /opt/nginx/nginx.conf
