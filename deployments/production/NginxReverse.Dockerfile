FROM nginx:alpine
COPY third_party/nginx/production.conf /etc/nginx/conf.d/default.conf
