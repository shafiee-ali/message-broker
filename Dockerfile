FROM golang:1.19

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /broker

EXPOSE 9091
EXPOSE 8080

CMD ["/broker"]
