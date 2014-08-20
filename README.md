## Example : Avro Schema Generator
Follow these steps to run avro schema generator:
* This examples assumes a running MySql setup. It also assumes the user root with no password.
* git clone aesop (ignore if done already)
* Run `mvn clean install` in `sample-mysql-relay` directory. This might take a while to download Aesop, Trooper and their dependencies from the various Maven repositories. 
* A `distribution` directory will be created in the `sample-mysql-relay` directory. This directory is the distribution which contains the built sample.
* To display help about schema generator, execute following command
```
java -cp "lib/*" com.flipkart.aesop.avro.schemagenerator.main.SchemaGeneratorCli -help
```


* Following command generates avro schema of all tables in database 'or_test' and saves the .avro files of each table in a folder at path '~/Desktop/schema'

```
java -cp "lib/*" com.flipkart.aesop.avro.schemagenerator.main.SchemaGeneratorCli -d or_test -f ~/Desktop/schema -v 1
```
* A message `Written Schema to folder <output-folder>` 




































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
