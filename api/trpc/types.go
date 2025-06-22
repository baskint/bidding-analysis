package trpc

// TRPCResponse represents a tRPC response structure
type TRPCResponse struct {
	Result *TRPCResult `json:"result,omitempty"`
	Error  *TRPCError  `json:"error,omitempty"`
}

// TRPCResult represents successful tRPC result
type TRPCResult struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

// TRPCError represents tRPC error
type TRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ProcessBidInput represents the input for bid processing
type ProcessBidInput struct {
	CampaignID            string   `json:"campaignId"`
	UserID                string   `json:"userId"`
	FloorPrice            float64  `json:"floorPrice"`
	DeviceType            string   `json:"deviceType"`
	OS                    string   `json:"os"`
	Browser               string   `json:"browser"`
	Country               string   `json:"country"`
	Region                string   `json:"region"`
	City                  string   `json:"city"`
	Keywords              []string `json:"keywords"`
	SegmentID             string   `json:"segmentId"`
	SegmentCategory       string   `json:"segmentCategory"`
	EngagementScore       float64  `json:"engagementScore"`
	ConversionProbability float64  `json:"conversionProbability"`
}

// CampaignStatsInput represents input for campaign statistics
type CampaignStatsInput struct {
	CampaignID string `json:"campaignId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

// BidHistoryInput represents input for bid history requests
type BidHistoryInput struct {
	CampaignID string `json:"campaignId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

// FraudAlertsInput represents input for fraud alerts
type FraudAlertsInput struct {
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	SeverityThreshold int    `json:"severityThreshold"`
}

// ModelAccuracyInput represents input for model accuracy
type ModelAccuracyInput struct {
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	ModelVersion string `json:"modelVersion"`
}
