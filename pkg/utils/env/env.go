package env

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/alfiejones/panda-ci/types"
	"github.com/rs/zerolog/log"
)

func getEnvWithFallback(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvWithError(key string, fallback *string) (*string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return &value, nil
	}

	if fallback != nil {
		return fallback, nil
	}

	return nil, fmt.Errorf("Missing env: %s", key)
}

func getIntEnvWithError(key string, fallback *int) (*int, error) {
	var fallbackStr *string
	if fallback != nil {
		conv := strconv.Itoa(*fallback)
		fallbackStr = &conv
	}

	str, err := getEnvWithError(key, fallbackStr)
	if err != nil {
		return fallback, err
	}

	value, err := strconv.Atoi(*str)
	if err != nil {
		log.Err(err)
		// TODO - we should probably have a docs page about our env variables & we can add a link to this message
		return fallback, fmt.Errorf("Unable to convert env: %s to int", key)
	}

	return &value, nil
}

type Stage string

const (
	StageLocal     Stage = "local"
	StageDev       Stage = "dev"
	StageCommunity Stage = "community"
	StageProd      Stage = "prod"
)

func GetStage() Stage {
	fallback := string(StageDev)
	stage, _ := getEnvWithError("STAGE", &fallback)

	return Stage(*stage)
}

func GetPostgresDSN() (*string, error) {
	var fallback *string

	return getEnvWithError("DSN", fallback)
}

func GetOryURL() (*string, error) {
	var fallback *string

	stage := GetStage()

	if stage == StageDev {
		fallback = types.Pointer("http://pandaci-kratos-dev.internal:4433")
	} else if stage == StageLocal {
		local := "http://localhost:4433"
		fallback = &local
	} else if stage == StageProd {
		fallback = types.Pointer("http://pandaci-kratos-prod.internal:4433")
	}

	return getEnvWithError("ORY_URL", fallback)
}

func GetOryAdminURL() (*string, error) {
	var fallback *string

	stage := GetStage()

	if stage == StageDev {
		fallback = types.Pointer("http://pandaci-kratos-dev.internal:4434")
	} else if stage == StageLocal {
		local := "http://localhost:4434"
		fallback = &local
	} else if stage == StageProd {
		fallback = types.Pointer("http://pandaci-kratos-prod.internal:4434")
	}

	return getEnvWithError("ORY_ADMIN_URL", fallback)
}

func GetOryAdminToken() string {
	fallback := ""
	val, _ := getEnvWithError("ORY_ADMIN_TOKEN", &fallback)
	return *val
}

func GetGithubAppID() (*int, error) {
	return getIntEnvWithError("GITHUB_APP_ID", nil)
}

func GetGithubAppPrivateKey() (*string, error) {
	str, err := getEnvWithError("GITHUB_APP_PRIVATE_KEY_BASE64", nil)
	if err != nil {
		return nil, err
	}

	envBytes, err := base64.StdEncoding.DecodeString(*str)
	return types.Pointer(string(envBytes)), err
}

func GetGithubAppClientSecret() (*string, error) {
	return getEnvWithError("GITHUB_APP_CLIENT_SECRET", nil)
}

func GetGithubAppClientID() (*string, error) {
	return getEnvWithError("GITHUB_APP_CLIENT_ID", nil)
}

func GetGithubAPIEndpoint() string {
	return getEnvWithFallback("GITHUB_API_ENDPOINT", "https://api.github.com")
}

func GetAPIHost() (*string, error) {
	var fallback *string

	stage := GetStage()

	if stage == StageDev {
		fallback = types.Pointer("0.0.0.0:5000")
	} else if stage == StageProd {
		fallback = types.Pointer("0.0.0.0:5000")
	} else if stage == StageLocal {
		fallback = types.Pointer("localhost:5000")
	}

	return getEnvWithError("API_HOST", fallback)
}

func GetBackendURL() (*string, error) {
	var fallback *string

	stage := GetStage()

	if stage == StageDev {
		branch, err := getDevBranch()
		if err != nil {
			return nil, err
		}
		fallback = types.Pointer(fmt.Sprintf("https://%s.api.dev.pandaci.com", *branch))
	} else if stage == StageProd {
		fallback = types.Pointer("https://api.pandaci.com")
	} else if stage == StageLocal {
		local := "http://localhost:5000"
		fallback = &local
	}

	return getEnvWithError("BACKEND_URL", fallback)
}

func getDevBranch() (*string, error) {
	return getEnvWithError("DEV_BRANCH", nil)
}

func GetFrontendURL() (*string, error) {
	var fallback *string

	stage := GetStage()

	if stage == StageDev {
		branch, err := getDevBranch()
		if err != nil {
			return nil, err
		}
		fallback = types.Pointer(fmt.Sprintf("https://%s.app.dev.pandaci.com", *branch))
	} else if stage == StageLocal {
		local := "http://localhost:5173"
		fallback = &local
	} else if stage == StageProd {
		fallback = types.Pointer("https://app.pandaci.com")
	}

	return getEnvWithError("FRONTEND_URL", fallback)
}

func GetOrchestratorGRPCURL() (*string, error) {
	var fallback *string

	stage := GetStage()

	if stage == StageDev {
		branch, err := getDevBranch()
		if err != nil {
			return nil, err
		}

		fallback = types.Pointer(fmt.Sprintf("https://%s.api.dev.pandaci.com/grpc", *branch))
	} else if stage == StageProd {
		fallback = types.Pointer("https://api.pandaci.com/grpc")
	} else {
		apiHost, err := GetAPIHost()
		if err != nil {
			return nil, err
		}

		local := fmt.Sprintf("http://%s/grpc", *apiHost)
		fallback = &local
	}

	return getEnvWithError("ORCHESTRATOR_GRPC_URL", fallback)
}

func GetRunnerAddress() (*string, error) {
	var fallback *string

	stage := GetStage()

	if stage == StageDev {
		fallback = types.Pointer("https://fly-runner.dev.pandaci.com/grpc")
	} else if stage == StageProd {
		fallback = types.Pointer("https://fly-runner.pandaci.com/grpc")
	} else {
		fallback = types.Pointer("http://localhost:5000/grpc")
	}

	return getEnvWithError("RUNNER_ADDRESS", fallback)
}

func GetJobsWorkflowID() (*string, error) {
	// TODO - maybe we can get rid of this in favor of the pb.WorkflowRunnerInitConfig
	// Stores the current workflow id for the running job
	return getEnvWithError("JOBS_WORKFLOW_ID", nil)
}

func GetRunnerPublicKey() (*string, error) {
	str, err := getEnvWithError("RUNNER_PUBLIC_KEY_BASE64", nil)
	if err != nil {
		return nil, err
	}

	envBytes, err := base64.StdEncoding.DecodeString(*str)
	return types.Pointer(string(envBytes)), err
}

func GetRunnerPrivateKey() (*string, error) {
	str, err := getEnvWithError("RUNNER_PRIVATE_KEY_BASE64", nil)
	if err != nil {
		return nil, err
	}

	envBytes, err := base64.StdEncoding.DecodeString(*str)
	return types.Pointer(string(envBytes)), err
}

func GetAllowedOrigins() ([]string, error) {
	var fallback string

	stage := GetStage()

	if stage == StageLocal {
		fallback = `["http://localhost:5173", "http://localhost:5000"]`
	} else if stage == StageDev {
		branch, err := getDevBranch()
		if err != nil {
			return []string{}, err
		}

		fallback = fmt.Sprintf(`["https://%s.app.dev.pandaci.com", "https://%s.api.dev.pandaci.com"]`, *branch, *branch)
	} else if stage == StageProd {
		fallback = `["https://app.pandaci.com", "https://api.pandaci.com"]`
	}

	originsJson := getEnvWithFallback("ALLOWED_ORIGINS", fallback)

	var origins []string

	err := json.Unmarshal([]byte(originsJson), &origins)
	if err != nil {
		return []string{}, err
	}

	return origins, nil
}

func GetCurrentEncryptionKeyID() (*string, error) {
	return getEnvWithError("ENCRYPTION_KEY_ID", nil)
}

func GetEncryptionKey(id string) (*string, error) {
	return getEnvWithError("ENCRYPTION_KEY_"+id, nil)
}

func GetGithubWebhookSecret() (*string, error) {
	return getEnvWithError("GITHUB_APP_WEBHOOK_SECRET", nil)
}

func GetPaddleWebhookSecret() (*string, error) {
	return getEnvWithError("PADDLE_WEBHOOK_SECRET", nil)
}

func GetPaddleAPIKey() (*string, error) {
	return getEnvWithError("PADDLE_API_KEY", nil)
}

func GetPosthogAPIKey() (*string, error) {
	return getEnvWithError("POSTHOG_API_KEY", nil)
}

func GetOryWebhookAPIKey() (*string, error) {
	stage := GetStage()

	var fallback *string
	if stage == "local" {
		fallback = types.Pointer("PLEASE-CHANGE-ME-I-AM-VERY-INSECURE")
	}
	return getEnvWithError("ORY_WEBHOOK_API_KEY", fallback)
}
func GetSMTPHost() (*string, error) {
	return getEnvWithError("SMTP_HOST", nil)
}

func GetSMTPPort() (*string, error) {
	return getEnvWithError("SMTP_PORT", nil)
}

func GetSMTPUsername() (*string, error) {
	return getEnvWithError("SMTP_USERNAME", nil)
}

func GetSMTPPassword() (*string, error) {
	return getEnvWithError("SMTP_PASSWORD", nil)
}
