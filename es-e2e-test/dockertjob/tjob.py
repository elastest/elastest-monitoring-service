import requests
import time
import sys
import os
import json

headers = {'content-type': 'application/json'}
event_i=0
def sendEv():
  global event_i
  url = os.environ["ET_EMS_HTTPINEVENTS_API"]
  evdata = { "index"+str(event_i) : 0, "channels": ["eschan"]}
  response = requests.post(url, headers=headers, data=json.dumps(evdata))
  # print(response)
  sys.stdout.flush()
  event_i+=1

print("Sending event with no subscriber. Nothing should happen")
sys.stdout.flush()
sendEv()
sendEv()
sendEv()
sendEv()
sendEv()
time.sleep(5)
# Adding elasticsearch endpoint
print("Adding elasticsearch endpoint")
url = os.environ["ET_EMS_API"] + "subscriber/elasticsearch"
evdata = { "ip" : os.environ["ET_SUT_HOST"],
        "port": int(os.environ["ET_SUT_PORT"]),
        "user": "esuser",
        "password": "espassword",
        "channel": "eschan"}
esid = requests.post(url, headers=headers, data=json.dumps(evdata)).text
esid=esid.replace('"', '')
sys.stdout.flush()
time.sleep(5)
print("Sending event with subscriber. Something should happen")
sys.stdout.flush()
sendEv()
sendEv()
sendEv()
sendEv()
sendEv()
time.sleep(5)
# Removing elasticsearch endpoint
print("Removing elasticsearch endpoint")
sys.stdout.flush()
url = os.environ["ET_EMS_API"] + "subscriber/"+esid
# requests.delete(url)
os.system('curl -o /dev/null -s -XDELETE '+url)
time.sleep(5)
print("Sending event with no subscriber. Nothing should happen")
sys.stdout.flush()
sendEv()
sendEv()
sendEv()
sendEv()
sendEv()
time.sleep(5)
os.system('echo Done')
