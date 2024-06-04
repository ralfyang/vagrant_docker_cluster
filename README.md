# vagrant_docker_cluster

## CLI mode version

### Prerequisites
- Go 1.16 or higher
- Vagrant
- VirtualBox
- Git

### Auto-Install & use(For the Macbook User or Linux)
* Clone the repository:

    ```sh
    git clone git@github.com:ralfyang/vagrant_docker_cluster.git
    cd vagrant_docker_cluster
    ```

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

## web-console mode
* setup the `.env` file as below
```
## Please make sure the Network information as below
Public_IP=111.111.111.2
Private_IP=192.168.100.13
PASSWORD=your_password_for_webconsole_login
```

* Running the webconsole for test
```
nohup go run ./main.go > output.log 2>&1 &
```

* Demo
![image](https://github.com/ralfyang/vagrant_docker_cluster/assets/4043594/eb3390cd-2551-4663-ad20-059d170e8786)



## Vagrant Web Console version

This project provides a web console for managing Vagrant VMs. The console allows you to start, stop, reload, reboot, and remove VMs, as well as view the current status of all VMs.


### Build for run
1. Build the `vvmanager` executable:

    ```sh
    go build -o vvmanager main.go
    ```

### Usage

The `webconsole.sh` script is used to manage the web console server. It supports `start`, `stop`, `restart`, and `run` commands.

- `start`: Starts the web console server in the background and rotates logs daily.
- `stop`: Stops the running web console server.
- `restart`: Restarts the web console server.
- `run`: Runs the web console server in the foreground.

To use the script, run the following commands:

1. Make the script executable:

    ```sh
    chmod +x webconsole.sh
    ```

2. Start the web console server:

    ```sh
    ./webconsole.sh start
    ```

3. Stop the web console server:

    ```sh
    ./webconsole.sh stop
    ```

4. Restart the web console server:

    ```sh
    ./webconsole.sh restart
    ```

5. Run the web console server in the foreground:

    ```sh
    ./webconsole.sh run
    ```

### Configuration

The `webconsole.sh` script handles log rotation and maintains a PID file to manage the server process. Logs are rotated daily, with a maximum of 4 log files retained (`history.log`, `history.log.1`, `history.log.2`, and `history.log.3`).

### Example

Here is an example session using the `webconsole.sh` script:

1. Start the server:

    ```sh
    ./webconsole.sh start
    ```

2. Access the web console by navigating to `http://<your_server_ip>:8080` in your web browser.

3. Stop the server:

    ```sh
    ./webconsole.sh stop
    ```

4. Restart the server:

    ```sh
    ./webconsole.sh restart
    ```

## Troubleshooting

- Ensure all dependencies are installed and properly configured.
- Check the log files (`history.log`) for any errors or warnings.
- If the server does not start, verify that no other process is using port 8080.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

