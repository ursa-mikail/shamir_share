package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/hashicorp/vault/shamir"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Directory   string
	TotalShares int
	Threshold   int
}

type KeyValues struct {
	KeyName  string
	KeyType  string
	KeyBytes []byte
}

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

func LoadConfig() (AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return AppConfig{}, err
	}

	totalShares, err := strconv.Atoi(os.Getenv("TOTAL_SHARES"))
	if err != nil {
		return AppConfig{}, err
	}

	threshold, err := strconv.Atoi(os.Getenv("THRESHOLD"))
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		Directory:   os.Getenv("DIRECTORY"),
		TotalShares: totalShares,
		Threshold:   threshold,
	}, nil
}

func CreateKeyShares(keyVals []KeyValues, dest string, totalShares int, threshold int) error {
	// Prepare empty list for each share
	shareGroups := make([][]KeyShare, totalShares)
	for _, kv := range keyVals {
		shares, err := shamir.Split(kv.KeyBytes, totalShares, threshold)
		if err != nil {
			return err
		}
		for i := 0; i < totalShares; i++ {
			shareGroups[i] = append(shareGroups[i], KeyShare{
				KeyName: kv.KeyName,
				KeyType: kv.KeyType,
				KeyShareVers: []KeyShareVers{
					{
						Version: 1,
						Share:   hex.EncodeToString(shares[i]),
					},
				},
			})
		}
	}

	// Save each share to its own file
	for i := 0; i < totalShares; i++ {
		filename := fmt.Sprintf("key_shares_%02d.json", i+1)
		fullPath := filepath.Join(dest, filename)
		allKeys := AllKeys{
			ShareNum: i + 1,
			Keys:     shareGroups[i],
		}
		if err := SaveKeyShareToJsonFile(fullPath, allKeys); err != nil {
			return err
		}
	}
	return nil
}

func SaveKeyShareToJsonFile(path string, keys AllKeys) error {
	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func RecoverKeyFromShares(source string, threshold int) (map[string][]byte, error) {
	files, err := filepath.Glob(filepath.Join(source, "key_shares_*.json"))
	if err != nil {
		return nil, err
	}
	if len(files) < threshold {
		return nil, fmt.Errorf("not enough shares: have %d, need %d", len(files), threshold)
	}
	sort.Strings(files)
	files = files[:threshold]

	shareMap := make(map[string][][]byte)

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		var ak AllKeys
		if err := json.Unmarshal(data, &ak); err != nil {
			return nil, err
		}

		for _, key := range ak.Keys {
			decoded, err := hex.DecodeString(key.KeyShareVers[0].Share)
			if err != nil {
				return nil, err
			}
			shareMap[key.KeyName] = append(shareMap[key.KeyName], decoded)
		}
	}

	result := make(map[string][]byte)
	for keyName, shares := range shareMap {
		combined, err := shamir.Combine(shares)
		if err != nil {
			return nil, err
		}
		result[keyName] = combined
	}
	return result, nil
}

// writeKeyToFile saves the key to a file in hex format under ./keys/
func writeKeyToFile(fullPath string, kv KeyValues) error {
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write hex-encoded key
	hexKey := hex.EncodeToString(kv.KeyBytes)
	_, err = f.WriteString(hexKey + "\n")
	return err
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	os.MkdirAll(cfg.Directory, fs.ModePerm)

	// Generate keys
	var keys []KeyValues
	for i := 0; i < 3; i++ {
		keyBytes := make([]byte, 32)
		if _, err := rand.Read(keyBytes); err != nil {
			panic(err)
		}
		keyName := fmt.Sprintf("key_%02d.key", i)
		keys = append(keys, KeyValues{
			KeyName:  keyName,
			KeyType:  "oct",
			KeyBytes: keyBytes,
		})
	}

	// Print and write each key
	for _, key := range keys {
		hexKey := hex.EncodeToString(key.KeyBytes)
		fullPath := filepath.Join("keys", key.KeyName)
		fmt.Printf("%s %s\n", hexKey, fullPath)

		if err := writeKeyToFile(fullPath, key); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", fullPath, err)
		}
	}

	err = CreateKeyShares(keys, cfg.Directory, cfg.TotalShares, cfg.Threshold)
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ Key shares saved as key_shares_0*.json")

	recovered, err := RecoverKeyFromShares(cfg.Directory, cfg.Threshold)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n✅ Recovered Keys:")
	for k, v := range recovered {
		fmt.Printf("%s %s\n", hex.EncodeToString(v), k)
	}
}

/*
go mod tidy
go run main.go

% go run main.go
f3f4a6d9b8486db95adbc2a50ea64f876fb128711f97347ba1b58dca14e59177 keys/key_00.key
dc958e7b38a842ac776db07b824901525c64a0628db0472401d2a8d268ba1491 keys/key_01.key
371cdc6028d7c849ae185249e7528ae87e4559d162d2f2144ec3a22e37c062aa keys/key_02.key
✅ Key shares saved as key_shares_0*.json

✅ Recovered Keys:
f3f4a6d9b8486db95adbc2a50ea64f876fb128711f97347ba1b58dca14e59177 key_00.key
dc958e7b38a842ac776db07b824901525c64a0628db0472401d2a8d268ba1491 key_01.key
371cdc6028d7c849ae185249e7528ae87e4559d162d2f2144ec3a22e37c062aa key_02.key

*/
