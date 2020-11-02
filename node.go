package main

import (
    "os/user"
    "path/filepath"
    "io/ioutil"
    
    "gopkg.in/yaml.v2"
)

type Node struct {
    Name     string `yaml:"name"`
    Host     string `yaml:"host"`
    User     string `yaml:"user"`
    Port     string `yaml:"port"`
    Password string `yaml:"password"`
    Referer  string `yaml:"referer"`
    CMD      string `yaml:"cmd"`
}

func getNodeFile() (string, error) {
    u, err := user.Current()
    if err != nil {
        return "", err
    }
    
    file := filepath.Join(u.HomeDir, ".ssh/assh.yaml")
    return file, nil
}

func loadNode() ([]*Node, error) {
    file, err := getNodeFile()
    if err != nil {
        return nil, err
    }
    
    body, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }
    
    data := []*Node{}
    err = yaml.Unmarshal(body, &data)
    if err != nil {
        return nil, err
    }
    
    return data, nil
}

func saveNode(list []*Node) error {
    file, err := getNodeFile()
    if err != nil {
        return err
    }
    
    b, err := yaml.Marshal(&list)
    if err != nil {
        return err
    }
    
    err = ioutil.WriteFile(file, b, 0644)
    return err
}