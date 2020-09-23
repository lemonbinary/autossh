package main

import (
    "encoding/binary"
    "os"
    "time"
    "fmt"
    "os/signal"
    "syscall"
    
    "golang.org/x/crypto/ssh"
    "github.com/moby/term"
)

func startTerminal(node *Node) error {
    config := &ssh.ClientConfig {
        User: node.User,
        Auth: []ssh.AuthMethod{ssh.Password(node.Password)},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout: time.Second * 10,
    }
    
    conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", node.Host, node.Port), config)
    if err != nil {
        return err
    }
    
    session, err := conn.NewSession()
    if err != nil {
        return err
    }
    defer session.Close()
    
    session.Stdout = os.Stdout
    session.Stderr = os.Stderr
    session.Stdin = os.Stdin
    
    modes := ssh.TerminalModes {
        ssh.ECHO: 1,
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }
    
    fd := os.Stdin.Fd()
    if !term.IsTerminal(fd) {
        return fmt.Errorf("not terminal fd")
    }
    
    state, err := term.MakeRaw(fd)
    if err != nil {
        return err
    }
    defer term.RestoreTerminal(fd, state)

    winsize, err := term.GetWinsize(fd)
    if err != nil {
        return err
    }
    
    termWidth  := int(winsize.Width)
    termHeight := int(winsize.Height)
    
    if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
        return err
    }
    
    if err := session.Shell(); err != nil {
        return err
    }
    
    go monitorWindow(session, os.Stdin.Fd())
    session.Wait()
    
    return nil
}

func monitorWindow(session *ssh.Session, fd uintptr) {
    sig := make(chan os.Signal, 1)
    
    signal.Notify(sig, syscall.SIGWINCH)
    defer signal.Stop(sig)
    
    for range sig {
        data, err := termSize(fd)
        if err != nil {
            return
        }
        session.SendRequest("window-change", false, data)
    }
}

func termSize(fd uintptr) ([]byte, error) {
    size := make([]byte, 16)
    
    winsize, err := term.GetWinsize(fd)
    if err != nil {
        return nil, err
    }
    
    binary.BigEndian.PutUint32(size, uint32(winsize.Width))
    binary.BigEndian.PutUint32(size[4:], uint32(winsize.Height))
    
    if debugLog != nil {
        debugLog.Println("winch", winsize.Width, winsize.Height)
    }
    
    return size, nil
}