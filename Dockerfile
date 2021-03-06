FROM golang:1.13-alpine AS certgen

RUN apk add git
WORKDIR /tmp
RUN cd /tmp && git clone https://github.com/square/certstrap && cd certstrap && go build
RUN /tmp/certstrap/certstrap --depot-path /tmp/certs init --passphrase "" --common-name localhost

FROM phusion/baseimage

CMD ["/sbin/my_init"]

RUN mkdir -p /app/certs
WORKDIR /app

RUN apt update
RUN apt install -y python3-dev gcc python3-pip libmysqlclient-dev

COPY requirements.txt requirements.txt
COPY history/history.proto history/history.proto

RUN pip3 install --upgrade pip
RUN pip3 install -r requirements.txt
RUN python3 -m grpc_tools.protoc -Ihistory/ --python_out=/app/ --grpc_python_out=/app/ history/history.proto

COPY --from=certgen /tmp/certs/* /app/certs/
COPY backend/* /app/
COPY start_server.sh /etc/service/server/run
COPY start_django.sh /etc/service/django/run

ENV PORT 50051

EXPOSE 80
EXPOSE 50051
