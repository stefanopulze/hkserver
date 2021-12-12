#!/usr/bin/env bash

RASP_IP=192.168.1.40
INSTALL_DIR=/opt/hkserver

echo "Stopping service"
ssh pi@$RASP_IP 'sudo systemctl stop hkserver'

echo "Updating binary"
scp -q ./output/hkserver pi@$RASP_IP:$INSTALL_DIR

echo "Restart service"
ssh pi@$RASP_IP 'sudo systemctl restart hkserver'

echo "🍺 App updated into $RASP_IP:$INSTALL_DIR"