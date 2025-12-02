<script>
	import CashflowChart from "$lib/components/cashflow_chart.svelte";
	import { ChevronDownIcon, ChevronLeftIcon, ChevronRightIcon, TrendingDownIcon, TrendingUpIcon } from "@lucide/svelte";
	import { Button, Dropdown, DropdownItem, Heading, Hr, P } from "flowbite-svelte";
	import { ChevronDownOutline } from "flowbite-svelte-icons";

</script>

<a
	href="/dashboard"
	class="flex flex-row gap-2 font-bold text-teal-600 hover:text-teal-800 hover:underline cursor-pointer"
	><ChevronLeftIcon />Kembali</a
>
<div class="flex-1">
	<Heading class="text-lg md:text-2xl ml-8 mb-8">Analisa & Statistik</Heading>
	<CashflowChart />

	<div class="mt-8 flex flex-col bg-white border border-teal-200 rounded-2xl">
		<div class="flex justify-between border-b-1 p-6 pb-3 border-teal-200">
			<div class="flex flex-col">
				<Heading class="text-xl mb-1">Riwayat Transaksi Harian</Heading>
				<span class="text-sm mb-2">Detail cashflow lengkap perhari</span>
			</div>
			<div class="flex flex-row gap-3 items-center">
				<button
					class="text-stone-500 hover:text-stone-800 active:text-stone-800 active:bg-stone-100 rounded-full cursor-pointer"
					><ChevronLeftIcon size={28} /></button
				>
				<Heading class="text-md mb-1">1 November 2025</Heading>
				<button
					class="text-stone-500 hover:text-stone-800 active:text-stone-800 active:bg-stone-100 rounded-full cursor-pointer"
					><ChevronRightIcon size={28} /></button
				>
			</div>
		</div>
		<div class="grid grid-cols-2 overflow-hidden">
			<div
				class="pt-4 p-6 border-r border-teal-200 bg-gradient-to-b from-emerald-50 via-white to-white"
			>
				<div class="pb-4 flex flex-row justify-between">
					<div>
						<Heading class="text-lg">Pemasukan</Heading>
						<Heading class="text-xl mb-2 text-emerald-500">Rp 40.000</Heading>
					</div>
					<TrendingUpIcon size={32} />
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
			<div class="pt-4 p-6 bg-gradient-to-b from-red-100 via-white to-white">
				<div class="pb-4 flex flex-row justify-between">
					<div>
						<Heading class="text-lg">Pengeluaran</Heading>
						<Heading class="text-xl mb-2 text-red-500">Rp 10.000</Heading>
					</div>
					<TrendingDownIcon size={32} />
				</div>
				<div>
					<div class="rounded-lg border border-teal-300 overflow-hidden bg-white">
						<div
							class="p-4 flex flex-row gap-6 items-center border-b border-teal-300 hover:bg-stone-100/50 active-bg-stone-50/30"
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
	</div>
</div>
