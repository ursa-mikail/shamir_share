package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/vault/shamir"
)

type AllKeys struct {
	ShareNum int        `json:"share_num"`
	Keys     []KeyShare `json:"keys"`
}

type KeyShare struct {
	KeyName      string         `json:"key_name"`
	KeyType      string         `json:"key_type"`
	KeyShareVers []KeyShareVers `json:"key_share_vers"`
}

type KeyShareVers struct {
	Version int    `json:"version"`
	Share   string `json:"share"`
}

func main() {
	// Define subcommands
	splitCmd := flag.NewFlagSet("split", flag.ExitOnError)
	recoverCmd := flag.NewFlagSet("recover", flag.ExitOnError)
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)

	// Split command flags
	splitKey := splitCmd.String("key", "", "Hex-encoded key to split (required)")
	splitKeyFile := splitCmd.String("keyfile", "", "Path to file containing hex-encoded key")
	splitKeyName := splitCmd.String("name", "secret.key", "Name for the key")
	splitTotal := splitCmd.Int("total", 5, "Total number of shares to create")
	splitThreshold := splitCmd.Int("threshold", 3, "Minimum shares needed to reconstruct")
	splitDir := splitCmd.String("dir", "shares", "Directory to save share files")
	splitRaw := splitCmd.Bool("raw", false, "Treat input as raw bytes, not hex-encoded")

	// Recover command flags
	recoverDir := recoverCmd.String("dir", "shares", "Directory containing share files")
	recoverKeyName := recoverCmd.String("name", "", "Specific key name to recover (recovers all if empty)")
	recoverOutput := recoverCmd.String("output", "", "Output file for recovered key (stdout if empty)")
	recoverRaw := recoverCmd.Bool("raw", false, "Output raw bytes, not hex-encoded")

	// Generate command flags
	generateSize := generateCmd.Int("size", 32, "Key size in bytes (default: 32 for 256-bit)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "split":
		splitCmd.Parse(os.Args[2:])
		handleSplit(*splitKey, *splitKeyFile, *splitKeyName, *splitTotal, *splitThreshold, *splitDir, *splitRaw)
	case "recover":
		recoverCmd.Parse(os.Args[2:])
		handleRecover(*recoverDir, *recoverKeyName, *recoverOutput, *recoverRaw)
	case "generate":
		generateCmd.Parse(os.Args[2:])
		handleGenerate(*generateSize)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Shamir Secret Sharing CLI Tool

Usage:
  shamir-cli <command> [options]

Commands:
  generate    Generate a random key
  split       Split a key into shares
  recover     Recover a key from shares

Examples:

  # Generate a 256-bit (32-byte) random key
  shamir-cli generate
  shamir-cli generate -size 16  # Generate 128-bit key

  # Split a hex key into 8 shares, requiring 3 to reconstruct
  shamir-cli split -key fe79f5ea5b2df8da348489c39a23fdb05ed67ca49ed8a55e19afe1af4f17e2a1 \
    -total 8 -threshold 3 -dir ./my_shares -name master.key

  # Split multiple keys to the same share directory
  shamir-cli split -keyfile ./key1.txt -total 5 -threshold 3 -dir ./shares -name key1
  shamir-cli split -keyfile ./key2.txt -total 5 -threshold 3 -dir ./shares -name key2
  shamir-cli split -keyfile ./key3.txt -total 5 -threshold 3 -dir ./shares -name key3

  # Recover all keys from shares
  shamir-cli recover -dir ./shares

  # Recover specific key
  shamir-cli recover -dir ./shares -name key2 -output key2_recovered.txt

Options:
  Run 'shamir-cli <command> -h' for command-specific options`)
}

func handleGenerate(size int) {
	if size <= 0 {
		fmt.Fprintf(os.Stderr, "Error: Key size must be positive\n")
		os.Exit(1)
	}

	keyBytes := make([]byte, size)
	if _, err := rand.Read(keyBytes); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating random key: %v\n", err)
		os.Exit(1)
	}

	hexKey := hex.EncodeToString(keyBytes)
	fmt.Print(hexKey) // No newline
}

func handleSplit(keyHex, keyFile, keyName string, total, threshold int, dir string, raw bool) {
	// Validate parameters
	if total < 2 {
		fmt.Fprintf(os.Stderr, "Error: Total shares must be at least 2\n")
		os.Exit(1)
	}
	if threshold < 2 || threshold > total {
		fmt.Fprintf(os.Stderr, "Error: Threshold must be between 2 and total shares\n")
		os.Exit(1)
	}

	// Get key bytes
	var keyBytes []byte
	var err error

	if raw {
		// Read as raw bytes
		if keyFile != "" {
			keyBytes, err = ioutil.ReadFile(keyFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading key file: %v\n", err)
				os.Exit(1)
			}
		} else if keyHex != "" {
			keyBytes = []byte(keyHex)
		} else {
			fmt.Fprintf(os.Stderr, "Error: Must provide -key or -keyfile\n")
			os.Exit(1)
		}
	} else {
		// Read as hex-encoded
		if keyFile != "" {
			// Read from file
			data, err := ioutil.ReadFile(keyFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading key file: %v\n", err)
				os.Exit(1)
			}
			keyHex = string(data)
		}

		if keyHex == "" {
			fmt.Fprintf(os.Stderr, "Error: Must provide -key or -keyfile\n")
			os.Exit(1)
		}

		// Remove all whitespace including newlines, spaces, tabs
		keyHex = strings.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return -1
			}
			return r
		}, keyHex)
		
		keyBytes, err = hex.DecodeString(keyHex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding hex key: %v\n", err)
			os.Exit(1)
		}
	}

	// Create directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Split the key
	shares, err := shamir.Split(keyBytes, total, threshold)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error splitting key: %v\n", err)
		os.Exit(1)
	}

	// Process each share
	for i := 0; i < total; i++ {
		filename := fmt.Sprintf("key_shares_%02d.json", i+1)
		fullPath := filepath.Join(dir, filename)

		var allKeys AllKeys

		// Check if share file already exists
		if data, err := ioutil.ReadFile(fullPath); err == nil {
			// File exists, parse it
			if err := json.Unmarshal(data, &allKeys); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing existing share file %s: %v\n", filename, err)
				os.Exit(1)
			}

			// Verify share numbers match
			if allKeys.ShareNum != i+1 {
				fmt.Fprintf(os.Stderr, "Error: Share number mismatch in %s\n", filename)
				os.Exit(1)
			}

			// Check if key name already exists
			for _, existingKey := range allKeys.Keys {
				if existingKey.KeyName == keyName {
					fmt.Fprintf(os.Stderr, "Error: Key '%s' already exists in share files\n", keyName)
					os.Exit(1)
				}
			}
		} else {
			// File doesn't exist, create new structure with empty keys array
			allKeys = AllKeys{
				ShareNum: i + 1,
				Keys:     make([]KeyShare, 0),
			}
		}

		// Create the new key share entry
		newKeyShare := KeyShare{
			KeyName: keyName,
			KeyType: "oct",
			KeyShareVers: []KeyShareVers{
				{
					Version: 1,
					Share:   hex.EncodeToString(shares[i]),
				},
			},
		}

		// Add the new key share
		allKeys.Keys = append(allKeys.Keys, newKeyShare)

		// Write the updated share file
		data, err := json.MarshalIndent(allKeys, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(fullPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing share file: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("✅ Successfully split key '%s' into %d shares\n", keyName, total)
	fmt.Printf("✅ Threshold: %d shares required to reconstruct\n", threshold)
	fmt.Printf("✅ Shares saved to: %s/\n", dir)
}

func handleRecover(dir, keyName, outputFile string, raw bool) {
	// Find share files
	files, err := filepath.Glob(filepath.Join(dir, "key_shares_*.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding share files: %v\n", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No share files found in %s\n", dir)
		os.Exit(1)
	}

	sort.Strings(files)

	// Read all shares
	shareMap := make(map[string][][]byte)

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not read %s: %v\n", file, err)
			continue
		}

		var ak AllKeys
		if err := json.Unmarshal(data, &ak); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not parse %s: %v\n", file, err)
			continue
		}

		for _, key := range ak.Keys {
			// Filter by key name if specified
			if keyName != "" && key.KeyName != keyName {
				continue
			}

			if len(key.KeyShareVers) > 0 {
				decoded, err := hex.DecodeString(key.KeyShareVers[0].Share)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Could not decode share: %v\n", err)
					continue
				}
				shareMap[key.KeyName] = append(shareMap[key.KeyName], decoded)
			}
		}
	}

	if len(shareMap) == 0 {
		if keyName != "" {
			fmt.Fprintf(os.Stderr, "Error: No shares found for key '%s'\n", keyName)
		} else {
			fmt.Fprintf(os.Stderr, "Error: No valid shares found\n")
		}
		os.Exit(1)
	}

	// Recover keys
	recoveredCount := 0
	for kName, shares := range shareMap {
		recovered, err := shamir.Combine(shares)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error recovering key '%s': %v\n", kName, err)
			continue
		}

		if outputFile != "" {
			// Save to file
			if raw {
				// Write raw bytes
				if err := os.WriteFile(outputFile, recovered, 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
					os.Exit(1)
				}
			} else {
				// Write as hex (no newline)
				hexKey := hex.EncodeToString(recovered)
				if err := os.WriteFile(outputFile, []byte(hexKey), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
					os.Exit(1)
				}
			}
			fmt.Printf("✅ Recovered key '%s' saved to: %s\n", kName, outputFile)
		} else {
			// Output to stdout
			if raw {
				// Write raw bytes to stdout
				os.Stdout.Write(recovered)
			} else {
				// Write as hex
				hexKey := hex.EncodeToString(recovered)
				fmt.Printf("✅ Recovered key '%s':\n%s\n", kName, hexKey)
			}
		}
		recoveredCount++
	}

	if recoveredCount == 0 {
		fmt.Fprintf(os.Stderr, "Error: Failed to recover any keys\n")
		os.Exit(1)
	}
}