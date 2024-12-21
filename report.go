package main

import (
	"fmt"
	"sort"
	"time"
)

type Report struct {
	Line int
	Path string
	Func string
}

var reports SortReport

type SortReport []Report

func (s SortReport) Less(i, j int) bool {
	return s[i].Line > s[j].Line
}

func (s SortReport) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortReport) Len() int {
	return len(s)
}

func (s SortReport) append(r Report) {
	reports = append(reports, r)
	sort.Sort(SortReport(reports))
}

func (s *SortReport) TryInsert(r Report) {
	if global_conf.RankNum == 0 {
		return
	}

	if s.Len() < global_conf.RankNum {
		s.append(r)
		return
	}

	min := (*s)[s.Len()-1].Line

	if r.Line > min {
		(*s) = (*s)[:s.Len()-1]
		s.append(r)
		return
	}
}

func (s SortReport) Print() {
	fmt.Println("********************")
	fmt.Println("Common:")
	fmt.Printf("\t   use time: %g s\n", time.Since(begin).Seconds())
	fmt.Printf("\tfile number: %d\n", fileNumber)
	fmt.Printf("\tfunc number: %d\n", funcNumber)
	fmt.Printf("\tsum   lines: %d\n", sumLines)
	fmt.Println("--------------------")

	for i, info := range s {
		fmt.Printf("#%d\n", i+1)
		fmt.Printf("\t path: %s\n", info.Path)
		fmt.Printf("\t func: %s\n", info.Func)
		fmt.Printf("\tlines: %d\n", info.Line)
	}
	fmt.Println("********************")
}
