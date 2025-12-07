<script lang="ts">
	import { Heading, P } from 'flowbite-svelte';
	import { TrendingUpIcon, TrendingDownIcon, ChevronDownIcon } from '@lucide/svelte';
	import { formatRupiah } from '$lib/utils/format_money';
	import { transactionsSources, type Transaction } from '$lib/types/transaction';

	let activeTab = $state<'income' | 'expense'>('income');

	let { transactions = [] }: { transactions: Transaction[] } = $props();

	// Split and calculate in one pass
	let incomes = $derived(transactions.filter((t) => t.type === 'income') ?? []);
	let expenses = $derived(transactions.filter((t) => t.type === 'expense') ?? []);

	let totalIncome = $derived(incomes.reduce((sum, t) => sum + t.amount, 0));
	let totalExpense = $derived(expenses.reduce((sum, t) => sum + t.amount, 0));
	let netTotal = $derived(totalIncome - totalExpense);
</script>

<!-- Summary Section-->
<div class="p-6 flex flex-col md:flex-row justify-between border-b border-teal-200">
	<Heading class="text-xl">Net Total</Heading>
	<Heading class={`text-xl mb-2 ${netTotal < 0 ? 'text-red-500' : 'text-emerald-500'}`}
		>{formatRupiah(totalIncome - totalExpense)}</Heading
	>
</div>

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
				<Heading class="text-2xl lg:text-xl mb-2 text-emerald-500"
					>{formatRupiah(totalIncome)}</Heading
				>
			</div>
			<TrendingUpIcon class="hidden lg:block" size={32} />
		</div>
		<div>
			{#if incomes.length > 0}
				<div class="rounded-lg border border-teal-300 overflow-hidden bg-white">
					{#each incomes as tx, i}
						<div
							class={[
								'p-4 flex flex-row gap-6 items-center hover:bg-stone-100/50 active:bg-stone-50/30',
								i < transactions.length - 1 ? 'border-b border-teal-300' : ''
							]}
						>
							<div class="flex-1">
								<Heading tag="h6" class="text-left">{transactionsSources[tx.source]}</Heading>
							</div>
							<Heading class="text-xl mb-2 text-emerald-500 self-start"
								>+ {formatRupiah(tx.amount)}</Heading
							>
						</div>
					{/each}
				</div>
			{:else}
				<div class="py-12"><P align="center">Tidak ada transaksi masuk</P></div>
			{/if}
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
				<Heading class="text-2xl lg:text-xl mb-2 text-red-500">{formatRupiah(totalExpense)}</Heading
				>
			</div>
			<TrendingDownIcon class="hidden lg:block" size={32} />
		</div>
		<div>
			{#if expenses.length > 0}
				<div class="rounded-lg border border-teal-300 overflow-hidden bg-white">
					{#each expenses as tx, i}
						<div
							class={[
								'p-4 flex flex-row gap-6 items-center hover:bg-stone-100/50 active:bg-stone-50/30',
								i < transactions.length - 1 ? 'border-b border-teal-300' : ''
							]}
						>
							<div class="flex-1">
								<Heading tag="h6" class="text-left">{transactionsSources[tx.source]}</Heading>
							</div>
							<Heading class="text-xl mb-2 text-red-500 self-start"
								>- {formatRupiah(tx.amount)}</Heading
							>
						</div>
					{/each}
				</div>
			{:else}
				<div class="py-12"><P align="center">Tidak ada transaksi keluar</P></div>
			{/if}
		</div>
	</div>
</div>
