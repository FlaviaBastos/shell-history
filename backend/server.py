from concurrent import futures
import time

import grpc

import history_pb2
import history_pb2_grpc

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


class Historian(history_pb2_grpc.HistorianServicer):

    def GetCommand(self, request, context):
      
        print(request)
        return history_pb2.Response(
          status=history_pb2._STATUS.values_by_name['OK'].name)

def serve():
    with open('certs/localhost.key') as f:
        private_key = f.read()
    with open('certs/localhost.crt') as f:
        certificate_chain = f.read()

    server_credentials = grpc.ssl_server_credentials(
      ((str.encode(private_key), str.encode(certificate_chain),),))
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    history_pb2_grpc.add_HistorianServicer_to_server(Historian(), server)
    server.add_secure_port('[::]:50051', server_credentials)
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()