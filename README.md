# log-collector

## 日志收集服务

### 基本原理
各个微服务将日志发送到`kafka`，再由日志收集服务从kafka中读取，可以输出到本地的日志文件，或者发送到ElasticSearch，进行日志分析的相关工作、

### 项目结构

#### reader
定义了通用的读取的接口
```go
type Reader interface {
	Read(ctx context.Context, ch chan<- []byte) error
}
```
其中`ctx context.Context`可以传入cancelCtx，可以由上层取消

`ch chan<- []byte`是只写的通道，将读取到的日志传入

目前的实现:
+ [kafka](./reader/kafka.go)

#### writer
定义了通用的写接口
```go
type Writer interface {
Write(data []byte) error

//程序结束时，注意close
Close() error
}
```
为了实现多写（写入多个），`Write`方法并没有将`ch <-chan []byte`作为参数，而是直接接受`data []byte`

`Close`方法是为了释放占用的资源，比如文件句柄

目前的实现:
+ [file](./writer/file.go)
+ [stdout](./writer/stdout.go)

#### collector
```go
type Collector struct {
	reader  []reader.Reader
	writer  []writer.Writer
	MsgChan chan []byte
}
```
`reader`是读取的数据的来源，而`writer`是写入的目的地，`MsgChan chan []byte`是作为读和写之间的中间件，reader将数据传入channel中，writer将其取出

#### config
借助于`viper`实现的，用来读取配置
