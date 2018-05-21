# ansible-london-demo
## setup
first of all start two containers by executing
```
docker-compose up -d --build
```
and forget about container technology for a minute (I know it's hard).   
Lets think about these two containers as a phisical machines (or vms somewhere in the cloud).

Double check if our machines are runnnig
```
docker-compose ps
docker ps
```

## demo time
open two terminals

- terminal 1:
```
# set NODE1_ID variable
NODE1_ID=$(docker ps | awk '/ansiblelondondemo_node1_1/ {print $1}')
# watch for ansible running on that machine
docker exec -it $NODE1_ID watch -n1 "ps -e -o command | grep '[a]nsible'"
```

- terminal 2:
```
# set NODE1_ID variable
NODE1_ID=$(docker ps | awk '/ansiblelondondemo_node1_1/ {print $1}')
# add app role to node1
docker exec -it -e CONSUL_HTTP_ADDR=consul:8500 $NODE1_ID consul kv put nodes/node1/roles/app
# hit /hello endpoint
curl localhost/hello

# add more roles
docker exec -it $NODE1_ID ls -la /tmp | grep dummy
docker exec -it -e CONSUL_HTTP_ADDR=consul:8500 $NODE1_ID consul kv put nodes/node1/roles/dummy_role1
docker exec -it -e CONSUL_HTTP_ADDR=consul:8500 $NODE1_ID consul kv put nodes/node1/roles/dummy_role2
docker exec -it $NODE1_ID ls -la /tmp | grep dummy

# remove dummy_role1
docker exec -it $NODE1_ID ls -la /tmp | grep dummy
docker exec -it -e CONSUL_HTTP_ADDR=consul:8500 $NODE1_ID consul kv delete nodes/node1/roles/dummy_role1
docker exec -it $NODE1_ID ls -la /tmp | grep dummy

# cleanup
docker-compose down
```
