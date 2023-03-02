## vagrant_docker_cluster

## Install the Requirements
### Source Clone first
```
git clone https://github.com/ralfyang/vagrant_docker_cluster.git
cd vagrant_docker_cluster
```

### Auto-Install & use(For the Macbook User or Linux)
[![asciicast](https://asciinema.org/a/78LKZjwwx0dMM595q0GqhDrvw.png)](https://asciinema.org/a/78LKZjwwx0dMM595q0GqhDrvw)

```
$ ./ctl.sh
===========================================================
 What do you want ?
===========================================================
[0] Install the Virtualbox & Vagrant <--- HERE!
[1] Start VM & login
[2] Login to VM
[3] Stop VM
[4] Reload VM
[5] Reboot VM
[RM] Remove VM
===========================================================
 Please insert a key as you need
===========================================================

```

## Install the Requirements - Manual
### Virtual Box install
* https://www.virtualbox.org/wiki/Downloads

### Vagrant install
* https://www.vagrantup.com/downloads.html


## Vagrant VM create
```
vagrant up docker01.dev
```
or

```
vagrant up docker0{1..9}.dev
```


### ssh-keygen - if you need to create ssh-key as below
```
ssh-keygen

Generating public/private rsa key pair.
Enter file in which to save the key (/Users/ralfyang/.ssh/id_rsa):
Created directory '/Users/ralfyang/.ssh'.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in /Users/ralfyang/.ssh/id_rsa.
Your public key has been saved in /Users/ralfyang/.ssh/id_rsa.pub.
The key fingerprint is:
```

### How to login
```
vagrant ssh docker01.dev
```
