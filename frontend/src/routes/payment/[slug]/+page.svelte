<script lang="ts">
	// Get the result of the form submission
	import { enhance } from '$app/forms';
	import { page } from '$app/stores';

	$: if ($page.form?.status === 'Success') {
		// The reactive statement (using $: ) runs whenever $page.form changes.
		console.log('Payment Success! Redirecting...');

		// 2. Use setTimeout for the delay, then perform the redirect.
		// If you want a full-page reload/refresh, set window.location.href = window.location.href;
		// Or for navigation without a full reload, use goto('/') from '$app/navigation'.
		setTimeout(() => {
			// Using window.location.href = window.location.href; causes a page refresh
			window.location.href = 'https://www.youtube.com/watch?v=-X-x10ntt0Y';
		}, 1500);
	}
</script>

<form method="POST" use:enhance>
	<button type="submit">Complete Mock Payment</button>
</form>

{#if $page.form?.success}
	<p style="color: green;">{$page.form.message}</p>
{:else if $page.form?.status === 'Error'}
	<p style="color: red;">Error: {$page.form.message}</p>
{/if}
