#from shamir_share import shamir_split, shamir_reconstruct, load_key_symmetric, get_key_hash
import json
import os

def test_shamir_without_inputs():
    """Test Shamir Secret Sharing with default values (no inputs)."""
    print("=== Test: Running Shamir Split with Default Values ===")

    key = b"default_secret_key:::[you should see this]"
    #key = b"default_secret_key:[you should see this]
    shares_required = 3
    shares_distributed = 5
    destination_folder = "./shamir_shares"

    shares = shamir_split(key, shares_required, shares_distributed)
    json_shares = {
        "key_root": {
            "hash": get_key_hash(key),
            "shares": shares
        }
    }

    os.makedirs(destination_folder, exist_ok=True)
    output_path = os.path.join(destination_folder, "shamir_shares.json")

    with open(output_path, "w") as outfile:
        json.dump(json_shares, outfile, indent=4)

    print(f"Shamir shares saved to {output_path}")

def test_shamir_reconstruct():
    """Test reconstructing the secret from saved shares."""
    print("\n=== Test: Running Shamir Reconstruction ===")

    shares_path = "./shamir_shares/shamir_shares.json"

    try:
        with open(shares_path, "r") as infile:
            data = json.load(infile)

        # Extract shares array from the nested structure
        shares_array = data["key_root"]

        shares = shares_array["shares"]# ["shares"]

        print("Shares:")
        for share in shares:
            print(share)

        # Reconstruct the key
        reconstructed_key = shamir_reconstruct(shares)
        print("\nReconstructed Key:")
        print(reconstructed_key)
        print("\nKey Hash (SHA-256):")
        print(get_key_hash(reconstructed_key))

        # Take the first 3 shares (or minimum required)
        #selected_shares = [
        #    {"share_id": str(i+1), "share_value": share}
        #    for i, share in enumerate(shares_array[:3])
        #]

        #print("Using Shares:", selected_shares)

        #recovered_key = shamir_reconstruct(selected_shares)
        #print(f"Reconstructed Key: {recovered_key}")

    except FileNotFoundError:
        print(f"Error: Share file '{shares_path}' not found! Run the split test first.")
    except Exception as e:
        print(f"Error reconstructing key: {e}")

if __name__ == "__main__":
    test_shamir_without_inputs()
    test_shamir_reconstruct()

"""
=== Test: Running Shamir Split with Default Values ===
Shamir shares saved to ./shamir_shares/shamir_shares.json

=== Test: Running Shamir Reconstruction ===
Shares:
required_shares
prime_mod
shares

Reconstructed Key:
b'default_secret_key:::[you should see this]'

Key Hash (SHA-256):
90e54e472e5b2f251b3f9aa17ae32441800555f9b1da3ca262fc94c7546eb673
"""
"""
# Test Script: test_usecase.py
import subprocess

def run_test():
    """Runs the shamir_share_demo script with no input arguments to test default behavior."""
    try:
        result = subprocess.run(
            ["python", "shamir_share_demo.py"],
            capture_output=True, text=True
        )
        # subprocess.run(["python", "shamir_share_demo.py", "--reconstruct"], capture_output=True, text=True)
        print("=== Test: Running `shamir_share_demo.py` with no arguments ===")
        print("Exit Code:", result.returncode)
        print("Output:\n", result.stdout)
        print("Errors:\n", result.stderr)


        print("=== Test: Running `shamir_share_demo.py` without arguments ===")
        print("Exit Code:", result.returncode)
        print("Output:\n", result.stdout)
        print("Errors:\n", result.stderr)

    except Exception as e:
        print(f"Test Failed: {e}")

if __name__ == "__main__":
    run_test()
"""