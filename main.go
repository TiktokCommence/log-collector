package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log-collector/collector"
	"log-collector/config"
	"log-collector/reader"
	"log-collector/writer"
	"os"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "config/config.yaml", "config path, eg: -conf config.yaml")
}
func main() {
	flag.Parse()
	appConf, err := config.GetConfig(flagconf)
	if err != nil {
		log.Fatalf("get config failed: %v", err)
	}
	var (
		readers []reader.Reader
		writers []writer.Writer
	)
	if appConf.Reader.Kafka != nil {
		KafkaBuilder := reader.NewKafkaReaderBuilder(appConf.Reader.Kafka.BrokersAddr, appConf.Reader.Kafka.Topic)
		kafka, err := KafkaBuilder.Build()
		if err != nil {
			log.Fatalf("create kafka reader failed: %v", err)
		}
		readers = append(readers, kafka)
	}
	if appConf.Writer.File != nil {
		//创建文件夹
		err := checkAndCreateDir(appConf.Writer.File.FilePath)
		if err != nil {
			panic(err)
		}

		FileBuilder := writer.NewFileWriterBuilder(appConf.Writer.File.FilePath, appConf.Writer.File.FileName, appConf.Writer.File.MaxSize, appConf.Writer.File.RotateByTime)
		file, err := FileBuilder.Build()
		if err != nil {
			log.Fatalf("create file writer failed: %v", err)
		}
		writers = append(writers, file)
	}
	if appConf.Writer.Stdout {
		stdoutBuilder := writer.NewStdoutWriterBuilder()
		stdout, err := stdoutBuilder.Build()
		if err != nil {
			log.Fatalf("create stdout writer failed: %v", err)
		}
		writers = append(writers, stdout)
	}
	c := collector.NewCollector(readers, writers, appConf.BuffSize)
	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("begin collect log........")
	err = c.Collect(cctx)
	if err != nil {
		log.Println("collect failed ", err)
		cancel()
	}
}

// checkAndCreateDir 检查文件夹是否存在，如果不存在则创建它
func checkAndCreateDir(dirPath string) error {
	// 检查文件夹是否存在
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		// 文件夹不存在，创建它
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		log.Println("Directory created:", dirPath)
	} else if err != nil {
		// 如果发生其他错误
		return fmt.Errorf("error checking directory: %v", err)
	} else {
		// 文件夹已存在
		log.Println("Directory already exists:", dirPath)
	}
	return nil
}
