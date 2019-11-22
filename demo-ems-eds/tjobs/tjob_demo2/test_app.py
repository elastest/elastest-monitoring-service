import requests
import os
from threading import Timer
import json
from websocket import create_connection
import sys
import unittest
import xmlrunner
import time
 
testSuccess = False
   
class TestCorrectness(unittest.TestCase):
  def test_correct(self):
    self.assertTrue(testSuccess, "Test Failed" )

testSuite = unittest.TestLoader().loadTestsFromTestCase(TestCorrectness)

class MonitoringTest():
  def __init__(self):

    self.condition = True
    self.ems = os.environ["ET_EMS_LSBEATS_HOST"]
    self.headers = {'content-type': 'text/plain'}
    self.stampers = ""
    self.monMachines = ""
    print("before sending requests")

    # get the stampers from file
    with open(os.environ['PWD'] + "/" + "stampers.txt") as f:
      self.stampers = f.read()

    # send stampers to EMS
    url = "http://" + self.ems + ":8888/stamper/tag0.1"
    response = requests.post(url, headers=self.headers, data=self.stampers)
    print(response.content)

    # get the monitoring machines from the file
    NUM_PAIRS = os.environ["TESTAPP_NUM_PAIRS"]
    with open(os.environ['PWD'] + "/" + "monitoring_machines.txt") as f:
      self.monMachines = f.read()
    self.monMachines = self.monMachines.replace("NUM_PAIRS", NUM_PAIRS)

    # send the monitoring machines to EMS
    url = "http://" + self.ems + ":8888/MonitoringMachine/signals0.1"
    response = requests.post(url, headers=self.headers, data=self.monMachines)
    print(response.content)

    print("after sending requests")
    sys.stdout.flush()

    self.start_test()

  def start_test(self):
    url = "ws://" + self.ems + ":3232"
    ws = create_connection(url)
    self.condition = True
    while(self.condition):
      result = ws.recv()
      result = json.loads(result)
      print result
      
      if "#terminatetest" in result["channels"]:
        print "test result found"
        self.condition = False
        if str(result.get('testCorrect', '')).lower() == 'true':
            global testSuccess
            testSuccess = True
        break
    print testSuccess
    xmlrunner.XMLTestRunner(verbosity=0, output='/tmp/test-reports').run(testSuite) 
    return True

if __name__ == "__main__":
  print("Starting the test")
  try:
      edstest = MonitoringTest()
  except Exception, e:
    print e
  print("Ending the test")

