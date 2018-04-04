#!/bin/bash
_NIC_NAME_=`ip route | grep default | awk '{print $5}'`
_GATEWAY_=`ip route | grep default | awk '{print $3}'`
_GW_NO_=`cat -n /etc/network/interfaces | grep gateway | grep -v '#' | awk '{print $1}'`
_RULE_NO_=`cat /etc/iproute2/rt_tables | grep -vE "^0|^#" | awk '{print $1}' | sort | head -1`

if ! grep -q switch /etc/iproute2/rt_tables
then
echo "$[ _RULE_NO_ -1 ]       switch" >> /etc/iproute2/rt_tables
fi

if ! grep -q switch /etc/network/interfaces
then
        sed -i "$_GW_NO_ a pre-down ip route del 0.0.0.0/0 via $_GATEWAY_ dev $_NIC_NAME_ table switch" /etc/network/interfaces
        sed -i "$_GW_NO_ a post-up ip route add 0.0.0.0/0 via $_GATEWAY_ dev $_NIC_NAME_ table switch" /etc/network/interfaces
        sed -i "$_GW_NO_ a pre-down ip rule del  from all to $_GATEWAY_ lookup switch" /etc/network/interfaces
        sed -i "$_GW_NO_ a post-up ip rule add  from all to $_GATEWAY_ lookup switch" /etc/network/interfaces
fi
