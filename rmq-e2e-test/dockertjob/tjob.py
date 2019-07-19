import requests
import time
import sys
import os
import json

jsonheaders = {'content-type': 'application/json'}
plainheaders = {'content-type': 'text/plain'}
def sendEv(evdata):
  url = os.environ["ET_EMS_HTTPINEVENTS_API"]
  response = requests.post(url, headers=jsonheaders, data=json.dumps(evdata))
  # print(response)
  sys.stdout.flush()

event_i=0
def eventsBurst():
  global event_i
  simpleoutdata = {"index": event_i, "channels": ["#outchan"]}
  event_i+=1
  simpleindata = {"index": event_i, "channels": ["#inchan"]}
  event_i+=1
  filterinilldata = {"index": event_i, "payload": 1, "channels": ["#filterinchan"]}
  event_i+=1
  filterindata = {"index": event_i, "payload": 10, "channels": ["#filterinchan"]}
  event_i+=1
  sendEv(simpleoutdata)
  sendEv(simpleindata)
  sendEv(filterinilldata)
  sendEv(filterindata)

# Adding rabbitmq endpoint
print("Adding rabbitmq endpoint")
 # curl -H "Content-Type: application/json" -XPOST "http://beastest.software.imdea.org:8888/subscriber/rabbitmq" -d '{"ip":"beastest.software.imdea.org", "port":5672, "user":"user", "password":"password", "key":"myexchange", "exchange_type":"fanout"}'
url = os.environ["ET_EMS_API"] + "subscriber/rabbitmq"
evdata = {
        "ip" : os.environ["ET_SUT_HOST"],
        "port": 5672,
        "user": "guest",
        "password": "guest",
        "channel": "#outchan",
        "key": "thekey",
        "exchange_type": "fanout"
        }
requests.post(url, headers=jsonheaders, data=json.dumps(evdata)).text
print("Deploying MoM and Stamper")
url = os.environ["ET_EMS_API"] + "MonitoringMachine/signals0.1"
with open('momdef.txt','rb') as payload:
    momid = requests.post(url, headers=plainheaders, data=payload).text
momid=momid.replace('"', '')
url = os.environ["ET_EMS_API"] + "stamper/tag0.1"
with open('stamperdef.txt','rb') as payload:
    stamperid = requests.post(url, headers=plainheaders, data=payload).text
stamperid=stamperid.replace('"', '')
time.sleep(5)
print("Sending event with stamper and mom")
sys.stdout.flush()
eventsBurst()
time.sleep(5)
# Removing MoM
print("Removing MoM")
sys.stdout.flush()
url = os.environ["ET_EMS_API"] + "MonitoringMachine/"+momid
# requests.delete(url)
os.system('curl -o /dev/null -s -XDELETE '+url)
time.sleep(5)
print("Sending event with no mom")
sys.stdout.flush()
eventsBurst()
time.sleep(5)
# Removing stamper
print("Removing stamper")
sys.stdout.flush()
url = os.environ["ET_EMS_API"] + "stamper/"+stamperid
# requests.delete(url)
os.system('curl -o /dev/null -s -XDELETE '+url)
time.sleep(5)
print("Sending event with no mom or stamper")
sys.stdout.flush()
eventsBurst()
time.sleep(5)
os.system('echo Done')
