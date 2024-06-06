#!/bin/bash
# Install minimal prerequisites (Ubuntu 18.04 as reference)
sudo apt update && sudo apt-get install build-essential cmake libgtk-3-dev libboost-all-dev libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg-turbo8-dev -y
# Download and unpack sources
make install