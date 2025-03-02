CREATE TABLE org (
    id CHAR(21) NOT NULL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255) NOT NULL,
    license JSONB,
    owner_user_id UUID NOT NULL -- Each org can only have 1 owner
);

CREATE TYPE org_role_enum AS ENUM ('admin', 'member');

CREATE TABLE org_users (
    user_id UUID NOT NULL,
    org_id CHAR(21) NOT NULL,
    role org_role_enum NOT NULL,
    PRIMARY KEY (user_id, org_id),
    FOREIGN KEY (org_id) REFERENCES org (id) ON DELETE CASCADE
);

CREATE TABLE pending_org_invites (
    org_id CHAR(21) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role org_role_enum NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (org_id, email),
    FOREIGN KEY (org_id) REFERENCES org (id) ON DELETE CASCADE
);
