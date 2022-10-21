package klog

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const (
	Red     = "\033[1;31m%v\033[0m"
	Green   = "\033[1;32m%v\033[0m"
	Yellow  = "\033[1;33m%v\033[0m"
	Blue    = "\033[1;34m%v\033[0m"
	Magenta = "\033[5;35m%v\033[0m"
)

// Printf take pattern(rd,gr,yl,bl,mg), varsString, varsValues
func Printf(pattern string, anything ...any) {
	pc, _, line, _ := runtime.Caller(1)
	var colorCode string
	var colorUsed = true
	switch pattern[:2] {
	case "rd":
		colorCode = "31"
	case "gr":
		colorCode = "32"
	case "yl":
		colorCode = "33"
	case "bl":
		colorCode = "34"
	case "mg":
		colorCode = "35"
	default:
		colorUsed = false
		colorCode = "34"
	}
	if colorUsed {
		pattern = pattern[2:]
	} 
	if strings.HasSuffix(pattern,"\n") {
		pattern = pattern[:len(pattern)-1]
		pattern = "\033[1;" + colorCode+"m"+runtime.FuncForPC(pc).Name()+"[line:"+strconv.Itoa(line)+"]:"+pattern+"\033[0m\n"
	} else {
		pattern = "\033[1;" + colorCode+"m"+runtime.FuncForPC(pc).Name()+"[line:"+strconv.Itoa(line)+"]:"+pattern+"\033[0m"
	}
	fmt.Printf(pattern, anything...)
}

func Printfs(pattern string, anything ...any) {
	var colorCode string
	var colorUsed = true
	switch pattern[:2] {
	case "rd":
		colorCode = "31"
	case "gr":
		colorCode = "32"
	case "yl":
		colorCode = "33"
	case "bl":
		colorCode = "34"
	case "mg":
		colorCode = "35"
	default:
		colorUsed = false
		colorCode = "34"
	}
	if colorUsed {
		pattern = "\033[1;" + colorCode +"m"+pattern[2:]
	} else {
		pattern = "\033[1;" + colorCode +"m"+pattern+"\033[0m"
	}
	if strings.HasSuffix(pattern,"\n") {
		pattern = pattern[:len(pattern)-1]+"\033[0m\n"
	} else {
		pattern = pattern+"\033[0m"
	}
	fmt.Printf(pattern, anything...)
}

// CheckError check if err not nil print it and return true
func CheckError(err error) bool {
	if err != nil {
		pc, _, line, _ := runtime.Caller(1)
		caller := runtime.FuncForPC(pc).Name()
		fmt.Printf("\033[1;31m [ERROR] %s [line:%d] : %v \033[0m \n", caller, line, err)
		return true
	}
	return false
}
