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
    Port     int    `yaml:"port"`
    Password string `yaml:"password"`
}

func loadConfig() ([]*Node, error) {
    u, err := user.Current()
    if err != nil {
        return nil, err
    }
    
    file := filepath.Join(u.HomeDir, ".ssh/assh.yaml")
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