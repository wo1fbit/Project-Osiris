package portscan

import (
    "osiris/modules/portscan/nmap"
    "fmt"
)

func Tool(tool string, ips []string) {
    if len(ips) == 0 {
        fmt.Println("No IPs entered, exiting...")
        return
    }
    switch tool{
    case "nmap":
        fmt.Println("Module: Portscan Tool:", tool)
        call_nmap(ips)
    default:
        fmt.Printf("%s not found\n", tool)
    }
}

func call_nmap(ips []string) {
    if err := nmap.Scan(ips, ""); err != nil {
        fmt.Println(err)
        return
    }
}
