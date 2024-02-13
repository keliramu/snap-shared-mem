#!/bin/bash

sudo snap remove uxx2
sudo snap install --dangerous ./uxx2_0.53_amd64.snap
#sudo snap connect uxx2:shmem-private
#sudo snap connect uxx2:network-bind
snap connections uxx2
