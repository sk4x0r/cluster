#GoCluster
This a library written in `go` to allow a group of servers to communicate with each other.

#Installation and Test

```
$ go get github.com/sk4x0r/cluster
$ go build github.com/sk4x0r/cluster
```

#Dependency
Library depends upon ZeroMQ 4, which can be installed from github.com/pebbe/zmq4


#Usage
New instance of server can be created by method `New()`.
```
s1:=cluster.New(serverId, configFile)
```

`New()` method takes two parameter, and returns object of type `Server`.

| Parameter		| Type		| Description  
| -------------|:---------:| -----------
| serverId		| `int` 	| unuque id assigned to each server
| configFile	| `string`  | path of the file containing configuration of all the servers

For example of configuration file, see _config.json_ file in source.

Running instance of a server can be killed using `StopServer()` method.
```
s1.StopServer()
```

# License

The package is available under GNU General Public License. See the _LICENSE_ file.
