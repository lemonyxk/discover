# How to use?

---
## run servers:

### 1.open a terminal and type:
```
$ cd client_1
$ go run ../main.go --config config.json
``` 

### 2.open a terminal and type:
```
$ cd client_2
$ go run ../main.go --config config.json
``` 

### 3.open a terminal and type:
```
$ cd client_3
$ go run ../main.go --config config.json
``` 

### 4.open a terminal and type:
```
$ cd client_4
$ go run ../main.go --config config.json
``` 

### now you got 4 nodes running in you machine.

---


## choose master node:
### 1.choose one of the nodes to be master, such as client_1: 
```
$ curl -X POST http://127.0.0.1:11001/bemaster
  OK
```

### 2.let other nodes to join the cluster:
```
$ curl -X POST -d "addr=127.0.0.1:11002" http://127.0.0.1:11001/join
  OK
$ curl -X POST -d "addr=127.0.0.1:11003" http://127.0.0.1:11001/join
  OK
$ curl -X POST -d "addr=127.0.0.1:11004" http://127.0.0.1:11001/join
  OK
```

### 3.check the cluster status(run in any node):
```
$ curl http://127.0.0.1:11003/serverlist
  [
    {
      "Addr": {
        "Addr": "127.0.0.1:11001",
        "Http": "127.0.0.1:11001",
        "Raft": "127.0.0.1:12001",
        "Tcp": "127.0.0.1:13001"
      },
      "IsMaster": true
    },
    {
      "Addr": {
        "Addr": "127.0.0.1:11002",
        "Http": "127.0.0.1:11002",
        "Raft": "127.0.0.1:12002",
        "Tcp": "127.0.0.1:13002"
      }
    },
    {
      "Addr": {
        "Addr": "127.0.0.1:11003",
        "Http": "127.0.0.1:11003",
        "Raft": "127.0.0.1:12003",
        "Tcp": "127.0.0.1:13003"
      }
    }
  ]
```

### 4.set a value(must run in master node):
```
$ curl -X POST -d "key=test_k&value=test_v" http://127.0.0.1:11001/Set
  OK
```

### 5.get a value:
```
$ curl http://127.0.0.1:11001/Get?key=test_k
  test_v
```

### 6.delete a value:
```
$ curl -X POST -d "key=test_k" http://127.0.0.1:11001/Delete
  OK
```


