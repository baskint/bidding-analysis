// scripts/generate-test-data.ts
/**
 * Test Data Generator for Bidding Analysis
 */

// const API_BASE_URL = process.env.API_URL || 'http://localhost:8080';
// const AUTH_TOKEN = process.env.AUTH_TOKEN || '';

const AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYmM2OWZkYjctZWM1ZS00NjVjLWE0YzItNDU0ODA4ZDlmMjdmIiwidXNlcm5hbWUiOiJiYXNraW5AdGFwa2FuLmNvbSIsImV4cCI6MTc2NjYzMTM3OCwiaWF0IjoxNzY2MDI2NTc4fQ.ZLoKUcN1XvsBx1JCBba-tF5XydBtlbDsD8DniATjIAo"
const API_BASE_URL="https://bidding-analysis-539382269313.us-central1.run.app"

console.log('Configuration:');
console.log('   API_URL: ' + API_BASE_URL);
console.log('   AUTH_TOKEN: ' + (AUTH_TOKEN ? 'Set' : 'Not set'));
console.log('');

if (!AUTH_TOKEN) {
  console.error('ERROR: Please set AUTH_TOKEN environment variable');
  console.error('   Example: export AUTH_TOKEN="your-jwt-token"');
  process.exit(1);
}

// Configuration
const NUM_EVENTS = 1000; // parseInt(process.argv[2] || '500', 10);
const BATCH_SIZE = 50;

// Data pools
const CAMPAIGNS = [
  { id: '1e6aad91-9499-434d-ab37-a61627b00c5b', name: 'Tech Products Q4' },
  { id: '65fc5974-10c1-4957-8a1b-7614a311292b', name: 'Holiday Shopping' },
  { id: 'a17e4042-45cb-4e12-b318-75376b74a1bb', name: 'Gaming Hardware' },
  { id: 'a1bb71c9-6daf-4300-8b59-6895986d2dc1', name: 'Mobile Apps' },
  { id: '7d001b14-0ddd-4de0-8b7d-5d0457ffd08e', name: 'E-commerce Summer' },
];

const DEVICE_TYPES = [
  { type: 'mobile', weight: 0.6 },
  { type: 'desktop', weight: 0.3 },
  { type: 'tablet', weight: 0.1 },
];

const COUNTRIES = [
  { code: 'US', weight: 0.4 },
  { code: 'UK', weight: 0.15 },
  { code: 'CA', weight: 0.1 },
  { code: 'AU', weight: 0.08 },
  { code: 'DE', weight: 0.12 },
  { code: 'FR', weight: 0.08 },
  { code: 'JP', weight: 0.07 },
];

const SEGMENTS = [
  { id: 'tech_enthusiasts', category: 'tech', weight: 0.25 },
  { id: 'shoppers', category: 'retail', weight: 0.30 },
  { id: 'gamers', category: 'gaming', weight: 0.20 },
  { id: 'business', category: 'b2b', weight: 0.15 },
  { id: 'students', category: 'education', weight: 0.10 },
];

const BROWSERS = ['Chrome', 'Safari', 'Firefox', 'Edge'];
const OS_LIST = ['Windows', 'macOS', 'iOS', 'Android', 'Linux'];

// Helper: Weighted random selection
function weightedRandom<T extends { weight: number }>(items: T[]): T {
  const totalWeight = items.reduce((sum, item) => sum + item.weight, 0);
  let random = Math.random() * totalWeight;

  for (const item of items) {
    random -= item.weight;
    if (random <= 0) return item;
  }

  return items[items.length - 1];
}

// Helper: Random between min and max
function randomBetween(min: number, max: number): number {
  return Math.random() * (max - min) + min;
}

// Helper: Random int between min and max (inclusive)
function randomInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

// Helper: Random element from array
function randomElement<T>(arr: T[]): T {
  return arr[Math.floor(Math.random() * arr.length)];
}

// Helper: Generate timestamp in the past N days
function randomTimestamp(daysAgo: number): Date {
  const now = new Date();
  const msAgo = daysAgo * 24 * 60 * 60 * 1000;
  const timestamp = now.getTime() - Math.random() * msAgo;
  return new Date(timestamp);
}

// Generate a single bid event
function generateBidEvent() {
  const campaign = randomElement(CAMPAIGNS);
  const device = weightedRandom(DEVICE_TYPES);
  const country = weightedRandom(COUNTRIES);
  const segment = weightedRandom(SEGMENTS);

  const baseFloorPrice = segment.category === 'b2b' ? 3.0 : 1.5;
  const floorPrice = randomBetween(baseFloorPrice * 0.8, baseFloorPrice * 1.2);

  const bidMultiplier = randomBetween(1.2, 2.5);
  const bidPrice = floorPrice * bidMultiplier;

  const ratio = bidPrice / floorPrice;
  let winProbability = 0.1 + (0.85 * Math.min((ratio - 1.0) / 2.0, 1.0));

  if (device.type === 'mobile') winProbability *= 0.9;
  if (device.type === 'tablet') winProbability *= 0.85;

  const won = Math.random() < winProbability;

  const winPrice = won ? bidPrice * randomBetween(0.85, 0.95) : undefined;

  const conversionRate = segment.category === 'b2b' ? 0.15 : 0.08;
  const converted = won && Math.random() < conversionRate;

  const keywordPools: Record<string, string[]> = {
    tech: ['technology', 'gadgets', 'electronics', 'innovation'],
    retail: ['shopping', 'deals', 'fashion', 'sale'],
    gaming: ['gaming', 'esports', 'pc', 'console'],
    b2b: ['business', 'enterprise', 'saas', 'software'],
    education: ['learning', 'courses', 'education', 'training'],
  };

  const keywords = [
    randomElement(keywordPools[segment.category] || ['general']),
    randomElement(keywordPools[segment.category] || ['general']),
  ];

  return {
    campaign_id: campaign.id,
    user_id: 'bc69fdb7-ec5e-465c-a4c2-454808d9f27f',
    bid_price: parseFloat(bidPrice.toFixed(4)),
    floor_price: parseFloat(floorPrice.toFixed(4)),
    won,
    ...(winPrice && { win_price: parseFloat(winPrice.toFixed(4)) }),
    converted,
    device_type: device.type,
    os: randomElement(OS_LIST),
    browser: randomElement(BROWSERS),
    country: country.code,
    segment_id: segment.id,
    segment_category: segment.category,
    keywords,
    timestamp: randomTimestamp(30).toISOString(),
  };
}

// Submit a batch of bid events
async function submitBatch(events: any[]) {
  const results = { success: 0, failed: 0 };

  for (const event of events) {
    try {
      const response = await fetch(API_BASE_URL + '/trpc/bidding.submit', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + AUTH_TOKEN,
        },
        body: JSON.stringify(event),
      });

      if (response.ok) {
        results.success++;
      } else {
        results.failed++;
        const error = await response.text();
        console.error('Failed to submit bid: ' + error.substring(0, 100));
      }
    } catch (error) {
      results.failed++;
      console.error('Error submitting bid:', error);
    }
  }

  return results;
}

// Main execution
async function main() {
  console.log('Starting test data generation...');
  console.log('Target: ' + NUM_EVENTS + ' bid events');
  console.log('Batch size: ' + BATCH_SIZE);
  console.log('API: ' + API_BASE_URL);
  console.log('');

  const totalResults = { success: 0, failed: 0 };
  const startTime = Date.now();

  // Generate and submit in batches
  for (let i = 0; i < NUM_EVENTS; i += BATCH_SIZE) {
    const batchSize = Math.min(BATCH_SIZE, NUM_EVENTS - i);
    const batch = Array.from({ length: batchSize }, () => generateBidEvent());

    const batchNum = Math.floor(i / BATCH_SIZE) + 1;
    const totalBatches = Math.ceil(NUM_EVENTS / BATCH_SIZE);
    console.log('Processing batch ' + batchNum + '/' + totalBatches + '...');

    const results = await submitBatch(batch);
    totalResults.success += results.success;
    totalResults.failed += results.failed;

    const progress = ((i + batchSize) / NUM_EVENTS * 100).toFixed(1);
    console.log('  Success: ' + results.success + ', Failed: ' + results.failed + ' (' + progress + '% complete)');

    // Small delay between batches
    if (i + BATCH_SIZE < NUM_EVENTS) {
      await new Promise(resolve => setTimeout(resolve, 500));
    }
  }

  const duration = ((Date.now() - startTime) / 1000).toFixed(1);

  console.log('');
  console.log('Test data generation complete!');
  console.log('   Success: ' + totalResults.success);
  console.log('   Failed: ' + totalResults.failed);
  console.log('   Duration: ' + duration + 's');
  console.log('   Rate: ' + (totalResults.success / parseFloat(duration)).toFixed(1) + ' events/sec');
  console.log('');
  console.log('Your dashboard should now have realistic data!');
}

main().catch((error) => {
  console.error('');
  console.error('Fatal error:');
  console.error(error);
  process.exit(1);
});
