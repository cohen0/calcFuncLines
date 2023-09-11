package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BlackListFiles []string
	BlackListDirs  []string
	Path           string
	RankNum        int
}

var global_conf Config

func init() {
	bs, err := os.ReadFile("./config.yaml")
	if err != nil {
		println(err)
		return
	}

	err = yaml.Unmarshal(bs, &global_conf)
	if err != nil {
		println(err)
		return
	}

	println("parse config success!!")
}
