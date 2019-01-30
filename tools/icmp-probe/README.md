## ping-probe

ping-probe读入环境变量IPLIST， 对其中的ip进行ping，如果丢包率不为0，则将时间戳、目的IP、发包数、收包数打到标准输出日志中。
ping的策略为：每秒1次，每次10个包，2秒没有返回即超时。
环境变量IPLIST为string类型，每个IP需要用空格分开。


