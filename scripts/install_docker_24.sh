#!/bin/bash

# Update the package index
sudo apt-get update

# Install required packages to enable the apt repository over HTTPS
sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release -y

# Add Docker’s official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up the stable repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Update the apt package index again
sudo apt-get update

# Find the specific version string for Docker version 24.0.2 (replace 'focal' with your version of Ubuntu if different)
VERSION_STRING=$(apt-cache madison docker-ce | grep '24.0.2' | head -1 | awk '{print $3}')

# Check if VERSION_STRING is not empty
if [ -z "$VERSION_STRING" ]; then
    echo "Docker version 24.0.2 not found"
    exit 1
fi

# Install Docker Engine
sudo apt-get install docker-ce="$VERSION_STRING" docker-ce-cli="$VERSION_STRING" containerd.io docker-compose-plugin -y

# Verify installation
sudo docker --version
