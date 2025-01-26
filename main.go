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
	VecStartingSize int     `json:"vec_starting_size"`
	VecGrowthFactor float64 `json:"vec_growth_factor"`
	DevMode         bool    `json:"dev_mode"`
	Debug           bool    `json:"debug"`
}

func defaultConfig() Config {
	return Config{
		VecStartingSize: 0,
		VecGrowthFactor: 1.5,
		DevMode:         false,
		Debug:           false,
	}
}

func findConfigFile() string {
	files, err := filepath.Glob("*.ndl")
	if err != nil || len(files) == 0 {
		return ""
	}
	return files[0]
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
	configFile := findConfigFile()
	config := parseConfig(configFile)
	jsonOutput, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonOutput))
}
