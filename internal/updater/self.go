package updater

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"syscall"
)

// CheckAndApply verifies if the binary on disk matches the one in the repo (or a new download)
// For this scaffolding, we simulate "new version" by checking a file in the git repo: bin/pivium
func CheckAndApply(currentPath string, repoBinaryPath string) error {
	// 1. Calculate Hashes
	currentHash, err := getFileHash(currentPath)
	if err != nil {
		return err
	}

	// Check if update source exists
	if _, err := os.Stat(repoBinaryPath); os.IsNotExist(err) {
		return nil // No update available locally
	}

	newHash, err := getFileHash(repoBinaryPath)
	if err != nil {
		return err
	}

	if currentHash == newHash {
		return nil // Already up to date
	}

	fmt.Println(">> New version detected. Updating...")

	// 2. Prepare new binary (Atomic Swap Pattern)
	tmpPath := currentPath + ".new"
	
	// Copy repo binary to tmp location
	if err := copyFile(repoBinaryPath, tmpPath); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, 0755); err != nil {
		return err
	}

	// 3. Rename (Atomic)
	if err := os.Rename(tmpPath, currentPath); err != nil {
		return err
	}

	// 4. Restart Process
	fmt.Println(">> Restarting into new binary...")
	execErr := syscall.Exec(currentPath, os.Args, os.Environ())
	if execErr != nil {
		return execErr
	}

	return nil
}

func getFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil { return err }
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil { return err }
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}