package klog

import (
	"fmt"
	"log"
	"path"
	"runtime"
)

const (
	Red     = "\033[1;31m%v\033[0m"
	Green   = "\033[1;32m%v\033[0m"
	Yellow  = "\033[1;33m%v\033[0m"
	Blue    = "\033[1;34m%v\033[0m"
	Magenta = "\033[5;35m%v\033[0m"
)

type Publisher interface {
	Publish(topic string, data map[string]any)
}

var (
	logger   = log.New(log.Writer(), "", 0)
	ss       = NewLimitedSlice[string](10)
	pub      Publisher
	topicPub = ""
)

func SaveToMem(size int) {
	if ss == nil {
		ss = NewLimitedSlice[string](size)
	}
}

func UsePublisher(publisher Publisher, topic string) {
	pub = publisher
	topicPub = topic
}

func GetLogs() *LimitedSlice[string] {
	return ss
}

// Printf takes pattern(rd,gr,yl,bl,mg), varsString, varsValues and prints the formatted log message
func Printfs(pattern string, anything ...interface{}) {
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
	msg := fmt.Sprintf(pattern, anything...)
	ss.Add(msg)
	if pub != nil {
		pub.Publish(topicPub, map[string]any{
			"log": msg,
		})
	}
	colorfulLogMessage := "\033[1;" + colorCode + "m" + msg + "\033[0m"
	fmt.Fprint(logger.Writer(), colorfulLogMessage)
}

// Printf takes pattern(rd,gr,yl,bl,mg), varsString, varsValues and prints the formatted log message
func Printf(pattern string, anything ...interface{}) {
	pf := "[INFO]"
	var colorCode string
	var colorUsed = true

	switch pattern[:2] {
	case "rd":
		colorCode = "31"
		pf = "[ERROR]"
	case "gr":
		colorCode = "32"
		pf = "[SUCCESS]"
	case "yl":
		colorCode = "33"
		pf = "[ERROR]"
	case "bl":
		colorCode = "34"
		pf = "[INFO]"
	case "mg":
		colorCode = "35"
		pf = "[DEBUG]"
	default:
		pf = "[INFO]"
		colorUsed = false
		colorCode = "34"
	}

	if colorUsed {
		pattern = pattern[2:]
	}
	pc, file, line, _ := runtime.Caller(1)
	caller := formatCaller(pc, file, line)
	msg := pf + caller + ":" + fmt.Sprintf(pattern, anything...)
	ss.Add(msg)
	if pub != nil {
		pub.Publish(topicPub, map[string]any{
			"log": msg,
		})
	}
	colorfulLogMessage := "\033[1;" + colorCode + "m" + msg + "\033[0m"
	fmt.Fprint(logger.Writer(), colorfulLogMessage)
}

func CheckError(err error) bool {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		caller := formatCaller(pc, file, line)
		msg := "[ERROR]" + caller + ":" + err.Error()
		ss.Add(msg)
		if pub != nil {
			pub.Publish(topicPub, map[string]any{
				"log": msg,
			})
		}
		colorfulLogMessage := "\033[1;31m" + msg + "\n" + "\033[0m"
		fmt.Fprint(logger.Writer(), colorfulLogMessage)
		return true
	}
	return false
}

// formatCaller formats the caller information as desired
func formatCaller(pc uintptr, file string, line int) string {
	_, filename := path.Split(file)
	return fmt.Sprintf("[%s:%d]", filename, line)
}
