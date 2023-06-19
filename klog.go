package klog

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
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
		ll := fmt.Sprintf(pattern[2:], anything...)
		ss.Add(ll)
		if pub != nil {
			pub.Publish(topicPub, map[string]any{
				"log": ll,
			})
		}
		pattern = "\033[1;" + colorCode + "m" + pattern[2:]
	} else {
		ll := fmt.Sprintf(pattern, anything...)
		ss.Add(ll)
		if pub != nil {
			pub.Publish(topicPub, map[string]any{
				"log": ll,
			})
		}
		pattern = "\033[1;" + colorCode + "m" + pattern + "\033[0m"
	}
	if strings.HasSuffix(pattern, "\n") {
		pattern = pattern[:len(pattern)-1] + "\033[0m\n"
	} else {
		pattern = pattern + "\033[0m"
	}

	fmt.Fprintf(logger.Writer(), pattern, anything...)
}

// Printf takes pattern(rd,gr,yl,bl,mg), varsString, varsValues and prints the formatted log message
func Printf(pattern string, anything ...interface{}) {
	pc, file, line := getCaller()
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

	logMessage := formatLogMessage(pf, pc, file, line, pattern, anything...)
	ss.Add(logMessage)
	if pub != nil {
		pub.Publish(topicPub, map[string]any{
			"log": logMessage,
		})
	}
	colorfulLogMessage := "\033[1;" + colorCode + "m" + logMessage + "\033[0m"
	fmt.Fprint(logger.Writer(), colorfulLogMessage)
}

// CheckError checks if err is not nil, prints it with caller information, and returns true
func CheckError(err error) bool {
	if err != nil {
		pc, file, line := getCaller()
		logMessage := formatLogMessage("[ERROR]", pc, file, line, "%v", err)
		ss.Add(logMessage)
		if pub != nil {
			pub.Publish(topicPub, map[string]any{
				"log": logMessage,
			})
		}
		colorfulLogMessage := fmt.Sprintf(Red, logMessage)
		fmt.Print(colorfulLogMessage)
		return true
	}
	return false
}

// formatCaller formats the caller information as desired
func formatCaller(pc uintptr, file string, line int) string {
	_, filename := path.Split(file)
	return fmt.Sprintf("[%s:%d]", filename, line)
}

// getCaller retrieves the caller information
func getCaller() (uintptr, string, int) {
	pc, file, line, _ := runtime.Caller(2) // Adjust the call depth based on your usage
	return pc, file, line
}

// formatLogMessage formats the log message with the appropriate prefix and caller information
func formatLogMessage(prefix string, pc uintptr, file string, line int, pattern string, anything ...interface{}) string {
	caller := formatCaller(pc, file, line)
	logMessage := fmt.Sprintf("%s%s: %s", prefix, caller, pattern)
	return logMessage
}
