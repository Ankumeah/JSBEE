FROM nginx:alpine@sha256:289decab414250121a93c3f1b8316b9c69906de3a4993757c424cb964169ad42

COPY  . /opt/nginx/
RUN chmod +x /opt/nginx/run_nginx.sh

ENTRYPOINT [ "/opt/nginx/run_nginx.sh" ]
