#etcd
##操作步骤
1. 下载最近版本etcd到/opt/etcd/下并解压
2. 写入上述配置文件的到/etc/systemd/system/etcd.service，注意其中--name=net1和IP地址需要修改
3. 执行systemctl启动命令，启动etcd服务
```
systemctl daemon-reload
systemctl reload etcd.service
systemctl start etcd.service
```
4. 可以通过systemctl status etcd.service查看etcd的启动状态
5. 拷贝etcdctl到/bin/目录中方便管理etcd集群，记得添加环境变量：ETCDCTL_API=3
