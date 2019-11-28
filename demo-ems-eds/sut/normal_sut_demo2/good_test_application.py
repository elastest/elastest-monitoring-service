from openmtc_app.onem2m import XAE
from openmtc_onem2m.model import Container
import gevent
import uuid
import os
import signal
import requests
import json
import time
from functools import partial

class TestApplication(XAE):

    def __init__(self, *args, **kw):
        super(TestApplication, self).__init__(*args, **kw)

        self.orch_path = 'onem2m/EDSOrch/edsorch/'
        self.sensor_temp_path = 'onem2m/TemperatureSensor/'
        self.actuator_simple_path = 'onem2m/SimpleActuator/'

        self.NUM_PAIRS = int(os.environ["TESTAPP_NUM_PAIRS"])
        self.stored_reply = {}
        self.sensor_requests = []
        self.actuator_requests = []
        self.app_ID = "testapplication"

        self.app_name = "TestApplication"
        self.ems = os.environ["ET_EMS_LSBEATS_HOST"]
        self.hostport = 'http://' + self.ems + ":8181"

        self.status = {}
        self.starttime = 0

    def __gen_ID(self):
        return uuid.uuid4().hex[:12]

    def __publish(self, message):
        print json.dumps(message)

    def __set_event(self, request):
        self.status = {'request': request, 'event': gevent.event.Event()}

    def __wait_event(self):
        self.status['event'].wait()
        self.status = {'request' : None, 'event' : None}

    def _on_register(self):

        # subscribe to the EDS orch response
        response_path = self.orch_path + 'response'
        self.add_container_subscription(response_path, self.handle_response)

        # subscribe to temperature sensor response
        response_path = self.sensor_temp_path + 'response'
        self.add_container_subscription(response_path, self.handle_response)

        # subscribe to the simple actuator response
        response_path = self.actuator_simple_path + 'response'
        self.add_container_subscription(response_path, self.handle_response)

        gevent.spawn_later(10,self.send_requests)
        self.run_forever()

    def _on_shutdown(self):
        # deregister the application - 4
        request_ID = str('deregister_'+ self.__gen_ID())
        request = [{'deregister': {'application': {'app_ID': self.app_ID, 'request_ID': request_ID}}}]
        request_path = self.orch_path + 'request'
        self.push_content(request_path, request)

    def send_requests(self):
        
        # World creation
        self.starttime = time.time()

        # Register the application
        request_ID = str('app_' + self.__gen_ID())
        self.__set_event(request_ID)
        request = [{'register': {'application': {'app_ID': self.app_ID, 'request_ID': request_ID}}}]
        request_path = self.orch_path + 'request'
        self.push_content(request_path, request)
        self.logger.info('sent request to register application')
        self.__wait_event()

        # Register NUM_PAIRS pairs of temp sensors - actuators
        for _ in range(self.NUM_PAIRS):
            # sensor
            request_ID = str('sensor_temp_' + self.__gen_ID())
            self.__set_event(request_ID)
            request = [{'register': {'sensor': {'app_ID': self.app_ID, 'request_ID': request_ID, 'sensor_type': 'temperature'}}}]
            self.push_content(request_path, request)
            self.sensor_requests.append(request_ID)
            self.logger.info('sent request to register sensor')
            self.__wait_event()
            # actuator
            request_ID = str('actuator_simple_' + self.__gen_ID())
            self.__set_event(request_ID)
            request = [{'register': {'actuator': {'app_ID': self.app_ID, 'request_ID': request_ID, 'actuator_type': 'simple'}}}]
            self.push_content(request_path, request)
            self.actuator_requests.append(request_ID)
            self.logger.info('sent request to register actuator')
            self.__wait_event()

        # Set up pairs
        for index in range(self.NUM_PAIRS):
            # switch on sensor
            request_ID = str('modify_' + self.__gen_ID())
            self.__set_event(request_ID)
            sensor_name = self.stored_reply[self.sensor_requests[index]]['conf']['name']
            request = [{'modify': {'app_ID': self.app_ID, 'request_ID': request_ID, 'name': sensor_name, 'conf': {'onoff':'ON', 'period':5, 'min':10, 'max':30}}}]
            request_path = self.sensor_temp_path + 'request'
            self.push_content(request_path, request)
            self.__wait_event()
            # config actuator
            request_ID = str('modify_' + self.__gen_ID())
            self.__set_event(request_ID)
            actuator_name = self.stored_reply[self.actuator_requests[index]]['conf']['name']
            request = [{'modify':{'app_ID':self.app_ID, 'request_ID': request_ID, 'name' : actuator_name, 'conf':{'delay':3}}}]
            request_path = self.actuator_simple_path + 'request'
            self.push_content(request_path, request)
            self.__wait_event()

        self.logger.info('System should be established...')

        # Subscribe to events
        for index in range(self.NUM_PAIRS):
            # sensor read
            sensor_request = self.sensor_requests[index] 
            self.add_container_subscription(self.stored_reply[sensor_request]['conf']['path'],
                partial(self.handle_temperature_sensor, index=index ))
            # actuator output
            actuator_request = self.actuator_requests[index]
            self.add_container_subscription(self.stored_reply[actuator_request]['conf']['out_path'],
               partial(self.handle_actuator_out, index=index))

        #stop the tjob after 1 minute
        gevent.spawn_later(60, self.app_shutdown)

    def app_shutdown(self):
        timestamp = time.time() - self.starttime
        json_message = {'ourmessage':'STOP_TEST', 'timestamp':timestamp}
        r = requests.post(self.hostport, json=json_message)
        gevent.sleep(1)
        os.kill(os.getpid(), signal.SIGTERM)

    def handle_actuator_out(self, cnt, con, index):
        timestamp = time.time() - self.starttime
        self.logger.info(':actuator: index %d value %s - time %s' % (index, float(con), timestamp))
        json_message = {'appname':'test1', 'type':'actuator', 'id':index, 'timestamp':timestamp }
        r = requests.post(self.hostport, json=json_message)

    def handle_temperature_sensor(self, cnt, con, index):
        timestamp = time.time() - self.starttime
        actuator_request = self.actuator_requests[index] 
        self.logger.info(':sensor: index %d value %s - time %s' % (index, float(con), timestamp))
        json_message = {'appname':'test1', 'type':'sensor', 'id':index, 'svalue':{'actual':float(con), 'threshold':20}, 'timestamp':timestamp}
        r = requests.post(self.hostport, json=json_message)
        if float(con) > 20:
            self.push_content(self.stored_reply[actuator_request]['conf']['in_path'], con)

    def handle_response(self, cnt, con):
        reply = con
        self.logger.info('EDS response from ' + cnt )
        # check if reply is for this application
        if reply.get('app_ID') == self.app_ID:
            # check the result in the reply
            request_ID = reply['request_ID']
            if reply.get('result') == 'SUCCESS':
                self.stored_reply[request_ID] = reply
                self.logger.info(request_ID + ' was a success')
            else:
                error = reply['error_string']
                self.logger.info(request_ID + ' did not succeed')
                self.logger.info('error ' + error_string)

        else:
            self.logger.info('received message not for this app')
            self.logger.info(reply)
        # if the program was waiting for the reply, notify it
        if self.status.get('request'):
            if self.status['request'] == reply.get('request_ID'):
                self.status['event'].set()


