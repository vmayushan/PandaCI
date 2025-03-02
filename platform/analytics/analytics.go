package analytics

import (
	"github.com/alfiejones/panda-ci/pkg/utils/env"
	"github.com/alfiejones/panda-ci/types"
	typesDB "github.com/alfiejones/panda-ci/types/database"
	"github.com/posthog/posthog-go"
	"github.com/rs/zerolog/log"
)

func trackEvent(message posthog.Capture) {
	apiKey, err := env.GetPosthogAPIKey()
	if err != nil {
		log.Error().Err(err).Msg("Error getting posthog api key")
		return
	}

	client, err := posthog.NewWithConfig(*apiKey, posthog.Config{Endpoint: "https://us.i.posthog.com"})
	if err != nil {
		log.Error().Err(err).Msg("Error creating posthog client")
		return
	}

	defer client.Close()

	if err := client.Enqueue(message); err != nil {
		log.Error().Err(err).Msg("Error sending analytics event")
	}
}

func TrackUserOrgEvent(user types.User, orgID string, message posthog.Capture) {
	message.DistinctId = user.ID

	message.Groups = posthog.NewGroups().
		Set("org", orgID)

	trackEvent(message)
}

func TrackOrgEvent(orgID string, message posthog.Capture) {
	message.Groups = posthog.NewGroups().
		Set("org", orgID)

	trackEvent(message)
}

func TrackUserProjectEvent(user types.User, project typesDB.Project, message posthog.Capture) {
	message.DistinctId = user.ID

	message.Groups = posthog.NewGroups().
		Set("org", project.OrgID).
		Set("project", project.ID)

	trackEvent(message)
}
