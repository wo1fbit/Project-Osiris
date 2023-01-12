package main

import (
    "osiris/modules/portscan"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "net"
    "bufio"
    "log"
)

func main() {
    fmt.Println("Osiris v0.1.0")
    if len(os.Args) < 3 {
        fmt.Println("Usage: osiris module1:tool1 module2:tool2... [ip_address/ip_file]")
        return
    }

    var ips = make(map[string]bool)
    lastArgument := os.Args[len(os.Args)-1]
    if net.ParseIP(lastArgument) != nil {
        ips[lastArgument] = true
    } else if strings.HasSuffix(lastArgument, ".txt") {
        _, err := os.Stat(lastArgument)
        if os.IsNotExist(err) {
            fmt.Println("File doesn't exist")
            return
        }
        ipFile, err := os.Open(lastArgument)
        if err != nil {
            log.Fatal(err)
        }
        defer ipFile.Close()
        ipScanner := bufio.NewScanner(ipFile)
        for ipScanner.Scan() {
            line := ipScanner.Text()
            if net.ParseIP(line) != nil {
                ips[line] = true
            }
        }
        if err := ipScanner.Err(); err != nil {
            log.Fatal(err)
        }
        }
        for i := 1; i < len(os.Args)-1; i++ {
            input := os.Args[i]
            parts := strings.Split(input, ":")

            if len(parts) != 2 {
                fmt.Println("Usage: osiris module1:tool1 module2:tool2... [ip_address/ip_file]")
                return
            }

            module := parts[0]
            tool := parts[1]

            _, err := exec.LookPath(tool)
            if err != nil {
                fmt.Println(tool, "not found in the system")
                return
            }

            switch module {
            case "portscan":
                var ipSlice []string
                for ip := range ips {
                    ipSlice = append(ipSlice, ip)
                }
                portscan.Tool(tool, ipSlice)
            default:
                fmt.Println("Invalid module")
            }
        }
        if len(ips) == 0 {
            fmt.Println("Please provide at least one valid IP address or a path to a file containing IP addresses")
            return
        }
}
