package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	// Vec
	VecStartingSize int     `json:"vec_starting_size"`
	VecGrowthFactor float64 `json:"vec_growth_factor"`
	// Debug & Dev
	DevMode bool `json:"dev_mode"`
	Debug   bool `json:"debug"`
	// Output
	OutputFile string `json:"output_file"`
}

func defaultConfig() Config {
	return Config{
		VecStartingSize: 0,
		VecGrowthFactor: 1.5,
		DevMode:         false,
		Debug:           false,
	}
}

func parseConfig(filename string) Config {
	if filename == "" {
		return defaultConfig()
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return defaultConfig()
	}
	defer file.Close()

	config := defaultConfig()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip comments and empty lines
		}
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
			parts := strings.SplitN(line, "]", 2)
			key := strings.TrimSpace(parts[0][1:])
			value := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))

			switch key {
			case "vec_starting_size":
				fmt.Sscanf(value, "%d", &config.VecStartingSize)
			case "vec_growth_factor":
				fmt.Sscanf(value, "%v", &config.VecGrowthFactor)
			case "output":
				fmt.Sscanf(value, "%s", &config.OutputFile)
			case "dev_mode":
				config.DevMode = strings.ToLower(value) == "true"
			case "debug":
				config.Debug = strings.ToLower(value) == "true"
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return config
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <config-directory>")
		os.Exit(1)
	}

	configDir := os.Args[1]
	files, err := filepath.Glob(filepath.Join(configDir, "*.ndl"))
	if err != nil || len(files) == 0 {
		fmt.Println("No config file found in the specified directory")
		os.Exit(1)
	}

	configFile := files[0]
	config := parseConfig(configFile)
	jsonOutput, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonOutput))
}
