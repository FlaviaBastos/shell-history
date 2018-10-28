FROM phusion/baseimage

CMD ["/sbin/my_init"]

RUN mkdir -p /app/certs 
WORKDIR /app

RUN apt-get update
RUN apt-get install -y python3-dev gcc python3-pip

COPY requirements.txt requirements.txt
COPY history/history.proto history/history.proto

RUN pip3 install --upgrade pip
RUN pip3 install -r requirements.txt
RUN python3 -m grpc_tools.protoc -Ihistory/ --python_out=/app/ --grpc_python_out=/app/ history/history.proto 

COPY certs/* /app/certs/
COPY backend/* /app/
COPY start_server.sh /etc/service/server/run
COPY start_django.sh /etc/service/django/run

ENV PORT 50051

EXPOSE 80
EXPOSE 50051
