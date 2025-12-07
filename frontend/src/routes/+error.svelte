<script lang="ts">
	import { page } from '$app/stores';
</script>

<div>
	<h1>{$page.status}</h1>
	<p>{$page.status === 404 ? 'Page not found' : 'Something went wrong'}</p>
	<a href="/">Go home</a>
</div>

<style>
	/* Inline styles only, no external CSS */
	div {
		padding: 2rem;
		text-align: center;
	}
</style>
