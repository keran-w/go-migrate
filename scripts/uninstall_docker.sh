sudo apt-get purge docker-ce docker-ce-cli containerd.io docker-compose-plugin
dpkg -l | grep -i docker
sudo rm -rf /var/lib/docker
sudo rm -rf /var/lib/containerd
sudo rm -rf /etc/docker
sudo apt-get update
sudo reboot