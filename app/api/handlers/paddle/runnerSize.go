package handlersPaddle

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/alfiejones/panda-ci/types"
)

func getMaxRunnerScale(plan types.CloudSubscriptionPlan, subscription *paddle.Subscription) int {

	if plan == types.CloudSubscriptionPlanPaused {
		// We don't have a valid subscription and this isn't an org which supports the
		// free plan, so we should return 0
		return 0
	}

	if plan == types.CloudSubscriptionPlanPro {
		return 8
	}

	if plan == types.CloudSubscriptionPlanEnterprise {
		return 16
	}

	return 4
}
