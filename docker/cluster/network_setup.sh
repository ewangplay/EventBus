#!/bin/bash
set -x

UP_DOWN="$1"
Eth0Addr=""

function printHelp () {
	echo "Usage: ./network_setup <up|down>"
}

function getETH0Addr() {
    Eth0Addr=`/sbin/ifconfig eth0|grep inet|grep -v inet6 | awk '{print $2}' | tr -d "addr:"`
    if [ "${Eth0Addr}" == "" ]; then
        Eth0Addr=`ifconfig eth0 | grep "inet addr:" | awk '{print $2}' | cut -c 6-`
    fi
    if [ "${Eth0Addr}" == "" ]; then
        Eth0Addr=`ifconfig eth0 | grep "inet" | grep -v "inet6" | awk '{print $2}'`
    fi 
    echo "eth0 net addr: ${Eth0Addr}"
}

function networkUp() {
    ## Set Env variables
    source ./env.sh

    ## get eth0 net addr
    getETH0Addr

    ## Copy eventbus configure
    mkdir -p /opt/ewangplay/eventbus/seeker0/etc
    cp ./config/ebseeker.cfg /opt/ewangplay/eventbus/seeker0/etc/ebseeker.cfg

    mkdir -p /opt/ewangplay/eventbus/node0/etc
    sed "s/# broadcast_address: \"\"/broadcast_address: \"${Eth0Addr}\"/g" ./config/ebnode.yaml > ./config/ebnode0.yaml
    sed -i "s/lookupd_tcp_addresses: \[\"127.0.0.1:4160\"\]/lookupd_tcp_addresses: \[\"eventbus-seeker0:4160\"\]/g" ./config/ebnode0.yaml
    cp ./config/ebnode0.yaml /opt/ewangplay/eventbus/node0/etc/ebnode.yaml
    rm -f ./config/ebnode0.yaml

    mkdir -p /opt/ewangplay/eventbus/node1/etc
    sed "s/http_address: \"0.0.0.0:8091\"/http_address: \"0.0.0.0:9091\"/g" ./config/ebnode.yaml > ./config/ebnode1.yaml
    sed -i "s/tcp_address: \"0.0.0.0:4150\"/tcp_address: \"0.0.0.0:5150\"/g" ./config/ebnode1.yaml
    sed -i "s/http_address: \"0.0.0.0:4151\"/http_address: \"0.0.0.0:5151\"/g" ./config/ebnode1.yaml
    sed -i "s/# broadcast_address: \"\"/broadcast_address: \"${Eth0Addr}\"/g" ./config/ebnode1.yaml
    sed -i "s/lookupd_tcp_addresses: \[\"127.0.0.1:4160\"\]/lookupd_tcp_addresses: \[\"eventbus-seeker0:4160\"\]/g" ./config/ebnode1.yaml
    cp ./config/ebnode1.yaml /opt/ewangplay/eventbus/node1/etc/ebnode.yaml
    rm -f ./config/ebnode1.yaml

    mkdir -p /opt/ewangplay/eventbus/node2/etc
    sed "s/http_address: \"0.0.0.0:8091\"/http_address: \"0.0.0.0:10091\"/g" ./config/ebnode.yaml > ./config/ebnode2.yaml
    sed -i "s/tcp_address: \"0.0.0.0:4150\"/tcp_address: \"0.0.0.0:6150\"/g" ./config/ebnode2.yaml
    sed -i "s/http_address: \"0.0.0.0:4151\"/http_address: \"0.0.0.0:6151\"/g" ./config/ebnode2.yaml
    sed -i "s/# broadcast_address: \"\"/broadcast_address: \"${Eth0Addr}\"/g" ./config/ebnode2.yaml
    sed -i "s/lookupd_tcp_addresses: \[\"127.0.0.1:4160\"\]/lookupd_tcp_addresses: \[\"eventbus-seeker0:4160\"\]/g" ./config/ebnode2.yaml
    cp ./config/ebnode2.yaml /opt/ewangplay/eventbus/node2/etc/ebnode.yaml
    rm -f ./config/ebnode2.yaml

    mkdir -p /opt/ewangplay/eventbus/admin/etc
    sed "s/127.0.0.1:4161/eventbus-seeker0:4161/g" ./config/ebadmin.cfg > ./config/ebadmin0.cfg
    cp ./config/ebadmin0.cfg /opt/ewangplay/eventbus/admin/etc/ebadmin.cfg
    rm -f ./config/ebadmin0.cfg

    ## Copy redis configure
    mkdir -p /opt/ewangplay/redis0/etc
    cp ./config/redis.conf /opt/ewangplay/redis0/etc/redis.conf

    docker-compose up -d
}

function networkDown() {
    source ./env.sh

    docker-compose down

    rm -rf /opt/ewangplay/eventbus
    rm -rf /opt/ewangplay/redis0
}

#Create the network using docker compose
if [ "${UP_DOWN}" == "up" ]; then
	networkUp
elif [ "${UP_DOWN}" == "down" ]; then ## Clear the network
	networkDown
elif [ "${UP_DOWN}" == "restart" ]; then ## Restart the network
	networkDown
	networkUp
else
	printHelp
	exit 1
fi
