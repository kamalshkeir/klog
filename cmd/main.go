package main

import (
	"fmt"

	"github.com/kamalshkeir/klog"
)

func main() {
	klog.SaveToMem(3)
	klog.Printfs("rdhello1\n")
	klog.Printfs("rdbye1\n")
	klog.Printf("rdhello2\n")
	klog.Printf("rdbye2\n")
	if v := klog.GetLogs(); v != nil {
		fmt.Println(v.Slice[2])
	}
}
