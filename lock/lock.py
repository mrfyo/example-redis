import redis
import uuid
import time
import math

def acquire_lock(rdb: redis.Redis, lockname: str, acquire_timeout=10):
    identifier = str(uuid.uuid4())

    end = time.time() + acquire_timeout
    while time.time() < end:
        if rdb.setnx('lock:' + lockname, identifier):
            return identifier
        time.sleep(0.001)
    
    return False


def release_lock(rdb: redis.Redis, lockname: str, identifier: str):
    pipe = rdb.pipeline(True)
    lockname = 'lock:' + lockname

    while True:
        try:
            pipe.watch(lockname)
            if pipe.get(lockname) == identifier:
                # 确认键没有被删除
                pipe.multi()
                pipe.delete(lockname)
                pipe.execute()
                return True
            pipe.unwatch()
            break
        except redis.exceptions.WatchError:
            pass
    
    return False


def acquire_lock_with_timeout(rdb: redis.Redis, lockname: str, acquire_timeout=10, lock_timeout=10):
    """
        带超时的锁（自动释放）
    """

    identifier = str(uuid.uuid4())
    lockname = 'lock:' + lockname
    lock_timeout = int(math.ceil(lock_timeout))

    end = time.time() + acquire_timeout

    while time.time() < end:
        if rdb.setnx(lockname, identifier):
            rdb.expire(lockname, lock_timeout)
        elif not rdb.ttl(lockname):
            rdb.expire(lockname, lock_timeout)
        
        time.sleep(0.001)
    
    return False

