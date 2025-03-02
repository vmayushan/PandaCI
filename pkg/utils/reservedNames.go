package utils

import (
	"regexp"
	"slices"
)

// Large list of reserved names.
// A lot of these probably aren't needed but it's better to be safe
// Most are ai generated
var reservedKeywords = []string{
	// Authentication & User Management
	"login", "logout", "signup", "signin", "signout", "registration",
	"register", "verify", "verification", "recovery", "password",
	"reset", "account", "profile", "settings", "preferences", "dashboard",
	"admin", "administrator", "manage", "users", "user-management",
	"roles", "permissions", "guest", "owner", "member", "visitor",
	"viewer", "maintainer", "collaborator", "superadmin", "moderator",
	"root", "members", "oauth2", "mfa", "totp", "sms", "two-factor",
	"2fa", "session", "sessions", "authenticator",

	// Support & Documentation
	"help", "support", "faq", "feedback", "contact", "about", "terms",
	"privacy", "legal", "documentation", "docs", "tutorials", "guide",
	"resources", "changelog", "status-page", "kb", "knowledgebase",
	"getting-started", "release-notes", "bug-report", "issue", "issues",
	"support-ticket", "forum",

	// API & Integration
	"api", "swagger", "graphql", "webhooks", "integrations",
	"oauth", "tokens", "auth", "auth-token", "jwt", "apikey", "api-key",
	"rest", "rest-api", "rpc", "sdk", "sdk-docs", "oauth-keys", "token-auth",
	"auth-code", "client-id", "client-secret",

	// Content & Media
	"blog", "news", "media", "images", "videos", "files", "downloads",
	"uploads", "assets", "static", "public", "gallery", "podcast",
	"video", "mp3", "mp4", "audio", "slides", "presentation", "ebook",
	"whitepaper", "story", "stories",

	// Site Navigation & Pages
	"home", "index", "services", "products", "pricing", "portfolio",
	"careers", "jobs", "press", "events", "sitemap", "robots.txt",
	"sitemap.xml", "pages", "feature", "features", "examples",

	// Functional Routes
	"search", "download", "upload", "export", "import", "analytics",
	"reports", "metrics", "notifications", "alerts", "logs", "audit",
	"monitor", "test", "beta", "trial", "demo", "sample", "bulk",
	"batch", "optin", "optout", "read", "write", "sync", "refresh",
	"ping", "version", "release-note", "experimental", "invite", "invites",

	// Technical & System
	"config", "configuration", "system", "internal", "maintenance",
	"health", "errors", "error", "500", "crash", "debug", "build", "builds",
	"test", "tests", "pipeline", "pipelines", "job", "jobs", "task", "tasks",
	"deploy", "deploys", "deployment", "deployments", "release", "releases",
	"rollback", "rollbacks", "stages", "step", "steps", "execute", "trigger",
	"triggers", "schedule", "scheduled", "workflow", "workflows", "artifact",
	"artifacts", "logs", "log", "cache", "caches", "queue", "queues",
	"runner", "runners", "agent", "agents", "executor", "executors",
	"environment", "environments", "branch", "branches", "tag", "tags",
	"schema", "schemas", "v1", "v2", "v3", "feature-flag", "env",
	"system-info", "ci", "runtime", "ci-cd", "cicd",

	// Version Control & Repository Management
	"repo", "repos", "repository", "repositories", "git", "github", "gitlab",
	"bitbucket", "commit", "commits", "pull-request", "pull-requests", "pr",
	"prs", "merge-request", "merge-requests", "mr", "mrs", "branch", "branches",
	"fork", "forks", "clone", "clones", "gitea", "perforce", "svn", "source-control",
	"patch", "merge", "tagging",

	// CI/CD Tools & Integrations
	"docker", "dockerfile", "kubernetes", "helm", "terraform", "ansible",
	"chef", "puppet", "salt", "jenkins", "circleci", "travis", "bamboo",
	"azure-pipelines", "aws-codebuild", "codepipeline", "gcp-cloudbuild",
	"code", "tools", "apps", "app", "cd", "orchestration", "auto-deploy",
	"automation", "triggered", "workflow-file",

	// Monitoring & Testing
	"monitor", "monitors", "health", "status", "uptime", "metrics", "performance",
	"unit-test", "integration-test", "functional-test", "acceptance-test",
	"load-test", "stress-test", "security-test", "sla", "logfiles", "uptime-monitor",
	"syslog", "trace", "tracing", "audit-log", "load", "error-tracking", "exception",
	"exceptions",

	// Miscellaneous CI/CD Terms
	"integration", "continuous-integration", "continuous-delivery", "continuous-deployment",
	"configuration", "config", "yml", "yaml", "runs", "run", "secrets", "secret",
	"release", "releases", "step-execution", "job-output", "pipeline-steps", "retries",

	// Marketing & SEO
	"promo", "offers", "discount", "deal", "free", "plans", "pricing",
	"subscribe", "unsubscribe", "newsletter", "trial", "marketing", "links",
	"link", "try", "get", "start", "begin", "end", "compare", "solutions",
	"solution", "explore", "plugin", "plugins", "action", "actions", "market",
	"marketplace", "catalog", "mobile", "amp", "board", "boards", "task",
	"tasks", "dash", "enterprise", "company", "invest", "investors",
	"contributors", "contriubute", "license", "oss", "open", "source",
	"open-source", "guides", "customers", "security", "partners", "experts",
	"community", "communities", "post", "posts", "question", "questions",
	"merch", "store", "handbook", "campaign", "campaigns", "school",
	"university", "startups", "roadmap", "press", "media", "teams", "team",
	"people", "education", "edu", "formum", "shop", "platform", "cloud",
	"dev", "develop", "developers", "developer", "cli", "topics", "topic",
	"ai", "event", "vs",

	// Commerce & Billing
	"billing", "invoice", "payment", "pricing", "checkout", "cart", "order",
	"purchase", "refund", "cancel", "upgrade", "downgrade", "shop", "quote",
	"subscription", "revenue", "renewal", "trial-period", "checkout-session",
	"discounts",

	// Common Abbreviations & Shortcuts
	"www", "http", "https", "ftp", "cdn", "ssl", "ipv4", "ipv6", "dns",
	"localhost", "127.0.0.1", "www2",

	// Time & Date
	"today", "tomorrow", "yesterday", "now", "future", "past", "history",
	"schedule", "calendar",

	// Social & Communication
	"message", "messages", "chat", "inbox", "conversation", "mail", "email",
	"notification", "mention", "follower", "following", "like", "comment",
	"share", "notifications", "dm", "direct-message", "follow", "unfollow",
	"friends", "subscribers", "share-link",

	// Geography & Location
	"location", "city", "state", "country", "region", "timezone",
	"coordinates", "map", "geolocation",

	// Miscellaneous
	"test", "example", "sample", "foo", "bar", "baz", "demo", "template",
	"default", "meta", "new", "edit", "delete", "create", "update", "duplicate",
	"copy", "archive", "trash", "recycle", "temp", "temporary", "sandbox",
	"private", "shared", "common", "public", "org", "orgs", "project",
	"projects", "activity", "bin", "templates", "changes", "updates",
	"improve", "improvements", "data", "panda-ci", "pandaci-com",
}

func IsURLNameValid(orgURL string) bool {
	urlNameRegex := regexp.MustCompile(`^([a-zA-Z0-9\-\_\%\~]+)*$`)
	return len(orgURL) > 2 && !slices.Contains(reservedKeywords, orgURL) && urlNameRegex.MatchString(orgURL)
}
