FROM nginx:alpine
COPY third_party/nginx/develop.conf /etc/nginx/conf.d/default.conf
