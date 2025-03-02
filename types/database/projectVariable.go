package typesDB

import (
	"time"
)

type ProjectEnvironment struct {
	ID            string    `db:"id"`
	ProjectID     string    `db:"project_id"`
	Name          string    `db:"name"`
	UpdatedAt     time.Time `db:"updated_at"`
	CreatedAt     time.Time `db:"created_at"`
	BranchPattern string    `db:"branch_pattern"`
}

type ProjectVariable struct {
	ID                   string    `db:"id"`
	ProjectID            string    `db:"project_id"`
	Key                  string    `db:"key"`
	Value                string    `db:"value"`
	UpdatedAt            time.Time `db:"updated_at"`
	CreatedAt            time.Time `db:"created_at"`
	EncryptionKeyID      string    `db:"encryption_key_id"`
	InitialisationVector string    `db:"initialisation_vector"`
	Sensitive            bool      `db:"sensitive"`
}

type ProjectVariableWithEnvironments struct {
	ID                   string
	ProjectID            string
	Key                  string
	Value                string
	UpdatedAt            time.Time
	CreatedAt            time.Time
	EncryptionKeyID      string
	InitialisationVector string
	Sensitive            bool

	Environments []ProjectEnvironment
}

type ProjectVariableOnProjectEnvironment struct {
	ProjectVariableID    string `db:"project_variable_id"`
	ProjectEnvironmentID string `db:"project_environment_id"`
}
