<script lang="ts" module>
	export interface CreateOrgFormData {
		slug: string;
		name: string;
	}
</script>

<script lang="ts">
	import { Description, Field, FieldGroup, Fieldset, Input, Label } from '$lib/components';

	let slug = $state<CreateOrgFormData['slug']>();
	let name = $state<CreateOrgFormData['name']>();

	let displaySlugDirty = $state(false);

	$effect(() => {
		if (!displaySlugDirty) slug = name?.replaceAll(' ', '-');
	});
</script>

<Fieldset>
	<FieldGroup>
		<Field>
			<Label>Organization Name</Label>
			<Input
				onchange={(e) => {
					name = e.currentTarget.value;
				}}
				required
				oninput={(e) => (name = e.currentTarget.value)}
				placeholder="Acme Inc"
				name="name"
				type="text"
			/>
		</Field>

		<Field>
			<Label>Organization URL</Label>
			<Input
				oninput={(e) => {
					slug = e.currentTarget.value.replaceAll(' ', '-');
					displaySlugDirty = slug !== '';
				}}
				value={slug}
				minlength={2}
				type="text"
				required
				placeholder="your-org"
				name="slug"
			/>
			<Description>
				https://app.pandaci.com/<b>{encodeURIComponent(slug || 'your-org')}</b>
			</Description>
		</Field>
	</FieldGroup>
</Fieldset>
