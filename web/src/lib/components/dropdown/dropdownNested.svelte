<script lang="ts">
	import { getContext, onMount, setContext, type Snippet } from 'svelte';
	import * as menu from '@zag-js/menu';
	import { normalizeProps, useMachine } from '@zag-js/svelte';
	import type { DropdownContext } from './dropdown.svelte';
	import { nanoid } from 'nanoid';

	type DropdownNestedProps = { children: Snippet; label: string };

	const { children, label }: DropdownNestedProps = $props();

	const ctx = getContext<DropdownContext>('dropdown');

	function createMenuAPI() {
		const id = nanoid(6);

		const service = useMachine(menu.machine, {
			id,
			positioning: { sameWidth: true },
			'aria-label': label
		});

		const api = $derived(menu.connect(service, normalizeProps));

		return {
			get api() {
				return api;
			},
			get service() {
				return service;
			}
		};
	}

	const subCtx = createMenuAPI();

	onMount(() => {
		ctx.api.setChild(subCtx.service);
		subCtx.api.setParent(ctx.service);
	});

	setContext<DropdownContext>('dropdown-sub', subCtx);
</script>

{@render children()}
