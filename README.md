# wfrpc

一款便捷的frp客户端（终端）开发运维远程维护调试利器，自动生成配置且支持后台运行。

wfrpc图形化客户端[wfrpc_gui](https://github.com/GavinTan/wfrpc_gui)。



![image-20220826093134131](https://raw.githubusercontent.com/GavinTan/files/master/picgo/image-20220826093134131.png)

## build

```shell
git clone https://github.com/GavinTan/wfrpc.git

cd wfrpc
make
```

## 使用

```shell
wfrpc -l 22 -r 2222

wfrpc -l 22 -r 2222 -l 80 -d
```



| 参数           | 描述                                                        | 默认值        |
| -------------- | ----------------------------------------------------------- | ------------- |
| --server_addr  | frp服务器地址                                               | nb33.3322.org |
| --server_port  | frp服务器端口                                               | 7000          |
| --server_token | frp服务器认证token                                          |               |
| --host         | 转发目标IP地址                                              | 本机ip        |
| --lport        | 转发端口, 支持指定多个                                      |               |
| --rport        | 转发后访问的端口, 不指定将随机生成, 支持指定多个与lport对应 |               |
| --type         | 指定端口类型 tcp or udp                                     | tcp           |
| --daemon       | 后台运行                                                    |               |
| --clear        | 清理配置                                                    |               |
| --showconfig   | 查看配置                                                    |               |
| --showdir      | 查看配置目录                                                |               |
| --noconfig     | 不自动生成配置文件                                          |               |
| status         | 查看后台运行状态                                            |               |
| stop           | 停止后台运行                                                |               |
