FROM golang:latest
RUN mkdir /fatalisa-public-api
COPY . /fatalisa-public-api
WORKDIR /fatalisa-public-api
RUN go build fatalisa-public-api && ls

FROM alpine:latest

RUN mkdir /app
COPY --from=0 /fatalisa-public-api/fatalisa-public-api /app
RUN chmod -R a+x /app
WORKDIR /app
ENV GIN_MODE release

EXPOSE 80

ENTRYPOINT ["./fatalisa-public-api"]