### docker部署consul集群

------

#### 启动consul参数说明：

```shell
–-network=host
# docker参数, 使得docker容器越过了net namespace的隔离，免去手动指定端口映射的步骤
-server
# consul支持以server或client的模式运行, server是服务发现模块的核心, client主要用于转发请求
-advertise
# 将本机私有IP传递到consul
-retry-join
# 指定要加入的consul节点地址，失败后会重试, 可多次指定不同的地址
-client
# 指定consul绑定在哪个client地址上，这个地址可提供HTTP、DNS、RPC等服务，默认是127.0.0.1
-bind
# 绑定服务器的ip地址；该地址用来在集群内部的通讯，集群内的所有节点到地址必须是可达的，默认是0.0.0.0
allow_stale
# 设置为true则表明可从consul集群的任一server节点获取dns信息, false则表明每次请求都会>经过consul的server leader
-bootstrap-expect
# 数据中心中预期的服务器数。指定后，Consul将等待指定数量的服务器可用，然后>启动群集。允许自动选举leader，但不能与传统-bootstrap标志一起使用, 需要在server模式下运行。
-data-dir
# 数据存放的位置，用于持久化保存集群状态
-node
# 群集中此节点的名称，这在群集中必须是唯一的，默认情况下是节点的主机名。
-config-dir
# 指定配置文件，当这个目录下有 .json 结尾的文件就会被加载，详细可参考https://www.consul.io/docs/agent/options.html#configuration_files
-enable-script-checks
# 检查服务是否处于活动状态，类似开启心跳
-datacenter
# 数据中心名称
-ui
# 开启ui界面
-join
# 指定ip，加入到已有的集群中
```

#### consul端口用途说明

```shell
8500 # http端口，用于http接口和web ui访问
8300 # server rpc端口，同一数据中心consul server之间通过该端口通信
8301 # serf lan端口，同一数据中心consul client通过该端口通信; 用于处理当前datacenter中LAN的gossip通信
8302 # serf wan端口，不同数据中心consul server通过该端口通信; agent Server使用，处理与其他datacenter的gossip通信
8600 # dns端口，用于已注册的服务发现
```

#### 开始搭建

`consul`分为`server`和`client`两种运行模式，也就是服务端和客户端模式。

- `server`模式主要负责参与共识协议，维护集群状态的集中视图，并响应集群中其他代理的查询。
- `client`模式负责转发所有请求给Server，参与gossip协议发现其他代理并检查它们的健康状态，然后将有关群集的查询转发给服务器代理。

~~为了保证服务的高可用性和奇数个方便选举，官方建议每个数据中心的`server`数量为3个或5个，这里我们就创建3个`server`，1个`client`。~~

##### 启动`node_1`节点，以`server`模式运行

```shell
docker run -d --name consul_1 -p 8501:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600 consul agent -server -node=server_1 -bootstrap-expect 3 -ui -client=0.0.0.0 -bind=0.0.0.0
```

##### 查看第一个`server`节点的ip地址

```shell
docker inspect --format '{{ .NetworkSettings.IPAddress }}' consul_1
172.18.0.2
```

##### 启动`node_2`节点，以`server`模式运行

```shell
docker run -d --name consul_2 -p 8502:8500 consul agent -server -node=server_2 -bootstrap-expect 3 -ui -client=0.0.0.0 -bind=0.0.0.0 -join 172.18.0.2
```

##### 启动`node_3`节点，以`server`模式运行

```shell
docker run -d --name consul_3 -p 8503:8500 consul agent -server -node=server_3 -bootstrap-expect 3 -ui -client=0.0.0.0 -bind=0.0.0.0 -join 172.18.0.2
```

##### ~~启动`node_4节点，以`client`模式运行~~

```shell
docker run -d --name consul_4 -p 8500:8500 consul agent -node=client -ui -client=0.0.0.0 -bind=0.0.0.0 -join 172.18.0.2
```

##### 查看consul集群成员信息

```shell
docker exec -t consul_1 consul members
```

```
Node      Address          Status  Type    Build  Protocol  DC   Segment
server_1  172.18.0.2:8301  alive   server  1.8.3  2         dc1  <all>
server_2  172.18.0.3:8301  alive   server  1.8.3  2         dc1  <all>
server_3  172.18.0.4:8301  alive   server  1.8.3  2         dc1  <all>
```

加入client再看看：

```
Node      Address          Status  Type    Build  Protocol  DC   Segment
server_1  172.18.0.2:8301  alive   server  1.8.3  2         dc1  <all>
server_2  172.18.0.3:8301  alive   server  1.8.3  2         dc1  <all>
server_3  172.18.0.4:8301  alive   server  1.8.3  2         dc1  <all>
client    172.18.0.5:8301  alive   client  1.8.3  2         dc1  <default>
```

