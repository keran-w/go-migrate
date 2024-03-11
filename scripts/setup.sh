sudo apt-get install python3-pip libprotobuf-dev libprotobuf-c-dev protobuf-c-compiler protobuf-compiler python3-protobuf libcap-dev libnl-3-dev xmlto asciidoc protobuf-c-compiler asciidoc libbsd-dev libdrm-dev libdrm-amdgpu1 libgnutls28-dev libnftables-dev libnet-dev libnl-3-dev libnl-genl-3-dev criu
cd ~
git clone https://github.com/checkpoint-restore/criu.git
cd criu
sudo make install
criu --version
