FROM nginx:alpine as build
COPY deployments/nginx/base.conf /etc/nginx/conf.d/default.conf

FROM build as prod
RUN sed -i 's/<profile>/prod/g' /etc/nginx/conf.d/default.conf

FROM build as test
RUN sed -i 's/<profile>/test/g' /etc/nginx/conf.d/default.conf

FROM build as dev
RUN sed -i 's/<profile>/dev/g' /etc/nginx/conf.d/default.conf

FROM dev as testing