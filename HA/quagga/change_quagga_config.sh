#!/bin/bash
#ospfd.conf
_NIC_NAME_=`ip route | grep default | awk '{print $5}'`
_GATEWAY_=`ip route | grep default | awk '{print $3}'`
_ROUTER_ID_=`ip addr show $_NIC_NAME_ | grep -w inet | awk '{print $2}' | awk -F "/" '{print $1}'`
_NETWORK_CIDR_=(`ip addr show lo | grep -w inet | awk '{print $2}' | awk -F "/" '{print $1}'`)
_NEIGHBOR_=`ip route | grep $_NIC_NAME_ | grep kernel | awk '{print $1}'`
#debian.conf
_OSPFD_IP_=$_ROUTER_ID_

sed -i "s@_OSPFD_IP_@$_OSPFD_IP_@g" /etc/quagga/debian.conf
sed -i "s@_NIC_NAME_@$_NIC_NAME_@g;s@_ROUTER_ID_@$_ROUTER_ID_@g;s@_GATEWAY_@$_GATEWAY_@g;s@_NEIGHBOR_@ $_NEIGHBOR_@g" /etc/quagga/ospfd.conf

if ! grep -q switch /etc/iproute2/rt_tables
then
echo "252	switch" >> /etc/iproute2/rt_tables
fi

if ! ip rule | grep -q switch
then
	ip rule add  from all to $_GATEWAY_ lookup switch
fi

if ! ip route show table switch | grep -q default
then
	ip route add 0.0.0.0/0 via $_GATEWAY_ dev $_NIC_NAME_ table switch
fi

if ! grep -q switch /etc/network/interfaces
then
	echo "post-up ip rule add  from all to $_GATEWAY_ lookup switch" >> /etc/network/interfaces
	echo "pre-down ip rule del  from all to $_GATEWAY_ lookup switch" >> /etc/network/interfaces
	echo "post-up ip route add 0.0.0.0/0 via $_GATEWAY_ dev $_NIC_NAME_ table switch" >> /etc/network/interfaces
	echo "pre-down ip route del 0.0.0.0/0 via $_GATEWAY_ dev $_NIC_NAME_ table switch" >> /etc/network/interfaces
fi

service quagga start
tail -f /dev/null
