FROM golang as backend
RUN apt update && apt install -y libglib2.0
WORKDIR /build
COPY server server
COPY vendor vendor
COPY config config
COPY go.mod go.sum Makefile ./
RUN make build_init && make build_backend

FROM node:12 as frontend

WORKDIR /build
COPY Makefile package.json webpack.config.js .eslintrc.json .babelrc ./
RUN npm i
COPY client client
RUN make build_frontend

FROM debian:stable-slim

WORKDIR /usr/src/app
RUN apt-get update; apt install -y libglib2.0
COPY --from=backend /build/dist/* ./
COPY --from=frontend /build/dist/* ./
COPY config ./data/state/config

EXPOSE 8334
VOLUME ["/usr/src/app/data/state/"]

ENTRYPOINT ["./filestash"]
