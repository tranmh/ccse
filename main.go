package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type ClaudeManager struct {
	cmd         *exec.Cmd
	stdin       *bufio.Writer
	stdout      *bufio.Scanner
	stderr      *bufio.Scanner
	printOutput bool
}

func printBanner() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Claude Code Session Extender                      ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════════╣")
	fmt.Println("║                                                                      ║")
	fmt.Println("║  Purpose: Automatically extend Claude Code usage beyond 5-hour       ║")
	fmt.Println("║           session limits without additional costs                    ║")
	fmt.Println("║                                                                      ║")
	fmt.Println("║  Strategy Explanation:                                               ║")
	fmt.Println("║  Claude Code has a 5-hour session limit. This tool starts sessions   ║")
	fmt.Println("║  3 hours BEFORE you typically begin work, ensuring fresh 5-hour      ║")
	fmt.Println("║  sessions are always available when needed.                          ║")
	fmt.Println("║                                                                      ║")
	fmt.Println("║  Example: If you work exactly 5 hours from 7:00 AM - 12:00 PM        ║")
	fmt.Println("║  → Session #1 starts at 4:02 AM (fresh 5-hour limit)                 ║")
	fmt.Println("║  → Session #2 starts at 9:02 AM (overlapping fresh limit)            ║")
	fmt.Println("║  → Result: You get 2 full sessions instead of just 1!                ║")
	fmt.Println("║  → This effectively doubles your available Claude Code time          ║")
	fmt.Println("║                                                                      ║")
	fmt.Println("║  Key Benefit: Multiple overlapping sessions = Extended usage         ║")
	fmt.Println("║                                                                      ║")
	fmt.Println("║  💡 CUSTOMIZATION TIP:                                               ║")
	fmt.Println("║  Adjust the schedule based on YOUR working hours for maximum         ║")
	fmt.Println("║  sessions. Start sessions 3 hours before you begin work.             ║")
	fmt.Println("║                                                                      ║")
	fmt.Println("║  Examples for different work schedules:                              ║")
	fmt.Println("║  • Work 9 AM - 2 PM → Use: \"2 6,11,16,21,2 * * *\"                    ║")
	fmt.Println("║  • Work 10 AM - 3 PM → Use: \"2 7,12,17,22,3 * * *\"                   ║")
	fmt.Println("║  • Work 6 AM - 11 AM → Use: \"2 3,8,13,18,23 * * *\"                   ║")
	fmt.Println("║                                                                      ║")
	fmt.Println("║  Default schedule: 4:02, 9:02, 14:02, 19:02, 0:02 (for 7AM start)    ║")
	fmt.Println("║                                                                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func main() {
	printBanner()

	var cronExpr string

	runAtStart := flag.Bool("run-at-start", false, "Send message immediately at startup")
	printResponses := flag.Bool("print-responses", true, "Print Claude's responses")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		// Default schedule: 4:02 AM and every 5 hours after (4:02, 9:02, 14:02, 19:02, 0:02)
		cronExpr = "2 4,9,14,19,0 * * *"
		fmt.Println("Using default schedule: 4:02 AM and every 5 hours after")
		fmt.Println("Times: 4:02, 9:02, 14:02, 19:02, 0:02")
	} else {
		cronExpr = args[0]
		fmt.Printf("Using custom schedule: %s\n", cronExpr)
	}

	c := cron.New()

	_, err := c.AddFunc(cronExpr, func() {
		runClaudeSession(*printResponses)
	})

	if err != nil {
		log.Fatalf("Invalid cron expression: %v", err)
	}

	fmt.Printf("Starting Claude Code automation with schedule: %s\n", cronExpr)
	if *runAtStart {
		fmt.Println("Will send initial message at startup")
	} else {
		fmt.Println("Will wait for first scheduled time")
	}
	fmt.Println("Press Ctrl+C to stop")

	// Run immediately once only if flag is set
	if *runAtStart {
		runClaudeSession(*printResponses)
	}

	c.Start()

	// Keep the program running
	select {}
}

func runClaudeSession(printResponses bool) {
	fmt.Printf("[%s] Starting new Claude Code session...\n", time.Now().Format("2006-01-02 15:04:05"))

	manager, err := startClaude(printResponses)
	if err != nil {
		fmt.Printf("Error starting Claude: %v\n", err)
		return
	}

	// Send "hi" message
	err = manager.sendMessage("hi")
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	} else {
		fmt.Println("Message 'hi' sent successfully")
	}

	// Wait and read response
	if printResponses {
		manager.readResponse()
	} else {
		time.Sleep(3 * time.Second)
	}

	// Stop Claude
	err = manager.stop()
	if err != nil {
		fmt.Printf("Error stopping Claude: %v\n", err)
	} else {
		fmt.Println("Claude Code session ended")
	}

	fmt.Printf("[%s] Session completed\n\n", time.Now().Format("2006-01-02 15:04:05"))
}

func startClaude(printOutput bool) (*ClaudeManager, error) {
	// Start Claude Code process
	cmd := exec.Command("claude", "code")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start claude code: %v", err)
	}

	manager := &ClaudeManager{
		cmd:         cmd,
		stdin:       bufio.NewWriter(stdin),
		stdout:      bufio.NewScanner(stdout),
		stderr:      bufio.NewScanner(stderr),
		printOutput: printOutput,
	}

	// Wait for Claude to initialize
	time.Sleep(3 * time.Second)

	return manager, nil
}

func (cm *ClaudeManager) sendMessage(message string) error {
	_, err := cm.stdin.WriteString(message + "\n")
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = cm.stdin.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush message: %v", err)
	}

	return nil
}

func (cm *ClaudeManager) readResponse() {
	if !cm.printOutput {
		time.Sleep(5 * time.Second)
		return
	}

	fmt.Println("Claude's response:")
	fmt.Println("---")

	// Read for a longer time to capture response
	timeout := time.After(8 * time.Second)
	responseReceived := false

	// Create channels for stdout and stderr
	stdoutChan := make(chan string, 10)
	stderrChan := make(chan string, 10)

	// Goroutine to read stdout
	go func() {
		for {
			if cm.stdout.Scan() {
				line := strings.TrimSpace(cm.stdout.Text())
				if line != "" {
					stdoutChan <- line
				}
			}
		}
	}()

	// Goroutine to read stderr
	go func() {
		for {
			if cm.stderr.Scan() {
				line := strings.TrimSpace(cm.stderr.Text())
				if line != "" {
					stderrChan <- line
				}
			}
		}
	}()

	// Read responses
	for {
		select {
		case <-timeout:
			fmt.Println("---")
			if !responseReceived {
				fmt.Println("No response received (timeout after 8s)")
			}
			return
		case line := <-stdoutChan:
			fmt.Printf("[stdout] %s\n", line)
			responseReceived = true
		case line := <-stderrChan:
			fmt.Printf("[stderr] %s\n", line)
			responseReceived = true
		}
	}
}

func (cm *ClaudeManager) stop() error {
	// Send exit command
	cm.stdin.WriteString("/exit\n")
	cm.stdin.Flush()

	// Give it time to exit gracefully
	time.Sleep(1 * time.Second)

	// Force kill if still running
	if cm.cmd.Process != nil {
		err := cm.cmd.Process.Kill()
		if err != nil && !strings.Contains(err.Error(), "process already finished") {
			return fmt.Errorf("failed to kill process: %v", err)
		}
	}

	return nil
}