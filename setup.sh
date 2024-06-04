#!/bin/bash

# Define variables
SERVICE_NAME="vvmanager"
SERVICE_USER=$(whoami)
SERVICE_GROUP=$(id -gn)
EXEC_PATH="/usr/local/bin/$SERVICE_NAME"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"
SCRIPT_PATH="/usr/local/bin/$SERVICE_NAME"

# Compile the Go application
echo "Compiling the Go application..."
go build -o $EXEC_PATH main.go

# Check if the build was successful
if [ ! -f "$EXEC_PATH" ]; then
    echo "Build failed. Please ensure you have Go installed and your main.go file is present."
    exit 1
fi

# Create systemd service file
echo "Creating systemd service file..."
cat <<EOL | sudo tee $SERVICE_FILE
[Unit]
Description=Vagrant VM Manager
After=network.target

[Service]
ExecStart=$EXEC_PATH
Restart=on-failure
User=$SERVICE_USER
Group=$SERVICE_GROUP

[Install]
WantedBy=multi-user.target
EOL

# Reload systemd to recognize the new service
echo "Reloading systemd..."
sudo systemctl daemon-reload

# Enable the service to start on boot
echo "Enabling the service..."
sudo systemctl enable $SERVICE_NAME

# Create the service management script
echo "Creating the service management script..."
cat <<EOL | sudo tee $SCRIPT_PATH
#!/bin/bash

case "\$1" in
    start)
        sudo systemctl start $SERVICE_NAME
        ;;
    stop)
        sudo systemctl stop $SERVICE_NAME
        ;;
    restart)
        sudo systemctl restart $SERVICE_NAME
        ;;
    status)
        sudo systemctl status $SERVICE_NAME
        ;;
    *)
        echo "Usage: $SERVICE_NAME {start|stop|restart|status}"
        exit 1
esac
exit 0
EOL

# Make the service management script executable
echo "Making the service management script executable..."
sudo chmod +x $SCRIPT_PATH

echo "Setup completed successfully!"

