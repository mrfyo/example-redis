import redis
import uuid
import time

from lock import acquire_lock, release_lock

def acquire_semaphore(rdb: redis.Redis, semname: str, limit: int, timeout=10):
    """
        计数信号量

        当多个进程运行在不同主机时，锁的稳定性受到主机时钟同步问题的影响
    """
    identifier = str(uuid.uuid4())
    now = time.time()

    pipe = rdb.pipeline(True)
    # 移除所有已经超时的进程
    pipe.zremrangebyscore(semname, '-inf', now - timeout)
    # 将当前进程加入集合中
    pipe.zadd(semname, identifier, now)
    # 返回当前线程的排名
    pipe.zrank(semname, identifier)
    # 如果排名小于限制量，则表示获取锁成功
    if pipe.execute()[-1] < limit:
        return identifier
    # 获取锁失败，则移除自身标识
    rdb.zrem(semname, identifier)
    return None


def release_semaphore(rdb: redis.Redis, semname: str, identifier: str):
    return rdb.zrem(semname, identifier)

##########################
# 公平信号量              
##########################


def acquire_fair_semaphore(rdb: redis.Redis, semname: str, limit: int, timeout=10):
    identifier = str(uuid.uuid4())
    czset = semname  + ":owner"
    ctr = semname + ':counter'

    now = time.time()
    p = rdb.pipeline(True)

    # 移除超时的信号量
    p.zremrangebyscore(semname, '-inf', now - timeout)
    # czset 更新为 两个有序集合的交集，分数仍然是 czset
    # 也就是说，时间戳对于排名没有决定关系，从而消除了时钟不一致带来的细微差距
    p.zinterstore(czset, {czset: 1, semname: 0})

    # 计数器自增
    p.incr(ctr)
    counter = p.execute()[-1]

    # 尝试获取信号量
    p.zadd(semname, identifier, now)
    p.zadd(czset, identifier, counter)

    # 检查排名，判断是否获得到信号量
    p.zrank(czset, identifier)
    if p.execute()[-1] < limit:
        return identifier
    
    # 未获得信号量，清除无用数据
    p.zrem(semname, identifier)
    p.zrem(czset, identifier)
    p.execute()

    return None


def release_fair_semaphore(rdb: redis.Redis, semname: str, identifier: str):
    p = rdb.pipeline(True)
    p.zrem(semname, identifier)
    p.zrem(semname + ':owner', identifier)
    return p.execute()[0]


def refresh_fair_semaphore(rdb: redis.Redis, semname: str, identifier: str):
    if rdb.zadd(semname, identifier, time.time()):
        # zadd的返回值表示成功被加入的新成员数量
        # 如果被加入而不是更新score的话，表示进程不在持有信号量的集合中
        # 也就是说该进程没有持有信号量，需要主动释放自己
        release_fair_semaphore(rdb, semname, identifier)
        return False

    return True


def acquire_semaphore_with_lock(rdb: redis.Redis, semname: str, limit: int, timeout=10):
    """
        用一个短暂的锁，来消除因时钟不同步导致的不公平竞争
    """
    identifier = acquire_lock(rdb, semname, acquire_timeout=0.01)
    if identifier:
        try:
            return acquire_fair_semaphore(rdb, semname, limit, timeout)
        finally:
            release_lock(rdb, semname, identifier)