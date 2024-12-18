package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type VipperSetting struct {
	*viper.Viper
}
type AppConfig struct {
	BuffSize uint         `yaml:"buffsize"`
	Reader   ReaderConfig `yaml:"reader"`
	Writer   WriterConfig `yaml:"writer"`
}
type ReaderConfig struct {
	Kafka *KafkaConfig `yaml:"kafka"`
}
type KafkaConfig struct {
	BrokersAddr []string `yaml:"brokersAddr"` //broker的地址
	Topic       string   `yaml:"topic"`       //topic
}
type WriterConfig struct {
	Stdout bool        `yaml:"stdout"`
	File   *FileConfig `yaml:"file"`
}
type FileConfig struct {
	FilePath     string `yaml:"filePath"`     //文件路径
	FileName     string `yaml:"fileName"`     //文件名称
	MaxSize      int64  `yaml:"maxSize"`      //分割的最大size(单位:字节)
	RotateByTime bool   `yaml:"rotateByTime"` //是否根据时间来进行切割
}

func (s *VipperSetting) ReadSection(k string, v interface{}) error {
	err := s.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(fileName string) (AppConfig, error) {
	var appConfig AppConfig
	vp := viper.New()
	vp.SetConfigFile(fileName)
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config file error:%w", err))
	}
	s := &VipperSetting{
		Viper: vp,
	}
	err = s.ReadSection("app", &appConfig)
	if err != nil {
		return AppConfig{}, err
	}
	return appConfig, nil
}
