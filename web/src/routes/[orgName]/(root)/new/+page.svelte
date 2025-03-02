<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	goto(`/${page.params.orgName}/new/github`, { replaceState: true });
</script>
