FROM golang:1.19

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY ./ ./

EXPOSE 9091
EXPOSE 8080

CMD ["./broker"]
