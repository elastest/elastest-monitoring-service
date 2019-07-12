FROM python
RUN pip install websocket-client
RUN pip install requests
COPY tjob.py /tjob.py
COPY wscli.py /wscli.py
COPY expected_out.txt /expected_out.txt
COPY check.sh /check.sh
ENTRYPOINT ["/check.sh"]
