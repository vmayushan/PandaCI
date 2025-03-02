CREATE TYPE git_integration_type_enum AS ENUM ('github');

CREATE TABLE git_integration (
    id CHAR(21) NOT NULL PRIMARY KEY,
    provider_id VARCHAR(36) NOT NULL,
    provider_account_id VARCHAR(36) NOT NULL,
    type git_integration_type_enum NOT NULL,
    UNIQUE (provider_id, type)
);
