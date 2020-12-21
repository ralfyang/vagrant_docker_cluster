OS_NAME="ubuntu/focal64"
NODE_COUNT = 10

Vagrant.configure("2") do |config|
  config.vm.synced_folder "~/.ssh/", "/tmp/conf.d/"
  config.vm.provision "shell", path: "./provisioning/docker.sh", args: ""
  NODE_COUNT.times do |i|
    node_id = "docker0#{i}.dev"
    config.vm.network "private_network", type: "dhcp", :adapter => 2
    config.vm.define node_id do |node|
      node.vm.box = OS_NAME
      node.vm.hostname = "#{node_id}"
    end
  end
  config.vm.provider :virtualbox do |vb|
      vb.customize ["modifyvm", :id, "--memory", "4096"]
      vb.customize ["modifyvm", :id, "--cpus", "2"]
  end
end
