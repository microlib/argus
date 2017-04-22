#!/bin/sh

# This script is called by the systemd.service script
# This needs to be fixed
export GOPATH=~/home/$USER/Projects
PATH=$PATH:/usr/local/go/bin

if [ "$1" = "start" ]
then
  cd /opt/microlib/argus/
  ./argus &
fi

if [ "$1" = "stop" ]
then
  cd /opt/microlib/argus/
  PID=$(ps -ef | grep argus | grep -v grep | awk '{print $2}')
  kill -s SIGTERM $PID
fi

exit 0
