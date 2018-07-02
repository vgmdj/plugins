# plugins
common plugins
- Redis
- Message Queue

# redis
对外提供一个全局的redis，配置好连接信息后，只需调用相应的函数即可

```
  import "github.com/vgmdj/plugins/redis"
  
  redis.NewRedis("127.0.0.1:6379","key",0)
  
  redis.store("key","item","value")
  
  v,_ := redis.GetString("key")
  
  //v == "value"

```

# MQ
[to be continued...](https://github.com/vgmdj/plugins/tree/master/mq)
