package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const targetOrigin = "sync_target"

type FromToConfig struct {
	FromOrigin string   `json:"fromOrigin"`
	ToOrigin   string   `json:"toOrigin"`
	Branches   []string `json:"branches"`
  Force      bool     `json:"force"`
}

type SyncConfig struct {
	TempDirectory string         `json:"tempDirectory"`
	FromToConfigs []FromToConfig `json:"fromToConfigs"`
}

func runCommand(dir string, command string, args ...string) {
  fmt.Println("Running command:", command, args)
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
}

func cloneRepo(dir string, origin string) {
	runCommand(dir, "git", "clone", origin)
}

func addRemote(dir string, origin string) {
	runCommand(dir, "git", "remote", "add", targetOrigin, origin)
}

func checkoutBranch(dir string, branch string) {
	runCommand(dir, "git", "checkout", branch)
  runCommand(dir, "git", "pull")
}

func syncBranch(dir string, branch string, force bool) {
  args := []string{"push"}
  if force {
    args = append(args, "-f")
  }
  args = append(args, targetOrigin, branch)
  runCommand(dir, "git", args...)
}

func getDirectoryFromOrigin(origin string) string {
	dir := origin[strings.LastIndex(origin, "/")+1:]
	if strings.HasSuffix(dir, ".git") {
		dir = dir[:len(dir)-4]
	}
	return dir
}

func runSync(tmpDir string, fromToConfig FromToConfig) {
	fromOrigin := fromToConfig.FromOrigin
	toOrigin := fromToConfig.ToOrigin
	branches := fromToConfig.Branches
  force := fromToConfig.Force
	fmt.Println("Syncing from", fromOrigin, "to", toOrigin, "branches", branches)
	fmt.Println("Temp dir:", tmpDir)
	cloneRepo(tmpDir, fromOrigin)
	projectDir := getDirectoryFromOrigin(fromOrigin)
	fullProjectDir := tmpDir + "/" + projectDir
	addRemote(fullProjectDir, toOrigin)
	for _, branch := range branches {
		checkoutBranch(fullProjectDir, branch)
		syncBranch(fullProjectDir, branch, force)
	}
}

func readConfigFromFile(filePath string) SyncConfig {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := SyncConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	return config
}

func prepareTempDirectory(config SyncConfig) string {
	if config.TempDirectory != "" {
    os.RemoveAll(config.TempDirectory)
		os.MkdirAll(config.TempDirectory, 0755)
		return config.TempDirectory
	}
	tmpDir := os.TempDir() + "/git-sync"
  if _, err := os.Stat(tmpDir); err == nil {
    os.RemoveAll(tmpDir)
  }
	os.MkdirAll(tmpDir, 0755)
	return tmpDir
}

func cleanupTempDirectory(tmpDir string) {
	os.RemoveAll(tmpDir)
}

func syncAll(config SyncConfig) {
	tmpDir := prepareTempDirectory(config)
	for _, fromToConfig := range config.FromToConfigs {
		runSync(tmpDir, fromToConfig)
	}
	cleanupTempDirectory(tmpDir)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: git-sync <config-file>")
		os.Exit(1)
	}
	configFilePath := os.Args[1]
	config := readConfigFromFile(configFilePath)
	syncAll(config)

}
