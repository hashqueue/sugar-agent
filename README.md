# sugar-agent

`sugar-agent`是一个抓取服务器性能数据的工具，它可以抓取服务器`某个时间段内`的`CPU`、`内存`、`磁盘`、`网络`等性能数据，并将数据通过`http API`发送回`sugar-server`。

## 性能数据抓取流程
1. `sugar-server`通过`rabbitmq`下发抓取服务器性能数据任务
2. `sugar-agent`作为`consumer`接收到任务后，会在本地创建一个抓取性能数据的任务，`consumer`设置了`prefetchCount=1`即`同一时刻`只能`运行一个`抓取性能数据的任务
3. 任务执行完毕后会通过`sugar-server`的`http API`回调，将性能数据通过`sugar-server`存入数据库中
4. 前端调用`sugar-server`的`http API`获取性能数据并展示

## Usage

```shell
hashqueue@hashqueue-pc:~/sugar-agent$ ./sugar-agent_amd64 -user guest -password guest -host 192.168.124.12 -port 5672 -exchange-name task_exchange -device-id 26
2023/03/11 15:02:36 Binding queue collect_device_26_perf_data_queue to exchange task_exchange
2023/03/11 15:02:36 [******] Started consumer [******] -> Waiting for messages. To exit press CTRL+C
2023/03/11 15:04:58 [x] Received a message [x] -> {"task_type": 0, "metadata": {"base_url": "http://192.168.124.12:8000", "task_uuid": "b107992c-f519-477e-ad91-e36956413f9a", "username": "consumer", "password": "88888888", "device_id": "26", "task_config": {"intervals": 10, "count": 10}}}
2023/03/11 15:04:58 [x] Start task [x]
2023/03/11 15:06:48 [x] Task is done [x]
2023/03/11 15:06:48 [x] Total use time: 110.026675 s [x]
2023/03/11 15:06:48 [x] Received a message [x] -> {"task_type": 0, "metadata": {"base_url": "http://192.168.124.12:8000", "task_uuid": "750aba77-dcec-45d2-a4a4-b645d33a1391", "username": "consumer", "password": "88888888", "device_id": "24", "task_config": {"intervals": 9, "count": 9}}}
2023/03/11 15:06:48 [x] Device id not match [x] -> deviceId from message: 24, my deviceId: 26
2023/03/11 15:06:48 Nothing to do, ack message and continue
2023/03/11 15:06:48 [x] Received a message [x] -> {"task_type": 0, "metadata": {"base_url": "http://192.168.124.12:8000", "task_uuid": "21bd57f2-d418-4206-8bc6-4f847dc8eee5", "username": "consumer", "password": "88888888", "device_id": "26", "task_config": {"intervals": 8, "count": 8}}}
2023/03/11 15:06:48 [x] Start task [x]
2023/03/11 15:08:00 [x] Task is done [x]
2023/03/11 15:08:00 [x] Total use time: 72.020304 s [x]
2023/03/11 15:08:00 [x] Received a message [x] -> {"task_type": 0, "metadata": {"base_url": "http://192.168.124.12:8000", "task_uuid": "aaef95dc-4583-4b2a-a7e4-7ac32ce434da", "username": "consumer", "password": "88888888", "device_id": "24", "task_config": {"intervals": 7, "count": 7}}}
2023/03/11 15:08:00 [x] Device id not match [x] -> deviceId from message: 24, my deviceId: 26
2023/03/11 15:08:00 Nothing to do, ack message and continue
```

## How to build this project
```shell
# Please install golang first, see https://go.dev/doc/install
# Use Git to clone this repo
git clone https://github.com/hashqueue/sugar-agent.git
# Install the project dependency package
cd sugar-agent/
go mod tidy
# Compile binary executable files
./build.sh
# Done.
```
