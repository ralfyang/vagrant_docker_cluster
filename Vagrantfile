# -*- mode: ruby -*-
# vi: set ft=ruby :

require "fileutils"

HOST_SSH_DIR = File.expand_path("~/.ssh")
WINDOWS_STORAGE_DIR = File.expand_path("./windows_storege", __dir__)

AUTO_GENERATE_HOST_SSH_KEY = true
HOST_SSH_KEY = File.join(HOST_SSH_DIR, "id_ed25519")

unless Dir.exist?(HOST_SSH_DIR)
  puts "==> Creating host SSH directory: #{HOST_SSH_DIR}"
  FileUtils.mkdir_p(HOST_SSH_DIR)
  FileUtils.chmod(0700, HOST_SSH_DIR)
end

if AUTO_GENERATE_HOST_SSH_KEY && !File.exist?(HOST_SSH_KEY)
  puts "==> Host SSH key not found. Generating: #{HOST_SSH_KEY}"

  system(
    "ssh-keygen",
    "-t", "ed25519",
    "-N", "",
    "-f", HOST_SSH_KEY,
    "-C", "vagrant-host-key"
  )

  FileUtils.chmod(0600, HOST_SSH_KEY) if File.exist?(HOST_SSH_KEY)
  FileUtils.chmod(0644, "#{HOST_SSH_KEY}.pub") if File.exist?("#{HOST_SSH_KEY}.pub")
end

unless Dir.exist?(WINDOWS_STORAGE_DIR)
  puts "==> Creating Windows shared storage directory: #{WINDOWS_STORAGE_DIR}"
  FileUtils.mkdir_p(WINDOWS_STORAGE_DIR)
end

Vagrant.configure("2") do |config|
  config.vm.synced_folder HOST_SSH_DIR, "/tmp/conf.d"

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
    vb.cpus = "2"
    # vb.gui = true
  end

  # Docker VM cluster
  (1..9).each do |i|
    node_id = "docker0#{i}.dev"

    config.vm.define node_id do |node|
      node.vm.box = "bento/ubuntu-24.04"
      node.vm.hostname = node_id

      node.vm.network "private_network",
        ip: "192.168.62.10#{i}",
        netmask: "255.255.255.0"

      node.vm.provision "shell",
        path: "./provisioning/docker.sh",
        args: "vagrant"

      node.vm.provider "virtualbox" do |vb|
        vb.memory = "4096"
        vb.cpus = "2"
      end
    end
  end

  # Test VM cluster
  (1..9).each do |i|
    node_id = "test0#{i}.dev"

    config.vm.define node_id do |node|
      node.vm.box = "bento/ubuntu-24.04"
      node.vm.hostname = node_id

      node.vm.network "private_network",
        ip: "192.168.62.20#{i}",
        netmask: "255.255.255.0"

      node.vm.provider "virtualbox" do |vb|
        vb.memory = "1024"
        vb.cpus = "1"
      end
    end
  end

  # Windows VM cluster
  (1..9).each do |i|
    node_id = "win0#{i}.dev"

    config.vm.define node_id do |node|
      node.vm.box = "gusztavvargadr/windows-10"
      node.vm.hostname = node_id

      node.vm.synced_folder WINDOWS_STORAGE_DIR, "c:\\tmp"

      node.vm.provision "shell",
        path: "./provisioning/provision.ps1",
        args: ""

      node.vm.network "forwarded_port",
        guest: 3389,
        host: "338#{i}",
        host_ip: "0.0.0.0"

      node.vm.network "private_network",
        ip: "192.168.62.21#{i}",
        netmask: "255.255.255.0"

      node.vm.provider "virtualbox" do |vb|
        vb.memory = "8196"
        vb.cpus = "4"
        vb.customize [
          "setextradata",
          :id,
          "GUI/LastGuestSizeHint",
          "1280x720"
        ]
      end
    end
  end

   # Windows VM cluster
   (1..9).each do |i|
     node_id = "win0#{i}.dev"
   
     config.vm.define node_id do |node|
       node.vm.box = "StefanScherer/windows_10"
       node.vm.hostname = node_id
       node.vm.communicator = "winrm"
   
       node.vm.synced_folder WINDOWS_STORAGE_DIR, "c:\\tmp"
   
       node.vm.provision "shell",
         path: "./provisioning/provision.ps1",
         privileged: false
   
       node.vm.network "forwarded_port",
         guest: 3389,
         host: "338#{i}",
         host_ip: "0.0.0.0"
   
       node.vm.network "private_network",
         ip: "192.168.62.21#{i}",
         netmask: "255.255.255.0"
   
       node.vm.provider "virtualbox" do |vb|
         vb.memory = "8196"
         vb.cpus = "4"
         vb.customize [
           "setextradata",
           :id,
           "GUI/LastGuestSizeHint",
           "1280x720"
         ]
       end
     end
   end
   
end
