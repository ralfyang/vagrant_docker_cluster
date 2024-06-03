# -*- mode: ruby -*-
# vi: set ft=ruby :


Vagrant.configure("2") do |config|
  config.vm.synced_folder "~/.ssh/", "/tmp/conf.d/"
  config.vm.provision "shell", path: "./provisioning/docker.sh", args: ""

      ## Docker VM cluster
      (1..9).each do |i|
      node_id = "docker0#{i}.dev"
         config.vm.define node_id do |node|
            node.vm.box = "ubuntu/focal64"
            node.vm.hostname = "#{node_id}"
#           node.vm.network "forwarded_port", guest: 8080, host: 808#{i}, host_ip: "127.0.0.1"
            node.vm.network "private_network", ip: "192.168.62.10#{i}", netmask: "255.255.255.0"
#            node.vm.synced_folder "./data", "/vagrant_data"
            node.vm.provider "virtualbox" do |vb|
              vb.memory = "4096"
              vb.cpus = "2"
             end
         end
      end

      ## Test VM cluster
      (1..9).each do |i|
      node_id = "test0#{i}.dev"
         config.vm.define node_id do |node|
            node.vm.box = "ubuntu/focal64"
            node.vm.hostname = "#{node_id}"
#           node.vm.network "forwarded_port", guest: 8080, host: 808#{i}, host_ip: "127.0.0.1"
            node.vm.network "private_network", ip: "192.168.62.20#{i}", netmask: "255.255.255.0"
#            node.vm.synced_folder "./data", "/vagrant_data"
            node.vm.provider "virtualbox" do |vb|
              vb.memory = "1024"
              vb.cpus = "1"
             end
         end
      end

## Win VM cluster
  (1..9).each do |i|
    node_id = "win0#{i}.dev"
    config.vm.define node_id do |node|
      config.vm.provision "shell", path: "./provisioning/provision.ps1", args: ""
      config.vm.synced_folder "./windows_storege", "c:\\tmp"
      node.vm.box = "gusztavvargadr/windows-10"
      node.vm.hostname = "#{node_id}"
      #node.vm.network "forwarded_port", guest: "3389", host: "338#{i}", host_ip: "127.0.0.1"
      node.vm.network "forwarded_port", guest: "3389", host: "338#{i}", host_ip: "0.0.0.0"
      node.vm.network "private_network", ip: "192.168.62.21#{i}", netmask: "255.255.255.0"
      #node.vm.synced_folder "./data", "/vagrant_data"
      node.vm.provider "virtualbox" do |vb|
        vb.memory = "8196"
        vb.cpus = "4"
        vb.customize ['setextradata', :id, 'GUI/LastGuestSizeHint', '1280x720']  # 해상도 설정 추가
      end
    end
  end


   config.vm.provider "virtualbox" do |vb|
     vb.memory = "2048"
     vb.cpus = "2"
     ## Display the VirtualBox GUI when booting the machine
     #vb.gui = true
   end
end
