<script lang="ts">
	import { Heading, P } from "flowbite-svelte";
    import { ChevronDownIcon } from "@lucide/svelte"
	import Logo from "$lib/components/logo.svelte";
	import Clock from "$lib/components/clock.svelte";
	let { children } = $props();
</script>

<div class="font-jakarta w-screen h-screen bg-stone-50">
	<div
		class="sticky top-0 z-50 h-14 lg:h-18 grid grid-cols-3 bg-white w-full items-center px-8 md:px-16 lg:px-24 border-b-1 border-teal-200 justify-between"
	>
		<a href="/dashboard" class="flex flex-row items-center gap-2">
			<Logo size="sm" />
			<Heading tag="h6" class="hidden md:inline">IMPHNEN</Heading>
		</a>
		<div class="text-center">
			<Heading tag="h6" class="md:hidden inline">IMPHNEN</Heading><span
				class="font-semibold hidden md:inline"><Clock /></span
			>
		</div>
		<div class="flex flex-row gap-2 justify-end items-center">
			<div class="hidden md:inline flex flex-col items-end">
				<P weight="semibold" size="sm" align="right">John Doe</P>
				<P size="xs" align="right">Toko Pak John</P>
			</div>
			<div class="shrink-0 w-10 h-10 bg-blue-500 rounded-full"></div>
			<button class="cursor:pointer" type="button" onclick={() => console.log('TEST')}>
				<ChevronDownIcon />
			</button>
		</div>
	</div>
	<div class="w-full px-8 md:px-16 lg:px-24 py-8">
		{@render children()}
	</div>
</div>
