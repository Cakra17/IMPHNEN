<script lang="ts">
	import { Heading } from 'flowbite-svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	$effect(() => {
		if (data?.status === 'Success') {
			setTimeout(() => {
				window.location.href = 'https://www.youtube.com/watch?v=-X-x10ntt0Y';
			}, 1500);
		}
	});
</script>

<div class="w-screen h-screen flex items-center justify-center">
	{#if data?.status === 'Success'}
		<Heading tag="h2">Pembayaran Berhasil.</Heading>
	{:else if data?.status === 'Error'}
		<Heading tag="h2">Pembayaran Gagal.</Heading>
	{/if}
</div>
