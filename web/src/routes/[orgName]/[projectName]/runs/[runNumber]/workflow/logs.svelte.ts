import { API } from '$lib/api';
import { type RunStatus } from '$lib/api/organization';
import { transformerNotationErrorLevel } from '@shikijs/transformers';
import { EventSource } from 'eventsource';
import { parse } from 'papaparse';
import { createHighlighterCore } from 'shiki/core';
import { createJavaScriptRegexEngine } from 'shiki/engine/javascript';

const highlighterPromise = createHighlighterCore({
	themes: [import('@shikijs/themes/snazzy-light'), import('@shikijs/themes/dracula-soft')],
	langs: [import('@shikijs/langs/jsonl')],
	engine: createJavaScriptRegexEngine()
});

function parseAsync(data: string, download = false) {
	return new Promise<string[][]>((resolve, reject) => {
		parse<string[]>(data, {
			download: download as false,
			fastMode: false,
			complete: (results) => {
				resolve(results.data);
			},
			error: (error: any) => {
				reject(error);
			}
		});
	});
}

export class Logs {
	#orgSlug: string;
	#projectSlug: string;

	#lang: 'ansi' | 'jsonl';

	#enabled: boolean;

	#status: string | null = null;
	#runID: string | null = null;
	#logURL: string | null = null;
	#stepID: string | null = null;

	#loading = $state(true);

	#stream: EventSource | null = null;
	#streamSourceController: AbortController | null = null;

	#logs = $state<string[][]>([]);
	#htmlLogs = $state<string>();

	constructor(
		orgSlug: string,
		projectSlug: string,
		lang: 'ansi' | 'jsonl',
		enabled: boolean = true
	) {
		this.#orgSlug = orgSlug;
		this.#projectSlug = projectSlug;
		this.#enabled = enabled;
		this.#loading = enabled;
		this.#lang = lang;
	}

	update(data: {
		runID?: string;
		status?: RunStatus;
		logURL?: string;
		stepID?: string;
		enabled?: boolean;
	}) {
		this.#runID = data.runID || null;
		this.#status = data.status || null;
		this.#logURL = data.logURL || null;
		this.#stepID = data.stepID || null;

		this.#enabled = data.enabled ?? true;

		this.#run();
	}

	async #streamLogs() {
		if (this.#stream || this.#streamSourceController || this.#status !== 'running' || !this.#runID)
			return;

		this.#loading = true;

		this.#streamSourceController = new AbortController();

		const { url, authorization } = await API.get(
			'/v1/orgs/{orgSlug}/projects/{projectSlug}/run/{workflowID}/stream/logs',
			{
				params: {
					orgSlug: this.#orgSlug,
					projectSlug: this.#projectSlug,
					workflowID: this.#runID
				},
				queries: this.#stepID
					? {
							step_id: this.#stepID
						}
					: undefined,
				signal: this.#streamSourceController.signal
			}
		);

		this.#stream = new EventSource(url, {
			fetch: (url, options) => {
				return fetch(url, {
					...options,
					headers: {
						...options?.headers,
						Authorization: authorization
					}
				});
			}
		});

		this.#stream.onopen = () => {
			this.#loading = false;
		};

		this.#stream.onmessage = async (event) => {
			const logs = await parseAsync('timestamp,type,data\n' + atob(event.data)).then((logs) => {
				return logs.slice(1);
			});
			this.#logs.push(...logs);
			await this.#renderLogs();
		};

		this.#stream.onerror = (event) => {
			console.error(event);
			// TODO - deal with stream close
		};
	}

	async #downloadLogs() {
		if (this.#status !== 'completed' || !this.#logURL) return;

		this.#loading = true;

		this.#stopStream();

		this.#logs = await parseAsync(this.#logURL, true);

		this.#renderLogs();
	}

	#stopStream() {
		if (this.#status !== 'running') {
			this.#streamSourceController?.abort();
			this.#streamSourceController = null;
			if (this.#stream) {
				this.#stream.close();
				this.#stream = null;
			}
		}
	}

	async #renderLogs() {
		if (this.#logs.length === 0) {
			this.#loading = false;
			this.#htmlLogs = '';
			return;
		}

		let logText = this.#logs.slice(1).map((row) => {
			return (row[2] ?? '').replace(/\r\n|\r|\n/g, '\n');
		});

		if (this.#lang === 'jsonl') {
			const [format, EstreePlugin, BabelPlugin] = await Promise.all([
				import('prettier/standalone').then((module) => module.format),
				import('prettier/plugins/estree'),
				import('prettier/plugins/babel')
			]);

			logText = await Promise.all(
				logText.map((log) =>
					format(log, { parser: 'json', plugins: [EstreePlugin, BabelPlugin] }).catch((e) => {
						console.error(e);
						return log;
					})
				)
			);
		}

		const highlighter = await highlighterPromise;

		const html = highlighter.codeToHtml(logText.join(''), {
			lang: this.#lang,
			themes: {
				light: 'snazzy-light',
				dark: 'dracula-soft'
			},
			colorReplacements: {
				'#282a36': '#242427'
			},
			transformers: [transformerNotationErrorLevel()]
		});

		this.#htmlLogs = html;
		this.#loading = false;
	}

	async #run() {
		if (!this.#enabled) {
			this.#stopStream();
			this.#logs = [];
			this.#loading = false;
			return;
		}
		this.#streamLogs();
		this.#downloadLogs();
	}

	$destroy() {
		this.#stopStream();
	}

	get logs() {
		return this.#logs;
	}

	get isLoading() {
		return this.#loading;
	}

	get htmlLogs() {
		return this.#htmlLogs;
	}
}
