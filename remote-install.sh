#!/bin/bash
git clone https://github.com/goody80/vagrant_docker_cluster.git
#curl -sL https://raw.githubusercontent.com/goody80/vagrant_docker_cluster/master/ctl.sh -o ./ctl.sh
chmod 755 ./ctl.sh
./ctl.sh
