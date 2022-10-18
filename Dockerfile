FROM golang:1.19.2-alpine3.16

RUN mkdir /app
ADD src/ /app
WORKDIR /app
ENV CGO_ENABLED=0

RUN go build -o main .

CMD ["/app/main"]