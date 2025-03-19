import type { Method } from 'api-typify';

export type SubscriptionPlan = 'free' | 'pro' | 'enterprise' | 'paused';
export type SubscriptionStatus = 'active' | 'canceled' | 'past_due' | 'trialing' | 'paused';

export interface Features {
	buildMinutes: number;
	maxBuidMinutes: number;
	buildMinutesPriceID: string;

	committers: number;
	maxCommitters: number;

	maxProjects: number;

	maxCloudRunnersScale: number;
}

export interface SubscriptionItem {
	productID: string;
	priceID: string;
	quantity: number;
}

export interface License {
	plan: SubscriptionPlan;
	features: Features;
	paddleData?: {
		customerID: string;
		lastNotificationOccurredAt: string;
		subscriptionID: string;
		subscriptionStatus: SubscriptionStatus;
		subscriptionItems: SubscriptionItem[];
		collectionMode: 'automatic' | 'manual';
		scheduledChange?: {
			action: 'cancel' | 'pause' | 'resume';
			effectiveAt: string;
			resumeAt?: string;
		};
		nextBillingAt?: string;
		billingPeriod?: {
			startsAt: string;
			endsAt: string;
		};
	};
}

export interface Organization {
	id: string;
	name: string;
	license?: License;
	ownerID: string;
	slug: string;
	avatarURL?: string;
	currentUsersRole?: 'member' | 'admin';
}

export interface Project {
	id: string;
	name: string;
	slug: string;
	orgID: string;
	avatarURL?: string;
	lastBuild?: string;
}

export type RunStatus = 'queued' | 'running' | 'completed' | 'pending';
export type RunConclusion = 'success' | 'failure' | 'skipped' | 'cancelled';

export type Paginated<T> = {
	data: T[];
	next: boolean;
};
export interface WorkflowRunAlert {
	type: 'error' | 'warning' | 'info';
	message: string;
	title: string;
}

export type WorkflowRunTrigger =
	| 'manual'
	| 'push'
	| 'pull_request-opened'
	| 'pull_request-synchronize'
	| 'pull_request-closed';

export interface WorkflowRun {
	id: string;
	name: string;
	projectID: string;
	status: RunStatus;
	conclusion?: RunConclusion;
	createdAt: string;
	finishedAt?: string;
	number: number;
	gitSha: string;
	gitBranch: string;
	prNumber?: number;
	trigger: WorkflowRunTrigger;
	outputURL?: string;
	gitTitle?: string;
	prURL?: string;
	commitURL?: string;
	committer: {
		name?: string;
		email: string;
		avatar?: string;
	};
	jobs?: JobRun[];
	alerts?: WorkflowRunAlert[];
}

export interface TaskRun {
	id: string;
	name: string;

	status: RunStatus;
	conclusion?: RunConclusion;

	createdAt: string;
	finishedAt?: string;

	dockerImage?: string;

	steps: StepRun[];
}

export interface JobRun {
	id: string;
	number: number;
	name: string;

	status: RunStatus;
	conclusion?: RunConclusion;

	runner: string;

	createdAt: string;
	finishedAt?: string;

	tasks: TaskRun[];
}

export interface StepRun {
	id: string;
	type: 'exec';

	name: string;

	createdAt: string;
	finishedAt?: string;

	status: RunStatus;
	conclusion?: RunConclusion;

	outputURL: string;

	meta: unknown;
}

export interface ProjectEnvironment {
	id: string;
	projectID: string;
	name: string;
	branchPattern: string;
	createdAt: string;
	updatedAt: string;
}

export interface ProjectVariable {
	id: string;
	projectID: string;
	key: string;
	value: string;
	updatedAt: string;
	createdAt: string;
	environments?: ProjectEnvironment[];
	sensitive: boolean;
}

export interface OrganizationUser {
	id: string;
	email: string;
	name: string;
	role: 'member' | 'admin';
}

export interface LogMeta {
	url: string;
}

export interface LogStream {
	url: string;
	authorization: string;
}

type DELETE = Method<{
	'/v1/orgs/{orgSlug}/projects/{projectSlug}/variables/{variableID}': object;
	'/v1/orgs/{orgSlug}/projects/{projectSlug}/environments/{environmentID}': object;
	'/v1/orgs/{orgName}/projects/{projectName}': object;
	'/v1/orgs/{orgSlug}': object;
	'/v1/orgs/{orgName}/users/{userID}': object;
}>;

type PUT = Method<{
	'/v1/orgs/{orgName}/projects/{projectName}': {
		req: {
			name: string;
			slug: string;
			avatarURL?: string;
		};
		res: Project;
	};

	'/v1/orgs/{orgName}': {
		req: {
			name: string;
			slug: string;
			avatarURL?: string;
		};
		res: Organization;
	};
	'/v1/orgs/{orgSlug}/projects/{projectSlug}/variables/{variableID}': {
		req: {
			key: string;
			value: string;
			environmentIDs: string[];
			sensitive: boolean;
		};
	};

	'/v1/orgs/{orgSlug}/projects/{projectSlug}/environments/{environmentID}': {
		req: {
			name: string;
			branchPattern: string;
		};
		res: ProjectEnvironment;
	};
}>;

type POST = Method<{
	'/v1/orgs': {
		res: Organization;
		req: {
			slug: string;
			name: string;
		};
	};
	'/v1/orgs/{orgName}/projects': {
		res: Project;
		req: {
			slug: string;
			name: string;
			gitProviderIntegrationID: string;
			gitProviderType: 'github';
			gitProviderRepoID: string;
		};
	};

	'/v1/orgs/{orgName}/users': {
		req: {
			email: string;
		};
	};

	'/v1/orgs/{orgSlug}/projects/{projectSlug}/environments': {
		req: {
			name: string;
			branchPattern: string;
		};
		res: ProjectEnvironment;
	};

	'/v1/orgs/{orgName}/projects/{projectName}/variables': {
		res: ProjectVariable;
		req: {
			value: string;
			key: string;
			environmentIDs: string[];
			sensitive: boolean;
		};
	};

	'/v1/orgs/{orgName}/projects/{projectName}/trigger': {
		req: {
			sha: string;
			branch: string;
		};
		res: WorkflowRun[];
	};

	'/v1/paddle/orgs/{orgSlug}/portal': {
		req: object;
		res: {
			generalURL: string;
		};
	};
}>;

type GET = Method<{
	'/v1/orgs': {
		res: Organization[];
	};
	'/v1/orgs/{orgSlug}': {
		res: Organization;
	};
	'/v1/orgs/{orgSlug}/usage': {
		res: {
			projectCount: number;
			usedBuildMinutes: number;
			usedCommitters: number;
		};
	};
	'/v1/orgs/{orgSlug}/users': {
		res: OrganizationUser[];
	};
	'/v1/orgs/{orgSlug}/projects': {
		res: Project[];
	};
	'/v1/orgs/{orgSlug}/projects/{projectSlug}': {
		res: Project;
	};
	'/v1/orgs/{orgSlug}/projects/{projectSlug}/runs': {
		res: Paginated<Omit<WorkflowRun, 'jobs' | 'alerts'>>;
		queries?: {
			page?: number;
			per_page?: number;
		};
	};

	'/v1/orgs/{orgSlug}/projects/{projectSlug}/environments': {
		res: ProjectEnvironment[];
	};

	'/v1/orgs/{orgSlug}/projects/{projectSlug}/variables': {
		res: ProjectVariable[];
	};

	'/v1/orgs/{orgSlug}/projects/{projectSlug}/variables/{variableID}': {
		res: ProjectVariable;
		queries?: {
			decrypt?: boolean;
		};
	};

	'/v1/orgs/{orgSlug}/projects/{projectSlug}/runs/{runNumber}': {
		res: WorkflowRun;
	};

	'/v1/orgs/{orgSlug}/projects/{projectSlug}/run/{workflowID}/stream/logs': {
		res: LogStream;
		queries?: {
			step_id?: string;
		};
	};

	'/v1/orgs/{orgSlug}/projects/{projectSlug}/run/{workflowID}/logs/{logID}': {
		res: LogMeta;
	};
}>;

export interface OrganizationAPI {
	GET: GET;
	POST: POST;
	DELETE: DELETE;
	PUT: PUT;
}
