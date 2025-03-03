package handlersPaddle

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/pandaci-com/pandaci/types"
)

func getMaxProjects(plan types.CloudSubscriptionPlan, subscription *paddle.Subscription) int {

	if plan == types.CloudSubscriptionPlanPaused {
		// We don't have a valid subscription and this isn't an org which supports the
		// free plan, so we should return 0
		return 0
	}

	if plan == types.CloudSubscriptionPlanPro {
		return 100
	}

	if plan == types.CloudSubscriptionPlanEnterprise {
		// We probably won't ever see this many but we needed a number
		return 999999
	}

	return 10
}
