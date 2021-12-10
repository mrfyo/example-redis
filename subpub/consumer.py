import redis

rdb = redis.Redis(host='127.0.0.1', port=6379, db=1, password='wait_123456')


def run():
    p = rdb.pubsub()
    p.subscribe(['mq-fruit']) # 订阅指定通道列表
    # p.psubscribe("mq-*")
    
    p.listen()

    for item in p.listen():
        print(item)
        if item['data'] == b'QUIT':
            break
    
    p.unsubscribe()
    p.close()

run()