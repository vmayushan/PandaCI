import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { sveltePhosphorOptimize } from 'phosphor-svelte/vite';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: [sveltePhosphorOptimize(), vitePreprocess()],
	kit: {
		adapter: adapter({
			fallback: 'index.html' // may differ from host to host,
		})
	}
};

export default config;
