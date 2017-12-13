# 个人Kubernetes实践

## 用virtualbox搭建本地k8s集群

    - - - - - -                   - - - - - - - -
    |    PC   | - - - - - - - - - |     k8s     |
    - - - - - -                   | deploy/ssh  |
                                  - - - - - - - -
                                        |
                        - - - - - - - - - - - - - - - - -
                        |                               |
                        |     switch(virtual box)       |
                        |                               |
                        - - - - - - - - - - - - - - - - -
                           |   |           |    |   |
           - - - - - - - - -   |           |    |   - - - - - - - - - - - - - - - - -
           |                   |           |    |                                   |
           |                   |           |    - - - - - - - - - - -|              |
    - - - - - - - - -  - - - - - - - - -  - - - - - - - - -   - - - - - - -  - -  -  - - - 
    |      net1     |  |     net2      |  |     net3      |   |    node1  |  |   node2   |
    | forward master|  | forward master|  | forward master|   |   compute |  |   compute |
    - - - - - - - - -  - - - - - - - - -  - - - - - - - - -   - - - - - - -  - - - - - - -   

## 需求：
1. k8s 可以免密码登录net1 net2 net3 node node2, 并搭建私有registry。
2. OS: Ubuntu 16.04 Kernel: 4.4.0-72-generic Docker: 17.09.0-ce Ansible: 2.4.1.0
3. 所有节点可以访问公网

## 准备工作
### k8s机器
1. 共享文件夹
2. Go 环境:go version go1.9.2 linux/amd64
3. 私有仓库:127.0.0.1:5000
4. 构建通用网络镜像(基于Ubuntu，安装了一些软件,见Dockerfiel)
