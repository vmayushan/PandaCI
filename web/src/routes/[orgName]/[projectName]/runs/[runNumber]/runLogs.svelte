<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { clsx } from 'clsx';

	type LogsProps = SvelteHTMLElements['div'];

	const { class: className, children, ...props }: LogsProps = $props();
</script>

<div class={clsx('run-logs', className)}>
	{@render children?.()}
</div>

<style>
	:global {
		.run-logs {
			@media (prefers-color-scheme: dark) {
				.shiki,
				.shiki span {
					color: var(--shiki-dark) !important;
					background-color: transparent !important;
					font-style: var(--shiki-dark-font-style) !important;
					font-weight: var(--shiki-dark-font-weight) !important;
					text-decoration: var(--shiki-dark-text-decoration) !important;
				}
			}

			pre.shiki {
				white-space: pre-wrap;
				word-wrap: break-word;
				text-indent: 3.4rem hanging each-line;
			}

			code {
				counter-reset: step;
				counter-increment: step 0;
			}

			code .line::before {
				content: counter(step);
				counter-increment: step;
				width: 1rem;
				margin-right: 4ch;
				display: inline-block;
				white-space: nowrap;
				text-align: right;
				color: var(--color-on-surface-variant);
			}

			code .line.error::before {
				color: var(--color-red-500);
			}
		}
	}
</style>
