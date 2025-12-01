<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>IMPHNEN (Ingin Menjadi Pengusaha Handal Namun Enggan Ngebuku)</title>
</svelte:head>

{@render children()}
