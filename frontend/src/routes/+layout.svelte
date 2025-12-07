<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import favicon_x96 from '$lib/assets/favicon-96x96.png';
	import favicon_ico from '$lib/assets/favicon.ico';
	import favicon_apple from '$lib/assets/apple-touch-icon.png';
	import manifest from '$lib/assets/site.webmanifest';
	import Logo from '$lib/components/logo.svelte';
	import { Button, Heading, P } from 'flowbite-svelte';
	import { LogInIcon } from '@lucide/svelte';
	import ResponsiveNav from '$lib/components/responsive_nav.svelte';
	import { page } from '$app/state';
	import type { NavItem } from '$lib/types/navitem';

	let { children } = $props();

	const navItems: NavItem[] = [
		{ href: '/pricing', label: 'Fitur & Pricing' },
		{ href: '/faq', label: 'FAQ' },
		{ href: 'http://202.155.95.111/api/v1/docs/index.html', label: 'API' },
		{ href: '/about', label: 'Tentang Kami' }
	];

	let scrollY = $state(0);
	let scrolled = $derived(scrollY > 5);

	let validHeader = $derived(
		!page.url.pathname.startsWith('/auth') &&
			!page.url.pathname.startsWith('/dashboard') &&
			!page.url.pathname.startsWith('/payment')
	);
</script>

<svelte:window bind:scrollY />

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="icon" type="image/png" href={favicon_x96} sizes="96x96" />
	<link rel="icon" type="image/svg+xml" href={favicon} />
	<link rel="shortcut icon" href={favicon_ico} />
	<link rel="apple-touch-icon" sizes="180x180" href={favicon_apple} />
	<link rel="manifest" href={manifest} />
	<title>IMPHNEN (Ingin Menjadi Pengusaha Handal Namun Enggan Ngebuku)</title>
</svelte:head>

<div class={validHeader ? 'bg-pattern' : ''}>
	{#if validHeader}
		<nav
			class="font-jakarta sticky top-0 z-50 w-full px-4 md:px-8 lg:px-16 xl:px-36 flex flex-row items-center justify-between transition-all duration-100 ease-in-out border-b {scrolled
				? 'h-[4rem] scrolled bg-white border-teal-200 xl:px-36'
				: 'h-[5rem] bg-transparent border-transparent xl:px-42!'}"
		>
			<!-- Logo Section -->
			<a href="/" class="flex flex-row items-center gap-2">
				<Logo size="sm" />
				<Heading tag="h6" class="hidden md:inline">IMPHNEN</Heading>
			</a>
			<!-- Navigation Links -->
			<div
				class="flex-1 flex flex-row py-2 h-full gap-4 justify-end items-center {scrolled
					? 'text-black'
					: 'text-white'}"
			>
				<ResponsiveNav
					pathname={page.url.pathname}
					{navItems}
					classes={scrolled ? '' : 'hover:text-emerald-200!'}
					topClass={scrolled ? 'top-[4rem]' : 'top-[5rem]'}
					justify="end"
				/>
				<Button
					class="rounded-full transition-all"
					size="sm"
					color={scrolled ? 'primary' : 'light'}
					href="/auth/register">Daftar Gratis</Button
				>
			</div>
		</nav>
	{/if}
	{@render children()}
</div>
