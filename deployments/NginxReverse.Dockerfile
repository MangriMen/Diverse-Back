FROM nginx:alpine as build
COPY third_party/nginx/base.conf /etc/nginx/conf.d/default.conf

FROM build as prod
RUN sed -i 's/<PROFILE>/prod/g' /etc/nginx/conf.d/default.conf

FROM build as test
RUN sed -i 's/<PROFILE>/test/g' /etc/nginx/conf.d/default.conf

FROM build as dev
RUN sed -i 's/<PROFILE>/dev/g' /etc/nginx/conf.d/default.conf
