
CREATE TYPE user_account_type_enum AS ENUM ('github');

CREATE TABLE user_account (
    user_id UUID NOT NULL,
    type user_account_type_enum NOT NULL,
    provider_account_id CHAR(36) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    refresh_token_expires_at TIMESTAMP NOT NULL,
    access_token VARCHAR(255) NOT NULL,
    access_token_expires_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, type),
    UNIQUE (type, provider_account_id)
);

CREATE TABLE user_account_oauth_state (
    id CHAR(21) NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    type user_account_type_enum NOT NULL,
    created_at TIMESTAMP NOT NULL
);
