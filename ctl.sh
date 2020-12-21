#!/bin/bash
### Code from https://github.com/ralfyang/vagrant_docker_cluster. Powered by Github.

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
  	arch=`uname -s`-`uname -m`
	os=`uname -s`

  	case $os in
    	Linux)
      		## VirtualBox install
      		sudo apt-add-repository "deb http://download.virtualbox.org/virtualbox/debian $(lsb_release -sc) contrib"
      		#wget -q https://www.virtualbox.org/download/oracle_vbox.asc -O- | sudo apt-key add -
			#wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
			curl https://www.virtualbox.org/download/oracle_vbox_2016.asc -o /tmp/virtualbox.key
			sudo apt-key add /tmp/virtualbox.key
      		sudo apt-get update
      		sudo apt-get install virtualbox-5.2  -y
   		;;

    	Darwin)
            	## VirtualBox Download & Install
            	#VirtualBox_installer="http://download.virtualbox.org/virtualbox/5.2.42/VirtualBox-5.2.42-137960-OSX.dmg"
            	VirtualBox_installer="https://download.virtualbox.org/virtualbox/6.1.16/VirtualBox-6.1.16-140961-OSX.dmg"
            	VirtualBox_file=$(echo "$VirtualBox_installer" | awk -F'/' '{print $NF}')
            	curl -L  $VirtualBox_installer -o ./$VirtualBox_file
            	sudo hdiutil attach $VirtualBox_file
            	sudo installer -pkg /Volumes/VirtualBox/VirtualBox.pkg -target /
            	hdiutil unmount /Volumes/VirtualBox/
            	rm -f ./$VirtualBox_file

		## Vagrant Download & Install
		#Vagrant_installer="https://releases.hashicorp.com/vagrant/2.0.1/vagrant_2.0.1_x86_64.dmg"
		Vagrant_installer="https://releases.hashicorp.com/vagrant/2.2.14/vagrant_2.2.14_x86_64.dmg"
		Vagrant_file=$(echo "$Vagrant_installer" | awk -F'/' '{print $NF}')
		curl -L $Vagrant_installer -o ./$Vagrant_file
		sudo hdiutil attach $Vagrant_file
		sudo installer -pkg /Volumes/Vagrant/vagrant.pkg -target /
		hdiutil unmount /Volumes/Vagrant
		rm -f ./$Vagrant_installer
	    	;;
    	esac


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

