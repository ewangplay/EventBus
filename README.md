# EventBus Service

- Support `Notifier` and `Queuer` type drivers for event.

- Provide Restful API to let user to send event -- either `Notifier` or `Queuer` type event.

## Setup

### Docker

- Build docker images

    ```
    make docker
    ```

- Clean docker images

    ```
    make docker-clean
    ```

- Run docker container manually

    - Run eventbus docker 

        ```
        docker run -d --name eventbus -v /opt/ewangplay/eventbus/log:/opt/ewangplay/eventbus/log -v /opt/ewangplay/eventbus/data:/opt/ewangplay/eventbus/data -v /opt/ewangplay/eventbus/etc:/opt/ewangplay/eventbus/etc -p 0.0.0.0:8091:8091 -p 0.0.0.0:4150:4150 -p 0.0.0.0:4151:4151 <eventbus-image-id>
        ```

    __Note:__ Volume option (/opt/ewangplay/eventbus/etc), the directory should contains the eventbus services configure files: [eventbus.yaml](./docker/config/eventbus.yaml.example).

### Manual Setup

- Build services

    Install following system packages before building the services:

    - Redhat/Centos/Fadora: snappy-devel, libzip-devel, zlib-devel, zlib, gcc-c++
    - Debian: libbz2-dev libsnappy-dev

        ```
        $ git clone ssh://git@192.168.250.3:10022/ewangplay/eventbus.git
        $ cd $GOPATH/github.com/ewangplay/eventbus
        $ make clean srvc
        ```

- Run unit test

    ```
    $ cd $GOPATH/github.com/ewangplay/eventbus
    $ make checks
    ```

- Run benchmark test

    ```
    $ cd $GOPATH/github.com/ewangplay/eventbus
    $ make bench
    ```

- Start eventbus services

	```
	$ ./build/bin/eventbus --config=/path/to/eventbus.yaml
	```

	__Note:__ Services can be obtained sample configure files from [example config](./docker/config) dir. Copy and mofify these files to fit your service instance.

## Sample Cluster

After build the eventbus docker images, you can use the sample cluster solution in [cluster](../docker/cluster) dir to build an all-in-one test cluster.

- Run the test cluster:

    ```
    cd ./docker/cluster 
    ./network_setup.sh up 
    ```

- Clean the test cluster:

    ```
    cd ./docker/cluster
    ./network_setup.sh down 
    ```

- How to test the cluster:

	In the `sdk` package, it provides some example programs such as [consumer-node](./sdk/examples/consumer-node/main.go), [consumer-seeker](./sdk/examples/consumer-seeker/main.go) [producer](./sdk/examples/producer/main.go), [http_producer](./sdk/examples/http_producer/main.go) to help to use the `EventBus SDK`. 

	These utils can also be used to test this sample cluster:

		- Producer
		
			```
			cd sdk/examples/producer
			go build .
			./producer 127.0.0.1:4150 test01
			```

		- Node Consumer
		
			```
			cd sdk/examples/consumer-node
			go build .
			./consumer-node 127.0.0.1:4150 test01
			```

		- Seeker Consumer
		
			```
			cd sdk/examples/consumer-seeker
			go build .
			./consumer-seeker 127.0.0.1:4161 test01
			```

		- HttpProducer
		
			```
			cd sdk/examples/http_producer
			go build .
			./http_producer 127.0.0.1:8091 ./body.json
			```

## Performance

In the [test/perf](./test/perf) dir, it provides some performance test util programs. 

- [multiple users parallel producer](./test/perf/producer-multi-user-paral.go)

	- Usage:

		```
		./producer-multi-user-paral --help

		Usage of ./producer-multi-user-paral:
		  -nodes string
				eventbus nodes addresses, sepatated by ',' (default "172.16.199.8:4150,172.16.199.8:5150,172.16.199.8:6150")
		  -topic string
				eventbus publish event topic (default "test01")
		  -user int
				concurrent user number (default 10)
		  -worker int
				concurrent worker number (default 10000)
		```

	- Test Sample:

		```
		./producer-multi-user-paral --nodes 172.16.199.8:4150 --topic test01 --user 10 --worker 10000
		```

		Cost: 5.396693401s
		TPS: 18529

- [multiple users parallel consumer](./test/perf/consumer-multi-user.go).

	- Usage:

		```
		./consumer-multi-user --help
		Usage of ./consumer-multi-user:
		  -nodes string
				eventbus nodes addresses, sepatated by ',' (default "172.16.199.8:4150,172.16.199.8:5150,172.16.199.8:6150")
		  -topic string
				eventbus publish event topic (default "test01")
		  -user int
				concurrent worker number (default 10)
		  -worker int
				concurrent worker number (default 10000)
		```

	- Test Sample:

		```
		./consumer-multi-user --nodes 172.16.199.8:4150 --topic test01 --user 10 --worker 10000
		```

		Cost: 11.244940496s
		TPS: 8892

## SDK

Please refer to [EventBus SDK Package](./sdk/README.md).

