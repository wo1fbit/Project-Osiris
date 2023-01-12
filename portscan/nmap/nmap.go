import (
    "fmt"
    "github.com/Ullaakut/nmap"
    "net"
    "os"
    "sync"
    "io"
    "golang.org/x/sync/semaphore"
    "context"
)

// Scan function
func Scan(ips []string, file string) error {
    rate := 20
    for _, ip := range ips {
        if net.ParseIP(ip) == nil {
            return fmt.Errorf("Invalid IP address: %s", ip)
        }
    }

    var f io.Writer
    var err error
    if file != "" {
        f, err = os.Create(file)
        if err != nil {
            return fmt.Errorf("Error creating file %s: %s", file, err)
        }
    } else {
        f = io.MultiWriter(os.Stdout)
    }

    var wg sync.WaitGroup
    sem := semaphore.NewWeighted(int64(rate))

    for _, ip := range ips {
        wg.Add(1)
        go scan(&wg, ip, sem, f)
    }

    wg.Wait()
    if file != "" {
                fmt.Printf("Scan complete! Results written to stdout and %s\n", file)
    }
    return nil
}

func scan(wg *sync.WaitGroup, ip string, sem *semaphore.Weighted, f io.Writer) {
    defer wg.Done()
    sem.Acquire(context.TODO(), 1)
    defer sem.Release(1)

    scanner, err := nmap.NewScanner(
        nmap.WithTargets(ip),
        nmap.WithPorts("1-65535"),
    )
    if err != nil {
        fmt.Println(err)
        return
    }

    result, warnings, err := scanner.Run()
    if err != nil {
        fmt.Println(err)
        return
    }

    if len(warnings) > 0 {
        fmt.Fprintf(f, "Warnings: %v\n", warnings)
    }

    for _, host := range result.Hosts {
        if len(host.Ports) == 0 || len(host.Addresses) == 0 {
            continue
        }

        fmt.Fprintf(f, "Host: %s (%s)\n", host.Addresses[0], host.Ports[0].State)
        for _, port := range host.Ports {
            fmt.Fprintf(f, "\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)
        }
    }
}
