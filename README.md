# sugar-agent

`sugar-agent`是一个抓取服务器性能数据的工具，它可以抓取服务器`某个时间段内`的`CPU`、`内存`、`磁盘`、`网络`等性能数据，并将数据通过`http API`发送回`sugar-server`。

## 性能数据抓取流程
1. `sugar-server`通过`rabbitmq`下发抓取服务器性能数据任务
2. `sugar-agent`作为`consumer`接收到任务后，会在本地创建一个抓取性能数据的任务，`consumer`设置了`prefetchCount=1`即`同一时刻`只能`运行一个`抓取性能数据的任务
3. 任务执行完毕后会通过`sugar-server`的`http API`回调，将性能数据通过`sugar-server`存入数据库中
4. 前端调用`sugar-server`的`http API`获取性能数据并展示

## Usage

```shell
hashqueue@hashqueue-pc:~/sugar-agent$ ./sugar-agent_amd64 -user guest -password guest -host localhost -port 5672 -exchange-name device_exchange -queue-name collect_device_perf_data_queue -routing-key device_perf_data
2023/03/05 13:18:22 Binding queue collect_device_perf_data_queue to exchange device_exchange with routing key device_perf_data
2023/03/05 13:18:22 [******] Started consumer [******] -> Waiting for messages. To exit press CTRL+C
2023/03/05 13:18:41 [x] Received a message [x] -> {"task_type": 0, "metadata": {"base_url": "http://127.0.0.1:8000", "task_id": "d5cfcfa5-266e-4faa-8200-3e5ad9fc8a4e", "username": "admin", "password": "admin3306", "task_config": {"intervals": 10, "count": 5}}}
2023/03/05 13:18:41 [x] Start task [x]
2023/03/05 13:19:37 [x] Task is done [x]
2023/03/05 13:19:37 [x] Total use time: 55.012806 s [x]
2023/03/05 13:31:07 [x] Received a message [x] -> {"task_type": 0, "metadata": {"base_url": "http://127.0.0.1:8000", "task_id": "5b03bc30-bf42-4008-97e0-29387dbbc24c", "username": "admin", "password": "admin3306", "task_config": {"intervals": 5, "count": 5}}}
2023/03/05 13:31:08 [x] Start task [x]
2023/03/05 13:31:38 [x] Task is done [x]
2023/03/05 13:31:38 [x] Total use time: 30.009797 s [x]
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
