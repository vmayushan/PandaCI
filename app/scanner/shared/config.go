package scannerShared

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/alfiejones/panda-ci/types"
	"github.com/rs/zerolog/log"
)

func ExtractWorkflowConfig(content []byte) (*types.WorkflowConfig, error) {
	rawConfig, err := extractWorkflowConfig(content)
	if err != nil {
		log.Error().Err(err).Msg("extracting config")
		return nil, err
	}

	return &types.WorkflowConfig{
		Config: *rawConfig,
	}, nil
}

func CleanJSON(input string) string {
	// Remove single-line comments (// ...)
	reSingleLine := regexp.MustCompile(`\/\/.*`)
	input = reSingleLine.ReplaceAllString(input, "")

	// Remove multi-line comments (/* ... */)
	reMultiLine := regexp.MustCompile(`\/\*[\s\S]*?\*\/`)
	input = reMultiLine.ReplaceAllString(input, "")

	// Remove trailing commas (before closing braces/brackets)
	reTrailingComma := regexp.MustCompile(`,\s*([\]}])`)
	input = reTrailingComma.ReplaceAllString(input, "$1")

	// Regular expression to match unquoted keys and add quotes around them
	re := regexp.MustCompile(`\b([a-zA-Z0-9_]+)\b\s*:`)

	// Replace matches with the quoted key
	input = re.ReplaceAllString(input, `"$1":`)

	// reNewLine := regexp.MustCompile(`/\r?\n/`)
	// input = reNewLine.ReplaceAllString(input, "")

	return strings.TrimSpace(input)
}

func extractWorkflowConfig(content []byte) (*types.WorkflowRawConfig, error) {
	startReg := `export\s+(?:const|let|var)\s+config(?:\s*:\s*\w+)?\s*=\s*`
	reg := GenerateRegex(types.WorkflowRawConfig{})
	fullReg := startReg + reg

	log.Info().Str("fullReg", fullReg).Msg("extractWorkflowConfig")

	configExists, err := regexp.Compile(fullReg)
	if err != nil {
		log.Error().Err(err).Msg("compiling configExists regexp")
		return nil, err
	}

	fullConfig := configExists.Find(content)
	if fullConfig == nil || len(fullConfig) == 0 {
		log.Info().Msg("no config object found")
		return &types.WorkflowRawConfig{}, nil
	}

	log.Info().Str("fullConfig", string(fullConfig)).Msg("extractWorkflowConfig")

	configContentRE, err := regexp.Compile(reg)
	if err != nil {
		log.Error().Err(err).Msg("compiling configContentRE regexp")
		return nil, err
	}

	match := configContentRE.Find(fullConfig)
	match = []byte(CleanJSON(string(match)))

	if match == nil || len(match) == 0 {
		// A config is defined but it doesn't match our regexp
		return nil, fmt.Errorf("bad config object")
	}

	config := &types.WorkflowRawConfig{}
	if err := json.Unmarshal(match, config); err != nil {
		return nil, err
	}

	return config, nil
}
