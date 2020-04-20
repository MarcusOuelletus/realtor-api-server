FROM golang:1.14-alpine

EXPOSE 9090

ENV RETS_SERVER_PATH=/rets-server

WORKDIR /rets-server

COPY . .

RUN GO111MODULE=on go get ./...

RUN go build -o main

CMD [ "./main", "-log" ]