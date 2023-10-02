FROM golang:1.21

LABEL maintainer="erickmayorgar"
LABEL repository="erickmayorgar/so1-p1-agent"

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]

ENV direccionServer 35.215.42.241