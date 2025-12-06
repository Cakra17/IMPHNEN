<script lang="ts">
	import {
		CalendarIcon,
		CameraIcon,
		ChevronRightIcon,
		CirclePlusIcon,
		FileTextIcon,
		PlusIcon,
		TrendingDownIcon,
		TrendingUpIcon
	} from '@lucide/svelte';
	import { Button, Card, Heading, Hr, Input, Label, Modal, P, Select } from 'flowbite-svelte';
	import telegramLogo from '../../lib/assets/telegram_logo.svg';
	import type { PageData } from '../$types';
	import { formatRupiah } from '$lib/utils/format_money';
	import { formatDate, formatDateSlash, formatRelativeTime } from '$lib/utils/date';

	let { data } = $props<{ data: PageData }>();

	// Access the summary data
	const summary = data.summary;
	const transactions = data.transactions;

	let txTypes: { value: string; name: string }[] = [
		{ value: 'expense', name: 'Pengeluaran' },
		{ value: 'income', name: 'Pemasukan' }
	];

	let inputModal = $state(false);
	let txValue: number = $state(0);
	let selectedTx: string = $state('expense');
	let selectedDate: string = $state(new Date().toISOString().split('T')[0]);

	console.log(data);
</script>

<div class="flex-1">
	<Heading class="text-xl md:text-2xl md:mb-2">Dashboard Keuangan</Heading>
	<P class="text-sm">Ringkasan hari ini, 1 Desember 2025</P>
</div>
<div class="w-full grid grid-cols-1 xl:grid-cols-4 gap-4 xl:gap-6 mt-8">
	<div
		class="xl:col-span-2 w-full flex flex-col rounded-xl gap-3 p-4 shadow-none hover:shadow-lg bg-teal-500 text-white"
	>
		<TrendingUpIcon />
		<div>
			<span class="text-sm">Pemasukan</span>
			<Heading class="text-lg text-white">{formatRupiah(summary?.total_income) ?? 'Rp 0'}</Heading>
		</div>
	</div>
	<div
		class="xl:col-span-2 w-full flex flex-col rounded-xl gap-3 p-4 shadow-none hover:shadow-lg bg-red-500 text-white"
	>
		<TrendingDownIcon />
		<div>
			<span class="text-sm">Pengeluaran</span>
			<Heading class="text-lg text-white">{formatRupiah(summary?.total_expense) ?? 'Rp 0'}</Heading>
		</div>
	</div>
	<div class="bg-white border border-teal-200 rounded-xl p-4 xl:p-5 flex flex-col gap-3">
		<span class="font-bold mb-1 md:mb-3">Aksi Cepat</span>
		<button
			onclick={() => (inputModal = true)}
			class="cursor-pointer flex-1 flex flex-col md:flex-row gap-4 p-4 items-center shadow-none ring-1 hover:ring-2 ring-teal-200 hover:ring-stone-300 rounded-lg transition-all"
		>
			<div class="w-12 h-12 flex items-center justify-center bg-stone-100 rounded-full">
				<PlusIcon size={24} />
			</div>
			<div class="flex flex-col items-center md:items-start">
				<p class="text-center md:text-left font-semibold">Input Manual</p>
				<p class="text-center md:text-left text-xs">Catat cashflow manual</p>
			</div>
		</button>
		<button
			onclick={() => console.log('TODO')}
			class="cursor-pointer flex-1 flex flex-col md:flex-row gap-4 p-4 items-center shadow-none ring-1 hover:ring-2 ring-teal-200 hover:ring-emerald-300 rounded-lg transition-all"
		>
			<div class="w-12 h-12 flex items-center justify-center bg-emerald-100 rounded-full">
				<CameraIcon size={24} />
			</div>
			<div class="flex flex-col items-center md:items-start">
				<p class="text-center md:text-left font-semibold">Scan Struk AI</p>
				<p class="text-center md:text-left text-xs">Tinggal foto, langsung kecatet</p>
			</div>
		</button>
		<button
			onclick={() => console.log('TODO')}
			class="cursor-pointer flex-1 flex flex-col md:flex-row gap-4 p-4 items-center shadow-none ring-1 hover:ring-2 ring-teal-200 hover:ring-sky-300 rounded-lg transition-all"
		>
			<div class="bg-emerald-100 rounded-full">
				<img src={telegramLogo} alt="Telegram" class="shrink-0 h-12 w-12" />
			</div>
			<div class="flex flex-col items-center md:items-start">
				<p class="text-center md:text-left font-semibold">Integrasi Telegram Bot</p>
				<p class="text-center md:text-left text-xs">Setup bot + terima pesanan via Telegram</p>
			</div>
		</button>
	</div>
	<div class="xl:col-span-3 bg-white border border-teal-200 rounded-xl flex flex-col">
		<div class="flex flex-row items-center gap-4 p-3 md:p-4 border-b-1 border-stone-200">
			<FileTextIcon class="text-teal-500 hidden md:inline" />
			<span class="flex-1 font-bold">Riwayat Transaksi</span>
			<a
				href="/dashboard/analytics"
				class="flex flex-row gap-1 md:gap-2 text-sm md:text-md text-right items-center font-bold text-teal-600 hover:text-teal-800 hover:underline cursor-pointer"
				><span class="hidden md:inline">Analisa Lengkap</span> <ChevronRightIcon /></a
			>
		</div>
		{#each transactions as transaction (transaction.id)}
			<div
				class="p-4 flex flex-row gap-6 items-center border-b border-teal-300 hover:bg-stone-50 active:bg-stone-50/30"
			>
				<div class="flex-1">
					<!-- TODO: Handle manual / bot / scan input -->
					<Heading tag="h6" class="text-left"
						>{transaction.source === 'manual' ? 'Input Manual' : 'TODO'}</Heading
					>
					<span class="text-sm mb-2">{formatRelativeTime(transaction.transaction_date, 1)}</span>
				</div>
				<Heading
					class={`text-xl mb-2 ${transaction.type === 'expense' ? 'text-red-600' : 'text-emerald-500'} self-start`}
					>{transaction.type === 'expense' ? '-' : '+'} {formatRupiah(transaction.amount)}</Heading
				>
			</div>
		{:else}
			<div class="flex w-full h-full justify-center items-center">
				<P>Belum ada transaksi hari ini.</P>
			</div>
		{/each}
	</div>
</div>

<Modal
	title="Tambah Transaksi Manual"
	form
	bind:open={inputModal}
	onaction={({ action }) => alert(`Handle "${action}"`)}
>
	<form method="POST" action="?/createManualTransaction">
		<div class="grid gap-6 grid-cols-1">
			<div>
				<Label for="amount" class="mb-2">Total</Label>
				<Input
					class="ps-10"
					type="number"
					id="amount"
					name="amount"
					placeholder={'100000'}
					bind:value={txValue}
					required
				>
					{#snippet left()}
						<P>Rp</P>
					{/snippet}
				</Input>
			</div>
			<div>
				<Label for="transaction_date" class="mb-2">Tanggal</Label>
				<Input
					class="ps-10"
					type="date"
					id="transaction_date"
					name="transaction_date"
					bind:value={selectedDate}
					required
				>
					{#snippet left()}
						<CalendarIcon />
					{/snippet}
				</Input>
			</div>
			<div>
				<Label for="transaction_date" class="mb-2">Tipe Transaksi</Label>
				<Select class="mt-2" items={txTypes} bind:value={selectedTx} id="type" name="type" />
			</div>
			<Button type="submit" disabled={txValue < 1 || selectedDate === '' || selectedTx === ''}
				>Tambah Transaksi</Button
			>
		</div>
	</form>

	{#snippet footer()}
		<div></div>
	{/snippet}
</Modal>
