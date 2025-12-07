<script lang="ts">
	import {
		CalendarIcon,
		CameraIcon,
		ChevronRightIcon,
		CirclePlusIcon,
		FileTextIcon,
		LoaderIcon,
		PlusIcon,
		TrendingDownIcon,
		TrendingUpIcon
	} from '@lucide/svelte';
	import {
		Button,
		ButtonGroup,
		Card,
		Datepicker,
		Fileupload,
		Heading,
		Hr,
		Input,
		Label,
		Modal,
		P,
		Select
	} from 'flowbite-svelte';
	import telegramLogo from '../../lib/assets/telegram_logo.svg';
	import type { PageData } from '../$types';
	import { formatRupiah } from '$lib/utils/format_money';
	import { formatDate, formatDateSeparator, formatRelativeTime } from '$lib/utils/date';
	import { ArrowDownToBracketOutline, ArrowUpFromBracketOutline } from 'flowbite-svelte-icons';
	import { page } from '$app/stores';
	import { applyAction, deserialize, enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import { transactionsSources } from '$lib/types/transaction';
	import type { Transaction } from '$lib/types/transaction';

	let { data } = $props<{ data: PageData }>();

	// Access the summary data
	const summary = data.summary;
	const transactions: Transaction[] = data.transactions;

	// Manual input modal
	let inputModal = $state(false);
	let txValue: number = $state(0);
	let selectedType: string = $state('expense');
	let selectedDate: string = $state(new Date().toISOString().split('T')[0]);

	// Receipt scan modal
	let scanModal = $state(false);
	let image: File | null = $state(null);
	let previewUrl: string | null = $state(null);

	let loading = $state(false);
	let error: string | null = $state(null);

	function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		image = target.files?.[0] ?? null;

		if (image) {
			if (previewUrl) URL.revokeObjectURL(previewUrl);
			previewUrl = URL.createObjectURL(image);
		}
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		if (!image) return;

		loading = true;
		error = null;

		try {
			const formData = new FormData();
			formData.append('image', image);

			// Call the SvelteKit form action properly
			const response = await fetch(window.location.pathname + '?/postReceiptImage', {
				method: 'POST',
				body: formData,
				headers: {
					'x-sveltekit-action': 'true'
				}
			});

			const result = deserialize(await response.text());

			if (result.type === 'success') {
				await applyAction(result);

				// Reset form
				image = null;
				previewUrl = null;
				scanModal = false;

				await invalidateAll();
			} else if (result.type === 'failure') {
				error = (result.data?.message as string) || 'Upload gagal';
			}
		} catch (err) {
			console.error('Upload error:', err);
			error = err instanceof Error ? err.message : 'Terjadi kesalahan';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex-1">
	<Heading class="text-xl md:text-2xl md:mb-2">Dashboard Keuangan</Heading>
	<P class="text-sm">Ringkasan hari ini, {formatDate(new Date())}</P>
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
	<div class="bg-white border border-teal-200 rounded-xl p-4 xl:p-5 flex flex-col gap-3 self-start">
		<span class="font-bold mb-1 md:mb-3">Aksi Cepat</span>
		<button
			onclick={() => (inputModal = true)}
			class="cursor-pointer flex-1 flex flex-col md:flex-row gap-4 p-4 items-center shadow-none ring-1 hover:ring-2 ring-teal-200 hover:ring-stone-300 rounded-lg transition-all"
		>
			<div class="shrink-0 w-12 h-12 flex items-center justify-center bg-stone-100 rounded-full">
				<PlusIcon size={24} />
			</div>
			<div class="flex flex-col items-center md:items-start">
				<p class="text-center md:text-left font-semibold">Input Manual</p>
				<p class="text-center md:text-left text-xs">Catat cashflow manual</p>
			</div>
		</button>
		<button
			onclick={() => (scanModal = true)}
			class="cursor-pointer flex-1 flex flex-col md:flex-row gap-4 p-4 items-center shadow-none ring-1 hover:ring-2 ring-teal-200 hover:ring-emerald-300 rounded-lg transition-all"
		>
			<div class="shrink-0 w-12 h-12 flex items-center justify-center bg-emerald-100 rounded-full">
				<CameraIcon size={24} />
			</div>
			<div class="flex flex-col items-center md:items-start">
				<p class="text-center md:text-left font-semibold">Scan Struk AI</p>
				<p class="text-center md:text-left text-xs">Tinggal foto, langsung kecatet</p>
			</div>
		</button>
		<button
			class="cursor-pointer flex-1 flex flex-col md:flex-row gap-4 p-4 items-center shadow-none ring-1 hover:ring-2 ring-teal-200 hover:ring-sky-300 rounded-lg transition-all"
		>
			<div class="shrink-0 bg-emerald-100 rounded-full">
				<img src={telegramLogo} alt="Telegram" class="shrink-0 h-12 w-12" />
			</div>
			<div class="flex flex-col items-center md:items-start">
				<p class="text-center md:text-left font-semibold">@tuankrebbot</p>
				<p class="text-center md:text-left text-xs">
					Bot integrasi Telegram. <br /><i>Integrasi custom Bot coming soon</i>
				</p>
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
		{#each transactions as tx (tx.id)}
			<div
				class="p-4 flex flex-row gap-6 items-center border-b border-teal-300 hover:bg-stone-50 active:bg-stone-50/30"
			>
				<div class="flex-1">
					<!-- TODO: Handle manual / bot / scan input -->
					<Heading tag="h6" class="text-left">{transactionsSources[tx.source]}</Heading>
					<span class="text-sm mb-2">{formatRelativeTime(tx.transaction_date, 1)}</span>
				</div>
				<Heading
					class={`text-xl mb-2 ${tx.type === 'expense' ? 'text-red-600' : 'text-emerald-500'} self-start`}
					>{tx.type === 'expense' ? '-' : '+'} {formatRupiah(tx.amount)}</Heading
				>
			</div>
		{:else}
			<div class="flex w-full h-full justify-center items-center">
				<P>Belum ada transaksi hari ini.</P>
			</div>
		{/each}
	</div>
</div>

<!-- Popup Modal Input Manual -->
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
				<ButtonGroup class="*:ring-primary-700! w-full">
					<Button
						class={`flex-1`}
						color={selectedType === 'income' ? 'green' : 'alternative'}
						onclick={() => (selectedType = 'income')}
					>
						<ArrowDownToBracketOutline />
						Pemasukan</Button
					>
					<Button
						class={`flex-1`}
						color={selectedType === 'expense' ? 'red' : 'alternative'}
						onclick={() => (selectedType = 'expense')}
					>
						<ArrowUpFromBracketOutline />
						Pengeluaran</Button
					>
				</ButtonGroup>
				<input type="hidden" name="type" value={selectedType} />
			</div>
			<Button type="submit" disabled={txValue < 1 || selectedDate === '' || selectedType === ''}
				>Tambah Transaksi</Button
			>
		</div>
	</form>

	{#snippet footer()}
		<div></div>
	{/snippet}
</Modal>

<!-- Popup Modal Scan Struk -->
<Modal
	title="Scan Struk Pembelian"
	form
	bind:open={scanModal}
	onaction={({ action }) => alert(`Handle "${action}"`)}
>
	<form method="POST" onsubmit={handleSubmit}>
		<div class="grid gap-6 grid-cols-1">
			<div>
				<Label for="amount" class="mb-2">Total</Label>
				<Fileupload
					id="image"
					name="image"
					accept="image/jpeg, image/png"
					onchange={handleFileSelect}
				/>
			</div>
			<div>
				<div class="h-64 flex items-center justify-center border border-teal-200 rounded-lg">
					{#if previewUrl}
						<img src={previewUrl} alt="preview" class="h-full" />
					{:else}
						<P>Tampilan gambar akan muncul disini</P>
					{/if}
				</div>
			</div>
			<Button type="submit" disabled={image === null || loading}>
				{#if loading}
					<div class="contents">
						<LoaderIcon class="animate-[spin_3s_linear_infinite]" /> Memproses...
					</div>
				{:else}
					Scan Transaksi
				{/if}
			</Button>
		</div>
	</form>

	{#snippet footer()}
		<div></div>
	{/snippet}
</Modal>
