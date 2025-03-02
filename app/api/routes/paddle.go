package routes

import (
	handlersPaddle "github.com/alfiejones/panda-ci/app/api/handlers/paddle"
	"github.com/alfiejones/panda-ci/app/api/middleware"
	middleware_loaders "github.com/alfiejones/panda-ci/app/api/middleware/loaders"
	"github.com/alfiejones/panda-ci/app/queries"
	"github.com/labstack/echo/v4"
)

func RegisterPaddleRoutes(e *echo.Echo, queries *queries.Queries) error {

	paddleHandler, err := handlersPaddle.NewHandler(queries)
	if err != nil {
		return err
	}

	v1 := e.Group("/v1/paddle")

	v1.POST("/webhook", paddleHandler.HandleWebhook)

	loaderMiddleware := middleware_loaders.New(queries)
	org := v1.Group("/orgs/:org_url_name", middleware.NewOryMiddleware().Session, loaderMiddleware.LoadOrg)

	org.POST("/portal", paddleHandler.HandlePortalRequest)

	return nil
}
