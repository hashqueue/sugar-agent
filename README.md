# sugar-agent

`sugar-agent`是一个抓取服务器性能数据的工具，它可以抓取服务器`某个时间段内`的`CPU`、`内存`、`磁盘`、`网络`等性能数据，并将数据通过`http API`发送回`sugar-server`。

## 性能数据抓取流程
1. `sugar-server`通过`rabbitmq`下发抓取服务器性能数据任务
2. `sugar-agent`作为`consumer`接收到任务后，会在本地创建一个抓取性能数据的任务，`consumer`设置了`prefetchCount=1`即`同一时刻`只能`运行一个`抓取性能数据的任务
3. 任务执行完毕后会通过`sugar-server`的`http API`回调，将性能数据通过`sugar-server`存入数据库中
4. 前端调用`sugar-server`的`http API`获取性能数据并展示

message body
```text
{"taskType":0,"metadata":{"durationTime":5,"count":5}}
# output usage
2023/03/04 13:22:18 Binding queue log_queue to exchange logs_direct with routing key info
2023/03/04 13:22:18 [******] Started consumer [******] -> Waiting for messages. To exit press CTRL+C
2023/03/04 13:22:39 [x] Received a message [x] -> {"taskType":0,"metadata":{"durationTime":1,"count":5}}
2023/03/04 13:22:49 [x] Task is done [x]
2023/03/04 13:22:49 [x] Total use time: 10.005037 s [x]
```
