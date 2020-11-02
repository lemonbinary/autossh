package main

import (
    "strings"
    "fmt"
    
    "github.com/manifoldco/promptui"
)

func search(list []*Node) (int, error) {
    templates := &promptui.SelectTemplates {
        Label:    "※ {{.|green}}",
        Active:   "➤ {{.Name|yellow}} {{.User|yellow}}{{`@`|yellow}}{{.Host|yellow}} {{.Referer|yellow}}",
        Inactive: "  {{.Name|faint}} {{.User|faint}}{{`@`|faint}}{{.Host|faint}} {{.Referer|yellow}}",
        Details: `
{{ "Name:" | faint }} {{ .Name }}
{{ "User:" | faint }} {{ .User }}
{{ "Host:" | faint }} {{ .Host }}
{{ "Port:" | faint }} {{ .Port }}
{{ "Referer:" | faint }} {{ .Referer }}
{{ "CMD:" | faint }} {{ .CMD }}`,
    }
    
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
        return -1, err
    }
    
    return index, nil
}

func option(nodes []*Node) []*Node {
    templates := &promptui.SelectTemplates {
        Label:    "※ {{.|green}}",
        Active:   "➤ {{.|yellow}}",
        Inactive: "  {{.|faint}}",
    }
    
    prompt := promptui.Select {
        Label: "OPTION",
        Items: []string{"ADD", "DEL"},
        Templates: templates,
    }
    
    _, result, _ := prompt.Run()
    if result == "ADD" {
        node := add()
        if node != nil {
            nodes = append(nodes, node)
            saveNode(nodes)
        }
    } else if result == "DEL" {
        nodes = del(nodes)
        saveNode(nodes)
    }
    return nodes
}

func add() *Node {
    temp := &promptui.PromptTemplates {
        Prompt:  "{{.}}",
        Valid:   "{{.|green}}",
        Invalid: "{{.|red}}",
        Success: "{{.|bold}}",
    }

    prompt := promptui.Prompt {
        Label: "ADD (name,host,user,port,password,referer,cmd)",
        Validate: func(input string) error {
            fields := strings.Split(input, ",")
            if len(fields) != 7 {
                return fmt.Errorf("error format")
            }
            
            return nil
        },
        Templates: temp,
        Default: "name,host,user,port,password,referer,cmd",
        AllowEdit: true,
    }
    
    result, err := prompt.Run()
    if err != nil {
        return nil
    }
    fields := strings.Split(result, ",")
    node := &Node {
        Name: fields[0],
        Host: fields[1],
        User: fields[2],
        Port: fields[3],
        Password: fields[4],
        Referer: fields[5],
        CMD: fields[6],
    }
    
    return node
}

func del(nodes []*Node) []*Node {
    index, err := search(nodes)
    if err != nil {
        return nodes
    }
    
    nodes = append(nodes[0:index], nodes[index+1:]...)
    return nodes
}

func startPrompt(nodes []*Node) {
start:
    list := []*Node{{Name:"OPTION"}}
    list = append(list, nodes...)

    index, err := search(list)
    if err != nil {
        fmt.Println("prompt failed:", err)
        return
    }
    
    // add del
    if index == 0 {
        nodes = option(nodes)
        goto start
    }
    
    err = startTerminal(list[index], nodes)
    if err != nil {
        fmt.Println("start terminal failed:", err)
    }
}