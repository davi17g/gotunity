#!/usr/bin/python3
import sys
import time
import socket
import argparse
import threading


MSG_SIZE = 1024
CONNECT_MSG = "CONNECT\r\n"
DISCONNECT_MSG = "DISCONNECT\r\n"
recv_bytes = 0
lock = threading.Lock()


class TrafficCapacityChecker(threading.Thread):
    def __init__(self):
        threading.Thread.__init__(self)

    def run(self):
        global recv_bytes
        global lock
        while True:
            time.sleep(1)
            lock.acquire()
            print("%d bytes\sec" % recv_bytes)
            recv_bytes = 0
            lock.release()


class Listener(threading.Thread):
    def __init__(self, sock):
        threading.Thread.__init__(self)
        self.sock = sock

    def run(self):
        global recv_bytes
        global lock
        while True:
            data, server = self.sock.recvfrom(MSG_SIZE)
            lock.acquire()
            recv_bytes += len(data)
            lock.release()


def check_args(args=None):
    parser = argparse.ArgumentParser(description="Connects to the server and receives server messages")
    parser.add_argument('-H', '--host',
                        help='Specify host address, if not specified default ip will be assigned',
                        default='127.0.0.1')
    parser.add_argument('-p', '--port',
                        help='Specify desired port address, if not specified default port will be assigned',
                        default='8089')

    result = parser.parse_args(args)
    return result.host, result.port


def main():
    host, port = check_args(sys.argv[1:])
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
    t1 = TrafficCapacityChecker()
    t2 = Listener(sock=sock)
    t1.setDaemon(True)
    try:
        sock.sendto(CONNECT_MSG.encode("utf-8"), (host, int(port)))
        t1.start()
        t2.start()
        t1.join()
        t2.join()
    except KeyboardInterrupt:
        sock.sendto(DISCONNECT_MSG.encode("utf-8"), (host, int(port)))

    finally:
        sock.close()

if __name__ == '__main__':
    main()