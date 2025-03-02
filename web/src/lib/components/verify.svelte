<script lang="ts">
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import { ArrowRight } from 'phosphor-svelte';

	interface VerifyProps {
		class?: string;
	}

	const { class: className }: VerifyProps = $props();

	const user = createQuery(() => queries.auth.session());

	const verifyAddress = $derived(
		user.data?.data.identity?.verifiable_addresses?.some((address) => address.verified === false)
	);
</script>

{#if verifyAddress}
	<div class={['flex w-full items-center justify-center gap-x-6 px-6 py-2.5', className]}>
		<a
			class="text-on-surface-variant hover:text-on-surface flex items-center justify-center space-x-1 text-sm"
			href={`/verification?email=${user.data?.data.identity?.traits.email}`}
		>
			<span>Please verify your email</span><span class="hidden sm:block">address</span>
			<ArrowRight data-slot="icon" />
		</a>
	</div>
{/if}
