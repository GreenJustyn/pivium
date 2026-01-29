package pkgs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Ensure takes a desired list of packages and installs what is missing
func Ensure(packages []string) error {
	var toInstall []string

	fmt.Println(">> Reconciling Packages...")
	
	// 1. Refresh package lists (optional, maybe run only if > 24h old in real prod)
	// exec.Command("apt-get", "update", "-q").Run()

	// 2. Check status
	for _, pkg := range packages {
		if !isInstalled(pkg) {
			fmt.Printf("   [+] Missing: %s\n", pkg)
			toInstall = append(toInstall, pkg)
		}
	}

	if len(toInstall) == 0 {
		fmt.Println("   All packages satisfied.")
		return nil
	}

	// 3. Install
	// DEBIAN_FRONTEND=noninteractive prevents prompts
	args := append([]string{"install", "-y", "--no-install-recommends"}, toInstall...)
	cmd := exec.Command("apt-get", args...)
	cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("apt install failed: %s\nOutput: %s", err, string(output))
	}

	fmt.Printf("   Successfully installed: %v\n", toInstall)
	return nil
}

// isInstalled uses dpkg-query to check if a package is cleanly installed
func isInstalled(pkg string) bool {
	// -W: show, -f: format status
	cmd := exec.Command("dpkg-query", "-W", "-f=${Status}", pkg)
	out, err := cmd.Output()
	if err != nil {
		// If command fails, package usually isn't installed
		return false
	}
	// We want "install ok installed"
	return strings.Contains(string(out), "install ok installed")
}