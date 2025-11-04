import json
import re

# Read the policy file
with open('data/policy.v1.json', 'r', encoding='utf-8') as f:
    policy = json.load(f)

# Sectors that should NOT have bearish screeners (incompatible with trend-following)
SKIP_BEARISH_SECTORS = ["Utilities", "Energy"]

# Function to convert bullish URL to bearish
def create_bearish_urls(bullish_urls, sector_name):
    """Create bearish versions of screeners with proper logic"""

    # Skip incompatible sectors
    if sector_name in SKIP_BEARISH_SECTORS:
        print(f"  [WARNING] Skipping {sector_name} (incompatible with trend-following)")
        return {}

    # Get sector filter string
    sector_map = {
        "Healthcare": "sec_healthcare",
        "Technology": "sec_technology",
        "Consumer Discretionary": "sec_consumercyclical",
        "Industrials": "sec_industrialgoods",
        "Communication Services": "sec_communication",
        "Consumer Defensive": "sec_consumergoods",
        "Financials": "sec_financial",
        "Real Estate": "sec_realestate",
        "Energy": "sec_energy",
        "Utilities": "sec_utilities"
    }

    sector_filter = sector_map.get(sector_name, "")
    bearish_urls = {}

    # Universe (bearish): REMOVE positive fundamental filters (they conflict with bearish logic)
    if 'universe' in bullish_urls:
        url = bullish_urls['universe']

        # Replace ta_sma200_pa with ta_sma200_pb (below SMA200)
        url = url.replace('ta_sma200_pa', 'ta_sma200_pb')

        # CRITICAL FIX: Remove positive fundamental filters
        # These create "value traps" - fundamentally strong stocks in downtrends are likely to reverse
        url = re.sub(r',fa_epsyoy_pos', '', url)
        url = re.sub(r',fa_epsyoy1_pos', '', url)
        url = re.sub(r',fa_sales5years_pos', '', url)
        url = re.sub(r',fa_roe_pos', '', url)
        url = re.sub(r',fa_roe_o15', '', url)

        bearish_urls['universe_bearish'] = url

    # Bounce (bearish): Price below SMA200, above SMA50 (temporary bounce), RSI overbought
    # This is a BEAR FLAG setup - solid logic for shorting bounces in downtrends
    if 'pullback' in bullish_urls:
        bearish_urls['bounce_bearish'] = (bullish_urls['pullback']
                                          .replace('ta_sma200_pa', 'ta_sma200_pb')
                                          .replace('ta_sma50_pb', 'ta_sma50_pa')
                                          .replace('ta_rsi_os40', 'ta_rsi_ob60'))

    # Breakdown (bearish): 52-week low + volume confirmation
    # IMPROVED: Add relative volume filter to avoid dead cat bounces
    if 'breakout' in bullish_urls:
        url = (bullish_urls['breakout']
               .replace('ta_sma200_pa', 'ta_sma200_pb')
               .replace('ta_highlow52w_nh', 'ta_highlow52w_nl'))

        # Add relative volume filter for conviction (avoid low-volume drifts)
        if 'sh_relvol' not in url:
            # Insert relative volume filter before &ft=4
            url = url.replace('&ft=4', ',sh_relvol_o2&ft=4')

        bearish_urls['breakdown_bearish'] = url

    # Death Cross (bearish): SMA50 below SMA200
    # Solid technical setup - confirmed downtrend on multiple timeframes
    if 'golden_cross' in bullish_urls:
        bearish_urls['death_cross_bearish'] = (bullish_urls['golden_cross']
                                               .replace('ta_sma200_pa', 'ta_sma200_pb')
                                               .replace('ta_sma50_pa200', 'ta_sma50_pb200')
                                               .replace('ta_pattern_tlsupport', 'ta_pattern_tlresistance'))

    return bearish_urls

# Update each sector with bearish screeners (skip incompatible sectors)
for sector in policy['sectors']:
    if 'screener_urls' in sector:
        print(f"\nProcessing {sector['name']}...")

        # Remove old bearish screeners first (clean slate)
        keys_to_remove = [k for k in sector['screener_urls'].keys() if 'bearish' in k]
        for key in keys_to_remove:
            del sector['screener_urls'][key]

        # Create new bearish screeners
        bearish_urls = create_bearish_urls(sector['screener_urls'], sector['name'])

        if bearish_urls:
            sector['screener_urls'].update(bearish_urls)
            print(f"  [OK] Added {len(bearish_urls)} bearish screeners")
        else:
            print(f"  [SKIP] No bearish screeners added")

# Write the updated policy back
with open('data/policy.v1.json', 'w', encoding='utf-8') as f:
    json.dump(policy, f, indent=2, ensure_ascii=False)

print("\n" + "="*60)
print("SUCCESS: FIXED bearish screeners added to policy.v1.json")
print("="*60)
print("\nKey fixes applied:")
print("  1. [OK] Removed positive fundamental filters from universe_bearish")
print("     (fa_epsyoy_pos, fa_sales5years_pos, fa_roe_pos)")
print("  2. [OK] Skipped bearish screeners for Utilities and Energy")
print("     (incompatible with trend-following)")
print("  3. [OK] Added volume confirmation to breakdown_bearish")
print("     (sh_relvol_o2 to avoid low-conviction breakdowns)")
print("  4. [OK] Kept bounce_bearish and death_cross_bearish (good logic)")
print("\nBearish screeners by type:")
print("  - universe_bearish: Stocks in downtrends (below SMA200)")
print("  - bounce_bearish: Bear flag setups (temporary bounces in downtrends)")
print("  - breakdown_bearish: 52-week lows with volume (conviction moves)")
print("  - death_cross_bearish: SMA50 crossed below SMA200 (confirmed downtrends)")
print("\n[WARNING] IMPORTANT: These are UNTESTED. Validate with paper trading first!")
