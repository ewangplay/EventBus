#!/bin/bash
set -x

UP_DOWN="$1"
Eth0Addr=""

COMPOSE_FILE_NSQ="docker-compose-nsq.yaml"
COMPOSE_FILE_EVENTBUS="docker-compose-eventbus.yaml"

function printHelp () {
	echo "Usage: ./network_setup <up|down>"
}

function getETH0Addr() {
    Eth0Addr=`/sbin/ifconfig enp2s0|grep inet|grep -v inet6 | awk '{print $2}' | tr -d "addr:"`
    if [ "${Eth0Addr}" == "" ]; then
        Eth0Addr=`ifconfig enp2s0 | grep "inet addr:" | awk '{print $2}' | cut -c 6-`
    fi
    if [ "${Eth0Addr}" == "" ]; then
        Eth0Addr=`ifconfig enp2s0 | grep "inet" | grep -v "inet6" | awk '{print $2}'`
    fi 
    echo "enp2s0 net addr: ${Eth0Addr}"
}

function networkUp() {
    ## Set Env variables
    source ./env.sh

    ## get eth0 net addr
    getETH0Addr

    ## Init nsqlookupd0 configure
    mkdir -p /opt/nsqio/nsq/nsqlookupd0/etc
    cp ./config/nsqlookupd.cfg /opt/nsqio/nsq/nsqlookupd0/etc/nsqlookupd.cfg

    ## Init nsqd0 configure
    mkdir -p /opt/nsqio/nsq/nsqd0/etc
    sed "s/# broadcast_address: \"\"/broadcast_address: \"${Eth0Addr}\"/g" ./config/nsqd.cfg > ./config/nsqd0.cfg
    sed -i "s/nsqlookupd_tcp_addresses: \[\"127.0.0.1:4160\"\]/nsqlookupd_tcp_addresses: \[\"nsqlookupd0:4160\"\]/g" ./config/nsqd0.cfg
    cp ./config/nsqd0.cfg /opt/nsqio/nsq/nsqd0/etc/nsqd.cfg
    rm -f ./config/nsqd0.cfg

    ## Init nsqd1 configure
    mkdir -p /opt/nsqio/nsq/nsqd1/etc
    sed "s/http_address: \"0.0.0.0:4151\"/http_address: \"0.0.0.0:5151\"/g" ./config/nsqd.cfg > ./config/nsqd1.cfg
    sed -i "s/tcp_address: \"0.0.0.0:4150\"/tcp_address: \"0.0.0.0:5150\"/g" ./config/nsqd1.cfg
    sed -i "s/# broadcast_address: \"\"/broadcast_address: \"${Eth0Addr}\"/g" ./config/nsqd1.cfg
    sed -i "s/nsqlookupd_tcp_addresses: \[\"127.0.0.1:4160\"\]/nsqlookupd_tcp_addresses: \[\"nsqlookupd0:4160\"\]/g" ./config/nsqd1.cfg
    cp ./config/nsqd1.cfg /opt/nsqio/nsq/nsqd1/etc/nsqd.cfg
    rm -f ./config/nsqd1.cfg

    ## Init nsqd2 configure
    mkdir -p /opt/nsqio/nsq/nsqd2/etc
    sed "s/http_address: \"0.0.0.0:4151\"/http_address: \"0.0.0.0:6151\"/g" ./config/nsqd.cfg > ./config/nsqd2.cfg
    sed -i "s/tcp_address: \"0.0.0.0:4150\"/tcp_address: \"0.0.0.0:6150\"/g" ./config/nsqd2.cfg
    sed -i "s/# broadcast_address: \"\"/broadcast_address: \"${Eth0Addr}\"/g" ./config/nsqd2.cfg
    sed -i "s/nsqlookupd_tcp_addresses: \[\"127.0.0.1:4160\"\]/nsqlookupd_tcp_addresses: \[\"nsqlookupd0:4160\"\]/g" ./config/nsqd2.cfg
    cp ./config/nsqd2.cfg /opt/nsqio/nsq/nsqd2/etc/nsqd.cfg
    rm -f ./config/nsqd2.cfg

    ## Init nsqadmin configure
    mkdir -p /opt/nsqio/nsq/nsqadmin/etc
    sed "s/127.0.0.1:4161/nsqlookupd0:4161/g" ./config/nsqadmin.cfg > ./config/nsqadmin0.cfg
    cp ./config/nsqadmin0.cfg /opt/nsqio/nsq/nsqadmin/etc/nsqadmin.cfg
    rm -f ./config/nsqadmin0.cfg

    ## Init redis configure
    mkdir -p /opt/ewangplay/redis0/etc
    cp ./config/redis.conf /opt/ewangplay/redis0/etc/redis.conf

    ## Init eventbus configure
    mkdir -p /opt/ewangplay/eventbus0/etc
    sed "s/http_address: \"0.0.0.0:4151\"/http_address: \"0.0.0.0:5151\"/g" ./config/eventbus.yaml > ./config/eventbus0.yaml
    sed -i "s/tcp_address: \"0.0.0.0:4150\"/tcp_address: \"nsqd0:4150\"/g" ./config/eventbus0.yaml
    sed -i "s/lookupd_tcp_addresses: \[\"127.0.0.1:4160\"\]/lookupd_tcp_addresses: \[\"nsqlookupd0:4160\"\]/g" ./config/eventbus0.yaml
    cp ./config/eventbus0.yaml /opt/ewangplay/eventbus0/etc/eventbus.yaml
    rm -f ./config/eventbus0.yaml

    docker-compose -f $COMPOSE_FILE_NSQ -f $COMPOSE_FILE_EVENTBUS up -d
}

function networkDown() {
    source ./env.sh

    docker-compose -f $COMPOSE_FILE_NSQ -f $COMPOSE_FILE_EVENTBUS down

    rm -rf /opt/ewangplay/eventbus0
    rm -rf /opt/ewangplay/redis0
    rm -rf /opt/nsqio/nsq
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
