# EventBus Service

`EventBus` is secure, reliable and efficient message forwarding service.

## Features

- Support `Notifier` and `Queuer` type drivers for event.

- Provide Restful API to let user to send event -- either `Notifier` or `Queuer` type event.

## Setup

### Manual Setup

- Build service

	```
	$ git clone https://github.com/ewangplay/eventbus.git
	$ cd eventbus
	$ make
	```

- Start eventbus service

	```
	$ /path/to/eventbus --config=/path/to/eventbus.yaml
	```

### Docker Setup

- Build docker image

    ```
    make docker
    ```

- Run docker container

	```
	docker run -d --name eventbus -v /opt/eventbus/log:/opt/eventbus/log -v /opt/eventbus/etc:/opt/eventbus/etc -p 0.0.0.0:8091:8091 ewangplay/eventbus
	```

    __Note:__ Volume option (/opt/eventbus/etc), the directory should contain the eventbus service configure file: [eventbus.yaml](./sampleconfig/eventbus.yaml).
