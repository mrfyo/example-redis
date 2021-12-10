## 发布订阅

#### PUBLISH

将消息发送到指定通道，返回订阅数量。

时间复杂度为 O(N + M)，N是订阅者数量，而 `M` 则是使用模式订阅(subscribed patterns)的客户端的数量。

```sh
# publish channal message
> publish mq-fruit Apple
```

#### SUBSCRIBE

订阅指定一个或多个通道的消息。

```sh
# subscribe channal [channal...]
> subscribe mq-fruit mq-vegetable
{'pattern': None, 'type': 'subscribe', 'channel': b'my-second-channel', 'data': 1}
```

- **type**: 反馈类型，包括两种模式
  - 通道订阅：'subscribe',  'unsubscribe',   'message'
  - 模式订阅：'psubscribe',  'punsubscribe', 'pmessage'
- **channel**: 消息来源的具体通道
- **pattern**: 订阅的模式
- **data**: 
  - 如果类型是 [un]subscribe，代表订阅客户端数量
  - 如果类型是 message，代表收到的具体消息

#### UNSUBSCRIBE

退订指定通道，如果没有指定通道，则退订所有通道。

```sh
# unsubscribe [channal...]
```

#### PSUBSCRIBE

订阅一个或多个符合指定模式的通道。时间复杂度是O(N)，N是订阅的模式数量。

```sh
# psubscribe pattern [pattern...]
> psubscribe mq-*
```

#### PUNSUBSCRIBE

退订指定模式，如果没有指定模式，则退订所有模式。

```sh
# psubscribe pattern [pattern...]
> pusubscribe mq-*
```

