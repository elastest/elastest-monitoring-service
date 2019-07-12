import os
import sys
'''
headers = {'content-type': 'text/plain'}

stampers = "when e.tag(#testresult) do #websocket\n when e.tag(#terminate) do #websocket\n".encode()

moms = "stream bool result_1 := e.strmatch(ourmessage,\"STOP_TEST\")\n stream bool result_2 := e.strmatch(value.integer,\"100\")\n stream bool result_3 := result_1 \/ result_2\n trigger result_3 do emit result_3 on #terminate".encode()

url = "http://" + ems + ":8888/stamper/tag0.1"
response = requests.post(url, headers=headers, data=stampers)
#response = requests.post(url, data=stampers)
print(response.content)

url = "http://" + ems + ":8888/MonitoringMachine/signals0.1"
response = requests.post(url, headers=headers, data=moms)
# # response = requests.post(url, data=moms)
print(response.content)

print "after sending requests"

i = 0

'''
from websocket import create_connection
ems = os.environ["ET_EMS_LSBEATS_HOST"]
url = "ws://" + ems + ":3232"
ws = create_connection(url)
sys.stdout.flush()
while True:
  result = ws.recv()
  if "_mapping [logs]" in result:
    print("Received event")
    sys.stdout.flush()
