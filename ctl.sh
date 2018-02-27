#!/bin/bash
hname=$1

  if [[ $hname = "" ]];then
    hname="docker01.dev"
  fi

sshkey_check(){
	if [[ ! -f $HOME/.ssh/id_rsa ]];then
		ssh-keygen -f $HOME/.ssh/id_rsa -t rsa -N ''
	fi
	return 0
}


start(){
	sshkey_check
	vagrant up $hname
	vagrant ssh $hname
}

stop(){
	vagrant halt $hname
}

connection(){
	vagrant ssh $hname
}

reload(){
	vagrant reload $hname
}

reboot(){
	vagrant halt $hname
	vagrant up $hname
}

remove(){
	echo " Are sure that remove the Virtual machine ? [ y ]"
	read sure
	if [[ $sure = "y" ]];then
		vagrant destroy -f $hname
	fi
}

application_install(){
	mkdir -p ~/tmp
	cd ~/tmp
	## VirtualBox Download & Install
	VirtualBox_installer="http://download.virtualbox.org/virtualbox/5.2.2/VirtualBox-5.2.2-119230-OSX.dmg"
	VirtualBox_file=$(echo "$VirtualBox_installer" | awk -F'/' '{print $NF}')
	wget $VirtualBox_installer
	sudo hdiutil attach $VirtualBox_file
	sudo installer -pkg /Volumes/VirtualBox/VirtualBox.pkg -target /
	hdiutil unmount /Volumes/VirtualBox/
	rm -f ./$VirtualBox_file

	## Vagrant Download & Install
	Vagrant_installer="https://releases.hashicorp.com/vagrant/2.0.1/vagrant_2.0.1_x86_64.dmg"
	Vagrant_file=$(echo "$Vagrant_installer" | awk -F'/' '{print $NF}')
	wget $Vagrant_installer
	sudo hdiutil attach $Vagrant_file
	sudo installer -pkg /Volumes/Vagrant/vagrant.pkg -target /
	hdiutil unmount /Volumes/Vagrant
	rm -f ./$Vagrant_installer
}

clear
BARR="==========================================================="
echo "$BARR"
echo " What do you want ?"
echo "$BARR"
echo "[0] Install the Virtualbox & Vagrant"
echo "[1] Start VM & login"
echo "[2] Login to VM"
echo "[3] Stop VM"
echo "[4] Reload VM"
echo "[5] Reboot VM"
echo "[RM] Remove VM"
echo "$BARR"
echo -n " Please insert a key as you need = "
read choice
echo "$BARR"

	case $choice in
		0)
			application_install;;
		1)
			start;;
		2)
			connection;;
		3)
			stop;;
		4)
			reload;;
		5)
			reboot;;
		RM|rm)
			remove;;
	esac

