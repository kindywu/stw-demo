# 执行以下命令，进行压力测试
* go test -bench=BenchmarkHttpHandler -benchmem -benchtime=10s -v
* go test -bench=BenchmarkHttpJson -benchmem -benchtime=10s -v
* go test -bench=BenchmarkHttpArray -benchmem -benchtime=10s -v
* go test -bench=BenchmarkHttpArrWithPool -benchmem -benchtime=10s -v

# Fiber is a Go web framework built on top of Fasthttp, the fastest HTTP engine for Go. It's designed to ease things up for fast development with zero memory allocation and performance in mind.
这个注意是用来做跟标准的Http库作比较。
