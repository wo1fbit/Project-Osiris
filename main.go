package main

import (
    "fmt"
    "osiris/portscan/nmap"
)

func main() {
    var ips []string
    var ip string
    fmt.Println("Enter IPs to scan, enter 'done' when finished")

    for {
        fmt.Scanln(&ip)
        if ip == "done" {
            break
        }

        ips = append(ips, ip)
    }

    if len(ips) == 0 {
        fmt.Println("No IPs entered, exiting...")
        return
    }

    if err := nmap.Scan(ips, ""); err != nil {
        fmt.Println(err)
        return
    }
}
