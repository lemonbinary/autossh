package main

import (
    "flag"
    "strings"
    "fmt"
    "log"
    "os"
    
    "github.com/manifoldco/promptui"
)

var (
    templates = &promptui.SelectTemplates {
        Label:    "※ {{.|green}}",
        Active:   "➤ {{.Name|yellow}} {{if .Host}}{{if .User}}{{.User|yellow}}{{`@`|yellow}}{{end}}{{.Host|yellow}}{{end}}",
        Inactive: "  {{.Name|faint}} {{if .Host}}{{if .User}}{{.User|faint}}{{`@`|faint}}{{end}}{{.Host|faint}}{{end}}",
    }
    
    debugLog *log.Logger
)

func search(list []*Node) (*Node, error) {
    prompt := promptui.Select {
        Label: "select host (↓ ↑ → ← move and / toggles search)",
        Items: list,
        Templates: templates,
        Size: 20,
        HideSelected: true,
        StartInSearchMode: true,
        Searcher: func(input string, index int) bool {
            n := list[index]
            content := fmt.Sprintf("%s %s %s", n.Name, n.User, n.Host)
            keys := strings.Split(input, " ")
            for _, k := range keys {
                if k == "" {
                    continue
                }
                if !strings.Contains(content, k) {
                    return false
                }
            }
            return true
        },
    }
    index, _, err := prompt.Run()
    if err != nil {
        return nil, err
    }
    
    return list[index], nil
}

var (
    debugFlag *bool
)

func main() {
    debugFlag = flag.Bool("debug", false, "debug to debug.log")
    flag.Parse()
    
    if *debugFlag {
        f, _ := os.OpenFile("./debug.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
        debugLog = log.New(f, "", log.Ldate)
    }
    
    nodes, err := loadConfig()
    if err != nil {
        fmt.Println("load config failed:", err)
        return
    }
    
    if debugLog != nil {
        debugLog.Println("node n", len(nodes))
    }
    
    n, err := search(nodes)
    if err != nil {
        fmt.Println("prompt failed:", err)
        return
    }
    
    err = startTerminal(n)
    if err != nil {
        fmt.Println("start terminal failed:", err)
    }
}