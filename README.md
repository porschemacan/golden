# Golden Server
a micro service framework.

## build

```
make
```
## Features
 - only support golang, keep framework light and simple;
 - support multiple serialization protocols, include json and protobuf, maybe thrift in future;
 - support service governance.

 	support CI, it's unittest friendly, it contain unit test framework, and mock framework; it's also system test friendly, it support gracefully exit to generate coverage data.

 	support operation and maintenance, it contain log, metrics, open tracing and dockerfile, which are closely related to operation and maintenance. log, metrics and open tracing are so called three pillars of service observability, and dockerfile helps with easy deployment.

 	support service call feature, include loadbalance and the mechanism of fuse and detect. also support limiter, include three different limiter: counter, leaky bucket, and token bucket.

 - support improvement, support low precision timer.

## RoadMap
 - √ golang server framework
 - √ rest api
 - √ json protocol
 - √ protobuf protocol
 - √ thrift protocol
 - √ redis support
 - √ limiter
 - √ mq producer
 - X mq comsumer
 - √ passport support
 - √ unit test framework
 - √ mock framework
 - √ gracefully exit
 - √ odin pprof monitor
 - √ auto tracing
 - √ log library
 - X metrics library
 - X Dockerfile support
 - X load balance
 - X fuse and detect
 - X low precision timer
