FROM node:21-alpine as build

WORKDIR /frontend

COPY . .
RUN npm install

RUN npm run build

FROM nginx:latest

COPY ./frontend.conf.template /etc/nginx/templates/frontend.conf.template
COPY --from=build /frontend/build/ /var/www/html/

