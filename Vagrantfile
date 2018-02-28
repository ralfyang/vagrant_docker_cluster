# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.a

######## VirtualBox Box image name
#OS_NAME="folimy/Ubuntu1604_with_docker"
#OS_URL="https://cloud-images.ubuntu.com/xenial/current/xenial-server-cloudimg-amd64-vagrant.box"
OS_NAME="ubuntu1604"
OS_URL="https://vagrantcloud.com/folimy/boxes/Ubuntu1604_with_docker/versions/0.1/providers/virtualbox.box"

VAGRANTFILE_API_VERSION = "2"
NODE_COUNT = 10
Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.provision "shell", path: "./provisioning/default_setting.sh", args: ""
  config.vm.synced_folder "~/.ssh/", "/tmp/conf.d/"
## For Docker instance  
  NODE_COUNT.times do |i|
    node_id = "docker0#{i}.dev"
    config.vm.define node_id do |node|
                                                                  ### Args= [service name] [master / slave] [master IP] [advertise IP]
    config.vm.provision "shell", path: "./provisioning/docker.sh", args: ""
#    config.vm.network "forwarded_port", guest: 5601, host: 5601
#    config.vm.network "forwarded_port", guest: 5000, host: 5000
#    config.vm.network "forwarded_port", guest: 9200, host: 9200
#    config.vm.network "private_network", ip: "192.168.10.1#{i}"
#    config.vm.network "public_network", :dev => "br0", :mode => "bridge", :type => "bridge", :ip => "192.168.133.10#{i}", :netmask => "255.255.255.0", :auto_config => "false"
      node.vm.box = OS_NAME
      node.vm.box_url = OS_URL
      node.vm.hostname = "#{node_id}"
    end
  end
    
  config.vm.provider :virtualbox do |vb|
      vb.customize ["modifyvm", :id, "--memory", "4096"]
      vb.customize ["modifyvm", :id, "--cpus", "2"]   
  end 
end
