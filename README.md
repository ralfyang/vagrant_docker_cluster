# vagrant_docker_cluster
## Vagrant install
 * Download: https://www.vagrantup.com/downloads.html
 
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
