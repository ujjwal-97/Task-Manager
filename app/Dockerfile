FROM golang:1.16

WORKDIR /go/src/app

COPY . .

ENV DB_HOST=db \
    DB_PORT=27017 \
    DB_USERNAME=admin \
    DB_PASSWORD=password \
    DB_NAME=taskmanager \
    hostip=192.168.29.34:22\
    hostusername=ujjwal\
    hostpassword=ujjwal123\
    healthcheckScript=healthCheck.sh

RUN go get -d -v ./...

RUN go install -v ./...

CMD ["app"]