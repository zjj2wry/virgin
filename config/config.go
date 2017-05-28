package config

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type Config struct {
	params map[string]string
}

func (this *Config) InitConfig(path string) {
	this.params = make(map[string]string)

	file, err := os.Open(path)
	log.Println(file)
	if err != nil {
		log.Fatal("app.conf can not find !")
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
