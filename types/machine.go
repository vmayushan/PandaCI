package types

import "fmt"

type CloudRunner string

const (
	CloudRunnerUbuntu1X  = "ubuntu-1x"
	CloudRunnerUbuntu2X  = "ubuntu-2x"
	CloudRunnerUbuntu4X  = "ubuntu-4x"
	CloudRunnerUbuntu8X  = "ubuntu-8x"
	CloudRunnerUbuntu16X = "ubuntu-16x"
)

func StringToCloudRunner(s string) (CloudRunner, error) {
	switch s {
	case "ubuntu-1x":
		return CloudRunnerUbuntu1X, nil
	case "ubuntu-2x":
		return CloudRunnerUbuntu2X, nil
	case "ubuntu-4x":
		return CloudRunnerUbuntu4X, nil
	case "ubuntu-8x":
		return CloudRunnerUbuntu8X, nil
	case "ubuntu-16x":
		return CloudRunnerUbuntu16X, nil
	default:
		return "", fmt.Errorf("unknown runner")
	}
}

func GetBuildMinutesScale(runner CloudRunner) int {
	switch runner {
	case CloudRunnerUbuntu1X:
		return 1
	case CloudRunnerUbuntu2X:
		return 2
	case CloudRunnerUbuntu4X:
		return 4
	case CloudRunnerUbuntu8X:
		return 8
	case CloudRunnerUbuntu16X:
		return 16
	default:
		return 1
	}
}
