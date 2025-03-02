package typesHTTP

import "time"

type ProjectVariable struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectID"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	Sensitive bool      `json:"sensitive"`

	Environments []ProjectEnvironment `json:"environments"`
}

type ProjectEnvironment struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectID"`
	Name          string    `json:"name"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
	BranchPattern string    `json:"branchPattern"`
}

type CreateProjectVariableBody struct {
	Key                   string   `json:"key" validate:"required,min=1,max=255"`
	Value                 string   `json:"value" validate:"required,min=1,max=65536"`
	ProjectEnvironmentIDs []string `json:"environmentIDs,omitempty"`
	Sensitive             bool     `json:"sensitive"`
}

type CreateProjectEnvironmentBody struct {
	Name          string `json:"name" validate:"required,min=1,max=255"`
	BranchPattern string `json:"branchPattern" validate:"required,min=1,max=255"`
}
