# vagrant_docker_cluster

## Requirement Install
### Virtual Box install
* https://www.virtualbox.org/wiki/Downloads

### Vagrant install
* https://www.vagrantup.com/downloads.html
* Box(Virtual Machine OS image) add
```
vagrant box add ubuntu1404 https://cloud-images.ubuntu.com/vagrant/trusty/current/trusty-server-cloudimg-amd64-vagrant-disk1.box
```


### Source Pull
```
git clone https://github.com/goody80/vagrant_docker_cluster.git
cd vagrant_docker_cluster
```

## Modify the `Vagrantfile`
* Need to change the box image as you have at below line in Vagrantfile


```
vi Vagrantfile

.
.
   node.vm.box = "ubuntu1404"
.
.
```

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
