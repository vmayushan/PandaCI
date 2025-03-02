<script lang="ts" module>
	import type { Api, Service } from '@zag-js/menu';

	export type DropdownContext = {
		readonly api: Api;
		readonly service: Service;
	};
</script>

<script lang="ts">
	import * as menu from '@zag-js/menu';
	import { nanoid } from 'nanoid';
	import { setContext, type Snippet } from 'svelte';
	import { normalizeProps, useMachine } from '@zag-js/svelte';

	interface DropdownProps {
		children: Snippet;
		triggerID?: string;
	}

	function createMenuAPI(id: string = nanoid(6)) {
		const service = useMachine(menu.machine, {
			id,
			positioning: { sameWidth: true },
			loopFocus: true
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

	const { children, triggerID }: DropdownProps = $props();

	const ctx = createMenuAPI(triggerID);
	setContext<DropdownContext>('dropdown', ctx);
</script>

{@render children()}
