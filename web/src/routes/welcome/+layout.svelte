<script lang="ts">
	import HideUntilAuthed from '../hideUntilAuthed.svelte';

	const { children } = $props();
</script>

<HideUntilAuthed />
{@render children()}
