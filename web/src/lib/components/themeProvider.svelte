<script lang="ts" module>
	export type Theme = 'light' | 'dark' | 'system';

	class ThemeContext {
		theme = $state<Theme>('light');

		constructor() {
			this.theme = (localStorage.getItem('theme') as Theme) || 'system';
		}

		resolvedTheme = () => {
			if (this.theme === 'system')
				return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
			return this.theme;
		};

		setTheme = (theme: Theme) => {
			this.theme = theme;
			localStorage.setItem('theme', theme);
		};
	}

	export function getTheme() {
		return getContext<ThemeContext>('theme');
	}
</script>

<script lang="ts">
	import { getContext, setContext, type Snippet } from 'svelte';

	interface ThemeContextProps {
		children: Snippet;
	}

	const { children }: ThemeContextProps = $props();

	setContext<ThemeContext>('theme', new ThemeContext());
</script>

{@render children()}
