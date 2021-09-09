FROM golang:latest
RUN mkdir /fatalisa-public-api
COPY . /fatalisa-public-api
WORKDIR /fatalisa-public-api
RUN go build fatalisa-public-api

FROM alpine:latest

RUN mkdir /app && chmod -r +x /app
COPY --from=0 /fatalisa-public-api/fatalisa-public-api /app
WORKDIR /app

EXPOSE 80

ENTRYPOINT ["fatalisa-public-api"]