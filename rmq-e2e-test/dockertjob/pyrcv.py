#!/usr/bin/env python
import sys
import pika
import json
import os

credentials = pika.PlainCredentials('user', 'password')
parameters = pika.ConnectionParameters(os.environ["ET_SUT_HOST"],
                                       5672,
                                       '/',
                                       credentials)
connection = pika.BlockingConnection(parameters)
channel = connection.channel()

channel.queue_declare(queue='thekey')

def callback(ch, method, properties, body):
    res = json.loads(body)
    if "index" in res:
        print(" [x] Received input event %r" % res["index"])
    else:
        print(" [x] Received generated event with value %r" % res["value"])
    sys.stdout.flush()

channel.basic_consume(
    queue='thekey', on_message_callback=callback, auto_ack=True)

print('Connected to RabbitMQ')
sys.stdout.flush()
channel.start_consuming()
