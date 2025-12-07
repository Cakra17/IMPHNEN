<script lang="ts">
	import CashflowChart from '$lib/components/cashflow_chart.svelte';
	import DailyCashflowTable from '$lib/components/daily_cashflow_table.svelte';
	import {
		ChevronDownIcon,
		ChevronLeftIcon,
		ChevronRightIcon,
		TrendingDownIcon,
		TrendingUpIcon
	} from '@lucide/svelte';
	import { Button, Datepicker, Dropdown, DropdownItem, Heading, Hr, P } from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	import type { PageData } from '../$types';
	import { formatDate, formatDateForBackend } from '$lib/utils/date';
	import { enhance } from '$app/forms';

	// Data
	let { data } = $props<{ data: PageData }>();

	// Cashflow table props
	let dateRangeTo = new Date();
	let dateRangeFrom = new Date();
	dateRangeFrom.setDate(dateRangeTo.getDate() - 7);

	// --- Daily Transactions table ---
	let dailyTransactions = $state(data.dailyTransactions);
	let selectedDate = $state(new Date());
	let formEl: HTMLFormElement;

	let disableNext = $derived(
		formatDateForBackend(selectedDate) === formatDateForBackend(new Date())
	);

	// Auto-submit form when date changes
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;
	$effect(() => {
		selectedDate; // Track this

		// Clear previous timer
		if (debounceTimer) clearTimeout(debounceTimer);

		// Set new timer
		debounceTimer = setTimeout(() => {
			if (formEl) {
				formEl.requestSubmit();
			}
		}, 500); // 500ms debounce

		// Cleanup on unmount
		return () => {
			if (debounceTimer) clearTimeout(debounceTimer);
		};
	});
</script>

<a
	href="/dashboard"
	class="flex flex-row gap-2 font-bold text-teal-600 block lg:hidden hover:text-teal-800 hover:underline cursor-pointer"
	><ChevronLeftIcon />Kembali</a
>
<div class="flex-1">
	<Heading class="text-lg md:text-2xl ml-8 lg:ml-0 mb-8">Analisa & Statistik</Heading>
	<CashflowChart />

	<div class="mt-8 flex flex-col bg-white border border-teal-200 rounded-2xl">
		<div
			class="flex flex-col md:flex-row items-center justify-between gap-2 border-b-1 p-6 pb-3 border-teal-200"
		>
			<div class="flex flex-col">
				<Heading class="text-xl mb-1">Riwayat Transaksi Harian</Heading>
				<span class="text-sm hidden md:block mb-2">Detail cashflow lengkap perhari</span>
			</div>
			<div class="flex-1 flex flex-row gap-3 items-center justify-end">
				<form
					method="POST"
					action="?/getDaily"
					bind:this={formEl}
					use:enhance={() => {
						return async ({ result }) => {
							if (result.type === 'success') {
								dailyTransactions = result.data?.dailyTransactions; // However your action returns it
								console.log(dailyTransactions);
							}
						};
					}}
					class="flex flex-row justify-end"
				>
					<input type="hidden" name="date" value={formatDateForBackend(selectedDate)} />

					<button
						type="button"
						class="text-stone-500 hover:text-stone-800 active:text-stone-800 active:bg-stone-100 rounded-full cursor-pointer"
						onclick={() => {
							selectedDate = new Date(selectedDate.getTime() - 24 * 60 * 60 * 1000);
						}}
					>
						<ChevronLeftIcon size={28} />
					</button>

					<Datepicker bind:value={selectedDate} availableTo={new Date()} />

					<button
						type="button"
						class={`${!disableNext ? 'hover:text-stone-800 active:text-stone-800 active:bg-stone-100 cursor-pointer text-stone-600' : 'text-stone-300'} rounded-full`}
						disabled={disableNext}
						onclick={() => {
							selectedDate = new Date(selectedDate.getTime() + 24 * 60 * 60 * 1000);
						}}
					>
						<ChevronRightIcon size={28} />
					</button>
				</form>
			</div>
		</div>
		<DailyCashflowTable transactions={dailyTransactions} />
	</div>
</div>
