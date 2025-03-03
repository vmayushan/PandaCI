package handlersPaddle

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/PaddleHQ/paddle-go-sdk/v3/pkg/paddlenotification"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	"github.com/pandaci-com/pandaci/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) HandleWebhook(c echo.Context) error {

	paddleWebhookSecret, err := env.GetPaddleWebhookSecret()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get paddle webhook secret")
		return err
	}

	verifier := paddle.NewWebhookVerifier(*paddleWebhookSecret)

	if ok, err := verifier.Verify(c.Request()); err != nil {
		log.Error().Err(err).Msg("Failed to verify paddle webhook")
		return err
	} else if !ok {
		log.Error().Msg("Failed to verify paddle webhook")
		return echo.NewHTTPError(401, "Failed to verify paddle webhook")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read body")
		return err
	}

	c.Request().Body = io.NopCloser(bytes.NewReader(body))

	var webhook paddlenotification.GenericNotificationEvent
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&webhook); err != nil {
		log.Error().Err(err).Msg("Failed to decode paddle webhook")
		return err
	}

	switch webhook.EventType {
	case paddlenotification.EventTypeNameSubscriptionCreated:
		return h.handleSubscriptionCreated(c)
	case paddlenotification.EventTypeNameSubscriptionUpdated:
		return h.handleSubscriptionUpdated(c)
	}

	return c.NoContent(200)
}

func (h *Handler) handleSubscriptionCreated(c echo.Context) error {
	body := paddlenotification.SubscriptionCreated{}
	if err := c.Bind(&body); err != nil {
		log.Error().Err(err).Msg("Failed to bind paddle webhook")
		return err
	}

	if err := h.handleSubscriptionData(c.Request().Context(), body.OccurredAt, body.Data.ID); err != nil {
		return err
	}

	return c.NoContent(200)
}

func (h *Handler) handleSubscriptionUpdated(c echo.Context) error {
	body := paddlenotification.SubscriptionUpdated{}
	if err := c.Bind(&body); err != nil {
		log.Error().Err(err).Msg("Failed to bind paddle webhook")
		return err
	}

	if err := h.handleSubscriptionData(c.Request().Context(), body.OccurredAt, body.Data.ID); err != nil {
		return err
	}

	return c.NoContent(200)
}

func (h *Handler) handleSubscriptionData(ctx context.Context, occuredAt string, subscriptionID string) error {

	occuredAtTime, err := time.Parse(time.RFC3339, occuredAt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse occured at")
		return err
	}

	subscription, err := h.paddleClient.GetSubscription(ctx, &paddle.GetSubscriptionRequest{
		SubscriptionID:         subscriptionID,
		IncludeNextTransaction: true,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get subscription")
		return err
	}

	orgID, ok := subscription.CustomData["orgID"]
	if !ok {
		log.Error().Msg("OrgId not found in custom data")
		return echo.NewHTTPError(400, "OrgId not found in custom data")
	}

	org, err := h.queries.Unsafe_GetOrgByID(ctx, orgID.(string))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get org")
		return err
	}

	oldLicense, err := org.GetLicense()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get license")
	}

	if oldLicense.PaddleData != nil && oldLicense.PaddleData.CustomerID != subscription.CustomerID && oldLicense.PaddleData.SubscriptionStatus != paddle.SubscriptionStatusCanceled {
		log.Error().Msg("CustomerID does not match, cancelling subscription")
		// We already have a customerID for this org.
		// This should only happen if someone is trying to buy a second subscription for an org.
		if _, err := h.paddleClient.CancelSubscription(ctx, &paddle.CancelSubscriptionRequest{
			SubscriptionID: subscriptionID,
		}); err != nil {
			return err
		}

		return nil
	}

	var nextBillingAt *time.Time
	if subscription.NextBilledAt != nil {
		nextBillingTime, err := time.Parse(time.RFC3339, *subscription.NextBilledAt)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse next billed at")
			return err
		}

		nextBillingAt = &nextBillingTime
	}

	subscriptionItems := make([]types.PaddleSubscriptionItem, len(subscription.Items))
	for i, item := range subscription.Items {
		subscriptionItems[i] = types.PaddleSubscriptionItem{
			ProductID: item.Product.ID,
			PriceID:   item.Price.ID,
			Quantity:  item.Quantity,
		}
	}

	plan := h.getPlan(ctx, *org, subscription)

	newLicense := types.CloudLicense{
		Plan: plan,
		PaddleData: &types.PaddleData{
			CustomerID:                 subscription.CustomerID,
			SubscriptionID:             subscription.ID,
			LastNotificationOccurredAt: occuredAtTime,
			SubscriptionStatus:         subscription.Status,
			BillingPeriod:              subscription.CurrentBillingPeriod,
			NextBillingAt:              nextBillingAt,
			ScheduledChange:            subscription.ScheduledChange,
			CollectionMode:             subscription.CollectionMode,
			SubscriptionItems:          subscriptionItems,
		},
		Features: types.Features{
			BuildMinutes:        getBuildMinutes(plan, subscription),
			MaxBuildMinutes:     getMaxBuildMinutes(plan, subscription),
			Committers:          getCommitters(plan, subscription),
			MaxCommitters:       getMaxCommitters(plan, subscription),
			MaxCloudRunnerScale: getMaxRunnerScale(plan, subscription),
			MaxProjects:         getMaxProjects(plan, subscription),
			BuildMinutesPriceID: getBuildMinutesPriceID(),
			CommittersPriceID:   getCommitterPriceID(),
		},
	}

	if err := org.SetLicense(newLicense); err != nil {
		log.Error().Err(err).Msg("Failed to set license")
		return err
	}

	if err := h.queries.UpdateOrgLicense(ctx, org); err != nil {
		log.Error().Err(err).Msg("Failed to update org license")
		return err
	}

	return nil
}
