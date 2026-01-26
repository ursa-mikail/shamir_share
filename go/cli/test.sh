# Manual test with a known key
#printf "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" > test.key

# Use -n flag with openssl (if supported) or pipe through tr
openssl rand -hex 32 | tr -d '\n' > test.key

./shamir-cli split -keyfile test.key -total 5 -threshold 3 -dir ./test_shares -name test.key
./shamir-cli recover -dir ./test_shares -name test.key -output test_recovered.key

echo "Original (hash):"
openssl sha256 test.key
echo ""
echo "Recovered (hash):"
openssl sha256 test_recovered.key
echo ""

diff test.key test_recovered.key && echo "✅ Files match!" || echo "❌ Files differ!"