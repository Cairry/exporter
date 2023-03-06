package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

var (
	// yaml 转 map
	resultMap = make(map[string]map[string]string)
	// DomainMap 获取最终子map
	DomainMap = make(map[string]string)
	ServerMap = make(map[string]string)
)

func Config() {

	// 打开文件
	file, err := os.Open("config/config.yaml")
	if err != nil {
		fmt.Println(err)
	}

	// 关闭文件
	defer file.Close()

	// 读取文件
	data, _ := ioutil.ReadAll(file)

	// 解析文件，输出给 resultMap, out: map[dominlist:map[baidu:www.baidu.com qq:www.qq.com]]
	if err := yaml.Unmarshal(data, &resultMap); err != nil {
		log.Fatalln(err)
	}

	// 获取 resultMap 中的value 子map「如：map[baidu:www.baidu.com qq:www.qq.com]」
	for key, value := range resultMap["domainlist"] {
		DomainMap[key] = value
	}

	// 获取端口号
	for portKey, portValue := range resultMap["server"] {
		ServerMap[portKey] = portValue
	}

}
