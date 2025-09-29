package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ANSI color codes
const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
	Cyan  = "\033[36m"
	Bold  = "\033[1m"
)

type PortResult struct {
	Port    int    `json:"port"`
	State   string `json:"state"`
	Banner  string `json:"banner,omitempty"`
	Elapsed int64  `json:"elapsed_ms"`
}

type ScanResult struct {
	Target   string       `json:"target"`
	Started  string       `json:"started"`
	Finished string       `json:"finished"`
	Elapsed  int64        `json:"elapsed_ms"`
	Ports    []PortResult `json:"ports"`
}

func classifyError(err error) string {
	if err == nil {
		return "open"
	}
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return "filtered"
	}
	lerr := strings.ToLower(err.Error())
	if strings.Contains(lerr, "refused") || strings.Contains(lerr, "connection refused") {
		return "closed"
	}
	if strings.Contains(lerr, "no route to host") || strings.Contains(lerr, "network is unreachable") {
		return "network-unreachable"
	}
	if strings.Contains(lerr, "i/o timeout") || strings.Contains(lerr, "deadline") {
		return "filtered"
	}
	return "error"
}

func worker(ctx context.Context, jobs <-chan int, results chan<- PortResult, wg *sync.WaitGroup, target string, timeout time.Duration, doBanner bool, lowResource bool) {
	defer wg.Done()
	dialer := &net.Dialer{}
	for {
		select {
		case <-ctx.Done():
			return
		case port, ok := <-jobs:
			if !ok {
				return
			}
			start := time.Now()
			address := net.JoinHostPort(target, strconv.Itoa(port))
			var conn net.Conn
			var err error
			if lowResource {
				connCtx, cancel := context.WithTimeout(ctx, timeout)
				conn, err = dialer.DialContext(connCtx, "tcp", address)
				cancel()
			} else {
				conn, err = dialer.DialContext(ctx, "tcp", address)
			}
			elapsed := time.Since(start).Milliseconds()

			if err != nil {
				results <- PortResult{Port: port, State: classifyError(err), Elapsed: elapsed}
				continue
			}

			var banner string
			if doBanner {
				_ = conn.SetReadDeadline(time.Now().Add(400 * time.Millisecond))
				buf := make([]byte, 256)
				n, _ := conn.Read(buf)
				if n > 0 {
					banner = strings.TrimSpace(string(buf[:n]))
					if len(banner) > 256 {
						banner = banner[:256]
					}
				}
			}
			_ = conn.Close()
			results <- PortResult{Port: port, State: "open", Banner: banner, Elapsed: elapsed}
		}
	}
}

func ensureDir(path string) error {
	d := filepath.Dir(path)
	if d == "." || d == "" {
		return nil
	}
	return os.MkdirAll(d, 0o755)
}

func main() {
	var target string
	var minPort, maxPort, threads int
	var timeoutMs int
	var doBanner, verbose, lowResource bool
	var outFile string

	flag.StringVar(&target, "t", "", "hadaf (IP ya hostname) - lazem")
	flag.IntVar(&minPort, "m", 1, "kamtarin port")
	flag.IntVar(&maxPort, "x", 1024, "bishtarin port")
	flag.IntVar(&threads, "r", 0, "tedad nokh (0 = auto)")
	flag.IntVar(&timeoutMs, "o", 1000, "timeout (ms)")
	flag.BoolVar(&doBanner, "b", false, "khandan banner")
	flag.BoolVar(&verbose, "v", false, "khoruji mofasal")
	flag.BoolVar(&lowResource, "l", false, "kam masraf baraye termux")
	flag.StringVar(&outFile, "u", "/sdcard/leshy_scan.json", "file JSON")
	flag.Parse()

	if target == "" {
		fmt.Printf("%sKheta: -t lazem%s\n", Red, Reset)
		fmt.Println("Rahnuma: ./leshy -h")
		os.Exit(1)
	}
	if minPort < 1 {
		minPort = 1
	}
	if maxPort > 65535 {
		maxPort = 65535
	}
	if minPort > maxPort {
		fmt.Printf("%sKheta: -m > -x nist%s\n", Red, Reset)
		os.Exit(1)
	}
	if lowResource {
		threads = 20
		timeoutMs = 1000
	} else if threads == 0 {
		threads = runtime.NumCPU() * 4
	}
	if threads < 1 {
		threads = 20
	}

	ip := target
	if net.ParseIP(target) == nil {
		ips, err := net.LookupIP(target)
		if err == nil && len(ips) > 0 {
			ip = ips[0].String()
		} else {
			fmt.Printf("%sKheta: %s resolve nashod%s\n", Red, target, Reset)
			os.Exit(1)
		}
	}

	if err := ensureDir(outFile); err != nil {
		fmt.Printf("%sKheta: pushe %s nasakht%s\n", Red, filepath.Dir(outFile), Reset)
		os.Exit(1)
	}

	fmt.Printf("%s>> Scan %s%s%s\n", Cyan, Bold, target, Reset)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		fmt.Printf("\n%s!! Ghat%s\n", Red, Reset)
		cancel()
	}()

	timeout := time.Duration(timeoutMs) * time.Millisecond
	totalPorts := maxPort - minPort + 1
	jobs := make(chan int, totalPorts)
	results := make(chan PortResult, totalPorts)

	var wg sync.WaitGroup
	var collectWg sync.WaitGroup
	found := make([]PortResult, 0, totalPorts)
	collectWg.Add(1)

	progress := 0
	progressMutex := &sync.Mutex{}
	openPorts := 0

	go func() {
		defer collectWg.Done()
		for r := range results {
			progressMutex.Lock()
			progress++
			if r.State == "open" || verbose {
				fmt.Printf("%s%d: %s%s%s\n", Green, r.Port, Bold, r.State, Reset)
				if r.Banner != "" {
					fmt.Printf("%s%s%s\n", Cyan, r.Banner, Reset)
				}
			}
			if r.State == "open" {
				openPorts++
			}
			progressMutex.Unlock()
			found = append(found, r)
		}
	}()

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(ctx, jobs, results, &wg, ip, timeout, doBanner, lowResource)
	}

	startTime := time.Now()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			progressMutex.Lock()
			percent := float64(progress) / float64(totalPorts) * 100
			fmt.Printf("\r%s%.0f%%%s", Cyan, percent, Reset)
			progressMutex.Unlock()
		}
	}()

totalLoop:
	for p := minPort; p <= maxPort; p++ {
		select {
		case <-ctx.Done():
			break totalLoop
		case jobs <- p:
		}
	}
	close(jobs)

	wg.Wait()
	close(results)
	collectWg.Wait()

	finishTime := time.Now()
	elapsed := finishTime.Sub(startTime).Milliseconds()

	sort.Slice(found, func(i, j int) bool { return found[i].Port < found[j].Port })

	scanRes := ScanResult{
		Target:   target,
		Started:  startTime.Format(time.RFC3339),
		Finished: finishTime.Format(time.RFC3339),
		Elapsed:  elapsed,
		Ports:    found,
	}

	tmp := outFile + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		fmt.Printf("%sKheta: file %s nasakht%s\n", Red, tmp, Reset)
		os.Exit(1)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(scanRes); err != nil {
		fmt.Printf("%sKheta: JSON nanevesht%s\n", Red, Reset)
		f.Close()
		os.Remove(tmp)
		os.Exit(1)
	}
	f.Close()
	if err := os.Rename(tmp, outFile); err != nil {
		fmt.Printf("%sKheta: rename %s nashod%s\n", Red, outFile, Reset)
		os.Exit(1)
	}

	fmt.Printf("\n%s>> Tamam!%s\n", Green, Reset)
	fmt.Printf("%s%s%s\n", Cyan, target, Reset)
	fmt.Printf("%sBaz: %d%s\n", Green, openPorts, Reset)
	fmt.Printf("%sZaman: %d ms%s\n", Cyan, elapsed, Reset)
	fmt.Printf("%sFile: %s%s%s\n", Cyan, Bold, outFile, Reset)
}
