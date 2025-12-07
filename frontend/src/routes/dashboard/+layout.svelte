<script lang="ts">
	import { Heading, Navbar, NavBrand, NavLi, NavUl, P, Tooltip } from 'flowbite-svelte';
	import { BellIcon, ChevronDownIcon, LogOutIcon } from '@lucide/svelte';
	import Logo from '$lib/components/logo.svelte';
	import { page } from '$app/state';
	import type { LayoutData } from '../$types';
	import ResponsiveNav from '$lib/components/responsive_nav.svelte';

	let { data, children } = $props<{ data: LayoutData; children: any }>();

	import { goto } from '$app/navigation';

	function handleLogout() {
		goto('/logout');
	}
</script>

<div class="font-jakarta w-[100%] h-screen overflow-x-hidden bg-stone-50">
	<nav
		class="sticky top-0 z-50 h-14 lg:h-18 bg-white w-full px-4 md:px-8 lg:px-16 xl:px-36 border-b border-teal-200 flex flex-row items-center justify-between"
	>
		<!-- Logo Section -->
		<a href="/dashboard" class="flex-1 flex flex-row items-center gap-2">
			<Logo size="sm" />
			<Heading tag="h6" class="hidden md:inline">IMPHNEN</Heading>
		</a>
		<!-- Navigation Links -->
		<ResponsiveNav pathname={page.url.pathname} />

		<!-- User Profile Section -->
		<div class="flex-1 flex flex-row py-2 h-full gap-4 justify-end items-center">
			<div class="hidden md:flex flex-col items-end">
				<P weight="semibold" size="sm" align="right">{data.user?.firstname} {data.user?.lastname}</P
				>
				<P size="xs" align="right">{data.user?.store_name}</P>
			</div>
			<button class="cursor-pointer" type="button" onclick={handleLogout}>
				<LogOutIcon />
			</button>
			<Tooltip>Keluar dari Akun</Tooltip>
		</div>
	</nav>
	<div class="w-full px-4 md:px-8 lg:px-16 xl:px-36 py-4 md:py-8">
		{@render children()}
	</div>
</div>
