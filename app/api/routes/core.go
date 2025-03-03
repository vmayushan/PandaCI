package routes

import (
	"github.com/pandaci-com/pandaci/app/api/handlers"
	handlersProject "github.com/pandaci-com/pandaci/app/api/handlers/project"
	"github.com/pandaci-com/pandaci/app/api/middleware"
	middleware_loaders "github.com/pandaci-com/pandaci/app/api/middleware/loaders"
	"github.com/pandaci-com/pandaci/app/git"
	"github.com/pandaci-com/pandaci/app/orchestrator"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/app/scanner"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	"github.com/pandaci-com/pandaci/platform/storage"
	"github.com/labstack/echo/v4"
)

func RegisterOrgRoutes(e *echo.Echo, queries *queries.Queries, gitHandler *git.Handler, scanner scanner.Handler, orchestrator orchestrator.Handler, bucketClient *storage.BucketClient, jwtHandler *jwt.JWTHandler) {
	orgHandler := handlers.NewHandler(queries)
	projectHandler := handlersProject.NewHandler(queries, gitHandler, scanner, orchestrator, bucketClient, jwtHandler)

	v1 := e.Group("/v1/orgs")
	v1.Use(middleware.NewOryMiddleware().Session)

	v1.POST("", orgHandler.CreateOrg)
	v1.GET("", orgHandler.GetCurrentUserOrgs)

	loaderMiddleware := middleware_loaders.New(queries)

	namedOrg := v1.Group("/:org_url_name", loaderMiddleware.LoadOrg)

	namedOrg.GET("", orgHandler.GetOrg)
	namedOrg.PUT("", orgHandler.UpdateOrg)
	namedOrg.DELETE("", orgHandler.DeleteOrg)

	namedOrg.GET("/usage", orgHandler.GetOrgUsage)

	namedOrg.POST("/users", orgHandler.InviteUserToOrg)
	namedOrg.GET("/users", orgHandler.GetOrgUsers)
	namedOrg.DELETE("/users/:user_id", orgHandler.RemoveUserFromOrg)

	namedOrg.GET("/projects", projectHandler.GetOrgProjects)
	namedOrg.POST("/projects", projectHandler.CreateProject)

	namedProject := namedOrg.Group("/projects/:project_url_name", loaderMiddleware.LoadProject)

	namedProject.GET("", projectHandler.GetProject)
	namedProject.PUT("", projectHandler.UpdateProject)
	namedProject.DELETE("", projectHandler.DeleteProject)
	namedProject.POST("/trigger", projectHandler.TriggerRun)
	namedProject.GET("/runs", projectHandler.GetWorkflowRuns)

	namedProject.POST("/environments", projectHandler.CreateProjectEnvironment)
	namedProject.GET("/environments", projectHandler.GetProjectEnvironments)
	namedProject.DELETE("/environments/:environment_id", projectHandler.DeleteProjectEnvironment)
	namedProject.PUT("/environments/:environment_id", projectHandler.UpdateProjectEnvironment)

	namedProject.GET("/variables", projectHandler.GetProjectVariables)
	namedProject.POST("/variables", projectHandler.CreateProjectVariable)
	namedProject.GET("/variables/:variable_id", projectHandler.GetProjectVariable)
	namedProject.DELETE("/variables/:variable_id", projectHandler.DeleteProjectVariable)
	namedProject.PUT("/variables/:variable_id", projectHandler.UpdateProjectVariable)

	workflowIDRun := namedProject.Group("/run/:run_id")
	workflowIDRun.Use(loaderMiddleware.LoadWorkflowRun)
	workflowIDRun.GET("/logs/:log_id", projectHandler.GetWorkflowRunLogs)
	workflowIDRun.GET("/stream/logs", projectHandler.GetLogStream)

	// TODO - move these into their own route handler as this is getting a little crazy

	workflowRun := namedProject.Group("/runs/:run_number")
	workflowRun.Use(loaderMiddleware.LoadWorkflowRun)

	workflowRun.GET("", projectHandler.GetWorkflowRunWithItems)
}
