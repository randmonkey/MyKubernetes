#!/bin/bash
#get somp snmp info
snmpwalk -v 2c -c Password 100.64.0.254 1.3.6.1.2.1.4.34.1.3.1.4 | cut -d "." -f 13-16 | cut -d " " -f 1,4
