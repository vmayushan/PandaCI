package grpcOrchestrator

import (
	"context"
	"math"
	"slices"

	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/alfiejones/panda-ci/types"
	typesDB "github.com/alfiejones/panda-ci/types/database"
	"github.com/rs/zerolog/log"
)

func (h *Handler) chargeForMinutes(ctx context.Context, org *typesDB.OrgDB) error {
	log.Debug().Str("orgID", org.ID).Msg("Charging for minutes")

	license, err := org.GetLicense()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get license")
		return err
	}

	if slices.Contains([]types.CloudSubscriptionPlan{types.CloudSubscriptionPlanPaused, types.CloudSubscriptionPlanFree}, license.Plan) {
		return nil
	}

	billingPeriod := license.GetBillingPeriod()

	buildMinutes, err := h.queries.GetBuildMinutes(ctx, org.ID, billingPeriod.StartsAt, billingPeriod.EndsAt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get build minutes")
		return err
	}

	if buildMinutes <= license.Features.BuildMinutes || buildMinutes >= license.Features.MaxBuildMinutes {
		// We are covered by bought minutes or we are over our limit, eitherway we don't need to charge
		return nil
	}

	quantity := int(math.Round((float64(buildMinutes-license.Features.BuildMinutes) / 500.0) + 0.5))

	// we want to update the license with the new build minutes asap to avoid double charging
	license.Features.BuildMinutes = buildMinutes + (quantity * 500)
	if err := org.SetLicense(license); err != nil {
		log.Error().Err(err).Msg("Failed to set license")
		return err
	}

	if err := h.queries.UpdateOrgLicense(ctx, org); err != nil {
		log.Error().Err(err).Msg("Failed to update org license")
		return err
	}

	if _, err := h.paddleClient.CreateSubscriptionCharge(ctx, &paddle.CreateSubscriptionChargeRequest{
		SubscriptionID: license.PaddleData.SubscriptionID,
		EffectiveFrom:  paddle.EffectiveFromNextBillingPeriod,
		Items: []paddle.CreateSubscriptionChargeItems{
			*paddle.NewCreateSubscriptionChargeItemsSubscriptionChargeItemFromCatalog(
				&paddle.SubscriptionChargeItemFromCatalog{
					PriceID:  license.Features.BuildMinutesPriceID,
					Quantity: quantity,
				}),
		},
	}); err != nil {
		log.Error().Err(err).Msg("Failed to create subscription charge")
		return err
	}

	return nil
}

func (h *Handler) chargeForCommitters(ctx context.Context, org *typesDB.OrgDB) error {
	log.Debug().Str("orgID", org.ID).Msg("Charging for committers")

	license, err := org.GetLicense()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get license")
		return err
	}

	billingPeriod := license.GetBillingPeriod()

	committersCount, err := h.queries.CountCommitters(ctx, org.ID, billingPeriod.StartsAt, billingPeriod.EndsAt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get committers")
		return err
	}

	if committersCount <= license.Features.Committers || committersCount >= license.Features.MaxCommitters {
		// We are covered by bought committers or we are over our limit, eitherway we don't need to charge
		return nil
	}

	quantity := committersCount - license.Features.Committers

	// we want to update the license with the new committers asap to avoid double charging
	license.Features.Committers = committersCount
	if err := org.SetLicense(license); err != nil {
		log.Error().Err(err).Msg("Failed to set license")
		return err
	}

	if err := h.queries.UpdateOrgLicense(ctx, org); err != nil {
		log.Error().Err(err).Msg("Failed to update org license")
		return err
	}

	if _, err := h.paddleClient.CreateSubscriptionCharge(ctx, &paddle.CreateSubscriptionChargeRequest{
		SubscriptionID: license.PaddleData.SubscriptionID,
		EffectiveFrom:  paddle.EffectiveFromNextBillingPeriod,
		Items: []paddle.CreateSubscriptionChargeItems{
			*paddle.NewCreateSubscriptionChargeItemsSubscriptionChargeItemFromCatalog(
				&paddle.SubscriptionChargeItemFromCatalog{
					PriceID:  license.Features.CommittersPriceID,
					Quantity: quantity,
				}),
		},
	}); err != nil {
		log.Error().Err(err).Msg("Failed to create subscription charge")
		return err
	}

	return nil
}
