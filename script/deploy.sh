#!/bin/bash
DIR=$(pwd)
SERVICE=/etc/systemd/system/shorters.service

go build

m4 -D_exe_=$DIR/shorters < "$DIR/script/shorters.service.m4" > $SERVICE

systemctl enable shorters.service
systemctl start shorters
