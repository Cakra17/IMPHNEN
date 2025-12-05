<script lang="ts">
	import { Heading, Navbar, NavBrand, NavLi, NavUl, P } from "flowbite-svelte";
    import { BellIcon, ChevronDownIcon } from "@lucide/svelte"
	import Logo from "$lib/components/logo.svelte";
	import Clock from "$lib/components/clock.svelte";
	import Separator from "$lib/components/primitives/separator.svelte";

	let { children } = $props();

	import { page } from "$app/state";
</script>

<div class="font-jakarta w-[99%] h-screen bg-stone-50">
	<Navbar
		class="sticky top-0 z-50 h-14 lg:h-18 bg-white w-full items-center px-4 md:px-12 lg:px-16 xl:px-24 border-b-1 border-teal-200 flex flex-row items-center"
	>
		<NavBrand href="/dashboard" class="flex flex-row items-center gap-2">
			<Logo size="sm" />
			<Heading tag="h6" class="hidden md:inline">IMPHNEN</Heading>
		</NavBrand>
		<NavUl>
			<NavLi href="/dashboard" class={page.url.pathname === '/dashboard' ? 'underline' : ''}
				>Dashboard</NavLi
			>
			<NavLi
				href="/dashboard/orders"
				class={page.url.pathname === '/dashboard/orders' ? 'underline' : ''}>Orderan</NavLi
			>
			<NavLi
				href="/dashboard/analytics"
				class={page.url.pathname === '/dashboard/analytics' ? 'underline' : ''}
				>Analisa Keuangan</NavLi
			>
			<NavLi
				href="/dashboard/config"
				class={page.url.pathname === '/dashboard/config' ? 'underline' : ''}>Pengaturan</NavLi
			>
		</NavUl>
		<div class="flex flex-row py-2 h-full gap-4 justify-end items-center">
			<div class="hidden md:inline flex flex-col items-end">
				<P weight="semibold" size="sm" align="right">John Doe</P>
				<P size="xs" align="right">Toko Pak John</P>
			</div>
			<div class="shrink-0 w-10 h-10 bg-blue-500 rounded-full"></div>
			<button class="cursor:pointer" type="button" onclick={() => console.log('TEST')}>
				<ChevronDownIcon />
			</button>
		</div>
	</Navbar>
	<div class="w-full px-4 md:px-12 lg:px-16 xl:px-24 py-8">
		{@render children()}
	</div>
</div>
