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
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
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

func worker(ctx context.Context, jobs <-chan int, results chan<- PortResult, wg *sync.WaitGroup, target string, timeout time.Duration, doBanner bool) {
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
			connCtx, cancel := context.WithTimeout(ctx, timeout)
			conn, err := dialer.DialContext(connCtx, "tcp", address)
			cancel()
			elapsed := time.Since(start).Milliseconds()

			if err != nil {
				results <- PortResult{Port: port, State: classifyError(err), Elapsed: elapsed}
				continue
			}

			var banner string
			if doBanner {
				_ = conn.SetReadDeadline(time.Now().Add(400 * time.Millisecond))
				buf := make([]byte, 512)
				n, _ := conn.Read(buf)
				if n > 0 {
					raw := string(buf[:n])
					raw = strings.Map(func(r rune) rune {
						if r == '\n' || r == '\r' || (r >= 0x20 && r <= 0x7e) {
							return r
						}
						return -1
					}, raw)
					if len(raw) > 512 {
						raw = raw[:512]
					}
					banner = raw
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
	var doBanner, verbose bool
	var outFile string

	flag.StringVar(&target, "target", "", "hadaf (IP ya hostname) - lazem")
	flag.IntVar(&minPort, "min", 1, "kamtarin port")
	flag.IntVar(&maxPort, "max", 1024, "bishtarin port")
	flag.IntVar(&threads, "threads", 100, "tedad nokh haye hamzaman")
	flag.IntVar(&timeoutMs, "timeout", 800, "timeout etesal (milli sanie)")
	flag.BoolVar(&doBanner, "banner", false, "khandan banner bad az etesal")
	flag.BoolVar(&verbose, "verbose", false, "nashun dadan khoruji mofasal")
	flag.StringVar(&outFile, "out", "scans/leshy_scan.json", "file khoruji JSON")
	flag.Parse()

	if target == "" {
		fmt.Printf("%sKheta: bayad hadaf (IP ya hostname) ro moshakhas konid. Mesal: --target 192.168.1.1%s\n", Red, Reset)
		fmt.Println("Baraye rahnuma: ./leshy --help")
		os.Exit(1)
	}
	if minPort < 1 {
		minPort = 1
	}
	if maxPort > 65535 {
		maxPort = 65535
	}
	if minPort > maxPort {
		fmt.Printf("%sKheta: port avval nemitune az port akhar bozorgtar bashe%s\n", Red, Reset)
		os.Exit(1)
	}
	if threads < 1 {
		threads = 100
	}

	ip := target
	if net.ParseIP(target) == nil {
		ips, err := net.LookupIP(target)
		if err == nil && len(ips) > 0 {
			ip = ips[0].String()
		} else {
			fmt.Printf("%sKheta: nemishe %s ro resolve kard: %v%s\n", Red, target, err, Reset)
			os.Exit(1)
		}
	}

	if err := ensureDir(outFile); err != nil {
		fmt.Printf("%sKheta: nemishe pushe khoruji ro sakht: %v%s\n", Red, err, Reset)
		os.Exit(1)
	}

	fmt.Printf("%sShuru scan rooye %s%s (%s) [port-ha: %d-%d]%s\n", Cyan, Bold, target, ip, minPort, maxPort, Reset)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		fmt.Printf("\n%sScan ghat shod tavasot karbar...%s\n", Red, Reset)
		cancel()
	}()

	timeout := time.Duration(timeoutMs) * time.Millisecond
	jobs := make(chan int, 1000)
	results := make(chan PortResult, 1000)

	var wg sync.WaitGroup
	var collectWg sync.WaitGroup
	found := make([]PortResult, 0, (maxPort-minPort+1))
	collectWg.Add(1)

	totalPorts := maxPort - minPort + 1
	progress := 0
	progressMutex := &sync.Mutex{}
	openPorts := 0

	go func() {
		defer collectWg.Done()
		for r := range results {
			progressMutex.Lock()
			progress++
			if r.State == "open" || verbose {
				banner := r.Banner
				if banner == "" {
					banner = "-"
				}
				fmt.Printf("%sPort %d: %s%s %s%s\n", Green, r.Port, Bold, r.State, func() string {
					if r.Banner != "" {
						return fmt.Sprintf("(banner: %s)", r.Banner)
					}
					return ""
				}(), Reset)
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
		go worker(ctx, jobs, results, &wg, ip, timeout, doBanner)
	}

	startTime := time.Now()

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			progressMutex.Lock()
			percent := float64(progress) / float64(totalPorts) * 100
			fmt.Printf("\r%s%d/%d port scan shod (%.0f%%)%s", Cyan, progress, totalPorts, percent, Reset)
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
		fmt.Printf("%sKheta dar sakht file moghat: %v%s\n", Red, err, Reset)
		os.Exit(1)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(scanRes); err != nil {
		fmt.Printf("%sKheta dar neveshtan JSON: %v%s\n", Red, err, Reset)
		f.Close()
		os.Remove(tmp)
		os.Exit(1)
	}
	f.Close()
	if err := os.Rename(tmp, outFile); err != nil {
		fmt.Printf("%sKheta dar taghir esm file: %v%s\n", Red, err, Reset)
		os.Exit(1)
	}

	fmt.Printf("\n%sScan tamam shod!%s\n", Green, Reset)
	fmt.Printf("%sHadaf: %s%s (%s)%s\n", Cyan, Bold, target, ip, Reset)
	fmt.Printf("%sPort-haye scan shode: %d%s\n", Cyan, totalPorts, Reset)
	fmt.Printf("%sPort-haye baz: %d%s\n", Green, openPorts, Reset)
	fmt.Printf("%sZaman kol: %d milli sanie%s\n", Cyan, elapsed, Reset)
	fmt.Printf("%sFile khoruji: %s%s%s\n", Cyan, Bold, outFile, Reset)
}
