#!/bin/bash
### Code from https://github.com/ralfyang/vagrant_docker_cluster. Powered by Github.

hname=$1
Vagrant_version="2.3.4"

if [[ $hname = "" ]]; then
    hname="docker01.dev"
fi

sshkey_check() {
    if [[ ! -f $HOME/.ssh/id_rsa ]]; then
        ssh-keygen -f $HOME/.ssh/id_rsa -t rsa -N ''
    fi
    return 0
}

resource_chk="./resource.status"
if [[ ! -f $resource_chk ]]; then
    touch $resource_chk
fi

start() {
    sshkey_check
    vagrant up $hname
    result_vm
    vagrant ssh $hname
}

stop() {
    vagrant halt $hname
    result_vm
}

connection() {
    vagrant ssh $hname
}

reload() {
    vagrant reload $hname
    result_vm
}

result_vm() {
    vagrant status | sed '$d' | sed '$d' | sed '1,2d' | sed '$d' | sed '$d' > $resource_chk
}

status() {
    result_vm
    cat $resource_chk
}

reboot() {
    vagrant halt $hname
    vagrant up $hname
    result_vm
}

remove() {
    echo " Are you sure you want to remove the Virtual machine? [ y ]"
    read sure
    if [[ $sure = "y" ]]; then
        vagrant destroy -f $hname
    fi
    result_vm
}

application_install() {
    mkdir -p ~/tmp
    cd ~/tmp
    arch=$(uname -s)-$(uname -m)
    os=$(uname -s)

    case $os in
        Linux)
            ## VirtualBox install
            sudo apt-add-repository "deb http://download.virtualbox.org/virtualbox/debian $(lsb_release -sc) contrib"
            curl https://www.virtualbox.org/download/oracle_vbox_2016.asc -o /tmp/virtualbox.key
            sudo apt-key add /tmp/virtualbox.key
            sudo apt-get update
            sudo apt-get install linux-headers-$(uname -r)
            sudo apt-get install virtualbox -y
            sudo /sbin/vboxconfig

            ## Vagrant install
            wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
            echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
            sudo apt update && sudo apt install vagrant -y
            ;;

        Darwin)
            ## VirtualBox Download & Install
            VirtualBox_installer="https://download.virtualbox.org/virtualbox/7.0.6/VirtualBox-7.0.6-155176-OSX.dmg"
            VirtualBox_Ext_pkg="https://download.virtualbox.org/virtualbox/7.0.6/Oracle_VM_VirtualBox_Extension_Pack-7.0.6.vbox-extpack"
            VirtualBox_file=$(echo "$VirtualBox_installer" | awk -F'/' '{print $NF}')
            curl -L $VirtualBox_installer -o ./$VirtualBox_file
            curl -L $VirtualBox_Ext_pkg -o Oracle_VM_VirtualBox_Extension_Pack.vbox-extpack
            sudo hdiutil attach $VirtualBox_file
            sudo installer -pkg /Volumes/VirtualBox/VirtualBox.pkg -target /
            hdiutil unmount /Volumes/VirtualBox/
            sudo vboxmanage extpack install ./Oracle_VM_VirtualBox_Extension_Pack.vbox-extpack
            rm -f ./$VirtualBox_file

            ## Vagrant Download & Install
            Vagrant_installer="https://releases.hashicorp.com/vagrant/$Vagrant_version/vagrant_${Vagrant_version}_darwin_amd64.dmg"
            Vagrant_file="vagrant_install.dmg"
            curl -L $Vagrant_installer -o ./$Vagrant_file
            sudo hdiutil attach $Vagrant_file
            sudo installer -pkg /Volumes/Vagrant/vagrant.pkg -target /
            hdiutil unmount /Volumes/Vagrant
            rm -f ./$Vagrant_installer
            ;;
    esac
}

case $1 in
    start)
        start
        ;;
    stop)
        stop
        ;;
    connection)
        connection
        ;;
    reload)
        reload
        ;;
    reboot)
        reboot
        ;;
    remove)
        remove
        ;;
    application_install)
        application_install
        ;;
    status)
        status
        ;;
    *)
        echo "Invalid command"
        ;;
esac

