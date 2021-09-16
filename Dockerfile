FROM golang:alpine
COPY . /fatalisa-public-api
WORKDIR /fatalisa-public-api
#ENV GIN_MODE release
#EXPOSE 80
RUN go build fatalisa-public-api
#RUN chmod +x /fatalisa-public-api/fatalisa-public-api
#ENTRYPOINT ["/fatalisa-public-api/fatalisa-public-api"]

FROM alpine:latest

COPY --from=0 /fatalisa-public-api/fatalisa-public-api /app
RUN chmod -R +x /app
ENV GIN_MODE release

COPY set-build-date.sh /set-build-date.sh
RUN chmod +x /set-build-date.sh && /set-build-date.sh

EXPOSE 80
ENTRYPOINT ["/app"]