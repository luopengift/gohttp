# gohttp
[![BuildStatus](https://travis-ci.org/luopengift/gohttp.svg?branch=master)](https://travis-ci.org/luopengift/gohttp)
[![GoDoc](https://godoc.org/github.com/luopengift/gohttp?status.svg)](https://godoc.org/github.com/luopengift/gohttp)
[![GoWalker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/luopengift/gohttp)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)
---

gohttp is used for RESTful APIs, Web apps, Http services in Golang.
It is used similar with [Tornado](http://www.tornadoweb.org).

### GO verion >=1.9

## Getting started

#### Sample example

* Simple Server Application: [server](https://github.com/luopengift/gohttp/blob/master/example/server/main.go)
* Simple Client Application: [client](https://github.com/luopengift/gohttp/blob/master/example/client/main.go)

#### Download and Install
```
go get github.com/luopengift/gohttp
```

#### Generate https tls
```
go run  $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```
#### Run
```
go run  $GOPATH/src/github.com/luopengift/gohttp/sample/server.go
```

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Added some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## Author
[@luopengift](luopengift@foxmail.com)

## License

gohttp source code is licensed under the [Apache Licence 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

