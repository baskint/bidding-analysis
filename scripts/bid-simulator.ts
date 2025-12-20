// scripts/bid-simulator.ts
/**
 * Bid Simulator - Continuously generates realistic bid events
 * 
 * Usage:
 *   ts-node scripts/bid-simulator.ts
 * 
 * This will generate 1 bid every 2-5 seconds until stopped (Ctrl+C)
 */

const API_BASE_URL = process.env.API_URL || 'http://localhost:8080';
const AUTH_TOKEN = process.env.AUTH_TOKEN || '';

if (!AUTH_TOKEN) {
  console.error('ERROR: Please set AUTH_TOKEN environment variable');
  console.error('   Example: export AUTH_TOKEN="your-jwt-token"');
  process.exit(1);
}

const CAMPAIGNS = [
  '69486fba-6b27-450e-bb05-8a1160ba1e10',
  '0e300360-d4e9-4492-82f9-e2118b819042',
  '687e0338-7d72-4df2-880a-f151048eabc2',
  '868f4419-9612-4c07-a897-457860c07fb5',
  'ba3f8fdc-4f76-4bfc-99fa-508bd454c7ba',
];

const DEVICE_TYPES = ['mobile', 'desktop', 'tablet'];
const COUNTRIES = ['US', 'UK', 'CA', 'AU', 'DE', 'FR', 'JP'];
const SEGMENTS = ['tech', 'retail', 'gaming', 'b2b', 'education'];
const BROWSERS = ['Chrome', 'Safari', 'Firefox', 'Edge'];
const OS_LIST = ['Windows', 'macOS', 'iOS', 'Android', 'Linux'];

function randomElement<T>(arr: T[]): T {
  return arr[Math.floor(Math.random() * arr.length)];
}

function randomBetween(min: number, max: number): number {
  return Math.random() * (max - min) + min;
}

function generateBid() {
  const floorPrice = randomBetween(1.0, 3.0);
  const bidPrice = floorPrice * randomBetween(1.2, 2.5);

  return {
    campaign_id: randomElement(CAMPAIGNS),
    user_id: '21d834cc-2f58-4897-9c87-7cc09e686110',
    bid_price: parseFloat(bidPrice.toFixed(4)),
    floor_price: parseFloat(floorPrice.toFixed(4)),
    device_type: randomElement(DEVICE_TYPES),
    os: randomElement(OS_LIST),
    browser: randomElement(BROWSERS),
    country: randomElement(COUNTRIES),
    segment_id: randomElement(SEGMENTS),
    segment_category: randomElement(SEGMENTS),
    keywords: ['technology', 'gadgets'],
  };
}

async function submitBid() {
  const bid = generateBid();

  try {
    const response = await fetch(API_BASE_URL + '/trpc/bidding.submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + AUTH_TOKEN,
      },
      body: JSON.stringify(bid),
    });

    if (response.ok) {
      const data: any = await response.json();
      const result = data.result?.data;
      console.log('✓ Bid submitted: $' + bid.bid_price.toFixed(4) +
        ' | ' + bid.device_type +
        ' | ' + bid.country +
        ' | Status: ' + (result?.status || 'unknown'));
    } else {
      const error = await response.text();
      console.error('✗ Failed:', error.substring(0, 100));
    }
  } catch (error) {
    console.error('✗ Error:', error);
  }
}

async function runSimulator() {
  console.log('🚀 Bid Simulator Started');
  console.log('   API: ' + API_BASE_URL);
  console.log('   Generating bids every 2-5 seconds...');
  console.log('   Press Ctrl+C to stop');
  console.log('');

  let count = 0;

  while (true) {
    count++;
    console.log('[' + count + '] Generating bid...');
    await submitBid();

    // Random delay between 2-5 seconds
    const delay = 2000 + Math.random() * 3000;
    await new Promise(resolve => setTimeout(resolve, delay));
  }
}

runSimulator().catch(console.error);
