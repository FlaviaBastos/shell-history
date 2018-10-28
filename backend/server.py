from concurrent import futures
import logging
import os
import sys
import time
import pytz
from datetime import datetime

import django
import grpc

import history_pb2
import history_pb2_grpc


from django.db import models

_ONE_DAY_IN_SECONDS = 60 * 60 * 24
logging.basicConfig(format='%(asctime)s - %(message)s', level=logging.INFO)
logger = logging.getLogger(__name__)


class Historian(history_pb2_grpc.HistorianServicer):

    def GetCommand(self, request):

        timestamp = datetime.fromtimestamp(request.timestamp, pytz.utc)
        command = ' '.join(request.command)
        new_command = Command(hostname=request.hostname, timestamp=timestamp,
                              username=request.username, oldpwd=request.oldpwd,
                              altusername=request.altusername, cwd=request.cwd,
                              command=command, exitcode=request.exitcode)
        new_command.save()
        return history_pb2.Response(
          status=history_pb2._STATUS.values_by_name['OK'].name)


def serve(port):
    logger.info('Starting shell-history backend service')
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    try:
        with open('certs/localhost.key') as f:
            private_key = f.read()
        with open('certs/localhost.crt') as f:
            certificate_chain = f.read()
        server_credentials = grpc.ssl_server_credentials(
        server.add_secure_port('[::]:{}'.format(port), server_credentials)
        ((str.encode(private_key), str.encode(certificate_chain),),))
        logger.info('Listening on secure port [::]:{}'.format(port))
    except:
        logger.info('Listening on insecure port [::]:{}'.format(port))
        server.add_insecure_port('[::]:{}'.format(port))

    history_pb2_grpc.add_HistorianServicer_to_server(Historian(), server)

    logger.info('Ready')
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    sys.path.append(os.path.join(os.getcwd(),"shell"))
    os.environ.setdefault("DJANGO_SETTINGS_MODULE", "shell.settings")

    # Setup django
    django.setup()
    from history.models import Command
    port = os.environ.get("PORT", 50051)
    serve(port=port)
