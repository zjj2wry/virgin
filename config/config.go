package config

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"config.1kf.com/server/base_go/log"
)

type Config struct {
	params map[string]string
}

func (this *Config) InitConfig(path string) {
	this.params = make(map[string]string)

	file, err := os.Open(path)
	log.Info(file)
	if err != nil {
		log.Error(errors.New("app.conf can not find !"))
		os.Exit(1)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		if strings.Index(s, "#") == 0 {
			continue
		}

		//n1 := strings.Index(s, "[")
		//n2 := strings.LastIndex(s, "]")
		//if n1 > -1 && n2 > -1 && n2 > n1+1 {
		//	this.strcet = strings.TrimSpace(s[n1+1 : n2])
		//	continue
		//}
		//if len(this.strcet) == 0 {
		//	continue
		//}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])
		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}
		key := frist
		this.params[key] = strings.TrimSpace(second)
	}
}

func (this *Config) GetString(key string) string {
	v, found := this.params[key]
	if !found {
		return ""
	}
	return v
}
