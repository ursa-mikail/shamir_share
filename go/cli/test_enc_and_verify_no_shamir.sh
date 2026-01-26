#!/bin/zsh

echo "===== Encryption and Verification Script ====="
echo

# Step 1: Create plaintext file with timestamp
echo "[Step 1] Creating plaintext file with timestamp..."
echo $(date '+%Y-%m-%d %H:%M:%S') >> plaintext.txt
echo "Plaintext created."
echo

# Step 2: Generate random key and IV
echo "[Step 2] Generating encryption key and IV..."
openssl rand -hex 32 > mykey.hex
openssl rand -hex 16 > myiv.hex

KEY=$(cat mykey.hex)
IV=$(cat myiv.hex)

echo "Key: $KEY"
echo "IV: $IV"
echo

# Step 3: Encrypt the file
echo "[Step 3] Encrypting plaintext.txt..."
openssl enc -aes-256-cbc -K "$KEY" -iv "$IV" \
  -in plaintext.txt \
  -out encrypted.bin

if [ $? -ne 0 ]; then
    echo "ERROR: Encryption failed!"
    exit 1
fi
echo "Encryption successful: encrypted.bin created."
echo

# Step 4: Save key and IV to secret.json
echo "[Step 4] Saving key and IV to secret.json..."
cat > secret.json <<EOF
{
  "key": "$KEY",
  "iv": "$IV"
}
EOF
echo "secret.json created."
echo

# Step 5: Extract key and IV from secret.json
echo "[Step 5] Extracting key and IV from secret.json..."
EXTRACTED_KEY=$(grep '"key"' secret.json | sed 's/.*"key": "\(.*\)".*/\1/')
EXTRACTED_IV=$(grep '"iv"' secret.json | sed 's/.*"iv": "\(.*\)".*/\1/')

echo "Extracted Key: $EXTRACTED_KEY"
echo "Extracted IV: $EXTRACTED_IV"
echo

# Step 6: Verify extraction matches original
echo "[Step 6] Verifying extracted values match original..."
if [ "$KEY" = "$EXTRACTED_KEY" ]; then
    echo "✓ Key matches!"
else
    echo "✗ Key MISMATCH!"
    echo "Original: $KEY"
    echo "Extracted: $EXTRACTED_KEY"
fi

if [ "$IV" = "$EXTRACTED_IV" ]; then
    echo "✓ IV matches!"
else
    echo "✗ IV MISMATCH!"
    echo "Original: $IV"
    echo "Extracted: $EXTRACTED_IV"
fi
echo

# Step 7: Decrypt using extracted values
echo "[Step 7] Decrypting using extracted key and IV..."
openssl enc -aes-256-cbc -d -K "$EXTRACTED_KEY" -iv "$EXTRACTED_IV" \
  -in encrypted.bin \
  -out decrypted.txt

if [ $? -ne 0 ]; then
    echo "ERROR: Decryption failed!"
    exit 1
fi
echo "Decryption successful: decrypted.txt created."
echo

# Step 8: Compare original and decrypted files
echo "[Step 8] Comparing original and decrypted files..."
if diff plaintext.txt decrypted.txt > /dev/null; then
    echo "✓✓✓ SUCCESS! Files match perfectly!"
    echo "The encryption and decryption process is verified."
else
    echo "✗✗✗ FAILURE! Files do not match!"
    echo "The encryption process has an error."
fi
echo

# Display file contents for verification
echo "[Step 9] Displaying file contents..."
echo
echo "--- Original plaintext.txt ---"
cat plaintext.txt
echo
echo "--- Decrypted decrypted.txt ---"
cat decrypted.txt
echo

echo "===== Process Complete ====="



echo ""
: <<'NOTE_BLOCK'

Creates plaintext with timestamp using date command
Generates random key and IV with OpenSSL
Encrypts the file using AES-256-CBC
Saves credentials to secret.json
Extracts key and IV from secret.json using grep and sed
Verifies the extracted values match the originals
Decrypts using the extracted credentials
Compares original and decrypted files with diff
Displays both file contents for visual verification

NOTE_BLOCK
echo "" 
