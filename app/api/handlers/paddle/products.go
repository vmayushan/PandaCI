package handlersPaddle

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/rs/zerolog/log"

	"github.com/pandaci-com/pandaci/pkg/utils/env"
	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
)

func getProProductID() string {
	stage := env.GetStage()

	if stage == "prod" {
		return "pro_01jkwz26eqajmca4e5xthx70c6"
	}

	return "pro_01jhzxfy5m5f1rbfs5yx1nej1q"
}

func getBuildMinutesPriceID() string {
	stage := env.GetStage()

	if stage == "prod" {
		return "pri_01jkwzam63n615aw0bhp46f2mf"
	}

	return "pri_01jjt5qw7m0pwjp0cx6scb2m3p"
}

func getBuildMinutesProductID() string {
	stage := env.GetStage()

	if stage == "prod" {
		return "pro_01jkwz3h1pzscxw4qqd7p0689g"
	}

	return "pro_01jjt4qzynkrgstad0q2b3vjez"
}

func getCommitterProductID() string {
	stage := env.GetStage()

	if stage == "prod" {
		return "pro_01jmx3863hfq453vdy4rp8v7kn"
	}

	return "pro_01jndtr8mxb16e3mnsz2yny4d1"
}

func getCommitterPriceID() string {
	stage := env.GetStage()

	if stage == "prod" {
		return "pri_01jmx3cbwsztpvpwp16rf8nkje"
	}

	return "pri_01jndts9h69jhehywdd8805w5x"
}

func (h *Handler) isPausedLicense(ctx context.Context, org typesDB.OrgDB, subscription *paddle.Subscription) bool {

	oldLicense, err := org.GetLicense()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get license")
		return false
	}

	validStatuses := []paddle.SubscriptionStatus{
		paddle.SubscriptionStatusActive,
		paddle.SubscriptionStatusTrialing,
		paddle.SubscriptionStatusPastDue,
	}

	isInvalidLicense := oldLicense.Plan == types.CloudSubscriptionPlanPaused
	isInvalidStatus := !slices.Contains(validStatuses, subscription.Status)

	if isInvalidLicense && isInvalidStatus {

		ownersFreeOrgs, err := h.queries.GetUsersFreeOrg(ctx, oldLicense.PaddleData.CustomerID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Msg("Failed to get users free org")
			return false
		}

		if ownersFreeOrgs != nil {
			// Owner already has a free org, so this one requires a paid plan
			return true
		}
	}

	return false
}

func (h *Handler) getPlan(ctx context.Context, org typesDB.OrgDB, subscription *paddle.Subscription) types.CloudSubscriptionPlan {

	if h.isPausedLicense(ctx, org, subscription) {
		return types.CloudSubscriptionPlanPaused
	}

	if subscription == nil {
		return types.CloudSubscriptionPlanFree
	}

	proProductID := getProProductID()

	for _, item := range subscription.Items {
		if item.Product.ID == proProductID {
			return types.CloudSubscriptionPlanPro
		}
	}

	return types.CloudSubscriptionPlanFree
}
