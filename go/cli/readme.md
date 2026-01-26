```
# Save the code as main.go
# Then initialize the module
go mod init shamir-cli
go get github.com/hashicorp/vault/shamir
go build -o shamir-cli
```

Generate a 256-bit key:
```
# Generate and display to terminal
./shamir-cli generate

# Generate and save to variable
KEY=$(./shamir-cli generate)
echo "Generated: $KEY"

# Generate with openssl and save to file
openssl rand -hex 32 > master.key

# Or generate with shamir-cli and save to file
./shamir-cli generate > master.key
```

Split a key into 8 shares (3 needed to reconstruct):
```
./shamir-cli split \
  -key fe79f5ea5b2df8da348489c39a23fdb05ed67ca49ed8a55e19afe1af4f17e2a1 \
  -total 8 \
  -threshold 3 \
  -dir ./shares_out \
  -name master.key

# Option 1: Use -keyfile to read from master.key
./shamir-cli split -keyfile master.key -total 8 -threshold 3 -dir ./shares_out -name master.key

# Option 2: Use the key from the file as a variable
KEY=$(cat master.key | tr -d '\n')
./shamir-cli split -key $KEY -total 8 -threshold 3 -dir ./shares_out -name master.key

# Then recover
./shamir-cli recover -dir ./shares_out -name master.key -output recovered.key

# Now they should match
if [ "$(cat master.key | tr -d '\n')" = "$(cat recovered.key | tr -d '\n')" ]; then
    echo "✅ Keys are equal"
else
    echo "❌ Keys are NOT equal"
fi
```

Recover specific key to a file:
```
./shamir-cli recover -dir ./shares_out -name master.key -output recovered.key
```

```
Without -name (recovers ALL keys in the shares):
./shamir-cli recover -dir ./shares_out -output recovered.key

This will recover every key found in the share files and output them all.

With -name (recovers ONLY the specified key):
bash./shamir-cli recover -dir ./shares_out -name master.key -output recovered.key

This filters to recover only the key named "master.key".

The -name parameter is useful when:
-You've split multiple different keys into the same share files
-You only want to recover one specific key, not all of them
```

```
# Step 1: Generate a 256-bit key and save it (with newline stripped)
./shamir-cli generate | tr -d '\n' > master.key
echo "Generated key:"
cat master.key
echo ""  # Add newline for readability

# Step 2: Split the key into 8 shares (3 needed to reconstruct)
./shamir-cli split -keyfile master.key -total 8 -threshold 3 -dir ./shares_out -name master.key

# Step 3: Verify the shares were created
ls -la ./shares_out/

# Step 4: Recover the key from shares
./shamir-cli recover -dir ./shares_out -name master.key -output recovered.key

# Step 5: Verify the keys match
echo "Original key:"
cat master.key
echo ""
echo "Recovered key:"
cat recovered.key

# Step 6: Compare them
if [ "$(openssl sha256 master.key | awk '{print $2}')" = "$(openssl sha256 recovered.key | awk '{print $2}')" ]; then
    echo "SHA256 hashes are equal"
else
    echo "SHA256 hashes are NOT equal"
fi
```

