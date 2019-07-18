FROM python
RUN pip install pika
RUN pip install requests
COPY tjob.py /tjob.py
COPY pyrcv.py /pyrcv.py
COPY momdef.txt /momdef.txt
COPY stamperdef.txt /stamperdef.txt
COPY expected_out.txt /expected_out.txt
COPY check.sh /check.sh
ENTRYPOINT ["/check.sh"]
