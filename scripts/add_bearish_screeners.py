import json

# Read the policy file
with open('data/policy.v1.json', 'r', encoding='utf-8') as f:
    policy = json.load(f)

# Function to convert bullish URL to bearish
def create_bearish_urls(bullish_urls, sector_name):
    """Create bearish versions of screeners"""

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

    # Universe (bearish): Replace ta_sma200_pa with ta_sma200_pb
    if 'universe' in bullish_urls:
        bearish_urls['universe_bearish'] = bullish_urls['universe'].replace('ta_sma200_pa', 'ta_sma200_pb')

    # Bounce (bearish): Price below SMA200, above SMA50 (temporary bounce), RSI overbought
    if 'pullback' in bullish_urls:
        # Replace ta_sma200_pa with ta_sma200_pb, ta_sma50_pb with ta_sma50_pa, ta_rsi_os40 with ta_rsi_ob60
        bearish_urls['bounce_bearish'] = (bullish_urls['pullback']
                                          .replace('ta_sma200_pa', 'ta_sma200_pb')
                                          .replace('ta_sma50_pb', 'ta_sma50_pa')
                                          .replace('ta_rsi_os40', 'ta_rsi_ob60'))

    # Breakdown (bearish): 52-week low
    if 'breakout' in bullish_urls:
        bearish_urls['breakdown_bearish'] = (bullish_urls['breakout']
                                             .replace('ta_sma200_pa', 'ta_sma200_pb')
                                             .replace('ta_highlow52w_nh', 'ta_highlow52w_nl'))

    # Death Cross (bearish): SMA50 below SMA200
    if 'golden_cross' in bullish_urls:
        bearish_urls['death_cross_bearish'] = (bullish_urls['golden_cross']
                                               .replace('ta_sma200_pa', 'ta_sma200_pb')
                                               .replace('ta_sma50_pa200', 'ta_sma50_pb200')
                                               .replace('ta_pattern_tlsupport', 'ta_pattern_tlresistance'))

    return bearish_urls

# Update each sector with bearish screeners
for sector in policy['sectors']:
    if 'screener_urls' in sector:
        bearish_urls = create_bearish_urls(sector['screener_urls'], sector['name'])
        sector['screener_urls'].update(bearish_urls)

# Write the updated policy back
with open('data/policy.v1.json', 'w', encoding='utf-8') as f:
    json.dump(policy, f, indent=2, ensure_ascii=False)

print("Added bearish screeners to all sectors in policy.v1.json")
print("\nBearish screeners added:")
print("  - universe_bearish: Stocks in downtrends (below SMA200)")
print("  - bounce_bearish: Temporary bounces in downtrends (RSI overbought)")
print("  - breakdown_bearish: 52-week lows")
print("  - death_cross_bearish: SMA50 crossing below SMA200")