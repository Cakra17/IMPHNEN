<script lang="ts">
	import { Heading } from 'flowbite-svelte';
	import { TrendingUpIcon, TrendingDownIcon, ChevronDownIcon } from '@lucide/svelte';

	let activeTab = $state<'income' | 'expense'>('income');
</script>

<!-- Mobile Tabs (visible on mobile only) -->
<div class="lg:hidden flex border-b border-teal-300">
	<button
		type="button"
		class={`flex-1 py-3 px-4 font-semibold transition-colors ${
			activeTab === 'income'
				? 'text-emerald-600 border-b-2 border-emerald-600 bg-emerald-50/50'
				: 'text-gray-500 hover:text-gray-700'
		}`}
		onclick={() => (activeTab = 'income')}
	>
		Pemasukan
	</button>
	<button
		type="button"
		class={`flex-1 py-3 px-4 font-semibold transition-colors ${
			activeTab === 'expense'
				? 'text-red-600 border-b-2 border-red-600 bg-red-50/50'
				: 'text-gray-500 hover:text-gray-700'
		}`}
		onclick={() => (activeTab = 'expense')}
	>
		Pengeluaran
	</button>
</div>

<!-- Desktop: Grid Layout (2 columns) -->
<!-- Mobile: Show only active tab -->
<div class="grid lg:grid-cols-2 overflow-hidden">
	<!-- Income Section -->
	<div
		class={`pt-4 p-6 md:border-r border-teal-200 bg-gradient-to-b from-emerald-50 via-white to-white ${
			activeTab === 'income' ? 'block' : 'hidden lg:block'
		}`}
	>
		<div class="pb-4 flex flex-row justify-between">
			<div>
				<Heading class="text-lg hidden lg:block">Pemasukan</Heading>
				<Heading class="text-2xl lg:text-xl mb-2 text-emerald-500">Rp 40.000</Heading>
			</div>
			<TrendingUpIcon class="hidden lg:block" size={32} />
		</div>
		<div>
			<div class="rounded-lg border border-teal-300 overflow-hidden bg-white">
				<div
					class="p-4 flex flex-row gap-6 items-center border-b border-teal-300 hover:bg-stone-50 active:bg-stone-50/30"
				>
					<div class="flex-1">
						<Heading tag="h6" class="text-left">Penjualan Produk A</Heading>
						<span class="text-sm mb-2">Terjual 10 Unit</span>
					</div>
					<ChevronDownIcon />
					<Heading class="text-xl mb-2 text-emerald-500 self-start">+ Rp 10.000</Heading>
				</div>
				<div
					class="p-4 flex flex-row gap-6 items-center border-b border-teal-300 hover:bg-stone-50 active:bg-stone-50/30"
				>
					<div class="flex-1">
						<Heading tag="h6" class="text-left">Penjualan Produk B</Heading>
						<span class="text-sm mb-2">Terjual 15 Unit</span>
					</div>
					<ChevronDownIcon />
					<Heading class="text-xl mb-2 text-emerald-500 self-start">+ Rp 15.000</Heading>
				</div>
				<div
					class="p-4 flex flex-row gap-6 items-center border-b border-teal-300 hover:bg-stone-50 active:bg-stone-50/30"
				>
					<div class="flex-1">
						<Heading tag="h6" class="text-left">Penjualan Produk C</Heading>
						<span class="text-sm mb-2">Terjual 3 Unit</span>
					</div>
					<ChevronDownIcon />
					<Heading class="text-xl mb-2 text-emerald-500 self-start">+ Rp 15.000</Heading>
				</div>
			</div>
		</div>
	</div>

	<!-- Expense Section -->
	<div
		class={`pt-4 p-6 bg-gradient-to-b from-red-50 via-white to-white ${
			activeTab === 'expense' ? 'block' : 'hidden lg:block'
		}`}
	>
		<div class="pb-4 flex flex-row justify-between">
			<div>
				<Heading class="hidden lg:block text-lg">Pengeluaran</Heading>
				<Heading class="text-2xl lg:text-xl mb-2 text-red-500">Rp 10.000</Heading>
			</div>
			<TrendingDownIcon class="hidden lg:block" size={32} />
		</div>
		<div>
			<div class="rounded-lg border border-teal-300 overflow-hidden bg-white">
				<div
					class="p-4 flex flex-row gap-6 items-center border-b border-teal-300 hover:bg-stone-100/50 active:bg-stone-50/30"
				>
					<div class="flex-1">
						<Heading tag="h6" class="text-left">Setoran Ormas</Heading>
						<span class="text-sm mb-2">Tagihan</span>
					</div>
					<Heading class="text-xl mb-2 text-red-500 self-start">- Rp 10.000</Heading>
				</div>
			</div>
		</div>
	</div>
</div>
