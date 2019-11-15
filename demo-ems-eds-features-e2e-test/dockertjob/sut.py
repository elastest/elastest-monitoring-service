import socket
import sys
import logging
import json
import requests

#request will have the form: operation_arg
def treatRequest(request, logged_users):
    logging.info("treating request %s with logged_users: %s", request, logged_users)
    j = json.loads(request)
    logging.info("json: %s", json.dumps(j))
    operation = j["operation"]
    args = j["args"]
    r = ""
    if (operation == "login"):
        r = login(args[0], logged_users)
    elif (operation == "test"):
        r = test(args[0], args[1], logged_users)
    elif (operation == "logout"):
        r = logout(args[0], logged_users)
    elif (operation == "exit"):
        r = {"result": "exit"}
    else:
        r = {"result": "Operation not supported"}
    logging.info("response: %s", json.dumps(r))
    return json.dumps(r) + "\n"

#idempotent
def login(usr, logged_users):
    logging.info("login usr: %s, %s", usr, logged_users)
    c = "cookie/" + usr
    logged_users[c] = usr
    response = {
        "usr" : usr,
        "cookie" : c,
        "result": "ok"
    }
    emsResponse = {
        "usr" : usr,
        "cookie" : c,
        "emptyStr": "",
        "op": "login",
        "arg": usr
    }
    informEMS(emsResponse)
    return response
#requires login
def test(test_name, cookie, logged_users):
    logging.info("testing test: %s with cookie: %s from logged_users: %s", test_name, cookie, logged_users)
    r = ""
    if cookie in logged_users:
        r = "ok"
    elif cookie:
        r = "Error, user with cookie: "+cookie+" not logged in"
    else:
        r = "Error, user without cookie"
    response = {
        "test": test_name,
        "cookie": cookie,
        "result" : r
    }
    emsResponse = {
        "op": "test",
        "arg": cookie
    }
    informEMS(emsResponse)
    return response

#requires login
def logout(cookie, logged_users):
    logging.info("deleting user with cookie: %s from logged_users: %s", cookie, logged_users)
    try:
        usr = logged_users[cookie]
        del logged_users[cookie]
        r = "ok" #success
    except KeyError as k:
        r = "Error: cannot logout user that has not logged in"
        usr = "not found"
    response = {
        "usr": usr,
        "cookie" : cookie,
        "result" : r
    }
    emsResponse = {
        "usr" : usr,
        "cookie" : "",
        "op": "logout",
        "arg": cookie
    }
    informEMS(emsResponse)
    return response

def informEMS(jsonObj):
    url =  "http://" + emsIp + ":"+str(emsPort) #8888 no! that's for specs!, to feed events use port 8181 for http events
    logging.info("new event %s to EMS in %s ", json.dumps(jsonObj), url)
    resp = requests.post(url, json=jsonObj)
    logging.info("EMS responded: %s : %s", resp, resp.content)
    #return resp
    

#MAIN
logging.basicConfig(format='%(levelname)s - %(asctime)s - %(message)s', level=logging.INFO, filename='./sut_log.txt', filemode='w')
encoding = "utf-8"
# Create a TCP/IP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
if len(sys.argv) < 4:
    print("Need arguments: port to listen, emsIp, emsPort")
    exit(1)
port = int(sys.argv[1])
emsIp = sys.argv[2]
emsPort = int(sys.argv[3]) #8181
# Bind the socket to the address given on the command line
try:
    server_address = ('', port)
    sock.bind(server_address)
    logging.info('starting up on %s port %s', sock.getsockname()[0], sock.getsockname()[1])
    sock.listen(1)
    logged_users = {} #map cookie id, state of the server
    while True: #always accept new connections
        logging.info('waiting for a connection')
        connection, client_address = sock.accept()
        logging.info('client connected: %s', client_address)
        try:
            while True: #always treat all data that a connection sends
                logging.info('waiting for data')
                data = ""
                last = ""
                end = "\n"
                while last != end: #each msg should end with "\n"
                    logging.info("expecting more data")
                    data += connection.recv(16).decode(encoding)
                    if data:
                        last = data[len(data)-1]
                    else: #there is no more data if data was empty
                        last = end
                logging.info('received "%s"', data)
                if data:
                    r = treatRequest(data, logged_users)
                    logging.info("sending response %s", r)
                    connection.sendall(r.encode(encoding))
                else:
                    break
        except Exception as e:
            logging.info("exception")
            logging.exception(e)
        finally:
            connection.close()
            logging.info("connection closed to: %s", client_address)

finally:
    #sock.shutdown(socket.SHUT_RDWR)
    sock.close()
    logging.info("socket closed")



"""
python -m py_compile sut.py
chmod +x ./sut.pyc
python ./sut.pyc
"""
#./rebuild_push.sh sut
#docker run --name orch_sut -p 10000:10000 -d --rm luismigueldanielsson/elastest-luismi:ems_orchestration_sut

#python sut.py 10000 localhost 8181
