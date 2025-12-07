<script lang="ts">
	import type { NavItem } from '$lib/types/navitem';
	import { ChevronDownIcon, ChevronUpIcon, MenuIcon, XIcon } from '@lucide/svelte';

	interface Props {
		pathname: string;
		navItems: NavItem[];
		classes?: string;
		topClass: string;
		justify?: 'start' | 'center' | 'end';
	}

	let { pathname, navItems, classes = '', topClass, justify: align = 'center' }: Props = $props();
	let mobileMenuOpen = $state(false);

	function isActive(href: string) {
		return pathname === href;
	}
</script>

<!-- Desktop Navigation (hidden on mobile) -->
<ul class={['hidden lg:flex flex-1 flex-row gap-4 items-center', `justify-${align}`]}>
	{#each navItems as item}
		<li>
			<a
				href={item.href}
				class={`hover:text-teal-700 rounded-full py-1.5 px-2.5 transition-colors ${classes} ${isActive(item.href) ? 'bg-teal-200' : ''}`}
			>
				{item.label}
			</a>
		</li>
	{/each}
</ul>

<!-- Mobile Menu Button -->
<div class="lg:hidden flex-1 flex justify-center">
	<button
		type="button"
		class="flex items-center gap-2 px-3 py-2 rounded-md hover:bg-teal-50 transition-colors"
		onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
		aria-label="Toggle menu"
	>
		<span class="font-medium text-sm">
			{navItems.find((item) => isActive(item.href))?.label || 'Menu'}
		</span>
		{#if mobileMenuOpen}
			<ChevronUpIcon class="w-6 h-6" />
		{:else}
			<ChevronDownIcon class="w-6 h-6" />
		{/if}
	</button>
</div>

<!-- Mobile Dropdown Menu -->
{#if mobileMenuOpen}
	<div
		class="lg:hidden absolute left-0 right-0 bg-white border-b border-teal-200 shadow-lg z-40 {topClass}"
	>
		<ul class="flex flex-col">
			{#each navItems as item}
				<li>
					<a
						href={item.href}
						class={`block px-6 py-3 hover:bg-teal-50 transition-colors ${isActive(item.href) ? 'bg-teal-100 border-l-4 border-teal-600' : ''}`}
						onclick={() => (mobileMenuOpen = false)}
					>
						{item.label}
					</a>
				</li>
			{/each}
		</ul>
	</div>
{/if}
