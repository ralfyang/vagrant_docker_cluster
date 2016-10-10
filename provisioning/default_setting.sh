#!/bin/bash
# Script for Default env. of zinst
#!/bin/bash
LANG=en_US.UTF-8
#sed -i '/^LANG=/d' /etc/sysconfig/i18n
#echo 'LANG=en_US.UTF-8' >> /etc/sysconfig/i18n
#sed -i 's/=enforcing/=disabled/g' /etc/selinux/config

#setenforce 0
apt update
apt-get install curl -y
apt-get install wget -y

mkdir -p /root/.ssh
cp -Rfv /tmp/conf.d/* /root/.ssh
