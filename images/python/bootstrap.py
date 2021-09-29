#!/usr/bin/python3
import sys
import os
import json
from importlib import import_module

handler_function = os.getenv("LAMBDA_HANDLER_FUNCTION")
function_name = os.getenv("LAMBDA_FUNCTION_NAME")
logfile = os.getenv("BOOTSTRAP_LOG_FILE")

log = open(logfile, 'w')

handler_file_name = "code." + ".".join(handler_function.split(".")[:-1])
handler_function_name = handler_function.split(".")[-1]

log.write(json.dumps({'module':handler_file_name, 'method': handler_function_name, 'name':function_name}) + "\n")

output = os.fdopen(3, 'w')

module = import_module(handler_file_name)
handler_method = getattr(module, handler_function_name)

for event in sys.stdin:
    response = {}
    try:
        eventData = json.loads(event)['eventData']
        response = handler_method(json.dumps(eventData, indent=4), "")
        output.write(json.dumps(response) + "\n")
    except Exception as e:
        log.write(json.dumps({'error': str(e)}) + "\n")
    finally:
        log.write(json.dumps({'event': event, 'response': response}) + "\n")

output.close()


