package handlersPaddle

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/alfiejones/panda-ci/types"
)

func getBuildMinutes(plan types.CloudSubscriptionPlan, subscription *paddle.Subscription) int {

	if plan == types.CloudSubscriptionPlanPaused {
		// We don't have a valid subscription and this isn't an org which supports the
		// free plan, so we should return 0
		return 0
	}

	productID := getBuildMinutesProductID()

	total := 6000
	if subscription.NextTransaction == nil {
		return total
	}

	for _, item := range subscription.NextTransaction.Details.LineItems {
		if item.Product.ID != nil && *item.Product.ID == productID {
			total += item.Quantity * 500
		}
	}

	return total
}

func getMaxBuildMinutes(plan types.CloudSubscriptionPlan, subscription *paddle.Subscription) int {
	if plan == types.CloudSubscriptionPlanPaused {
		// We don't have a valid subscription and this isn't an org which supports the
		// free plan, so we should return 0
		return 0
	}

	planDefault := 6000

	if plan == types.CloudSubscriptionPlanPro || plan == types.CloudSubscriptionPlanEnterprise {
		planDefault = 99999999
	}

	return planDefault
}
