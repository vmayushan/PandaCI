<script lang="ts">
	import {
		Dialog,
		DialogTitle,
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
	import type { Project, Organization, ProjectEnvironment } from '$lib/api/organization';
	import { API, type apiDefs } from '$lib/api';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { handleForm } from '$lib/utils';
	import { queries } from '$lib/queries';
	import DialogCloseButton from '$lib/components/dialog/dialogCloseButton.svelte';
	import Dropdown from '$lib/components/dropdown/dropdown.svelte';
	import DropdownButton from '$lib/components/dropdown/dropdownButton.svelte';
	import DropdownMenu from '$lib/components/dropdown/dropdownMenu.svelte';
	import DropdownCheckbox from '$lib/components/dropdown/dropdownCheckbox.svelte';
	import { CaretDown, Eye, EyeSlash, Stack } from 'phosphor-svelte';
	import Description from '$lib/components/fieldset/description.svelte';
	import DropdownSection from '$lib/components/dropdown/dropdownSection.svelte';
	import DropdownDivider from '$lib/components/dropdown/dropdownDivider.svelte';
	import DropdownItem from '$lib/components/dropdown/dropdownItem.svelte';
	import InputGroup from '$lib/components/inputGroup.svelte';
	import Switch from '$lib/components/switch.svelte';

	interface CreateVariableModal {
		open: boolean;
		project: Project;
		org: Organization;
		environments: ProjectEnvironment[];
	}

	let { open = $bindable(), org, project, environments }: CreateVariableModal = $props();

	let showValue = $state(false);

	const queryClient = useQueryClient();

	const projectVariableMutation = createMutation(() => ({
		mutationFn: (
			data: apiDefs['POST']['/v1/orgs/{orgName}/projects/{projectName}/variables']['req']
		) =>
			API.post('/v1/orgs/{orgName}/projects/{projectName}/variables', {
				body: { ...data },
				params: { orgName: org.slug, projectName: project.slug }
			}),
		onSettled: () => {
			queryClient.invalidateQueries(
				queries.variables.projectVariables(org.slug, project.slug)._ctx.list()
			);
		}
	}));

	const checkedItems = $state<
		Record<
			string,
			{
				checked: boolean;
				name: string;
			}
		>
	>(
		environments.reduce(
			(acc, env) => {
				acc[env.id] = { checked: false, name: env.name };
				return acc;
			},
			{} as Record<string, { checked: boolean; name: string }>
		)
	);
</script>

<Dialog bind:open>
	<DialogTitle>Create variable</DialogTitle>
	{#if projectVariableMutation.error}
		<Text variant="error">{projectVariableMutation.error.message}</Text>
	{/if}
	<form
		class="flex flex-col space-y-8"
		onsubmit={(e) => {
			const { data } =
				handleForm<apiDefs['POST']['/v1/orgs/{orgName}/projects/{projectName}/variables']['req']>(
					e
				);

			data.environmentIDs = Object.entries(checkedItems)
				.filter(([, { checked }]) => checked)
				.map(([id]) => id);

			data.sensitive = (data.sensitive as unknown as string) === 'on';

			projectVariableMutation.mutate(
				{
					...data
				},
				{
					onSuccess: () => {
						(e.target as HTMLFormElement)?.reset();
						open = false;
					}
				}
			);
		}}
	>
		<DialogBody>
			<Fieldset>
				<FieldGroup>
					<Field class="flex flex-col">
						<Label>Environments</Label>
						<Description class="mt-1">Restrict this variable to specific branches.</Description>
						<Dropdown triggerID="env-selector">
							<DropdownButton outline full data-slot="control" type="button">
								<Stack data-slot="icon" />
								{@const selectedEnvs = Object.values(checkedItems).filter(({ checked }) => checked)}
								<span>
									{#if selectedEnvs.length}
										{selectedEnvs.map(({ name }) => name).join(', ')}
									{:else}
										All environments
									{/if}
								</span>

								{#snippet indicator(iconProps)}
									<CaretDown data-slot="icon" class="ml-auto!" {...iconProps} />
								{/snippet}
							</DropdownButton>

							<DropdownMenu disablePortal>
								<DropdownSection>
									{#if environments.length === 0}
										<Text class="mx-4 py-2">No environments</Text>
									{/if}
									{#each environments as env (env.id)}
										<DropdownCheckbox bind:checked={checkedItems[env.id].checked} value={env.id}>
											{env.name}
										</DropdownCheckbox>
									{/each}
								</DropdownSection>
								<DropdownDivider />
								<DropdownItem value="manage" href={`/${org.slug}/${project.slug}/environments`}>
									Manage environments
								</DropdownItem>
							</DropdownMenu>
						</Dropdown>
					</Field>

					<Field>
						<Switch label="Sensitive" name="sensitive" />
						<Description>Sensitive values won't be accessible from our dashboard</Description>
					</Field>

					<Field>
						<Label>Name</Label>
						<Input type="text" name="key" />
					</Field>

					<Field>
						<Label>Value</Label>
						<InputGroup class="flex space-x-2">
							<Input type={showValue ? 'text' : 'password'} name="value" />
							<Button
								tooltip={showValue ? 'Hide value' : 'Show value'}
								type="button"
								aria-label={showValue ? 'Hide value' : 'Show value'}
								aria-pressed={showValue}
								outline
								onclick={() => (showValue = !showValue)}
							>
								{#if showValue}
									<EyeSlash data-slot="icon" />
								{:else}
									<Eye data-slot="icon" />
								{/if}
							</Button>
						</InputGroup>
					</Field>
				</FieldGroup>
			</Fieldset>
		</DialogBody>
		<DialogActions>
			<DialogCloseButton plain>Cancel</DialogCloseButton>
			<Button
				loading={projectVariableMutation.isPending}
				color="dark/white"
				class="self-end "
				type="submit"
			>
				Save
			</Button>
		</DialogActions>
	</form>
</Dialog>
