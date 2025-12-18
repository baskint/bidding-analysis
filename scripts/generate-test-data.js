// scripts/generate-test-data.ts
/**
 * Test Data Generator for Bidding Analysis
 */
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g = Object.create((typeof Iterator === "function" ? Iterator : Object).prototype);
    return g.next = verb(0), g["throw"] = verb(1), g["return"] = verb(2), typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var API_BASE_URL = process.env.API_URL || 'http://localhost:8080';
var AUTH_TOKEN = process.env.AUTH_TOKEN || '';
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
var NUM_EVENTS = 2; // parseInt(process.argv[2] || '500', 10);
var BATCH_SIZE = 50;
// Data pools
var CAMPAIGNS = [
    { id: '69486fba-6b27-450e-bb05-8a1160ba1e10', name: 'Tech Products Q4' },
    { id: '0e300360-d4e9-4492-82f9-e2118b819042', name: 'Holiday Shopping' },
    { id: '687e0338-7d72-4df2-880a-f151048eabc2', name: 'Gaming Hardware' },
    { id: '868f4419-9612-4c07-a897-457860c07fb5', name: 'Mobile Apps' },
    { id: 'ba3f8fdc-4f76-4bfc-99fa-508bd454c7ba', name: 'E-commerce Summer' },
];
var DEVICE_TYPES = [
    { type: 'mobile', weight: 0.6 },
    { type: 'desktop', weight: 0.3 },
    { type: 'tablet', weight: 0.1 },
];
var COUNTRIES = [
    { code: 'US', weight: 0.4 },
    { code: 'UK', weight: 0.15 },
    { code: 'CA', weight: 0.1 },
    { code: 'AU', weight: 0.08 },
    { code: 'DE', weight: 0.12 },
    { code: 'FR', weight: 0.08 },
    { code: 'JP', weight: 0.07 },
];
var SEGMENTS = [
    { id: 'tech_enthusiasts', category: 'tech', weight: 0.25 },
    { id: 'shoppers', category: 'retail', weight: 0.30 },
    { id: 'gamers', category: 'gaming', weight: 0.20 },
    { id: 'business', category: 'b2b', weight: 0.15 },
    { id: 'students', category: 'education', weight: 0.10 },
];
var BROWSERS = ['Chrome', 'Safari', 'Firefox', 'Edge'];
var OS_LIST = ['Windows', 'macOS', 'iOS', 'Android', 'Linux'];
// Helper: Weighted random selection
function weightedRandom(items) {
    var totalWeight = items.reduce(function (sum, item) { return sum + item.weight; }, 0);
    var random = Math.random() * totalWeight;
    for (var _i = 0, items_1 = items; _i < items_1.length; _i++) {
        var item = items_1[_i];
        random -= item.weight;
        if (random <= 0)
            return item;
    }
    return items[items.length - 1];
}
// Helper: Random between min and max
function randomBetween(min, max) {
    return Math.random() * (max - min) + min;
}
// Helper: Random int between min and max (inclusive)
function randomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}
// Helper: Random element from array
function randomElement(arr) {
    return arr[Math.floor(Math.random() * arr.length)];
}
// Helper: Generate timestamp in the past N days
function randomTimestamp(daysAgo) {
    var now = new Date();
    var msAgo = daysAgo * 24 * 60 * 60 * 1000;
    var timestamp = now.getTime() - Math.random() * msAgo;
    return new Date(timestamp);
}
// Generate a single bid event
function generateBidEvent() {
    var campaign = randomElement(CAMPAIGNS);
    var device = weightedRandom(DEVICE_TYPES);
    var country = weightedRandom(COUNTRIES);
    var segment = weightedRandom(SEGMENTS);
    var baseFloorPrice = segment.category === 'b2b' ? 3.0 : 1.5;
    var floorPrice = randomBetween(baseFloorPrice * 0.8, baseFloorPrice * 1.2);
    var bidMultiplier = randomBetween(1.2, 2.5);
    var bidPrice = floorPrice * bidMultiplier;
    var ratio = bidPrice / floorPrice;
    var winProbability = 0.1 + (0.85 * Math.min((ratio - 1.0) / 2.0, 1.0));
    if (device.type === 'mobile')
        winProbability *= 0.9;
    if (device.type === 'tablet')
        winProbability *= 0.85;
    var won = Math.random() < winProbability;
    var winPrice = won ? bidPrice * randomBetween(0.85, 0.95) : undefined;
    var conversionRate = segment.category === 'b2b' ? 0.15 : 0.08;
    var converted = won && Math.random() < conversionRate;
    var keywordPools = {
        tech: ['technology', 'gadgets', 'electronics', 'innovation'],
        retail: ['shopping', 'deals', 'fashion', 'sale'],
        gaming: ['gaming', 'esports', 'pc', 'console'],
        b2b: ['business', 'enterprise', 'saas', 'software'],
        education: ['learning', 'courses', 'education', 'training'],
    };
    var keywords = [
        randomElement(keywordPools[segment.category] || ['general']),
        randomElement(keywordPools[segment.category] || ['general']),
    ];
    return __assign(__assign({ campaign_id: campaign.id, user_id: '21d834cc-2f58-4897-9c87-7cc09e686110', bid_price: parseFloat(bidPrice.toFixed(4)), floor_price: parseFloat(floorPrice.toFixed(4)), won: won }, (winPrice && { win_price: parseFloat(winPrice.toFixed(4)) })), { converted: converted, device_type: device.type, os: randomElement(OS_LIST), browser: randomElement(BROWSERS), country: country.code, segment_id: segment.id, segment_category: segment.category, keywords: keywords, timestamp: randomTimestamp(30).toISOString() });
}
// Submit a batch of bid events
function submitBatch(events) {
    return __awaiter(this, void 0, void 0, function () {
        var results, _i, events_1, event_1, response, error, error_1;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    results = { success: 0, failed: 0 };
                    _i = 0, events_1 = events;
                    _a.label = 1;
                case 1:
                    if (!(_i < events_1.length)) return [3 /*break*/, 9];
                    event_1 = events_1[_i];
                    _a.label = 2;
                case 2:
                    _a.trys.push([2, 7, , 8]);
                    return [4 /*yield*/, fetch(API_BASE_URL + '/trpc/bidding.submit', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                                'Authorization': 'Bearer ' + AUTH_TOKEN,
                            },
                            body: JSON.stringify(event_1),
                        })];
                case 3:
                    response = _a.sent();
                    if (!response.ok) return [3 /*break*/, 4];
                    results.success++;
                    return [3 /*break*/, 6];
                case 4:
                    results.failed++;
                    return [4 /*yield*/, response.text()];
                case 5:
                    error = _a.sent();
                    console.error('Failed to submit bid: ' + error.substring(0, 100));
                    _a.label = 6;
                case 6: return [3 /*break*/, 8];
                case 7:
                    error_1 = _a.sent();
                    results.failed++;
                    console.error('Error submitting bid:', error_1);
                    return [3 /*break*/, 8];
                case 8:
                    _i++;
                    return [3 /*break*/, 1];
                case 9: return [2 /*return*/, results];
            }
        });
    });
}
// Main execution
function main() {
    return __awaiter(this, void 0, void 0, function () {
        return __generator(this, function (_a) {
            console.log('Starting test data generation...');
            console.log('Target: ' + NUM_EVENTS + ' bid events');
            console.log('Batch size: ' + BATCH_SIZE);
            console.log('API: ' + API_BASE_URL);
            console.log('');
            // const totalResults = { success: 0, failed: 0 };
            // const startTime = Date.now();
            // // Generate and submit in batches
            // for (let i = 0; i < NUM_EVENTS; i += BATCH_SIZE) {
            //   const batchSize = Math.min(BATCH_SIZE, NUM_EVENTS - i);
            //   const batch = Array.from({ length: batchSize }, () => generateBidEvent());
            //   const batchNum = Math.floor(i / BATCH_SIZE) + 1;
            //   const totalBatches = Math.ceil(NUM_EVENTS / BATCH_SIZE);
            //   console.log('Processing batch ' + batchNum + '/' + totalBatches + '...');
            //   const results = await submitBatch(batch);
            //   totalResults.success += results.success;
            //   totalResults.failed += results.failed;
            //   const progress = ((i + batchSize) / NUM_EVENTS * 100).toFixed(1);
            //   console.log('  Success: ' + results.success + ', Failed: ' + results.failed + ' (' + progress + '% complete)');
            //   // Small delay between batches
            //   if (i + BATCH_SIZE < NUM_EVENTS) {
            //     await new Promise(resolve => setTimeout(resolve, 500));
            //   }
            // }
            // const duration = ((Date.now() - startTime) / 1000).toFixed(1);
            // console.log('');
            // console.log('Test data generation complete!');
            // console.log('   Success: ' + totalResults.success);
            // console.log('   Failed: ' + totalResults.failed);
            // console.log('   Duration: ' + duration + 's');
            // console.log('   Rate: ' + (totalResults.success / parseFloat(duration)).toFixed(1) + ' events/sec');
            // console.log('');
            console.log('Your dashboard should now have realistic data!');
            return [2 /*return*/];
        });
    });
}
main().catch(function (error) {
    console.error('');
    console.error('Fatal error:');
    console.error(error);
    process.exit(1);
});
