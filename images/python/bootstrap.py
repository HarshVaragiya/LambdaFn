import sys
import os

log = open("../../bin/logfile.txt", 'a')

from example import lambda_handler

for event in sys.stdin:
    log.write("Input Event: " + event)
    response = lambda_handler(event, "")
    log.write(f"Response : {response}")
    print(response)


