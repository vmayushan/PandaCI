<script lang="ts">
	import {
		Dialog,
		DialogTitle,
		DialogDescription,
		DialogBody,
		DialogActions,
		Text,
		Fieldset,
		Field,
		FieldGroup,
		Label,
		Input,
		Button
	} from '$lib/components';
	import { API, type apiDefs } from '$lib/api';
	import { createMutation } from '@tanstack/svelte-query';
	import { handleForm } from '$lib/utils';
	import { page } from '$app/state';
	import DialogCloseButton from '$lib/components/dialog/dialogCloseButton.svelte';

	interface CreateVariableModal {
		open: boolean;
	}

	let { open = $bindable() }: CreateVariableModal = $props();

	const inviteUserMutation = createMutation(() => ({
		mutationFn: ({ email }: { email: string }) =>
			API.post('/v1/orgs/{orgName}/users', {
				body: {
					email
				},
				params: {
					orgName: page.params.orgName
				}
			})
	}));
</script>

<Dialog bind:open>
	<DialogTitle>Invite member</DialogTitle>
	<DialogDescription>Invite users to your organization to collaborate on projects</DialogDescription
	>
	{#if inviteUserMutation.error}
		<Text variant="error">{inviteUserMutation.error.message}</Text>
	{/if}
	<form
		class="flex flex-col space-y-8"
		onsubmit={(e) => {
			const { data } = handleForm<apiDefs['POST']['/v1/orgs/{orgName}/users']['req']>(e);

			inviteUserMutation.mutate(data, {
				onSuccess: () => {
					(e.target as HTMLFormElement)?.reset();
					open = false;
				}
			});
		}}
	>
		<DialogBody>
			<Fieldset>
				<FieldGroup>
					<Field>
						<Label>Email</Label>
						<Input autofocus={true} type="email" name="email" />
					</Field>
				</FieldGroup>
			</Fieldset>
		</DialogBody>
		<DialogActions>
			<DialogCloseButton plain>Cancel</DialogCloseButton>
			<Button
				loading={inviteUserMutation.isPending}
				color="dark/white"
				class="self-end"
				type="submit"
			>
				Invite
			</Button>
		</DialogActions>
	</form>
</Dialog>
