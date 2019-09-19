#!/bin/bash
#ospfd.conf
if [[ ${_NIC_NAME_} == "" ]]; then
	echo "environment variable _NIC_NAME_ not specified, use default route to configure ospf..."
	_NIC_NAME_=`ip route | grep default | awk '{print $5}'`
fi
if [[ ${_GATEWAY_} == "" ]]; then
	echo "environment variable _GATEWAY_ not found, configure ospf by interface ${_NIC_NAME_}..."
	_GATEWAY_=`ip route | grep default | awk '{print $3}'`
	_GATEWAY_CIDR_=(`echo $_GATEWAY_ | awk -F "." '{print $1"."$2"."$3}'`)
	_ROUTER_ID_=`ip addr show $_NIC_NAME_ | grep $_GATEWAY_CIDR_ | awk '{print $2}' | awk -F "/" '{print $1}'`
	_NEIGHBOR_=`ip route | grep kernel  | grep $_GATEWAY_CIDR_ | awk '{print $1}'`
fi
#debian.conf
_OSPFD_IP_=$_ROUTER_ID_

sed -i "s@_OSPFD_IP_@$_OSPFD_IP_@g" /etc/quagga/debian.conf
sed -i "s@_NIC_NAME_@$_NIC_NAME_@g;s@_ROUTER_ID_@$_ROUTER_ID_@g;s@_GATEWAY_@$_GATEWAY_@g;s@_NEIGHBOR_@ $_NEIGHBOR_@g" /etc/quagga/ospfd.conf

service quagga start

tail -f /dev/null
