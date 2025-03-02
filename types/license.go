package types

import (
	"time"

	"github.com/PaddleHQ/paddle-go-sdk/v3"
)

type CloudSubscriptionPlan string

const (
	CloudSubscriptionPlanPaused     CloudSubscriptionPlan = "paused"
	CloudSubscriptionPlanFree       CloudSubscriptionPlan = "free"
	CloudSubscriptionPlanPro        CloudSubscriptionPlan = "pro"
	CloudSubscriptionPlanEnterprise CloudSubscriptionPlan = "enterprise"
)

type CloudSubscriptionStatus string

const (
	CloudSubscriptionStatusActive   CloudSubscriptionStatus = "active"
	CloudSubscriptionStatusTrialing CloudSubscriptionStatus = "trialing"
	CloudSubscriptionStatusPastDue  CloudSubscriptionStatus = "past_due"
	CloudSubscriptionStatusCanceled CloudSubscriptionStatus = "canceled"
	CloudSubscriptionStatusPaused   CloudSubscriptionStatus = "paused"
)

type Features struct {
	BuildMinutes        int    `json:"buildMinutes"`
	MaxBuildMinutes     int    `json:"maxBuildMinutes"`
	BuildMinutesPriceID string `json:"buildMinutesPriceID"`

	Committers        int    `json:"committers"`
	MaxCommitters     int    `json:"maxCommitters"`
	CommittersPriceID string `json:"committersPriceID"`

	MaxProjects int `json:"maxProjects"`

	MaxCloudRunnerScale int `json:"maxCloudRunnerScale"`
}

type PaddleSubscriptionItem struct {
	ProductID string `json:"productID"`
	PriceID   string `json:"priceID"`
	Quantity  int    `json:"quantity"`
}

type PaddleData struct {
	CustomerID                 string                              `json:"customerID"`
	LastNotificationOccurredAt time.Time                           `json:"lastNotificationOccurredAt"`
	SubscriptionID             string                              `json:"subscriptionID"`
	SubscriptionStatus         paddle.SubscriptionStatus           `json:"subscriptionStatus"`
	SubscriptionItems          []PaddleSubscriptionItem            `json:"subscriptionItems"`
	CollectionMode             paddle.CollectionMode               `json:"collectionMode"`
	ScheduledChange            *paddle.SubscriptionScheduledChange `json:"scheduledChange"`
	NextBillingAt              *time.Time                          `json:"nextBillingAt"`
	BillingPeriod              *paddle.TimePeriod                  `json:"billingPeriod"`
}

type CloudLicense struct {
	PaddleData *PaddleData           `json:"paddleData"`
	Plan       CloudSubscriptionPlan `json:"plan"`
	Features   Features              `json:"features"`
}

type BillingPeriod struct {
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}

func (l *CloudLicense) GetBillingPeriod() BillingPeriod {
	if l.PaddleData == nil || l.PaddleData.BillingPeriod == nil {
		return BillingPeriod{
			StartsAt: time.Now().AddDate(0, -1, 0),
			EndsAt:   time.Now(),
		}
	}

	startsAt, err := time.Parse(time.RFC3339, l.PaddleData.BillingPeriod.StartsAt)
	if err != nil {
		return BillingPeriod{}
	}

	endsAt, err := time.Parse(time.RFC3339, l.PaddleData.BillingPeriod.EndsAt)
	if err != nil {
		return BillingPeriod{}
	}

	return BillingPeriod{
		StartsAt: startsAt,
		EndsAt:   endsAt,
	}
}
