package handlersProject

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pandaci-com/pandaci/app/api/middleware"
	middlewareLoaders "github.com/pandaci-com/pandaci/app/api/middleware/loaders"
	"github.com/pandaci-com/pandaci/pkg/encryption"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	utilsValidator "github.com/pandaci-com/pandaci/pkg/utils/validator"
	"github.com/pandaci-com/pandaci/platform/analytics"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	typesHTTP "github.com/pandaci-com/pandaci/types/http"
	"github.com/labstack/echo/v4"
	"github.com/posthog/posthog-go"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GetProjectVariables(c echo.Context) error {

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	variablesDB, err := h.queries.GetProjectVariablesWithEnvironments(c.Request().Context(), project)
	if err != nil {
		return err
	}

	variables := make([]typesHTTP.ProjectVariable, len(variablesDB))
	for i, variable := range variablesDB {
		projectVariable := typesHTTP.ProjectVariable{
			ID:           variable.ID,
			ProjectID:    variable.ProjectID,
			Key:          variable.Key,
			UpdatedAt:    variable.UpdatedAt,
			CreatedAt:    variable.CreatedAt,
			Sensitive:    variable.Sensitive,
			Environments: []typesHTTP.ProjectEnvironment{},
		}

		for _, env := range variable.Environments {
			projectVariable.Environments = append(projectVariable.Environments, typesHTTP.ProjectEnvironment{
				ID:            env.ID,
				ProjectID:     env.ProjectID,
				Name:          env.Name,
				UpdatedAt:     env.UpdatedAt,
				CreatedAt:     env.CreatedAt,
				BranchPattern: env.BranchPattern,
			})
		}

		variables[i] = projectVariable
	}

	return c.JSON(http.StatusOK, variables)
}

func (h *Handler) UpdateProjectVariable(c echo.Context) error {
	id := c.Param("variable_id")

	user := middleware.GetUser(c)

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	variableReq := typesHTTP.EditProjectVariableBody{}

	if err := c.Bind(&variableReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	validate := utilsValidator.NewValidator()
	if err := validate.Struct(variableReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utilsValidator.ValidatorErrors(err))
	}

	if strings.HasPrefix(variableReq.Key, "PANDACI_") {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Variable key cannot start with 'PANDACI_' prefix"})
	}

	variable, err := h.queries.GetProjectVariableByID(c.Request().Context(), project, id)
	if err != nil {
		return err
	}

	// Encrypt the variable value
	keyID, err := env.GetCurrentEncryptionKeyID()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current encryption key ID")
		return err
	}

	key, err := env.GetEncryptionKey(*keyID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get encryption key")
		return err
	}

	encrypted, iv, err := encryption.Encrypt(variableReq.Value, *key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt variable")
		return err
	}

	variable.Key = variableReq.Key
	variable.Value = encrypted
	variable.InitialisationVector = iv
	variable.EncryptionKeyID = *keyID
	variable.Sensitive = variableReq.Sensitive

	// Check for duplicate variables with same key in same environment
	currentVariables, err := h.queries.GetProjectVariablesWithEnvironments(c.Request().Context(), project)
	if err != nil {
		return err
	}

	for _, currentVariable := range currentVariables {
		// Skip the current variable being updated
		if currentVariable.ID == variable.ID {
			continue
		}

		if currentVariable.Key != variable.Key {
			continue
		}

		if len(currentVariable.Environments) == 0 && len(variableReq.ProjectEnvironmentIDs) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("A variable with the same key has already been created for the default environment")})
		}

		for _, env := range currentVariable.Environments {
			for _, reqEnvID := range variableReq.ProjectEnvironmentIDs {
				if env.ID == reqEnvID {
					return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("A variable with the same key is already attached to the %s environment", env.Name)})
				}
			}
		}
	}

	if err := h.queries.UpdateProjectVariable(c.Request().Context(), variable, variableReq.ProjectEnvironmentIDs); err != nil {
		return err
	}

	analytics.TrackUserProjectEvent(user, *project, posthog.Capture{
		Event: "project_variable_updated",
	})

	return c.NoContent(http.StatusOK)
}

func (h *Handler) DeleteProjectVariable(c echo.Context) error {
	id := c.Param("variable_id")

	user := middleware.GetUser(c)

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	if err := h.queries.DeleteProjectVariableByID(c.Request().Context(), *project, id); err != nil {
		return err
	}

	analytics.TrackUserProjectEvent(user, *project, posthog.Capture{
		Event: "project_variable_deleted",
	})

	return c.NoContent(http.StatusOK)
}

func (h *Handler) CreateProjectVariable(c echo.Context) error {

	user := middleware.GetUser(c)

	variableReq := typesHTTP.CreateProjectVariableBody{}

	if err := c.Bind(&variableReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	validate := utilsValidator.NewValidator()
	if err := validate.Struct(variableReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utilsValidator.ValidatorErrors(err))
	}

	if strings.HasPrefix(variableReq.Key, "PANDACI_") {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Variable key cannot start with 'PANDACI_' prefix"})
	}

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	variable := typesDB.ProjectVariable{
		ProjectID: project.ID,
		Key:       variableReq.Key,
		Sensitive: variableReq.Sensitive,
	}

	// Encrypt the variable value
	keyID, err := env.GetCurrentEncryptionKeyID()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current encryption key ID")
		return err
	}

	key, err := env.GetEncryptionKey(*keyID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get encryption key")
		return err
	}
	encrypted, iv, err := encryption.Encrypt(variableReq.Value, *key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt variable")
		return err
	}

	variable.Value = encrypted
	variable.InitialisationVector = iv
	variable.EncryptionKeyID = *keyID

	currentVariables, err := h.queries.GetProjectVariablesWithEnvironments(c.Request().Context(), project)
	if err != nil {
		return err
	}

	// Check if we already have an existing variable with the same key & same environment
	for _, currentVariable := range currentVariables {
		if currentVariable.Key != variable.Key {
			continue
		}

		if len(currentVariable.Environments) == 0 && len(variableReq.ProjectEnvironmentIDs) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("A variable with the same key has already been created for the default environment")})
		}

		for _, env := range currentVariable.Environments {
			for _, reqEnvID := range variableReq.ProjectEnvironmentIDs {
				if env.ID == reqEnvID {
					return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("A variable with the same key is already attached to the %s environment", env.Name)})
				}
			}
		}
	}

	if err := h.queries.CreateProjectVariable(c.Request().Context(), &variable, variableReq.ProjectEnvironmentIDs); err != nil {
		return err
	}

	analytics.TrackUserProjectEvent(user, *project, posthog.Capture{
		Event: "project_variable_created",
	})

	res := typesHTTP.ProjectVariable{
		ID:        variable.ID,
		ProjectID: variable.ProjectID,
		Key:       variable.Key,
		UpdatedAt: variable.UpdatedAt,
		CreatedAt: variable.CreatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProjectVariable(c echo.Context) error {
	id := c.Param("variable_id")
	decrypt := c.QueryParam("decrypt") == "true"

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	variable, err := h.queries.GetProjectVariableByID(c.Request().Context(), project, id)
	if err != nil {
		return err
	}

	res := typesHTTP.ProjectVariable{
		ID:        variable.ID,
		ProjectID: variable.ProjectID,
		Key:       variable.Key,
		UpdatedAt: variable.UpdatedAt,
		CreatedAt: variable.CreatedAt,
		Sensitive: variable.Sensitive,
	}

	if decrypt {

		if variable.Sensitive {
			return c.String(http.StatusForbidden, "cannot show the decrypted value of a sensitive variable")
		}

		// Decrypt the variable value
		key, err := env.GetEncryptionKey(variable.EncryptionKeyID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get encryption key")
			return err
		}

		decrypted, err := encryption.Decrypt(variable.Value, variable.InitialisationVector, *key)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decrypt variable")
			return err
		}

		res.Value = decrypted
	}

	return c.JSON(http.StatusOK, res)
}
