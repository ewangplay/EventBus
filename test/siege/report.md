# siege test

## Test
`siege -c 255 -r 1000  "http://172.16.199.8:8091/v1/event POST < ../data/payment-notify.json"`

## Result

### first time
- log level: info

```
Transactions:		      255000 hits
Availability:		      100.00 %
Elapsed time:		      631.42 secs
Data transferred:	       34.25 MB
Response time:		        0.35 secs
Transaction rate:	      403.85 trans/sec
Throughput:		        0.05 MB/sec
Concurrency:		      143.27
Successful transactions:      255000
Failed transactions:	           0
Longest transaction:	       16.83
Shortest transaction:	        0.01
```


### second time
- log level: error

```
Transactions:		      255000 hits
Availability:		      100.00 %
Elapsed time:		      438.42 secs
Data transferred:	       34.04 MB
Response time:		        0.16 secs
Transaction rate:	      581.63 trans/sec
Throughput:		        0.08 MB/sec
Concurrency:		       93.46
Successful transactions:      255000
Failed transactions:	           0
Longest transaction:	       16.93
Shortest transaction:	        0.01
```
