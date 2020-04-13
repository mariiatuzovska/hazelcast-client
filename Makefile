start-server:
	../hazelcast-3.12.6/bin/start.sh

# do not work
start-management-center:
	../hazelcast-3.12.6/management-center/start.sh 

map-example:
	./client map-example

queue-example:
	./client queue-example

help: 
	go build -o client
	./client

build:
	go build -o client