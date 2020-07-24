# All-in-one test

## Procedure

### Start Sample Cluster
	Refer to [Sample Cluster](../../README.md)

### Start All-in-one test

- Start single seeker-consumer

	```
	cd ./consumer-seeker/
	go build .
	./consumer-seeker -topic top001
	```

- Start pushing message with single producer and single node

	```
	cd ./producer
	go build .
	./producer -topic top001 -message-count 10000
	```

- Pushing message with multiple producers

	```
	cd ./producer
	go build .
	./producer -topic top001 -message-count 10000 -producer-count 2
	```

- Pushing message with multiple topics

	```
	./producer -topic top001 top002 -message-count 10000
	```

