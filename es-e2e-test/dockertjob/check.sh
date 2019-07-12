#!/bin/bash
python3 wscli.py >> out.txt &
python3 tjob.py >> out.txt
cat out.txt
echo -=-=-=-=-=-=-=-=-=-=-=-=-
diff expected_out.txt out.txt
