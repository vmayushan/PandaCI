<script lang="ts">
	import Field from '../fieldset/field.svelte';
	import Label from '../fieldset/label.svelte';
	import Input from '../input.svelte';
	import Messages from './messages.svelte';
	import Button from '../button.svelte';
	import type { UiNode } from '@ory/client';
	import { omit } from 'lodash-es';
	import Text from '../text/text.svelte';
	import Code from '../text/code.svelte';
	import { GithubLogo } from 'phosphor-svelte';

	interface RenderProps {
		nodes: UiNode[];
		filters?: {
			groups?: string[];
			attributes?: Record<string, string | boolean | number>[];
			excludeAttributes?: Record<string, string | boolean | number>[];
		};
		smallButtons?: boolean;
		defaultValues?: Record<string, string>;
		disabled?: boolean;
	}

	const {
		nodes: unsortedNodes,
		disabled,
		filters,
		smallButtons,
		defaultValues
	}: RenderProps = $props();

	let filteredNodes = unsortedNodes;

	if (filters?.groups) {
		filteredNodes = filteredNodes.filter(({ group }) => filters.groups?.includes(group));
	}

	if (filters?.excludeAttributes) {
		filteredNodes = filteredNodes.filter(
			({ attributes }) =>
				!filters.excludeAttributes?.some((attr) =>
					Object.entries(attr).every(
						([key, value]) => attributes[key as keyof typeof attributes] === value
					)
				)
		);
	}

	if (filters?.attributes) {
		filteredNodes = filteredNodes.filter(({ attributes }) =>
			filters.attributes?.some((attr) =>
				Object.entries(attr).every(
					([key, value]) => attributes[key as keyof typeof attributes] === value
				)
			)
		);
	}

	const visible: UiNode[] = [];
	const hidden: UiNode[] = [];

	for (const node of filteredNodes) {
		if (node.attributes.node_type === 'input' && node.attributes.type === 'hidden')
			hidden.push(node);
		else visible.push(node);
	}

	const nodes = [...visible, ...hidden];
</script>

{#snippet input({ node }: { node: UiNode })}
	<Field>
		{#if node.meta.label}
			<Label>
				{node.meta.label.text}
			</Label>
		{/if}
		<Input
			onchange={(e) => {
				if ((node.attributes as any).name === 'code') {
					e.currentTarget.value = e.currentTarget.value.trim();
				}
			}}
			defaultValue={defaultValues?.[(node.attributes as any).name as string]}
			disabled={disabled || (node.attributes as any).disabled || undefined}
			{...omit(node.attributes as any, ['node_type', 'disabled'])}
		/>
		<Messages messages={node.messages} />
	</Field>
{/snippet}

<!-- eslint-disable-next-line svelte/require-each-key -->
{#each nodes as node}
	{#if node.type === 'input' && node.attributes.node_type === 'input'}
		{#if node.attributes.type === 'hidden'}
			<input {...node.attributes as any} />
		{:else if node.attributes.type === 'submit' && node.group === 'oidc'}
			<Button
				name={node.attributes.name}
				type={node.attributes.type}
				value={node.attributes.value}
				disabled={disabled || node.attributes.disabled}
				outline
				full
			>
				<GithubLogo data-slot="icon" />
				{node.meta.label?.text || 'Submit'}
			</Button>
		{:else if node.attributes.type === 'submit'}
			<Button
				name={node.attributes.name}
				type={node.attributes.type}
				value={node.attributes.value}
				disabled={disabled || node.attributes.disabled}
				full={!smallButtons}
			>
				{node.meta.label?.text || 'Submit'}
			</Button>
		{:else if node.attributes.type === 'button'}
			<Button
				name={node.attributes.name}
				type={node.attributes.type}
				value={node.attributes.value}
				disabled={node.attributes.disabled}
				onclick={() => eval((node.attributes as any).onclick)}
				full={!smallButtons}
			>
				{node.meta.label?.text}
			</Button>
		{:else}
			{@render input({ node })}
		{/if}
	{:else if node.type === 'text' && node.attributes.node_type === 'text'}
		{#if node.meta.label}
			<Text>{node.meta.label.text}</Text>
		{/if}
		<Code>{node.attributes.text}</Code>
	{:else if node.attributes.node_type === 'a'}
		<Button href={node.attributes.href} id={node.attributes.id} {disabled}>
			{node.attributes.title.text}
		</Button>
	{:else if node.attributes.node_type === 'script'}
		<svelte:element
			this={'script'}
			src={node.attributes.src}
			nonce={node.attributes.nonce}
			integrity={node.attributes.integrity}
			async={node.attributes.async}
			crossorigin={node.attributes.crossorigin as any}
			type={node.attributes.type}
			id={node.attributes.id}
			referrerpolicy="no-referrer"
		></svelte:element>
	{:else if node.type === 'img'}
		{#if node.meta.label}
			<Text>{node.meta.label.text}</Text>
		{/if}
		<img {...omit(node.attributes, ['node_type']) as any} />
	{/if}
{/each}
