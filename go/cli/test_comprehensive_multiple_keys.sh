#!/bin/bash
set -e

echo "════════════════════════════════════════════════════════════════"
echo "  Shamir Secret Sharing - MULTI-KEY TEST"
echo "════════════════════════════════════════════════════════════════"

# Cleanup
rm -rf test_multi_shares
rm -f key1.txt key2.txt key3.txt secret.json
rm -f key1_recovered.txt key2_recovered.txt key3_recovered.txt secret_recovered.json
rm -f plaintext.txt encrypted.bin decrypted.txt

echo ""
echo "🔨 Building shamir-cli..."
go build -o shamir-cli main.go

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  TEST 1: Multiple Hex Keys to Same Share Directory"
echo "════════════════════════════════════════════════════════════════"

echo ""
echo "🔑 Step 1: Generate three different hex keys..."
openssl rand -hex 32 | tr -d '\n' > key1.txt
openssl rand -hex 32 | tr -d '\n' > key2.txt
openssl rand -hex 16 | tr -d '\n' > key3.txt

echo "Generated key1.txt (256-bit)"
echo "Generated key2.txt (256-bit)"
echo "Generated key3.txt (128-bit)"

echo ""
echo "🔪 Step 2: Split all three keys to the same share directory..."
echo "  Splitting key1..."
./shamir-cli split -keyfile key1.txt -total 5 -threshold 3 -dir ./test_multi_shares -name key1

echo "  Splitting key2..."
./shamir-cli split -keyfile key2.txt -total 5 -threshold 3 -dir ./test_multi_shares -name key2

echo "  Splitting key3..."
./shamir-cli split -keyfile key3.txt -total 5 -threshold 3 -dir ./test_multi_shares -name key3

echo ""
echo "📋 Step 3: Inspect share file structure..."
echo "Contents of key_shares_01.json:"
cat test_multi_shares/key_shares_01.json | jq '.keys | map(.key_name)'

echo ""
echo "🔥 Step 4: Delete original keys (simulate loss)..."
cp key1.txt key1_backup.txt
cp key2.txt key2_backup.txt
cp key3.txt key3_backup.txt
rm key1.txt key2.txt key3.txt
echo "❌ Deleted all original keys"

echo ""
echo "🔧 Step 5: Recover all keys from shares..."
./shamir-cli recover -dir ./test_multi_shares

echo ""
echo "🔧 Step 6: Recover individual keys to files..."
./shamir-cli recover -dir ./test_multi_shares -name key1 -output key1_recovered.txt
./shamir-cli recover -dir ./test_multi_shares -name key2 -output key2_recovered.txt
./shamir-cli recover -dir ./test_multi_shares -name key3 -output key3_recovered.txt

echo ""
echo "✅ Step 7: Verify all recoveries..."
echo "Key 1:"
if diff key1_backup.txt key1_recovered.txt; then
    echo "  ✅ KEY1 MATCHES!"
else
    echo "  ❌ KEY1 DIFFERS!"
    exit 1
fi

echo ""
echo "Key 2:"
if diff key2_backup.txt key2_recovered.txt; then
    echo "  ✅ KEY2 MATCHES!"
else
    echo "  ❌ KEY2 DIFFERS!"
    exit 1
fi

echo ""
echo "Key 3:"
if diff key3_backup.txt key3_recovered.txt; then
    echo "  ✅ KEY3 MATCHES!"
else
    echo "  ❌ KEY3 DIFFERS!"
    exit 1
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  TEST 2: Mixed Key Types (Hex + Raw JSON)"
echo "════════════════════════════════════════════════════════════════"

echo ""
echo "🔑 Step 1: Generate encryption key/IV as JSON..."
KEY=$(openssl rand -hex 32 | tr -d '\n')
IV=$(openssl rand -hex 16 | tr -d '\n')

cat > secret.json <<EOF
{
  "key": "$KEY",
  "iv": "$IV",
  "description": "Database encryption credentials"
}
EOF

echo "Created secret.json:"
cat secret.json

echo ""
echo "🔪 Step 2: Add JSON to existing shares (with hex keys)..."
rm -rf test_multi_shares
./shamir-cli split -keyfile key1_recovered.txt -total 5 -threshold 3 -dir ./test_multi_shares -name key1
./shamir-cli split -keyfile key2_recovered.txt -total 5 -threshold 3 -dir ./test_multi_shares -name key2
./shamir-cli split -keyfile secret.json -raw -total 5 -threshold 3 -dir ./test_multi_shares -name secret.json

echo ""
echo "📋 Step 3: Inspect mixed share structure..."
echo "Share file now contains:"
cat test_multi_shares/key_shares_01.json | jq '.keys | map(.key_name)'

echo ""
echo "🔥 Step 4: Delete originals..."
cp secret.json secret_backup.json
rm secret.json
echo "❌ Deleted secret.json"

echo ""
echo "🔧 Step 5: Recover JSON from shares with mixed keys..."
./shamir-cli recover -dir ./test_multi_shares -name secret.json -raw -output secret_recovered.json

echo ""
echo "✅ Step 6: Verify JSON recovery..."
if diff secret_backup.json secret_recovered.json; then
    echo "  ✅ JSON FILE MATCHES!"
else
    echo "  ❌ JSON FILE DIFFERS!"
    exit 1
fi

echo ""
echo "🔓 Step 7: Test encryption with recovered JSON credentials..."
KEY_RECOVERED=$(jq -r '.key' secret_recovered.json)
IV_RECOVERED=$(jq -r '.iv' secret_recovered.json)

echo "Recovered Key: $KEY_RECOVERED"
echo "Recovered IV: $IV_RECOVERED"

# Create test data
echo "Testing encryption with recovered credentials from multi-key share" > plaintext.txt

# Encrypt
openssl enc -aes-256-cbc -K "$KEY_RECOVERED" -iv "$IV_RECOVERED" \
  -in plaintext.txt -out encrypted.bin

# Decrypt
openssl enc -aes-256-cbc -d -K "$KEY_RECOVERED" -iv "$IV_RECOVERED" \
  -in encrypted.bin -out decrypted.txt

if diff plaintext.txt decrypted.txt > /dev/null; then
    echo "  ✅ ENCRYPTION/DECRYPTION WITH RECOVERED CREDENTIALS WORKS!"
else
    echo "  ❌ ENCRYPTION/DECRYPTION FAILED!"
    exit 1
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  TEST 3: Threshold Recovery (using subset of shares)"
echo "════════════════════════════════════════════════════════════════"

echo ""
echo "🔪 Step 1: Create temporary directory with only 3 of 5 shares..."
mkdir -p test_threshold_shares
cp test_multi_shares/key_shares_01.json test_threshold_shares/
cp test_multi_shares/key_shares_03.json test_threshold_shares/
cp test_multi_shares/key_shares_05.json test_threshold_shares/

echo "Copied shares 1, 3, and 5 (threshold=3)"

echo ""
echo "🔧 Step 2: Recover all keys using only threshold shares..."
./shamir-cli recover -dir ./test_threshold_shares

echo ""
echo "🔧 Step 3: Verify specific key recovery..."
./shamir-cli recover -dir ./test_threshold_shares -name key2 -output key2_threshold.txt

if diff key2_backup.txt key2_threshold.txt; then
    echo "  ✅ THRESHOLD RECOVERY WORKS!"
else
    echo "  ❌ THRESHOLD RECOVERY FAILED!"
    exit 1
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  TEST 4: Error Handling - Duplicate Key Name"
echo "════════════════════════════════════════════════════════════════"

echo ""
echo "🔪 Step 1: Try to add duplicate key name..."
if ./shamir-cli split -keyfile key1_backup.txt -total 5 -threshold 3 -dir ./test_multi_shares -name key1 2>&1 | grep -q "already exists"; then
    echo "  ✅ DUPLICATE KEY REJECTION WORKS!"
else
    echo "  ❌ DUPLICATE KEY SHOULD HAVE BEEN REJECTED!"
    exit 1
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  🎉 ALL MULTI-KEY TESTS PASSED!"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "Summary:"
echo "  ✅ Multiple hex keys can be stored in same share files"
echo "  ✅ Mixed key types (hex + raw JSON) work together"
echo "  ✅ Individual key recovery from multi-key shares works"
echo "  ✅ Threshold recovery works with multiple keys"
echo "  ✅ Duplicate key names are properly rejected"
echo "  ✅ Encryption credentials from multi-key shares work"
echo ""
echo "Share file structure supports:"
echo "  - Multiple keys per share file"
echo "  - Different key types (hex-encoded, raw bytes)"
echo "  - Selective recovery by key name"
echo "  - Recovery of all keys at once"
echo ""

# Cleanup
rm -rf test_threshold_shares