package main

import (
    "flag"
    "fmt"
    "log"
    "os"
)

var (
    debugLog *log.Logger
)

func main() {
    debugFlag := flag.Bool("debug", false, "debug to debug.log")
    flag.Parse()
    
    if *debugFlag {
        f, _ := os.OpenFile("./debug.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
        debugLog = log.New(f, "", log.Ldate)
    }
    
    nodes, err := loadNode()
    if err != nil {
        fmt.Println("load node failed:", err)
        return
    }
    
    if debugLog != nil {
        debugLog.Println("node n", len(nodes))
    }
    
    startPrompt(nodes)
}