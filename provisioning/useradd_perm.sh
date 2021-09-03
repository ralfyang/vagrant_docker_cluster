#!/bin/bash

## function for user add & sudo permission
User_add_comm(){
	user=$1
	adduser $user
	usermod -aG sudo $user
	groups $user
	echo "$user  ALL=(ALL) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/$user
}

user_list=$1
## user_list empty check
if [[ $user_list = "" ]];then
	echo "# "
	echo "# need a file of userlist"
	echo "# "
	exit 0
fi


## For loop for each user
users=`cat $user_list`
for i in $users
do
	echo "i = $i"
	User_add_comm $i
done
