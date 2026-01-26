#!/bin/bash
set -e

echo "════════════════════════════════════════════════════════════════"
echo "  Shamir Secret Sharing - Comprehensive Test"
echo "════════════════════════════════════════════════════════════════"

# Cleanup
rm -rf test_shares test_shares_hex
rm -f secret.json secret_recovered.json test.key test_recovered.key
rm -f plaintext.txt encrypted.bin decrypted.txt

echo ""
echo "🔨 Building shamir-cli..."
go build -o shamir-cli

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  TEST 1: JSON File Splitting (with -raw flag)"
echo "════════════════════════════════════════════════════════════════"

echo ""
echo "🔑 Step 1: Generate key and IV as JSON..."
KEY=$(openssl rand -hex 32 | tr -d '\n')
IV=$(openssl rand -hex 16 | tr -d '\n')

cat > secret.json <<EOF
{
  "key": "$KEY",
  "iv": "$IV"
}
EOF

echo "Created secret.json:"
cat secret.json

echo ""
echo "🔪 Step 2: Split JSON file using Shamir..."
./shamir-cli split -keyfile secret.json -raw -total 5 -threshold 3 -dir ./test_shares -name secret.json

echo ""
echo "🔥 Step 3: Delete original JSON (simulate loss)..."
cp secret.json secret_backup.json
rm secret.json
echo "❌ Deleted secret.json"

echo ""
echo "🔧 Step 4: Recover JSON from shares..."
./shamir-cli recover -dir ./test_shares -name secret.json -raw -output secret_recovered.json

echo ""
echo "✅ Step 5: Verify recovery..."
echo "Original (hash):"
openssl sha256 secret_backup.json

echo ""
echo "Recovered (hash):"
openssl sha256 secret_recovered.json

echo ""
if diff secret_backup.json secret_recovered.json; then
    echo "✅ JSON FILES MATCH!"
else
    echo "❌ JSON FILES DIFFER!"
    exit 1
fi

echo ""
echo "🔓 Step 6: Extract credentials and test encryption..."
KEY_RECOVERED=$(jq -r '.key' secret_recovered.json)
IV_RECOVERED=$(jq -r '.iv' secret_recovered.json)

echo "Recovered Key: $KEY_RECOVERED"
echo "Recovered IV: $IV_RECOVERED"

# Create test data
echo "This is secret data for encryption test" > plaintext.txt

# Encrypt with recovered credentials
openssl enc -aes-256-cbc -K "$KEY_RECOVERED" -iv "$IV_RECOVERED" \
  -in plaintext.txt -out encrypted.bin

# Decrypt
openssl enc -aes-256-cbc -d -K "$KEY_RECOVERED" -iv "$IV_RECOVERED" \
  -in encrypted.bin -out decrypted.txt

echo ""
if diff plaintext.txt decrypted.txt > /dev/null; then
    echo "✅ ENCRYPTION/DECRYPTION TEST PASSED!"
else
    echo "❌ ENCRYPTION/DECRYPTION TEST FAILED!"
    exit 1
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  TEST 2: Hex Key Splitting (traditional method)"
echo "════════════════════════════════════════════════════════════════"

echo ""
echo "🔑 Step 1: Generate hex key..."
openssl rand -hex 32 | tr -d '\n' > test.key
echo "Generated test.key"

echo ""
echo "🔪 Step 2: Split hex key..."
./shamir-cli split -keyfile test.key -total 5 -threshold 3 -dir ./test_shares_hex -name test.key

echo ""
echo "🔥 Step 3: Delete original key..."
cp test.key test_backup.key
rm test.key
echo "❌ Deleted test.key"

echo ""
echo "🔧 Step 4: Recover key from shares..."
./shamir-cli recover -dir ./test_shares_hex -name test.key -output test_recovered.key

echo ""
echo "✅ Step 5: Verify recovery..."
echo "Original (hash):"
openssl sha256 test_backup.key

echo ""
echo "Recovered (hash):"
openssl sha256 test_recovered.key

echo ""
if diff test_backup.key test_recovered.key; then
    echo "✅ HEX KEY FILES MATCH!"
else
    echo "❌ HEX KEY FILES DIFFER!"
    exit 1
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  🎉 ALL TESTS PASSED!"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "Summary:"
echo "  ✅ JSON file split and recovery works"
echo "  ✅ Recovered credentials can encrypt/decrypt data"
echo "  ✅ Traditional hex key splitting still works"
echo ""

