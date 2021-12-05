## Redis 练习场



## 1. 文章投票

该例子来源于《Redis 实战》第一章。希望通过 Redis 的集合，实现对文章进行快速排名，其中文章的评分与发布时间呈负相关，与点赞数呈正相关。也就说，某一文章的评分会随着时间而降低，而点赞数会增加它。

#### 1.1 数据结构设计

**计数器couter**，采用`hash`存储，key值为`couter`。用来存储文章和用户信息表的主键。

```json
{
    "user": 1,
    "article": 5
}
```

**文章信息表article**，采用 `hash`存储，key值为`article:id`。用户信息包括

```json
{
    "id": 1,
    "nickname": "Jack Me",
    "username": "Jack",
    "password": "b0804ec967f48520697662a204f5fe72"
}
```

**文章信息表article**采用 `hash`存储，key值为`article:id`。文章中包括的字段如下所示

```json
{
    "id": 1,
    "title": "Go to statement considered harmful",
    "link": "http://goo.gl/kZUSu",
    "poster": "user:83271",
    "time": 1638673699,
    "votes": 5
}
```

当然，为了避免使用`SCAN`命令来获取已经保存的用户和文章的 key ，我们还需要用集合保存表的索引。另外，考虑到我们可能需要按照时间段获取一部分文章。这里我们使用 `zset` 来作为索引表。虽然这里没有针对用户按时间查找的需求，但为了统一，都采用`zset`。下面是两个有序集合，分别作为**用户信息ID集合**，**文章ID集合**。

```json
{
    "record:user": [
        {
    		"score": 1638666257,
            "member": "user:4"
        }
    ],
    "record:user": [
        {
    		"score": 1638673699,
            "member": "article:3"
        }
    ]
}
```

通过上面几个存储对象，已经基本完成了用户和文章的记录。下面，需要通过使用其他的存储对象来关联用户和文章，来实现真正的功能。

**文章发布记录集合**

显而易见，未发布的文章是不应该被投票的。为了快速判断文章是否发布和过滤文章，我们这里还需要使用一个集合来记录已发布文章的ID。如果我们需要按照时间段来查找文章的话，我们可以继续使用`zset`，就像保存创建记录一样。**ZSCORE** 和 **SISMEMBER** 都是 O(1) 操作，所以set 和 zset 都可以满足需要。这里继续使用 zset，基数是发布时间。

```json
{
    "publish:article": [
       {
    		"score": 1638666300,
            "member": "user:4"
        }
    ]
}
```

**投票记录集合**

为了防止刷票，我们需要为每一篇文章都增加一个投票记录集合，来记录那些用户向其投上一个宝贵的票。key格式为`vote:article:articleID`

```json
{
    "vote:article:3": [
       "user:6",
       "user:5"
    ]
}
```

**分数记录有序集合**

完成需求的关键结构出现了，为了实现文章可以按照分数进行排名，我们需要使用一个有序集合来记录。其中基数是分数，成员是文章的key。

```
{
	"score:article": [
		{
    		"score": 1296,
            "member": "article:3"
        },
        {
    		"score": 3024,
            "member": "article:4"
        }
	]
}
```

#### 1.2 操作

**创建用户**

```sh
> HINCRBY counter user 1
> HSET user:4 id 4 nickname 'Jack Me' username Jack password 4ff9fc6e4e5d5f590c4f2134a8cc96d1
> ZADD record:user 1638666257 user:4
```

**创建文章**

```sh
> HINCRBY counter article 1
> HSET article:3 id 3 title 'Go to statement considered harmful' link 'http://goo.gl/kZUSu' poster 'user:83271' time 1638673699 votes 0
> ZADD record:article 1638673699 article:3
```

**发布文章**

```sh
> ZADD publish:article 1638673999 article:3
> ZADD publis
```

