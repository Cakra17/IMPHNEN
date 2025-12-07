<script>
	import { goto } from '$app/navigation';
	import Logo from '$lib/components/logo.svelte';
	import { formatDate } from '$lib/utils/date';
	import { BriefcaseIcon, TrendingDownIcon, TrendingUpIcon } from '@lucide/svelte';
	import { Button, Heading, P } from 'flowbite-svelte';
</script>

<div class="font-jakarta smooth-scroll">
	<div
		class="flex flex-row justify-center items-center px-4 md:px-8 lg:px-16 xl:px-36 py-4 md:py-8 lg:py-16 xl:py-36 gap-36"
	>
		<div class="flex-1 flex flex-col items-center lg:items-start">
			<div class="flex-1 flex flex-col py-16 justify-center">
				<Heading tag="h1" class="text-white">Kelola Bisnis,</Heading>
				<Heading tag="h1" class="text-emerald-400">Tanpa Mager.</Heading>
			</div>
			<Heading tag="h4" class="text-white italic mb-8"
				>Ingin Menjadi Pengusaha Handal Namun Enggan Ngebuku</Heading
			>
			<P class="text-white mb-8" size="lg"
				>Aplikasi pembukuan UMKM yang memahami kamu males ngetik. Tinggal foto struk, langsung
				tercatat. Orderan dari Telegram? Auto masuk. Semua cashflow terlihat jelas tanpa ribet.</P
			>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<Button pill size="xl" color="light" class="transition-all" href="/auth/register"
					>Mulai Sekarang</Button
				>
				<Button pill size="xl" color="blue" class="transition-all">Lihat Demo</Button>
			</div>
		</div>
		<div
			class="hidden lg:block flex-1 flex flex-col items-start bg-stone-50 rounded-2xl shadow-lg hover:shadow-2xl transition-all p-8"
		>
			<Heading class="text-xl md:text-2xl md:mb-2">Dashboard Keuangan</Heading>
			<P class="text-sm mb-8">Ringkasan hari ini, {formatDate(new Date())}</P>
			<div
				class="xl:col-span-2 w-full flex flex-col rounded-xl gap-3 p-4 shadow-none hover:shadow-lg bg-teal-500 text-white mb-4"
			>
				<TrendingUpIcon />
				<div>
					<span class="text-sm">Pemasukan</span>
					<Heading class="text-lg text-white">Rp 800.000</Heading>
				</div>
			</div>
			<div
				class="xl:col-span-2 w-full flex flex-col rounded-xl gap-3 p-4 shadow-none hover:shadow-lg bg-red-500 text-white"
			>
				<TrendingDownIcon />
				<div>
					<span class="text-sm">Pengeluaran</span>
					<Heading class="text-lg text-white">Rp 400.000</Heading>
				</div>
			</div>
		</div>
	</div>
	<!-- Problem Statement -->
	<section class="py-20 bg-white">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center mb-16">
				<h2 class="text-4xl font-bold text-gray-900 mb-4">Capek Ngebuku Manual? Kami Paham.</h2>
				<p class="text-xl text-gray-600">Masalah yang sering dialami UMKM Indonesia</p>
			</div>
			<div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
				<div
					class="bg-gradient-to-br from-red-50 to-orange-50 rounded-2xl p-6 card-hover border border-red-200"
				>
					<div class="text-4xl mb-4">😫</div>
					<h3 class="text-xl font-bold mb-2">Mager Nyatet</h3>
					<p class="text-gray-700 mb-2 italic">
						"Udah capek jualan, masih harus ngetik semua transaksi..."
					</p>
					<p class="text-gray-600">
						Setiap hari puluhan transaksi, tapi nulis manual? Males banget.
					</p>
				</div>
				<div
					class="bg-gradient-to-br from-yellow-50 to-amber-50 rounded-2xl p-6 card-hover border border-yellow-200"
				>
					<div class="text-4xl mb-4">😰</div>
					<h3 class="text-xl font-bold mb-2">Bingung Untung-Rugi</h3>
					<p class="text-gray-700 mb-2 italic">"Akhir bulan baru sadar, uangnya kemana?"</p>
					<p class="text-gray-600">
						Gak tahu bisnis untung atau rugi karena pencatatan acak-acakan.
					</p>
				</div>
				<div
					class="bg-gradient-to-br from-blue-50 to-cyan-50 rounded-2xl p-6 card-hover border border-blue-200"
				>
					<div class="text-4xl mb-4">😤</div>
					<h3 class="text-xl font-bold mb-2">Orderan Berantakan</h3>
					<p class="text-gray-700 mb-2 italic">"Pelanggan order via Chat, sering kelewat..."</p>
					<p class="text-gray-600">Chat campur jadi satu, orderan nyangkut, pelanggan komplain.</p>
				</div>
				<div
					class="bg-gradient-to-br from-purple-50 to-pink-50 rounded-2xl p-6 card-hover border border-purple-200"
				>
					<div class="text-4xl mb-4">🤯</div>
					<h3 class="text-xl font-bold mb-2">Waktu Terbuang</h3>
					<p class="text-gray-700 mb-2 italic">"Ngitung-ngitung manual sampai malem..."</p>
					<p class="text-gray-600">Harusnya bisa ngurusin bisnis, malah sibuk ngebuku.</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Solution Statement -->
	<section class="py-20 bg-stone-100">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center mb-16">
				<h2 class="text-4xl font-bold text-gray-900 mb-4">
					Biarin AI yang Kerja. Kamu Fokus Jualan Aja.
				</h2>
				<p class="text-xl text-gray-600">Fitur lengkap yang bikin hidup kamu lebih mudah</p>
			</div>
			<div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
				<div class="bg-white rounded-2xl p-8 card-hover border border-teal-200">
					<div class="w-16 h-16 bg-teal-100 rounded-full flex items-center justify-center mb-6">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="32"
							height="32"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="text-teal-600"
							><path
								d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"
							></path><circle cx="12" cy="13" r="3"></circle></svg
						>
					</div>
					<h3 class="text-2xl font-bold mb-3">Scan Struk AI</h3>
					<p class="text-gray-600 mb-4">
						Foto aja, langsung kecatet. AI otomatis ekstrak data dari struk belanja. Akurasi tinggi,
						cepat banget.
					</p>
					<ul class="space-y-2 text-sm text-gray-600">
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Upload foto struk
						</li>
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							AI ekstrak otomatis
						</li>
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Langsung masuk pembukuan
						</li>
					</ul>
				</div>

				<!-- Feature 2 -->
				<div class="bg-white rounded-2xl p-8 card-hover border border-teal-200">
					<div class="w-16 h-16 bg-emerald-100 rounded-full flex items-center justify-center mb-6">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="32"
							height="32"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="text-emerald-600"
							><line x1="12" y1="1" x2="12" y2="23"></line><path
								d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"
							></path></svg
						>
					</div>
					<h3 class="text-2xl font-bold mb-3">Dashboard Cashflow</h3>
					<p class="text-gray-600 mb-4">
						Lihat untung-rugi dalam sekejap. Real-time analytics dengan grafik visual yang mudah
						dibaca.
					</p>
					<ul class="space-y-2 text-sm text-gray-600">
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Pemasukan vs pengeluaran
						</li>
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Grafik visual interaktif
						</li>
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Filter harian/mingguan/bulanan
						</li>
					</ul>
				</div>

				<!-- Feature 3 -->
				<div class="bg-white rounded-2xl p-8 card-hover border border-teal-200">
					<div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mb-6">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="32"
							height="32"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="text-blue-600"
							><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg
						>
					</div>
					<h3 class="text-2xl font-bold mb-3">Telegram Bot</h3>
					<p class="text-gray-600 mb-4">
						Terima order otomatis via Telegram. Orderan langsung masuk sistem, gak ada yang
						terlewat.
					</p>
					<ul class="space-y-2 text-sm text-gray-600">
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Integrasi Telegram Business
						</li>
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Order masuk otomatis
						</li>
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Notifikasi real-time
						</li>
					</ul>
				</div>

				<!-- Feature 4 -->
				<div class="bg-white rounded-2xl p-8 card-hover border border-teal-200">
					<div
						class="w-16 h-16 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center mb-6"
					>
						<BriefcaseIcon size={28} />
					</div>
					<h3 class="text-2xl font-bold mb-3">Manajemen Produk</h3>
					<p class="text-gray-600 mb-4">
						Atur stok & katalog produk. Tambah produk dengan foto, atur harga dan stok dengan mudah.
					</p>
					<ul class="space-y-2 text-sm text-gray-600">
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Tambah stok dengan mudah
						</li>
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Integrasi dengan bot
						</li>
						<li class="flex items-start gap-2">
							<svg class="w-5 h-5 text-teal-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20"
								><path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								></path></svg
							>
							Form simpel
						</li>
					</ul>
				</div>
			</div>
		</div>
	</section>

	<section
		class="flex flex-col gap-8 justify-center items-center px-4 md:px-8 lg:px-16 xl:px-36 py-16"
	>
		<Heading tag="h4" class="text-white"
			>Tunggu apa lagi? Ayo mulai kemudahan perjalanan bisnis <span class="text-emerald-400"
				>Sekarang!</span
			></Heading
		>
		<Button pill size="xl" color="light" class="transition-all">Mulai Sekarang</Button>
	</section>

	<!-- Footer -->
	<section class="bg-stone-100">
		<P size="sm" class="text-teal-500 text-center py-4"
			>© 2025 IMPHNEN. All rights reserved. Powered by <a
				href="https://kolosal.ai"
				class="hover:text-teal-900">Kolosal.ai</a
			></P
		>
	</section>
</div>

<style>
	.card-hover {
		transition: all 0.3s ease;
	}

	.card-hover:hover {
		transform: translateY(-8px);
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
	}

	.smooth-scroll {
		scroll-behavior: smooth;
	}
</style>
