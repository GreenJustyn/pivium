package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"pivium/internal/ceph"
	"pivium/internal/config"
	"pivium/internal/pkgs"
	"pivium/internal/proxmox"
	"pivium/internal/updater"
)

// Hardcoded for bootstrapping, in prod pass via -ldflags
var Version = "0.1.0"

func main() {
	rootDir := flag.String("root", "/opt/pivium", "Root directory of the repo")
	mode := flag.String("mode", "reconcile", "Operation mode")
	flag.Parse()

	hostname, _ := os.Hostname()
	log.Printf("Starting pivium v%s on %s", Version, hostname)

	// 1. Self-Update Check
	// Assuming the git repo has a compiled binary in bin/pivium
	executable, _ := os.Executable()
	err := updater.CheckAndApply(executable, *rootDir+"/bin/pivium")
	if err != nil {
		log.Printf("Update check failed: %v", err)
	}

	if *mode == "reconcile" {
		reconcile(*rootDir, hostname)
	}
}

func reconcile(rootDir, hostname string) {
	// 1. Pull latest Git changes
	fmt.Println(">> Pulling latest configuration...")
	gitCmd := exec.Command("git", "-C", rootDir, "pull")
	if out, err := gitCmd.CombinedOutput(); err != nil {
		log.Printf("Git pull warning: %v (Output: %s)", err, string(out))
		// Continue anyway, maybe we are offline but have local config
	}

	// 2. Load Config
	cfg, err := config.Load(rootDir, hostname)
	if err != nil {
		log.Fatalf("Critical: Failed to load config: %v", err)
	}

	fmt.Printf("   Active Role: Proxmox=%v, Ceph=%v\n", cfg.Proxmox.Enabled, cfg.Ceph.Enabled)

	// 3. Package Management
	if err := pkgs.Ensure(cfg.System.Packages); err != nil {
		log.Printf("Error ensuring packages: %v", err)
	}

	// 4. Proxmox Reconciliation
	if err := proxmox.Reconcile(*cfg); err != nil {
		log.Printf("Error reconciling Proxmox: %v", err)
	}

	// 5. Ceph Reconciliation
	if err := ceph.Reconcile(*cfg); err != nil {
		log.Printf("Error reconciling Ceph: %v", err)
	}

	log.Println(">> Reconciliation Complete.")
}