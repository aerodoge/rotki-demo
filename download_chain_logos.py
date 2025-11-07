#!/usr/bin/env python3
"""
Download chain logos from DeBank API
"""
import os
import sys
import requests
import json
from pathlib import Path
from urllib.parse import urlparse

# DeBank API endpoint
API_ENDPOINT = "https://pro-openapi.debank.com/v1/chain/list"

# Output directory
OUTPUT_DIR = Path("frontend/public/images/chains")

def download_chain_logos():
    """Download all chain logos from DeBank"""

    # Create output directory if it doesn't exist
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    print(f"📁 Created directory: {OUTPUT_DIR}")

    # Try to get chain list from DeBank API
    print("\n🔍 Fetching chain list from DeBank API...")

    headers = {
        'accept': 'application/json'
    }

    try:
        response = requests.get(API_ENDPOINT, headers=headers, timeout=30)

        if response.status_code == 401:
            print("⚠️  API requires authentication. Trying alternative method...")
            # Try the public endpoint without auth
            alt_endpoint = "https://api.debank.com/chain/list"
            response = requests.get(alt_endpoint, headers=headers, timeout=30)

        if response.status_code != 200:
            print(f"❌ Failed to fetch chain list: HTTP {response.status_code}")
            print(f"Response: {response.text}")
            return

        data = response.json()

        # Debug: print first 500 chars of response
        print(f"📋 Response preview: {str(data)[:500]}")

        # Handle different response formats
        if isinstance(data, dict):
            # If response has nested data.chains structure
            if 'data' in data and isinstance(data['data'], dict) and 'chains' in data['data']:
                chains = data['data']['chains']
            # If response is a dict with 'data' key
            elif 'data' in data:
                chains = data['data']
            # If response is a dict with chain IDs as keys
            else:
                chains = []
                for chain_id, chain_data in data.items():
                    if isinstance(chain_data, dict):
                        chain_data['id'] = chain_id
                        chains.append(chain_data)
                    else:
                        # Simple string value, create basic object
                        chains.append({'id': chain_id, 'name': chain_data})
        else:
            chains = data

        print(f"✅ Found {len(chains)} chains")
        print(f"📋 First chain: {chains[0] if chains else 'None'}")

        # Download each logo
        success_count = 0
        skip_count = 0
        fail_count = 0

        for chain in chains:
            chain_id = chain.get('id')
            chain_name = chain.get('name', chain_id)
            logo_url = chain.get('logo_url')

            if not logo_url:
                print(f"⏭️  Skipping {chain_name} ({chain_id}): No logo URL")
                skip_count += 1
                continue

            # Determine file extension from URL
            parsed_url = urlparse(logo_url)
            ext = os.path.splitext(parsed_url.path)[1] or '.png'

            # Save with chain ID as filename
            output_path = OUTPUT_DIR / f"{chain_id}{ext}"

            try:
                print(f"📥 Downloading {chain_name} ({chain_id})...", end=" ")
                logo_response = requests.get(logo_url, timeout=30)

                if logo_response.status_code == 200:
                    with open(output_path, 'wb') as f:
                        f.write(logo_response.content)
                    print(f"✅ Saved to {output_path}")
                    success_count += 1
                else:
                    print(f"❌ Failed: HTTP {logo_response.status_code}")
                    fail_count += 1

            except Exception as e:
                print(f"❌ Error: {e}")
                fail_count += 1

        print(f"\n{'='*60}")
        print(f"📊 Summary:")
        print(f"   ✅ Successfully downloaded: {success_count}")
        print(f"   ⏭️  Skipped (no logo): {skip_count}")
        print(f"   ❌ Failed: {fail_count}")
        print(f"   📁 Total chains: {len(chains)}")
        print(f"{'='*60}")

        # Save chain list metadata
        metadata_path = OUTPUT_DIR / "chains.json"
        with open(metadata_path, 'w') as f:
            json.dump(chains, f, indent=2)
        print(f"\n💾 Saved chain metadata to {metadata_path}")

    except Exception as e:
        print(f"❌ Error fetching chain list: {e}")
        sys.exit(1)

if __name__ == "__main__":
    download_chain_logos()
