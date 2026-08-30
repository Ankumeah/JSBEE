#! /bin/sh

envsubst '$HOSTNAME $PORT $API_VERSION' < /opt/nginx/nginx.conf.template > /opt/nginx/nginx.conf

nginx -g "daemon off;" -c /opt/nginx/nginx.conf
