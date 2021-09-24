
def lambda_handler(event, context):
    print(event)
    return {
        'statusCode': 200,
        'message': 'Hello from the test lambda!'
    }