import redis
import time
import random

rdb = redis.Redis(host='127.0.0.1', port=6379, db=1, password='wait_123456')


fruits = ['Apple', 'Pear', 'Grape', 'Strawberry', 'Banan']

def run():
    for i in range(10):

        i = random.randint(0, len(fruits) - 1)
        message = fruits[i]
        print("count: ", i, " --> ", message)
        rdb.publish(channel='mq-fruit', message=message)
        
        time.sleep(1)
    rdb.publish(channel='mq-fruit', message='QUIT')
    
run()