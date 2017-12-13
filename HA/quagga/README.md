# k8s + docker + quagga
## 需求分析
ingress高可用，kube-router的简化，参考：https://cloudnativelabs.github.io/post/2017-11-01-kube-high-available-ingress/

## 文件解释
### Dockerfile
定制了一个quagga镜像，顺带安装了一些便于排查网络问题的软件
   
### change_quagga_config.sh
修改部分配置文件并重启quagga进程，然后用tail保持容器不退出(待优化)

### quagga.yaml
部署yaml

### zebra.conf等
quagga配置文件

## 踩到的坑
1. quagga进程在容器中重启的时候需要容器以特权模式启动，docker命令行加--privileged=true参数，k8s中见quagga.yaml文件。
2. ospfd会监听一个业务IP，有可能导致2604端口非预期暴露
3. 该示例中，是将lo口上配置的IP地址宣告给对端，假定lo口和互联地址在同一个段，会导致icmp redirect
4. 重分发路由，过滤端口为lo，并deny掉127.0.0.0/8
5. 防火墙需要放行一些端口:2604(quagga ospf) 59(ospf)，可以限定源IP地址
   
## 其他详细信息等踩完坑再来补充
