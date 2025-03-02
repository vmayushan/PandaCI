<script lang="ts">
	import * as popover from '@zag-js/popover';
	import { portal, useMachine, normalizeProps } from '@zag-js/svelte';
	import SidebarItem from './sidebar/sidebarItem.svelte';
	import { ChatCircleText } from 'phosphor-svelte';
	import SidebarLabel from './sidebar/sidebarLabel.svelte';
	import posthog from 'posthog-js';
	import Button from './button.svelte';
	import Field from './fieldset/field.svelte';
	import Fieldset from './fieldset/fieldset.svelte';
	import FieldGroup from './fieldset/fieldGroup.svelte';
	import Label from './fieldset/label.svelte';
	import TextArea from './textArea.svelte';
	import { handleForm } from '$lib/utils';
	import { createMutation } from '@tanstack/svelte-query';
	import { PUBLIC_STAGE } from '$env/static/public';

	interface FeedbackProps {
		mobile: boolean;
	}

	const { mobile }: FeedbackProps = $props();

	const service = useMachine(popover.machine, { id: `feedback-${mobile}` });

	const api = $derived(popover.connect(service, normalizeProps));

	let prevOpen = $state(false);

	const survey_id =
		PUBLIC_STAGE === 'prod'
			? '0194f4eb-0ffc-0000-2d14-3c378a31fdca'
			: '0194f4d7-d0ea-0000-6770-f374c2415fa3';

	$effect(() => {
		if (api.open) {
			if (!prevOpen) {
				posthog.capture('survey shown', {
					$survey_id: survey_id
				});
				prevOpen = true;
			}
		} else if (prevOpen) {
			posthog.capture('survey dismissed', {
				$survey_id: survey_id
			});
			prevOpen = false;
		}
	});

	const sendFeedback = createMutation(() => ({
		mutationFn: (data: { feedback: string }) => {
			posthog.capture('survey sent', {
				$survey_id: survey_id,
				$survey_response: data.feedback
			});
			return new Promise((resolve) => setTimeout(resolve, 500));
		},
		onSettled: () => {
			api.setOpen(false);
		}
	}));
</script>

<SidebarItem tooltip="Feedback" {...api.getTriggerProps()}>
	<ChatCircleText data-slot="icon" /><SidebarLabel>Feedback</SidebarLabel>
</SidebarItem>

{#snippet content()}
	<div
		{...api.getContentProps()}
		class="bg-surface-high border-outline relative rounded-lg border p-4 shadow"
	>
		<form
			class="space-y-8"
			onsubmit={(e) => {
				const { data } = handleForm<{ feedback: string }>(e);
				sendFeedback.mutate(data, {
					onSuccess: () => {
						e.target.reset();
					}
				});
			}}
		>
			<Fieldset>
				<FieldGroup>
					<Field>
						<Label>What can we do to improve our product?</Label>
						<TextArea placeholder="Start typing..." name="feedback" />
					</Field>
				</FieldGroup>
			</Fieldset>

			<Button loading={sendFeedback.isPending} color="dark/white" type="submit">Submit</Button>
		</form>
	</div>
{/snippet}
{#if mobile}
	<div {...api.getPositionerProps()}>
		{@render content()}
	</div>
{:else}
	<div use:portal {...api.getPositionerProps()}>
		{@render content()}
	</div>
{/if}
