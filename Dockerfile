FROM golang:alpine
COPY . /fatalisa-public-api
WORKDIR /fatalisa-public-api
ENV GIN_MODE release
EXPOSE 80
RUN go run fatalisa-public-api

#FROM alpine:latest

#COPY --from=0 /fatalisa-public-api/fatalisa-public-api /app
#RUN chmod -R +x /app
#ENV GIN_MODE release
#ENV TZ Asia/Jakarta

#EXPOSE 80

#ENTRYPOINT ["/app"]