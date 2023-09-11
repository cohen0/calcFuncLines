package main

import (
	"bufio"
	"flag"

	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	fileNumber int
	funcNumber int
	inpath     string
	begin      time.Time
	sumLines   int
)

func init() {
	flag.StringVar(&inpath, "path", "E:\\workspace\\test\\", "input path")
	begin = time.Now()
}

func getFuncName(str string) string {
	arr := strings.Split(str, "func ")

	return arr[1]
}

func blacklist(path string) bool {
	ext := filepath.Ext(path)
	if ext != ".go" {
		return true
	}

	_, f := filepath.Split(path)

	for _, file := range global_conf.BlackListFiles {
		matchd, _ := filepath.Match(file, f)
		if matchd {
			return true
		}
	}

	return false
}

func processFile(path string) {
	if blacklist(path) {
		return
	}

	file, err := os.Open(path)
	if err != nil {
		println(err)
		return
	}
	defer file.Close()

	fileNumber++

	scanner := bufio.NewScanner(file)
	count := 0
	var funcname string
	var record bool
	var multiAnnotate bool

	for scanner.Scan() {
		line := scanner.Text()

		tmp := strings.TrimSpace(line)
		if len(tmp) == 0 { //空行
			continue
		}
		if strings.Index(tmp, "//") == 0 { //单行注释
			continue
		}
		if strings.Count(tmp, "/*") > 0 { //多行注释开始
			multiAnnotate = true
		}
		if strings.Count(tmp, "*/") > 0 { //多行注释结束
			multiAnnotate = false
		}

		if multiAnnotate {
			continue
		}

		sumLines++

		//func start
		if strings.Index(line, "func ") == 0 {
			funcname = getFuncName(line)
			record = true
			funcNumber++
		}

		//func end
		if record == true && strings.IndexByte(line, '}') == 0 {
			record = false
			reports.TryInsert(count-1, path, funcname)
			count = 0
		}

		if record {
			count++
		}
	}

	return
}

func blacklistDir(path string) bool {
	for _, dir := range global_conf.BlackListDirs {
		if dir == path {
			return true
		}
	}

	return false
}

func processDir(path string) {
	file, err := os.Open(path)
	if err != nil {
		println(err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if !fileInfo.IsDir() {
		println("input error: path is not a dir!!")
		return
	}

	fileinfos, err := file.ReadDir(0)
	if err != nil {
		println(err)
		return
	}

	for _, info := range fileinfos {
		if info.IsDir() {
			if blacklistDir(info.Name()) {
				continue
			}
			subpath := filepath.Join(path, info.Name())
			processDir(subpath)
			continue
		}

		processFile(filepath.Join(path, info.Name()))
	}
}

func main() {
	flag.Parse()
	processDir(global_conf.Path)
	reports.Print()
}
