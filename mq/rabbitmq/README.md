# RabbitMQ

## 创建连接
```
rabbit, _ := NewRabbit("127.0.0.1:5672", "/", "user", "pwd")

```

## 发送消息
```
//send to mq
rabbit.SendToQue("exchange", "key", []byte("OK"))

```

## 接收消息
```
//receive from mq
msgs, _ := rabbit.ReceiveFromMQ("exchange", "key", "queue", nil)
for msg := range msgs {
    log.Println("Receive a message from mq: ", string(msg.Body))

    //do some thing

}
```

## 限流
```
每次ack后，再接收新消息
rabbit.SetQos(1,0,true)

```
