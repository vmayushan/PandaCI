CREATE TABLE project (
    id CHAR(21) NOT NULL PRIMARY KEY,
    org_id CHAR(21) NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    git_integration_id CHAR(21) NOT NULL,
    git_provider_repo_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (org_id) REFERENCES org (id) ON DELETE CASCADE,
    FOREIGN KEY (git_integration_id) REFERENCES git_integration (id),
    UNIQUE (org_id, slug)
);
