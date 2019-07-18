#!/usr/bin/env python
import sys
import pika

credentials = pika.PlainCredentials('user', 'password')
parameters = pika.ConnectionParameters('localhost',
                                       5672,
                                       '/',
                                       credentials)
connection = pika.BlockingConnection(parameters)
channel = connection.channel()

channel.queue_declare(queue='thekey')

def callback(ch, method, properties, body):
    print(" [x] Received %r" % body)
    sys.stdout.flush()

channel.basic_consume(
    queue='thekey', on_message_callback=callback, auto_ack=True)

print('Connected to RabbitMQ')
sys.stdout.flush()
channel.start_consuming()
