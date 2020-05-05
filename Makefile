build:
	go build -o client

start-server-1:
	../hazelcast-3.12.7-server-1/bin/start.sh

start-server-2:
	../hazelcast-3.12.7-server-2/bin/start.sh

start-server-3:
	../hazelcast-3.12.7-server-3/bin/start.sh