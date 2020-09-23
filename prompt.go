package main

import (
    "strconv"
    "strings"
    "fmt"
    
    "github.com/manifoldco/promptui"
)

func search(list []*Node) (int, error) {
    templates := &promptui.SelectTemplates {
        Label:    "※ {{.|green}}",
        Active:   "➤ {{.Name|yellow}} {{if .Host}}{{if .User}}{{.User|yellow}}{{`@`|yellow}}{{end}}{{.Host|yellow}}{{end}}",
        Inactive: "  {{.Name|faint}} {{if .Host}}{{if .User}}{{.User|faint}}{{`@`|faint}}{{end}}{{.Host|faint}}{{end}}",
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
        Label: "ADD (name,host,user,port,password)",
        Validate: func(input string) error {
            fields := strings.Split(input, ",")
            for _, f := range fields {
                if f == "" {
                    return fmt.Errorf("empty field is not allow")
                }
            }
            if len(fields) != 5 {
                return fmt.Errorf("error format")
            }
            
            return nil
        },
        Templates: temp,
        Default: "name,host,user,port,password",
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
        Password: fields[4],
    }
    node.Port, _ = strconv.Atoi(fields[3])
    
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
    list := []*Node{&Node{Name:"OPTION"}}
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
    
    err = startTerminal(list[index])
    if err != nil {
        fmt.Println("start terminal failed:", err)
    }
}