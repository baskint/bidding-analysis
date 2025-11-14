-- Generate test fraud alerts
-- Run this SQL to populate your fraud dashboard with sample data

-- Insert test fraud alerts
INSERT INTO fraud_alerts (
    id, 
    campaign_id, 
    alert_type, 
    severity, 
    description, 
    affected_user_ids, 
    detected_at, 
    status
)
SELECT 
    uuid_generate_v4(),
    c.id,
    (ARRAY['click_velocity', 'geo_anomaly', 'device_anomaly', 'ip_anomaly', 'bot_detection'])[floor(random() * 5 + 1)],
    floor(random() * 7 + 3)::int, -- Severity 3-10
    CASE floor(random() * 5)::int
        WHEN 0 THEN 'Abnormal click velocity detected from IP ' || (floor(random() * 255) + 1)::text || '.' || (floor(random() * 255) + 1)::text || '.1.100'
        WHEN 1 THEN 'Geographic location mismatch detected - unusual activity from multiple countries'
        WHEN 2 THEN 'Suspicious device fingerprint detected - possible bot activity'
        WHEN 3 THEN 'Multiple rapid clicks from same source within 100ms window'
        ELSE 'IP reputation check failed - known fraud network detected'
    END,
    ARRAY[]::TEXT[],
    NOW() - (random() * interval '7 days'), -- Random time in last 7 days
    (ARRAY['active', 'investigating', 'resolved'])[floor(random() * 3 + 1)]
FROM campaigns c
CROSS JOIN generate_series(1, 15) -- 15 alerts per campaign
WHERE c.user_id = (SELECT id FROM users LIMIT 1) -- Your user ID
LIMIT 50; -- Total of 50 alerts

-- Update some predictions to have fraud_risk = true
-- This populates the "blocked bids" and "amount saved" metrics
UPDATE predictions p
SET fraud_risk = true
FROM bid_events be
WHERE p.bid_event_id = be.id
  AND be.user_id = (SELECT id FROM users LIMIT 1)
  AND random() < 0.05 -- Mark 5% of bids as fraud
  AND be.timestamp >= NOW() - interval '30 days';

-- Verify the data was created
SELECT 
    alert_type,
    COUNT(*) as count,
    AVG(severity) as avg_severity,
    COUNT(*) FILTER (WHERE status = 'active') as active_count
FROM fraud_alerts fa
JOIN campaigns c ON fa.campaign_id = c.id
WHERE c.user_id = (SELECT id FROM users LIMIT 1)
GROUP BY alert_type
ORDER BY count DESC;

-- Check fraud predictions
SELECT 
    COUNT(*) as total_predictions,
    COUNT(*) FILTER (WHERE fraud_risk = true) as fraud_predictions,
    SUM(CASE WHEN fraud_risk = true THEN be.bid_price ELSE 0 END) as amount_saved
FROM predictions p
JOIN bid_events be ON p.bid_event_id = be.id
WHERE be.user_id = (SELECT id FROM users LIMIT 1)
  AND be.timestamp >= NOW() - interval '30 days';
  