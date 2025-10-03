package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
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

// Noskhe barname
const Noskhe = "0.2.16"

type NatijePort struct {
	Port     int    `json:"port"`
	Protokol string `json:"protokol"`
	Vaziyat  string `json:"vaziyat"`
	Baner    string `json:"baner,omitempty"`
	Noskhe   string `json:"noskhe,omitempty"`
	Zaman    int64  `json:"zaman_ms"`
}

type NatijeScan struct {
	Hadaf    string       `json:"hadaf"`
	Shuru    string       `json:"shuru"`
	Tamam    string       `json:"tamam"`
	ZamanKol int64        `json:"zaman_ms"`
	PortHa   []NatijePort `json:"port_ha"`
}

func tasnifKheta(err error, protokol string) string {
	if err == nil {
		return "baz"
	}
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return "baz|filtered"
	}
	lerr := strings.ToLower(err.Error())
	if strings.Contains(lerr, "refused") || strings.Contains(lerr, "connection refused") {
		return "baste"
	}
	if strings.Contains(lerr, "no route to host") || strings.Contains(lerr, "network is unreachable") {
		return "shabake-napadid"
	}
	if strings.Contains(lerr, "i/o timeout") || strings.Contains(lerr, "deadline") {
		return "baz|filtered"
	}
	return "kheta"
}

// decodeDNSResponse: Extract TXT record from DNS response
func decodeDNSResponse(data []byte) string {
	if len(data) < 12 {
		return ""
	}
	// Skip DNS header (12 bytes)
	pos := 12
	// Read query name
	for pos < len(data) && data[pos] != 0 {
		pos += int(data[pos]) + 1
	}
	pos++ // Skip null byte
	if pos >= len(data) {
		return ""
	}
	// Skip QTYPE and QCLASS (4 bytes)
	pos += 4
	if pos >= len(data) {
		return ""
	}
	// Read answer section
	for pos < len(data) && data[pos] != 0 {
		pos += int(data[pos]) + 1
	}
	pos++ // Skip null byte
	if pos+10 >= len(data) {
		return ""
	}
	// Skip TYPE, CLASS, TTL, RDLENGTH (10 bytes)
	pos += 10
	if pos >= len(data) {
		return ""
	}
	// Read TXT record length
	txtLen := int(data[pos])
	pos++
	if pos+txtLen > len(data) {
		return ""
	}
	// Extract TXT record
	txt := string(data[pos : pos+txtLen])
	return strings.TrimSpace(txt)
}

func girBaner(conn net.Conn, port int, protokol string, noepayload string) (baner, noskhe string) {
	if protokol == "udp" {
		if payload := GirPayloadUDP(noepayload); payload != nil {
			fmt.Fprintf(os.Stderr, "%sErsal payload UDP baraye port %d (%s): %v%s\n", Cyan, port, noepayload, payload, Reset)
			_, err := conn.Write(payload)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%sKheta dar ersal payload UDP baraye port %d: %v%s\n", Red, port, err, Reset)
				return "", ""
			}
		} else {
			fmt.Fprintf(os.Stderr, "%sPayload UDP baraye %s namojud ast%s\n", Red, noepayload, Reset)
			return "", ""
		}
		_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	} else {
		_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sKheta dar khandan baner baraye port %d: %v%s\n", Red, port, err, Reset)
		return "", ""
	}
	if n == 0 {
		fmt.Fprintf(os.Stderr, "%sHich dadei az port %d (%s) daryaft nashod%s\n", Red, port, protokol, Reset)
		return "", ""
	}

	// Extract ASCII-only for cleaner banner
	baner = ""
	if noepayload == "version.bind" && protokol == "udp" && port == 53 {
		baner = decodeDNSResponse(buf[:n])
		noskhe = baner
	} else {
		for _, b := range buf[:n] {
			if b >= 32 && b <= 126 { // Printable ASCII
				baner += string(b)
			}
		}
		baner = strings.TrimSpace(baner)
		if len(baner) > 1024 {
			baner = baner[:1024]
		}
	}

	if protokol == "tcp" && (port == 80 || port == 443) {
		_, err := conn.Write([]byte("GET / HTTP/1.0\r\nHost: localhost\r\n\r\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sKheta dar ersal HTTP request baraye port %d: %v%s\n", Red, port, err, Reset)
			return baner, ""
		}
		_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		n, err := conn.Read(buf)
		if err == nil && n > 0 {
			response := string(buf[:n])
			lines := strings.Split(response, "\r\n")
			for _, line := range lines {
				if strings.HasPrefix(strings.ToLower(line), "server:") {
					noskhe = strings.TrimSpace(strings.TrimPrefix(line, "Server:"))
					if noskhe != "" {
						baner = noskhe
					}
				}
			}
		}
	}

	if noskhe == "" && baner != "" && noepayload != "version.bind" {
		lowerBaner := strings.ToLower(baner)
		if strings.Contains(lowerBaner, "ssh-") || strings.Contains(lowerBaner, "ftp") || strings.Contains(lowerBaner, "220 ") {
			noskhe = strings.SplitN(baner, "\n", 2)[0]
		} else if strings.Contains(lowerBaner, "bind") {
			if idx := strings.Index(baner, "BIND"); idx >= 0 {
				noskhe = strings.TrimSpace(baner[idx:])
			}
		} else if strings.Contains(lowerBaner, "ntpd") || strings.Contains(lowerBaner, "ntp") {
			if idx := strings.Index(lowerBaner, "ntpd"); idx >= 0 {
				noskhe = strings.TrimSpace(baner[idx:])
			} else if idx := strings.Index(lowerBaner, "ntp"); idx >= 0 {
				noskhe = strings.TrimSpace(baner[idx:])
			}
		} else if strings.Contains(lowerBaner, "snmp") {
			noskhe = strings.TrimSpace(baner)
		}
	}

	fmt.Fprintf(os.Stderr, "%sBaner daryaft shod baraye port %d (%s): %s%s\n", Cyan, port, protokol, baner, Reset)
	return baner, noskhe
}

func kargar(ctx context.Context, jobs <-chan int, natayej chan<- NatijePort, wg *sync.WaitGroup, hadaf string, timeout time.Duration, doBaner bool, kamMasraf bool, protokol string, noepayload string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case port, ok := <-jobs:
			if !ok {
				return
			}
			start := time.Now()
			addr := net.JoinHostPort(hadaf, strconv.Itoa(port))
			var conn net.Conn
			var err error
			if protokol == "tcp" {
				dialer := &net.Dialer{}
				if kamMasraf {
					ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
					conn, err = dialer.DialContext(ctxTimeout, "tcp", addr)
					cancel()
				} else {
					conn, err = dialer.DialContext(ctx, "tcp", addr)
				}
			} else {
				udpAddr, errResolve := net.ResolveUDPAddr("udp", addr)
				if errResolve != nil {
					natayej <- NatijePort{Port: port, Protokol: protokol, Vaziyat: tasnifKheta(errResolve, protokol), Zaman: time.Since(start).Nanoseconds() / 1e6}
					continue
				}
				var udpConn *net.UDPConn
				udpConn, err = net.DialUDP("udp", nil, udpAddr)
				if err == nil {
					if kamMasraf {
						_ = udpConn.SetReadDeadline(time.Now().Add(timeout))
					} else {
						_ = udpConn.SetReadDeadline(time.Now().Add(3 * time.Second))
					}
				}
				conn = udpConn
			}
			if err != nil {
				natayej <- NatijePort{Port: port, Protokol: protokol, Vaziyat: tasnifKheta(err, protokol), Zaman: time.Since(start).Nanoseconds() / 1e6}
				continue
			}

			var baner, noskhe string
			if doBaner {
				baner, noskhe = girBaner(conn, port, protokol, noepayload)
			}
			if conn != nil {
				if closeErr := conn.Close(); closeErr != nil {
					fmt.Fprintf(os.Stderr, "%sKheta dar baste shodan port %d: %v%s\n", Red, port, closeErr, Reset)
				}
			}
			vaziyat := tasnifKheta(err, protokol)
			if baner == "" && noskhe == "" && protokol == "udp" {
				vaziyat = "baz|filtered"
			}
			natayej <- NatijePort{Port: port, Protokol: protokol, Vaziyat: vaziyat, Baner: baner, Noskhe: noskhe, Zaman: time.Since(start).Nanoseconds() / 1e6}
		}
	}
}

func sakhtPushe(path string) error {
	d := filepath.Dir(path)
	if d == "." || d == "" {
		return nil
	}
	return os.MkdirAll(d, 0o755)
}

func main() {
	var hadaf string
	var minPort, maxPort, nokh int
	var timeoutMs int
	var doBaner, mofasal, kamMasraf, makhfi, neshanNoskhe bool
	var fileJSON, protokol, noepayload string
	var udpTimeout int

	flag.StringVar(&hadaf, "t", "", "hadaf (IP ya hostname) - lazem")
	flag.IntVar(&minPort, "m", 1, "kamtarin port")
	flag.IntVar(&maxPort, "x", 1024, "bishtarin port")
	flag.IntVar(&nokh, "r", 0, "tedad nokh (0 = auto)")
	flag.IntVar(&timeoutMs, "o", 1000, "timeout (ms, faghat ba -l)")
	flag.IntVar(&udpTimeout, "w", 3000, "timeout UDP (ms, baraye baner)")
	flag.BoolVar(&doBaner, "b", false, "khandan baner va noskhe")
	flag.BoolVar(&mofasal, "v", false, "khoruji mofasal")
	flag.BoolVar(&kamMasraf, "l", false, "kam masraf baraye termux")
	flag.BoolVar(&makhfi, "f", false, "makhfi az firewall")
	flag.StringVar(&fileJSON, "u", "/sdcard/leshy_scan.json", "file JSON")
	flag.StringVar(&protokol, "p", "tcp", "protokol (tcp ya udp)")
	flag.StringVar(&noepayload, "y", "none", "payload baraye udp (dns, version.bind, ntp, snmp, test, ya none)")
	flag.BoolVar(&neshanNoskhe, "V", false, "neshan dadan noskhe barname")
	flag.Parse()

	if neshanNoskhe {
		fmt.Printf("Leshy Scanner Port, Noskhe: %s\n", Noskhe)
		os.Exit(0)
	}

	if hadaf == "" {
		fmt.Printf("%sKheta: -t lazem ast%s\n", Red, Reset)
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
		fmt.Printf("%sKheta: kamtarin port (-m) bayad az bishtarin port (-x) kochektar bashad%s\n", Red, Reset)
		os.Exit(1)
	}
	if protokol != "tcp" && protokol != "udp" {
		fmt.Printf("%sKheta: protokol (-p) bayad 'tcp' ya 'udp' bashad%s\n", Red, Reset)
		os.Exit(1)
	}
	if protokol == "udp" && noepayload != "dns" && noepayload != "ntp" && noepayload != "snmp" && noepayload != "test" && noepayload != "version.bind" && noepayload != "none" {
		fmt.Printf("%sKhta: payload (-y) bayad 'dns', 'ntp', 'snmp', 'test', 'version.bind', ya 'none' bashad%s\n", Red, Reset)
		os.Exit(1)
	}
	if kamMasraf {
		nokh = 20
		timeoutMs = 1000
	} else if nokh == 0 {
		nokh = runtime.NumCPU() * 8
		if nokh > 100 {
			nokh = 100
		}
	}
	if nokh < 1 {
		nokh = 20
	}
	if udpTimeout < 500 {
		udpTimeout = 500
	}

	ip := hadaf
	if net.ParseIP(hadaf) == nil {
		ips, err := net.LookupIP(hadaf)
		if err != nil || len(ips) == 0 {
			fmt.Printf("%sKheta: %s resolve nashod%s\n", Red, hadaf, Reset)
			os.Exit(1)
		}
		ip = ips[0].String()
	}

	if err := sakhtPushe(fileJSON); err != nil {
		fmt.Printf("%sKheta: sakht pushe %s momken nashod: %v%s\n", Red, filepath.Dir(fileJSON), err, Reset)
		os.Exit(1)
	}

	fmt.Printf("%s>> Scan %s%s (%s)%s\n", Cyan, Bold, hadaf, protokol, Reset)
	if makhfi {
		fmt.Printf("%s>> Makhfi: port-haye tasadofi + dirang%s\n", Cyan, Reset)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		fmt.Printf("\n%s!! Ghat shod%s\n", Red, Reset)
		cancel()
	}()

	timeout := time.Duration(timeoutMs) * time.Millisecond
	totalPorts := maxPort - minPort + 1
	jobs := make(chan int, totalPorts)
	natayej := make(chan NatijePort, totalPorts)

	var wg sync.WaitGroup
	var collectWg sync.WaitGroup
	natayeje := make([]NatijePort, 0, totalPorts)
	collectWg.Add(1)

	progress := 0
	progressMutex := &sync.Mutex{}
	bazPorts := 0

	go func() {
		defer collectWg.Done()
		for r := range natayej {
			progressMutex.Lock()
			progress++
			if r.Vaziyat == "baz" || r.Vaziyat == "baz|filtered" || mofasal {
				fmt.Printf("%s%d/%s: %s%s%s\n", Green, r.Port, r.Protokol, Bold, r.Vaziyat, Reset)
				if r.Baner != "" {
					fmt.Printf("%sBaner: %s%s\n", Cyan, r.Baner, Reset)
				}
				if r.Noskhe != "" {
					fmt.Printf("%sNoskhe: %s%s\n", Cyan, r.Noskhe, Reset)
				}
			}
			if r.Vaziyat == "baz" || r.Vaziyat == "baz|filtered" {
				bazPorts++
			}
			progressMutex.Unlock()
			natayeje = append(natayeje, r)
		}
	}()

	for i := 0; i < nokh; i++ {
		wg.Add(1)
		go kargar(ctx, jobs, natayej, &wg, ip, timeout, doBaner, kamMasraf, protokol, noepayload)
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

	// Enqueue ports
	ports := make([]int, 0, totalPorts)
	for p := minPort; p <= maxPort; p++ {
		ports = append(ports, p)
	}
	if makhfi {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(ports), func(i, j int) { ports[i], ports[j] = ports[j], ports[i] })
	}
totalLoop:
	for _, p := range ports {
		select {
		case <-ctx.Done():
			break totalLoop
		case jobs <- p:
			if makhfi {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			}
		}
	}
	close(jobs)

	wg.Wait()
	close(natayej)
	collectWg.Wait()

	finishTime := time.Now()
	zamanKol := finishTime.Sub(startTime).Milliseconds()

	sort.Slice(natayeje, func(i, j int) bool { return natayeje[i].Port < natayeje[j].Port })

	scanRes := NatijeScan{
		Hadaf:    hadaf,
		Shuru:    startTime.Format(time.RFC3339),
		Tamam:    finishTime.Format(time.RFC3339),
		ZamanKol: zamanKol,
		PortHa:   natayeje,
	}

	tmp := fileJSON + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		fmt.Printf("%sKheta: sakht file %s momken nashod: %v%s\n", Red, tmp, err, Reset)
		os.Exit(1)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(scanRes); err != nil {
		fmt.Printf("%sKheta: neveshtan JSON be %s momken nashod: %v%s\n", Red, tmp, err, Reset)
		f.Close()
		os.Remove(tmp)
		os.Exit(1)
	}
	if err := f.Close(); err != nil {
		fmt.Printf("%sKheta: baste shodan file %s momken nashod: %v%s\n", Red, tmp, err, Reset)
		os.Remove(tmp)
		os.Exit(1)
	}
	if err := os.Rename(tmp, fileJSON); err != nil {
		fmt.Printf("%sKheta: rename %s be %s momken nashod: %v%s\n", Red, tmp, fileJSON, err, Reset)
		os.Remove(tmp)
		os.Exit(1)
	}

	fmt.Printf("\n%s>> Tamam!%s\n", Green, Reset)
	fmt.Printf("%s%s (%s)%s\n", Cyan, hadaf, protokol, Reset)
	fmt.Printf("%sBaz: %d%s\n", Green, bazPorts, Reset)
	fmt.Printf("%sZaman: %d ms%s\n", Cyan, zamanKol, Reset)
	fmt.Printf("%sFile: %s%s%s\n", Cyan, Bold, fileJSON, Reset)
}
