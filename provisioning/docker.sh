#!/bin/bash
### Script for docker env. 
#sudo apt-get install libxslt-dev libxml2-dev libvirt-dev zlib1g-dev ruby-dev -y
#vagrant box add ubuntu1404 http://
#vagrant plugin install vagrant-libvirt 
#vagrant plugin install vagrant-mutate

wget -qO- https://get.docker.com/ | sh
sudo usermod -aG docker root
service docker start

# Docker-compose install
curl -L https://github.com/docker/compose/releases/download/1.8.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
chmod 755 /usr/local/bin/docker-compose



curl -sL bit.ly/ralf_dcs -o ./dcs
sudo chmod 755 ./dcs
sudo mv ./dcs /usr/bin/dcs
