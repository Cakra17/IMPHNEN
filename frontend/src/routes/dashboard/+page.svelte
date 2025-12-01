<script>
	import { CameraIcon, ChevronRightIcon, CirclePlusIcon, FileTextIcon, PlusIcon, TrendingDownIcon, TrendingUpIcon } from "@lucide/svelte";
	import { Card, Heading, Hr, P } from "flowbite-svelte";
    import { WhatsappSolid } from "flowbite-svelte-icons";
</script>

<div class="flex-1">
	<Heading class="text-lg md:text-2xl mb-2">Dashboard Keuangan</Heading>
	<P class="text-sm">Ringkasan hari ini, 1 Desember 2025</P>
</div>
<div class="w-full grid grid-cols-1 xl:grid-cols-4 gap-4 xl:gap-8 mt-8">
	<div
		class="xl:col-span-2 w-full flex flex-col rounded-xl gap-3 p-4 shadow-none hover:shadow-lg bg-teal-500 text-white"
	>
		<TrendingUpIcon />
		<div>
			<span class="text-sm">Pemasukan</span>
			<Heading class="text-lg text-white">Rp 400.000</Heading>
		</div>
	</div>
	<div
		class="xl:col-span-2 w-full flex flex-col rounded-xl gap-3 p-4 shadow-none hover:shadow-lg bg-red-500 text-white"
	>
		<TrendingDownIcon />
		<div>
			<span class="text-sm">Pemasukan</span>
			<Heading class="text-lg text-white">Rp 400.000</Heading>
		</div>
	</div>
	<div class="bg-white border border-teal-200 rounded-xl p-6 flex flex-col gap-3">
		<span class="font-bold mb-3">Aksi Cepat</span>
		<button
			onclick={() => console.log('TODO')}
			class="cursor-pointer flex-1 flex flex-row gap-4 p-4 items-center shadow-none ring-1 hover:ring-2 ring-teal-200 hover:ring-teal-300 rounded-lg transition-all"
		>
			<div class="p-3 bg-stone-100 rounded-full">
				<PlusIcon size={24} />
			</div>
			<div class="flex flex-col items-start">
				<p class="text-left font-semibold">Input Manual</p>
				<p class="text-left text-xs">Catat cashflow manual</p>
			</div>
		</button>
		<button
			onclick={() => console.log('TODO')}
			class="cursor-pointer flex-1 flex flex-row gap-4 p-4 items-center shadow-none ring-1 hover:ring-2 ring-teal-200 hover:ring-teal-300 rounded-lg transition-all"
		>
			<div class="p-3 bg-sky-100 rounded-full">
				<CameraIcon size={24} />
			</div>
			<div class="flex flex-col items-start">
				<p class="text-left font-semibold">Scan Struk AI</p>
				<p class="text-left text-xs">Tinggal foto, langsung kecatet</p>
			</div>
		</button>
		<button
			onclick={() => console.log('TODO')}
			class="cursor-pointer flex-1 flex flex-row gap-4 p-4 items-center shadow-none ring-1 hover:ring-2 ring-teal-200 hover:ring-teal-300 rounded-lg transition-all"
		>
			<div class="p-3 bg-emerald-100 rounded-full">
				<WhatsappSolid class="shrink-0 h-6 w-6" />
			</div>
			<div class="flex flex-col items-start">
				<p class="text-left font-semibold">Integrasi WA Business</p>
				<p class="text-left text-xs">Setup bot + Auto-Catat via WA</p>
			</div>
		</button>
	</div>
	<div class="xl:col-span-3 bg-white border border-teal-200 rounded-xl flex flex-col gap-3">
		<div class="flex flex-row gap-4 p-6 border-b-1 border-stone-200">
			<FileTextIcon class="text-teal-500" />
			<span class="flex-1 font-bold">Riwayat Transaksi</span>
			<a
				href="/dashboard/analytics"
				class="flex flex-row gap-2 font-bold text-teal-600 hover:text-teal-800 hover:underline cursor-pointer"
				>Analisa Lengkap <ChevronRightIcon /></a
			>
		</div>
	</div>
</div>
