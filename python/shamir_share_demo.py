#!pip install sslib
from sslib import shamir
import json
import argparse
from hashlib import sha256
import os
import sys

def get_key_hash(key_bytes):
    """Returns the SHA-256 hash of the given key bytes."""
    sha256_hash = sha256(key_bytes)
    return sha256_hash.hexdigest()

def load_key_symmetric(keyfilepath):
    """Loads a symmetric key from a file, or generates a default key if the file is missing."""
    try:
        with open(keyfilepath, "rb") as keyfile:
            return keyfile.read()
    except FileNotFoundError:
        print(f"Warning: Key file '{keyfilepath}' not found. Using default key.")
        return b"default_secret_key:[you should see this]"

def shamir_split(key, shares_required, shares_distributed):
    """Splits a key into Shamir secret shares."""
    shares = shamir.to_base64(shamir.split_secret(key, shares_required, shares_distributed))
    return shares

def shamir_reconstruct(shares):
    """Reconstructs the key from Shamir shares."""
    return shamir.recover_secret(shamir.from_base64(shares))

def main(args=None):
    parser = argparse.ArgumentParser(description="Shamir Secret Sharing for Key Splitting and Reconstruction")
    parser.add_argument("--minimum-shares", type=int, default=3, help="Minimum shares required to reconstruct the key (default: 3)")
    parser.add_argument("--number-to-distribute", type=int, default=5, help="Total shares to distribute (default: 5)")
    parser.add_argument("--path-to-key", type=str, default="default_key.bin", help="Path to the key file to be split (default: 'default_key.bin')")
    parser.add_argument("--path-to-destination-folder", type=str, default="./shamir_shares", help="Destination folder to store shares (default: './shamir_shares')")
    parser.add_argument("--reconstruct", action="store_true", help="Reconstruct the key from shares")

    # Parse known arguments and ignore unrecognized ones
    if args is None:
        args, _ = parser.parse_known_args()
    else:
        args, _ = parser.parse_known_args(args)

    if args.reconstruct:
        # Reconstruct the key
        json_file_path = os.path.join(args.path_to_destination_folder, "shamir_shares.json")
        try:
            with open(json_file_path, "r") as file:
                data = json.load(file)
                shares_data = data["key_root"]["shares"]
                shares = shares_data["shares"]# ["shares"]

                print("Shares:")
                for share in shares:
                    print(share)

                # Reconstruct the key
                reconstructed_key = shamir_reconstruct(shares)
                print("\nReconstructed Key:")
                print(reconstructed_key)
                print("\nKey Hash (SHA-256):")
                print(get_key_hash(reconstructed_key))
        except FileNotFoundError:
            print(f"Error: File '{json_file_path}' not found.")
        except KeyError as e:
            print(f"Error: Invalid JSON structure. Missing key: {e}")
    else:
        # Split the key
        key = load_key_symmetric(args.path_to_key)
        shares = shamir_split(key, args.minimum_shares, args.number_to_distribute)

        json_shares = {
            "key_root": {
                "hash": get_key_hash(key),
                "shares": {
                    "required_shares": args.minimum_shares,
                    "prime_mod": "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAhQ==",
                    "shares": shares
                }
            }
        }

        os.makedirs(args.path_to_destination_folder, exist_ok=True)
        output_path = os.path.join(args.path_to_destination_folder, "shamir_shares.json")

        with open(output_path, "w") as outfile:
            json.dump(json_shares, outfile, indent=4)

        print(f"Shamir shares saved to {output_path}")

if __name__ == "__main__":
    # Example usage:
    # To split the key:
    main()
    # To reconstruct the key:
    main(["--reconstruct"])
    #main()

"""
Warning: Key file 'default_key.bin' not found. Using default key.
Shamir shares saved to ./shamir_shares/shamir_shares.json
Shares:
required_shares
prime_mod
shares

Reconstructed Key:
b'default_secret_key:[you should see this]'

Key Hash (SHA-256):
f21525b3cea8d85d6566eb6d72bac07060d9df2f385560aae7bbc7ebf30a21ea
"""