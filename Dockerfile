FROM golang:1.18.1-alpine3.15 as builder

RUN apk update
RUN apk add --update nodejs yarn
RUN mkdir /app

ADD . /app/
WORKDIR /app

RUN yarn install --frozen-lockfile 
RUN yarn build
RUN yarn build:server

FROM alpine:latest  as production
WORKDIR /app/
RUN ls
COPY --from=0 /app/main ./
COPY --from=0 /app/dist ./dist

CMD ["/app/main"]