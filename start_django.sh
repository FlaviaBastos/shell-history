#!/bin/sh
# Starts the django server

cd /app
python3 manage.py migrate
python3 manage.py runserver 0.0.0.0:80