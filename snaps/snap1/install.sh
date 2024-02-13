#!/bin/bash

sudo snap remove uxx
sudo snap install --dangerous ./uxx_0.23_amd64.snap
#sudo snap connect uxx:shmem-private
#sudo snap connect uxx:network-bind
snap connections uxx
