package main

import (
	"os/signal"
	"syscall"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type PortInfo struct {
	Port     int
	Protocol string
	Process  string
	PID      string
	Command  string
	Address  string
}

func main() {
	portmap := scanPorts()
	displayMap(portmap)
}

func scanPorts() []PortInfo {
	var portmap []PortInfo

	// Get listening ports
	cmd := exec.Command("ss", "-tuln")
	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error running ss: %v", err)
		return portmap
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 4 {
			proto := parts[0]
			localAddr := parts[3]
			
			port := extractPort(localAddr)
			if port == 0 {
				continue
			}

			// Try to get process info
			process, pid, command := getProcessInfo(port)

			portmap = append(portmap, PortInfo{
				Port:     port,
				Protocol: proto,
				Process:  process,
				PID:      pid,
				Command:  command,
				Address:  localAddr,
			})
		}
	}

	return portmap
}

func extractPort(addr string) int {
	// Handle IPv6 addresses like [::1]:8080
	re := regexp.MustCompile(`\]:?(\d+)$`)
	match := re.FindStringSubmatch(addr)
	if len(match) > 1 {
		port, _ := strconv.Atoi(match[1])
		return port
	}

	// Handle IPv4 addresses like 0.0.0.0:8080
	parts := strings.Split(addr, ":")
	if len(parts) > 0 {
		port, err := strconv.Atoi(parts[len(parts)-1])
		if err == nil {
			return port
		}
	}

	return 0
}

func getProcessInfo(port int) (string, string, string) {
	// Use lsof to find process using the port
	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t", "-P", "-n")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", ""
	}

	pid := strings.TrimSpace(string(output))
	if pid == "" {
		return "", "", ""
	}

	// Get process name
	cmd = exec.Command("ps", "-p", pid, "-o", "comm=", "-o", "args=")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return pid, pid, ""
	}

	parts := strings.Fields(string(output))
	if len(parts) >= 2 {
		return parts[0], pid, strings.Join(parts[1:], " ")
	}

	return pid, pid, ""
}

func displayMap(portmap []PortInfo) {
	if len(portmap) == 0 {
		color.Yellow("No listening ports found")
		return
	}

	fmt.Println(color.CyanString("\n=== PORT MAP ==="))
	fmt.Println()

	sort.Slice(portmap, func(i, j int) bool {
		return portmap[i].Port < portmap[j].Port
	})

	for _, p := range portmap {
		portColor := getColorForPort(p.Port)
		fmt.Printf("%-6s %-10s %-15s %s\n",
			portColor(fmt.Sprintf("%d", p.Port)),
			p.Protocol,
			p.Address,
			color.HiWhiteString(p.Process),
		)
		
		if p.Command != "" {
			fmt.Printf("    PID: %s | %s\n", color.HiYellowString(p.PID), p.Command)
		}
		fmt.Println()
	}

	fmt.Println(color.YellowString("\nCommands:"))
	fmt.Println("  lsof -i :<port>  # Find more details")
	fmt.Println("  kill <pid>       # Kill process")
	fmt.Println("  ss -tuln         # List all ports")
}

func getColorForPort(port int) func(string) string {
	if port < 1024 {
		return func(s string) string { return color.HiRedString(s) }
	} else if port < 8000 {
		return func(s string) string { return color.HiYellowString(s) }
	}
	return func(s string) string { return color.HiGreenString(s) }
}