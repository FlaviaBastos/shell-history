FROM python:3-alpine

RUN mkdir -p /app/certs
WORKDIR /app

RUN apk add --no-cache python-dev gcc build-base

COPY requirements.txt requirements.txt

RUN pip install -r requirements.txt


COPY certs/* /app/certs/
COPY backend/* /app/

CMD python server.py
