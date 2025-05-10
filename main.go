package main

import (
	"log"

	"github.com/ridehalgh/go-perfmon/tui"
	"github.com/ridehalgh/go-perfmon/utils"
)

// func getMemUsage() (float64, error) {

// }

// func showProcessInfo() {
// 	// --- Configuration ---
// 	// refreshInterval defines how often (in seconds) the stats are updated.
// 	refreshInterval := 2 * time.Second

// 	// --- Setup Signal Handling for Graceful Shutdown ---
// 	// Create a channel to listen for OS signals.
// 	sigs := make(chan os.Signal, 1)
// 	// Notify this channel for Interrupt (Ctrl+C) and SIGTERM signals.
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

// 	// --- Main Monitoring Loop ---
// 	// This loop will run until a shutdown signal is received.
// 	time.Sleep(1 * time.Second) // Give user a moment to read

// 	// Create a ticker that fires at the refreshInterval.
// 	ticker := time.NewTicker(refreshInterval)
// 	defer ticker.Stop() // Ensure the ticker is stopped when main exits.

// 	for {
// 		select {
// 		case <-ticker.C: // Triggered at each refreshInterval
// 			// Get CPU Usage
// 			// Percent calculates the CPU usage percentage.
// 			// The first argument is the interval over which to calculate usage.
// 			// The second argument (percpu) is false, meaning we want total CPU usage, not per core.

// 			// Get Memory Usage
// 			// VirtualMemory returns statistics about system memory usage.
// 			vmStat, err := mem.VirtualMemory()
// 			if err != nil {
// 				log.Printf("Error getting memory usage: %v\n", err)
// 				continue // Skip this iteration if there's an error
// 			}

// 			// Get Process Information
// 			// Get the current process ID
// 			pid := os.Getpid()
// 			// Get the process object for the current process
// 			p, err := process.NewProcess(int32(pid))
// 			if err != nil {
// 				log.Printf("Error getting process info: %v\n", err)
// 				continue // Skip this iteration if there's an error
// 			}
// 			currMem, err := p.MemoryInfo()
// 			if err != nil {
// 				log.Printf("Error getting current process memory info: %v\n", err)
// 				continue // Skip this iteration if there's an error
// 			}

// 			fmt.Printf("Current Process Information: %v \n", formatBytesAuto(currMem.RSS))

// 			// --- Display Information ---
// 			fmt.Println("--- System Performance Monitor ---")
// 			fmt.Printf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05"))
// 			fmt.Println("----------------------------------")

// 			// CPU Usage
// 			//if len(cpuPercentages) > 0 {
// 			//		fmt.Printf("CPU Usage: %.2f%%\n", cpuPercentages[0])
// 			//		} else {
// 			//				fmt.Println("CPU Usage: N/A")
// 			//			}

// 			// Memory Usage
// 			fmt.Printf("Memory Usage:\n")
			// fmt.Printf("  Total:     %v \n", formatBytesAuto(vmStat.Total))
			// fmt.Printf("  Available: %v \n", formatBytesAuto(vmStat.Available))
			// fmt.Printf("  Used:      %v (%.2f%%)\n", formatBytesAuto(vmStat.Used), vmStat.UsedPercent)
			// fmt.Printf("  Free:      %v \n", formatBytesAuto(vmStat.Free)) // Free is not always the same as Available

// 			fmt.Println("----------------------------------")
// 			fmt.Println("Press Ctrl+C to exit.")

// 		case <-sigs: // Triggered when Ctrl+C or SIGTERM is received
// 			fmt.Println("\nShutting down System Performance Monitor...")
// 			fmt.Println("Goodbye!")
// 			return // Exit the main function, and thus the program
// 		}
// 	}

// }

func main() {

	if err := utils.InitLog("perfmon.log"); err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}
	defer utils.CloseLog()

	tui.InitTui()
}
