CREATE TABLE project_environment (
    id CHAR(21) NOT NULL PRIMARY KEY,
    project_id CHAR(21) NOT NULL,
    name VARCHAR(255) NOT NULL,
    branch_pattern VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (project_id) REFERENCES project (id) ON DELETE CASCADE,
    UNIQUE (project_id, name)
);

CREATE TABLE project_variable (
    id CHAR(21) NOT NULL PRIMARY KEY,
    project_id CHAR(21) NOT NULL,
    key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    encryption_key_id CHAR(8) NOT NULL,
    initialisation_vector TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    sensitive BOOLEAN NOT NULL,
    FOREIGN KEY (project_id) REFERENCES project (id) ON DELETE CASCADE
);

CREATE TABLE project_variable_on_project_environment (
    project_environment_id CHAR(21) NOT NULL,
    project_variable_id CHAR(21) NOT NULL,
    PRIMARY KEY (project_environment_id, project_variable_id),
    FOREIGN KEY (project_environment_id) REFERENCES project_environment (id) ON DELETE CASCADE,
    FOREIGN KEY (project_variable_id) REFERENCES project_variable (id) ON DELETE CASCADE
);
