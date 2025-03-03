package main

import (
	"encoding/base64"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/pandaci-com/pandaci/app/grpc"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	"github.com/pandaci-com/pandaci/pkg/stream"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	pbConnect "github.com/pandaci-com/pandaci/proto/go/v1/v1connect"
	"github.com/rs/zerolog/log"
)

func StartHandler(logs *ArrayWriter, workflowID string, orchestratorClient pbConnect.OrchestratorServiceClient, jwtHandler jwt.JWTHandler, config *pb.WorkflowRunnerInitConfig, port int, worstJobConclusion map[string]pb.Conclusion) error {

	stepLogs := make(map[string]*stream.LogStream)

	eventStream := stream.NewEventStream()

	log.Info().Msg("starting grpc server")

	e := echo.New()

	e.Debug = true // TODO - remove	this

	e.Use(middleware.Secure())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*.pandaci.com"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	if err := grpc.RegisterWorkflowRPC(e, jwtHandler, orchestratorClient, config.WorkflowMeta, worstJobConclusion, stepLogs, eventStream); err != nil {
		return err
	}

	streamGroup := e.Group("/v1/stream", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.NewHTTPError(401, "unauthorized")
			}

			workflowJWT := token[7:]

			if claims, err := jwtHandler.ValidateWorkflowToken(workflowJWT); err != nil {
				return echo.NewHTTPError(401, "unauthorized")
			} else if claims.WorkflowID != workflowID {
				return echo.NewHTTPError(401, "unauthorized")
			}

			return next(c)
		}
	})

	// streamGroup.GET("/events", func(c echo.Context) error {
	// 	w := c.Response()
	// 	w.Header().Set("Content-Type", "text/event-stream")
	// 	w.Header().Set("Cache-Control", "no-cache")
	// 	w.Header().Set("Connection", "keep-alive")

	// 	w.Flush()

	// 	eventChan := eventStream.Subscribe()
	// 	defer eventStream.Unsubscribe(eventChan)

	// 	for {
	// 		select {
	// 		case event := <-eventChan:
	// 			_, err := fmt.Fprintf(w.Writer, "data: %s\n\n", event)
	// 			if err != nil {
	// 				log.Error().Err(err).Msg("writing event")
	// 				return fmt.Errorf("error writing event")
	// 			}
	// 			w.Flush()
	// 		case <-c.Request().Context().Done():
	// 			{
	// 				log.Info().Msg("client disconnected")
	// 				return nil
	// 			}
	// 		}
	// 	}
	// })

	streamGroup.GET("/logs", func(c echo.Context) error {
		w := c.Response()
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		w.Flush()

		var logChan chan []string
		if c.QueryParam("step_id") != "" {
			logChan = stepLogs[c.QueryParam("step_id")].Subscribe()
			defer stepLogs[c.QueryParam("step_id")].Unsubscribe(logChan)
		} else {
			logChan = logs.logStream.Subscribe()
			defer logs.logStream.Unsubscribe(logChan)
		}

		for {
			select {
			case entries, ok := <-logChan:
				if !ok {
					return nil // channel closed
				}
				for _, entry := range entries {
					data := base64.StdEncoding.EncodeToString([]byte(entry))
					_, err := fmt.Fprintf(w.Writer, "data: %s\n\n", data)
					if err != nil {
						log.Error().Err(err).Msg("writing log")
						return c.String(500, "error writing log")
					}
				}
				w.Flush()
			case <-c.Request().Context().Done():
				{
					log.Info().Msg("client disconnected")
					return nil
				}
			}
		}
	})

	go func() {
		if err := e.Start(fmt.Sprintf("%s:%d", config.Host, port)); err != nil {
			log.Fatal().Err(err).Msg("starting echo server")
		}
	}()

	log.Info().Msg("grpc server started")

	return nil
}
