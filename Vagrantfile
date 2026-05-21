# -*- mode: ruby -*-
# vi: set ft=ruby :

require "fileutils"
require "rbconfig"

HOST_ARCH = RbConfig::CONFIG["host_cpu"]
IS_ARM_HOST = HOST_ARCH =~ /arm|aarch64/i

HOST_SSH_DIR = File.expand_path("~/.ssh")
WINDOWS_STORAGE_DIR = File.expand_path("./windows_storege", __dir__)

AUTO_GENERATE_HOST_SSH_KEY = true
HOST_SSH_KEY = File.join(HOST_SSH_DIR, "id_ed25519")

ENABLE_WINDOWS_VMS = ENV.fetch("ENABLE_WINDOWS_VMS", "false") == "true"
REQUESTED_WINDOWS_VM = ARGV.any? { |arg| arg =~ /^win\d+\.dev$/ }

# Windows VM은 기본 비활성화
if REQUESTED_WINDOWS_VM && !ENABLE_WINDOWS_VMS
  abort <<~MSG
    ERROR: Windows VMs are disabled.

    Requested command:
      vagrant #{ARGV.join(" ")}

    To enable Windows VMs, run:

      ENABLE_WINDOWS_VMS=true vagrant up win01.dev --provider=virtualbox

    Linux Docker VMs can still be used normally:

      vagrant up docker01.dev
      vagrant up docker02.dev
      vagrant up mid01.dev
  MSG
end

# Ensure host ~/.ssh exists
unless Dir.exist?(HOST_SSH_DIR)
  puts "==> Creating host SSH directory: #{HOST_SSH_DIR}"
  FileUtils.mkdir_p(HOST_SSH_DIR)
  FileUtils.chmod(0700, HOST_SSH_DIR)
end

# Ensure host SSH key exists
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

# Ensure Windows shared storage exists
unless Dir.exist?(WINDOWS_STORAGE_DIR)
  puts "==> Creating Windows shared storage directory: #{WINDOWS_STORAGE_DIR}"
  FileUtils.mkdir_p(WINDOWS_STORAGE_DIR)
end

Vagrant.configure("2") do |config|
  # Host SSH directory is shared into Linux VMs.
  # docker.sh should copy only safe files from /tmp/conf.d if needed.
  config.vm.synced_folder HOST_SSH_DIR, "/tmp/conf.d"

  # Global VirtualBox defaults
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

  # Medium Linux VM cluster
  (1..9).each do |i|
    node_id = "mid0#{i}.dev"

    config.vm.define node_id do |node|
      node.vm.box = "bento/ubuntu-24.04"
      node.vm.hostname = node_id

      node.vm.network "private_network",
        ip: "192.168.62.20#{i}",
        netmask: "255.255.255.0"

      node.vm.provider "virtualbox" do |vb|
        vb.memory = "8196"
        vb.cpus = "4"
      end
    end
  end

  # Windows VM cluster
  #
  # Disabled by default.
  #
  # Apple Silicon + VirtualBox + Windows ARM:
  #
  #   ENABLE_WINDOWS_VMS=true vagrant up win01.dev --provider=virtualbox
  #
  # Intel/x86 + VirtualBox + Windows x86:
  #
  #   ENABLE_WINDOWS_VMS=true vagrant up win01.dev --provider=virtualbox
  #
  if ENABLE_WINDOWS_VMS
    (1..9).each do |i|
      node_id = "win0#{i}.dev"

      config.vm.define node_id do |node|
        node.vm.box = "stromweld/windows-11"
        node.vm.box_architecture = IS_ARM_HOST ? "arm64" : "amd64"

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
          vb.gui = true

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
end
