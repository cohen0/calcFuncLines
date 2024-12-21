package main

import (
	"bufio"
	"sync/atomic"

	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	fileNumber int32
	funcNumber int32
	sumLines   int32
	begin      time.Time
)

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

func ProcessFile(param interface{}) []Report {
	path := param.(string)
	var dates []Report

	if blacklist(path) {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		println(err)
		return nil
	}
	defer file.Close()

	atomic.AddInt32(&fileNumber, 1)

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

		atomic.AddInt32(&sumLines, 1)

		//func start
		if strings.Index(line, "func ") == 0 {
			funcname = getFuncName(line)
			record = true
			atomic.AddInt32(&funcNumber, 1)
		}
		//func end
		if record && strings.IndexByte(line, '}') == 0 {
			record = false
			r := Report{count - 1, path, funcname}
			dates = append(dates, r)
			count = 0
		}

		if record {
			count++
		}
	}

	return dates
}

func blacklistDir(path string) bool {
	for _, dir := range global_conf.BlackListDirs {
		if dir == path {
			return true
		}
	}

	return false
}

func processDir(path string, pool *TaskPool) {
	var err error
	file, err := os.Open(path)
	if err != nil {
		println(err)
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
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
			processDir(subpath, pool)
			continue
		}

		file := filepath.Join(path, info.Name())
		pool.AddTask(file, ProcessFile)
	}
}

func main() {
	pool := NewTaskPool()
	go pool.Run()
	processDir(global_conf.Path, pool)
	pool.Stop()
	reports.Print()
}
