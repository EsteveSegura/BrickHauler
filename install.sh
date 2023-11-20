#!/bin/bash

# Check if we're root
if [ "$(id -u)" != "0" ]; then
   echo "Exute as root" 1>&2
   exit 1
fi

# Download
wget https://github.com/EsteveSegura/BrickHauler/releases/download/0.1.0/BrickHauler-linux-amd64 -O /var/tmp/BrickHauler

# Move binaries to /usr/local/bin
sudo cp /var/tmp/BrickHauler /usr/local/bin/BrickHauler
sudo chmod +x /usr/local/bin/BrickHauler

# Allow to be called in lowercase
sudo ln -s /usr/local/bin/BrickHauler /usr/local/bin/brickhauler

# Remove temporary files
sudo rm /var/tmp/BrickHauler

echo "BrickHauler installed successfully!"