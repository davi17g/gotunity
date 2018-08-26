#!/usr/bin/python3
import sys
import socket
import time
import random
import argparse


MSG_SIZE = 1024


def check_args(args=None):
    parser = argparse.ArgumentParser(description="Generates random UDP every 100 miliseconds")
    parser.add_argument('-H', '--host',
                        help='Specify host address, if not specified default ip will be assigned',
                        default='localhost')
    parser.add_argument('-p', '--port',
                        help='Specify desired port address, if not specified default port will be assigned',
                        default='8888')

    result = parser.parse_args(args)
    return result.host, result.port


def main():
    host, port = check_args(sys.argv[1:])
    soc = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    while True:
        msg = bytearray(random.getrandbits(8) for _ in range(MSG_SIZE))
        time.sleep(0.1)
        soc.sendto(msg, (host, int(port)))


if __name__ == '__main__':
    main()