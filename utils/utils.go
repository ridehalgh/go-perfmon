package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var LogFile *os.File

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

var initOnce sync.Once
var isDebug bool

func InitLog(logPath string) error {
	var err error
	initOnce.Do(func() {
		if len(os.Getenv("DEBUG")) > 0 {
			isDebug = true
		}

		logDir := filepath.Dir(logPath)
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			err = os.MkdirAll(logDir, os.ModePerm)
			if err != nil {
				log.Fatalf("Failed to create log directory: %s: %w", logDir, err)
				return
			}
		}

		LogFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			err = fmt.Errorf("failed to open log file: %s: %w", logPath, err)
			return
		}

		logFlags := log.Ldate | log.Ltime
		var defaultWriter io.Writer = LogFile

		DebugLogger = log.New(defaultWriter, "[DEBUG] ", logFlags)
		InfoLogger = log.New(defaultWriter, "[INFO] ", logFlags)
		WarnLogger = log.New(defaultWriter, "[WARN] ", logFlags)
		ErrorLogger = log.New(defaultWriter, "[ERROR] ", logFlags)

		InfoLogger.Println("Logging initialized...")
		if isDebug {
			DebugLogger.Println("Debug logging enabled...")
		}
	})

	return err
}

func CloseLog() {
	if LogFile != nil {
		InfoLogger.Println("Closing log file...")
		err := LogFile.Close()
		if err != nil {
			// Fallback
			log.Printf("Error closing log file: %v\n", err)
		}
	}
}

const (
	B  uint64 = 1
	KB        = 1024 * B
	MB        = 1024 * KB
	GB        = 1024 * MB
	TB        = 1024 * GB
	PB        = 1024 * TB
)

type Unit int

const (
	UnitBytes Unit = iota
	UnitKilobytes
	UnitMegabytes
	UnitGigabytes
	UnitTerabytes
	UnitPetabytes
)

func convertBytesToUnit(bytes uint64, unit Unit) (float64, string) {
	switch unit {
	case UnitBytes:
		return float64(bytes), "B"
	case UnitKilobytes:
		return float64(bytes) / float64(KB), "KB"
	case UnitMegabytes:
		return float64(bytes) / float64(MB), "MB"
	case UnitGigabytes:
		return float64(bytes) / float64(GB), "GB"
	case UnitTerabytes:
		return float64(bytes) / float64(TB), "TB"
	case UnitPetabytes:
		return float64(bytes) / float64(PB), "PB"
	default:
		return float64(bytes), "<unknown>"
	}
}

func formatBytesString(bytes uint64, unit Unit) string {
	val, unitStr := convertBytesToUnit(bytes, unit)
	if unit == UnitBytes {
		return fmt.Sprintf("%.0f %s", val, unitStr)
	} else {
		return fmt.Sprintf("%.2f %s", val, unitStr)
	}
}

func FormatBytesAuto(bytes uint64) string {
	var unit Unit
	switch {
	case bytes >= PB:
		unit = UnitPetabytes
	case bytes >= TB:
		unit = UnitTerabytes
	case bytes >= GB:
		unit = UnitGigabytes
	case bytes >= MB:
		unit = UnitMegabytes
	case bytes >= KB:
		unit = UnitKilobytes
	default:
		unit = UnitBytes
	}
	return formatBytesString(bytes, unit)
}
