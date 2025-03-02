package handlersPaddle

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/alfiejones/panda-ci/types"
)

func getCommitters(plan types.CloudSubscriptionPlan, subscription *paddle.Subscription) int {

	if plan == types.CloudSubscriptionPlanPaused {
		// We don't have a valid subscription and this isn't an org which supports the
		// free plan, so we should return 0
		return 0
	}

	productID := getCommitterProductID()

	total := 5

	if subscription.NextTransaction != nil {
		for _, item := range subscription.NextTransaction.Details.LineItems {
			if item.Product.ID != nil && *item.Product.ID == productID {
				total += item.Quantity
			}
		}
	}

	return total
}

func getMaxCommitters(plan types.CloudSubscriptionPlan, subscription *paddle.Subscription) int {

	if plan == types.CloudSubscriptionPlanPaused {
		// We don't have a valid subscription and this isn't an org which supports the
		// free plan, so we should return 0
		return 0
	}

	planDefault := 5

	if plan == types.CloudSubscriptionPlanPro || plan == types.CloudSubscriptionPlanEnterprise {
		planDefault = 999999
	}

	return planDefault
}
