<script lang="ts">
	import {
		Button,
		Heading,
		Input,
		Label,
		Modal,
		P,
		PaginationNav,
		Toast,
		Tooltip
	} from 'flowbite-svelte';
	import type { PageData } from '../$types';
	import type { UserData } from '$lib/types/user';
	import { ChevronLeftIcon, LoaderIcon, PencilIcon, PlusIcon, TrashIcon } from '@lucide/svelte';
	import { enhance } from '$app/forms';
	import { CheckCircleSolid, CloseCircleSolid } from 'flowbite-svelte-icons';
	import { formatRupiah } from '$lib/utils/format_money';
	import { invalidate } from '$app/navigation';
	import { redirect } from '@sveltejs/kit';

	let { data } = $props<{ data: PageData }>();

	let user: UserData = data?.user;
	let products: Product[] = data?.products;
	let product_meta: PaginationMeta = data?.products_meta;

	// Profile edit vars
	let profileEditModal = $state(false);
	let editStatus: boolean | null = $state(null);

	// Profile delete vars
	let userDeleteModal = $state(false);
	let deleteVerify = $state('');

	let submitting = $state(false);
	let submitStatus: boolean | null = $state(null);
	// Product add vars
	let productsAddModal = $state(false);

	// Single product edit vars
	let productEditModal = $state(false);
	let selectedProductIndex: number | null = $state(null);
</script>

<a
	href="/dashboard"
	class="flex flex-row gap-2 font-bold text-teal-600 block lg:hidden hover:text-teal-800 hover:underline cursor-pointer"
	><ChevronLeftIcon />Kembali</a
>
<Heading class="text-lg md:text-2xl ml-8 lg:ml-0 mb-8">Pengaturan UMKM</Heading>

<div class="w-full grid grid-cols-1 xl:grid-cols-4 gap-4 xl:gap-6">
	<!-- User Profile Section -->
	<div class="flex flex-col gap-4 xl:gap-6 self-start">
		<div
			class="w-full bg-white border border-teal-200 rounded-xl flex flex-col self-start overflow-hidden"
		>
			<div
				class="p-4 xl:p-5 h-18 border-b border-teal-200 flex flex-row justify-between items-center"
			>
				<Heading tag="h5" class="">Profil Pengguna</Heading>
				<button
					onclick={() => (profileEditModal = true)}
					class="text-stone-600 hover:text-stone-900 active:text-emerald-700 rounded-full"
					><PencilIcon size={20} /></button
				>
				<Tooltip>Edit Profil</Tooltip>
			</div>
			<div class="bg-white p-4 flex flex-col w-full gap-6 break-words">
				<div>
					<Heading tag="h6" class="pb-1">Nama Lengkap</Heading>
					<P>{user.firstname} {user.lastname}</P>
				</div>
				<div>
					<Heading tag="h6" class="pb-1">Nama UMKM</Heading>
					<P>{user.store_name}</P>
				</div>
				<div>
					<Heading tag="h6" class="pb-1">Email</Heading>
					<P>{user.email}</P>
				</div>
			</div>
		</div>
		<button
			onclick={() => (userDeleteModal = true)}
			class="cursor-pointer flex-1 flex flex-col md:flex-row gap-4 p-4 text-red-900 bg-red-50 items-center shadow-none ring-1 hover:ring-2 ring-red-300 rounded-lg transition-all"
		>
			<div class="w-12 h-12 flex items-center justify-center bg-red-600 text-red-50 rounded-full">
				<TrashIcon size={24} />
			</div>
			<div class="flex flex-col items-center md:items-start">
				<p class="text-center md:text-left font-semibold">Hapus Akun</p>
				<p class="text-center md:text-left text-xs">Hapus akun secara permanen</p>
			</div>
		</button>
	</div>

	<!-- Store Product Section -->
	<div
		class="bg-white col-span-1 xl:col-span-3 border border-teal-200 rounded-xl flex flex-col gap-3 self-start overflow-hidden"
	>
		<div class="p-4 h-18 border-b border-teal-200 flex flex-row justify-between items-center">
			<div class="flex flex-row gap-2 items-center">
				<Heading tag="h5">Daftar Produk</Heading>
				<P class="text-lg">({product_meta.total_data})</P>
			</div>
			<button
				onclick={() => (productsAddModal = true)}
				class="text-stone-600 hover:text-stone-900 active:text-emerald-700 rounded-full"
				><PlusIcon size={20} /></button
			>
			<Tooltip>Tambah Produk</Tooltip>
		</div>
		<div class="flex flex-row justify-center">
			<PaginationNav
				currentPage={product_meta.page}
				totalPages={product_meta.total_page}
				onPageChange={() => {
					/*TODO*/
				}}
			/>
		</div>
		{#if products.length > 0}
			<div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-4 items-center place-items-stretch">
				{#each products as product, i}
					<div class="w-full h-full flex flex-col border border-teal-200 rounded-lg items-center">
						<div class="flex-1">
							<img
								src={product.image_url}
								alt={`Gambar ${product.name}`}
								class="h-36 m:h-56 xl:h-64"
							/>
						</div>
						<div class="flex flex-row w-full bg-stone-50 border-t border-teal-200">
							<div class="p-4 flex flex-col items-start">
								<Heading tag="h5">{product.name}</Heading>
								<P class="font-semibold">{formatRupiah(product.price)}</P>
								<P class="text-sm">Sisa stok: {product.stock}</P>
							</div>
							<div class="flex-1 p-4 flex flex-col items-end">
								<button
									onclick={() => {
										productEditModal = true;
										selectedProductIndex = i;
									}}
									class="text-stone-600 hover:text-stone-900 active:text-emerald-700 rounded-full"
									><PencilIcon size={20} /></button
								>
								<Tooltip>Edit {product.name}</Tooltip>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{:else}
			<div class="py-12 flex flex-col gap-4 items-center">
				<P align="center">Toko mu belum memiliki produk.</P><Button
					size="sm"
					onclick={() => (productsAddModal = true)}><PlusIcon class="me-2" /> Tambah Produk</Button
				>
			</div>
		{/if}
	</div>
</div>

<!-- Profile Edit Modal -->
<Modal
	bind:open={profileEditModal}
	title="Edit Profil"
	form
	onaction={({ action }) => alert(`Handle "${action}"`)}
>
	<form
		method="POST"
		action="?/updateProfile"
		use:enhance={() => {
			return async ({ result }) => {
				if (result.type === 'success') {
					editStatus = true;
					setTimeout(() => window.location.reload(), 1000);
				} else if (result.type === 'failure') {
					editStatus = false;
				}
			};
		}}
	>
		<div class="flex flex-col items-center gap-6">
			<div class="w-full">
				<Label for="email" class="mb-2 text-md">Email</Label>
				<P>{user.email}</P>
			</div>
			<div class="w-full grid grid-cols-1 md:grid-cols-2 gap-6">
				<div class="w-full">
					<Label for="firstname" class="mb-2">Nama Depan</Label>
					<Input type="text" name="firstname" defaultValue={user.firstname} required></Input>
				</div>
				<div class="w-full">
					<Label for="lastname" class="mb-2">Nama Belakang</Label>
					<Input type="text" name="lastname" defaultValue={user.lastname} required></Input>
				</div>
			</div>
			<div class="w-full">
				<Label for="store_name" class="mb-2">Nama UMKM</Label>
				<Input type="text" name="store_name" defaultValue={user.store_name} required></Input>
			</div>
			{#if editStatus === true}
				<Toast color="green" class="w-full" dismissable={false}>
					{#snippet icon()}
						<CheckCircleSolid class="h-5 w-5" />
						<span class="sr-only">Check icon</span>
					{/snippet}
					Profile Telah Dirubah.
				</Toast>
			{:else if editStatus === false}
				<Toast color="red" class="w-full" dismissable={false}>
					{#snippet icon()}
						<CloseCircleSolid class="h-5 w-5" />
						<span class="sr-only">Close icon</span>
					{/snippet}
					Gagal mmerubah profil.
				</Toast>
			{/if}
			<Button type="submit" disabled={editStatus === true}>Edit Profil</Button>
		</div>
	</form>
</Modal>

<!-- Profile Delete Modal -->
<Modal title="Penghapusan Akun" bind:open={userDeleteModal}>
	<form
		method="POST"
		action="?/deleteProfile"
		use:enhance={() => {
			return async ({ result }) => {
				if (result.type === 'redirect') {
					// Let the redirect happen
					window.location.href = result.location;
				} else if (result.type === 'failure') {
					alert(result.data?.error || 'Failed to delete account');
				}
			};
		}}
	>
		<P>
			Semua data terhadap UMKM <strong>{user.store_name}</strong> akan dihapus termasuk: Daftar transaksi,
			Daftar produk, Daftar order, dan Integrasi Telegram Bot
		</P>
		<P>
			Semua orderan, termasuk yang belum kamu selesaikan akan terhapus. Pastikan untuk menyelesaikan
			atau mencatat manual orderan jika perlu!
		</P>
		<P class="text-red-900 font-semibold pb-4">Tindakan ini tidak dapat dibatalkan.</P>
		<P>Untuk menghapus akun, ketik nama UMKM kamu pada kolom dibawah.</P>
		<div class="w-full pb-6">
			<Input type="text" name="verification" bind:value={deleteVerify}></Input>
		</div>
		<Button type="submit" color="red" class="me-2" disabled={deleteVerify !== user.store_name}>
			Hapus Akun
		</Button>
		<Button
			type="button"
			onclick={() => {
				deleteVerify = '';
				userDeleteModal = false;
			}}
			color="alternative"
		>
			Batalkan
		</Button>
	</form>
</Modal>

<!-- Product Add Modal -->
<Modal
	bind:open={productsAddModal}
	title="Tambah Produk"
	form
	onaction={({ action }) => alert(`Handle "${action}"`)}
>
	<form
		method="POST"
		action="?/submitProduct"
		enctype="multipart/form-data"
		use:enhance={() => {
			submitting = true;

			return async ({ result }) => {
				submitting = false;

				if (result.type === 'success') {
					submitStatus = true;
					setTimeout(() => window.location.reload(), 1000);
				} else if (result.type === 'failure') {
					submitStatus = false;
				}
			};
		}}
	>
		<div class="flex flex-col items-center gap-6">
			<div class="w-full">
				<Label for="name" class="mb-2">Nama</Label>
				<Input type="text" name="name" placeholder="Produk Saya" required></Input>
			</div>
			<div class="w-full">
				<Label for="price" class="mb-2">Harga</Label>
				<Input class="ps-10" type="number" name="price" placeholder="100000" required>
					{#snippet left()}
						<P>Rp</P>
					{/snippet}
				</Input>
			</div>
			<div class="w-full">
				<Label for="stock" class="mb-2">Stok</Label>
				<Input type="number" name="stock" placeholder="100" required></Input>
			</div>
			<div class="w-full">
				<Label for="image" class="mb-2">Gambar</Label>
				<Input type="file" name="image" accept="image/jpeg, image/png" required></Input>
			</div>
			{#if submitStatus === true}
				<Toast color="green" class="w-full" dismissable={false}>
					{#snippet icon()}
						<CheckCircleSolid class="h-5 w-5" />
						<span class="sr-only">Check icon</span>
					{/snippet}
					Produk berhasil ditambahkan.
				</Toast>
			{:else if submitStatus === false}
				<Toast color="red" class="w-full" dismissable={false}>
					{#snippet icon()}
						<CloseCircleSolid class="h-5 w-5" />
						<span class="sr-only">Close icon</span>
					{/snippet}
					Gagal menambahkan produk.
				</Toast>
			{/if}
			<Button type="submit" disabled={submitting || submitStatus === true}>
				{#if submitting}
					<div class="contents">
						<LoaderIcon class="animate-[spin_3s_linear_infinite]" /> Menambahkan produk...
					</div>
				{:else}
					Tambahkan Produk
				{/if}
			</Button>
		</div>
	</form>
</Modal>

<!-- Product Edit Modal -->
<Modal
	bind:open={productEditModal}
	title="Edit Produk"
	form
	onaction={({ action }) => alert(`Handle "${action}"`)}
>
	{#if selectedProductIndex !== null}
		<form
			method="POST"
			action="?/updateProduct"
			enctype="multipart/form-data"
			use:enhance={() => {
				submitting = true;

				return async ({ result }) => {
					submitting = false;

					if (result.type === 'success') {
						submitStatus = true;
						setTimeout(() => window.location.reload(), 1000);
					} else if (result.type === 'failure') {
						submitStatus = false;
					}
				};
			}}
		>
			<input type="text" name="id" value={products[selectedProductIndex].id} hidden />
			<div class="flex flex-col items-center gap-6">
				<div class="w-full">
					<Label for="name" class="mb-2">Nama</Label>
					<Input type="text" name="name" placeholder={products[selectedProductIndex].name}></Input>
				</div>
				<div class="w-full">
					<Label for="price" class="mb-2">Harga</Label>
					<Input
						class="ps-10"
						type="number"
						name="price"
						placeholder={products[selectedProductIndex].price.toString()}
					>
						{#snippet left()}
							<P>Rp</P>
						{/snippet}
					</Input>
				</div>
				<div class="w-full">
					<Label for="stock" class="mb-2">Stok*</Label>
					<Input
						type="number"
						name="stock"
						defaultValue={products[selectedProductIndex].stock.toString()}
						required
					></Input>
				</div>
				<div class="w-full">
					<Label for="image" class="mb-2">Gambar</Label>
					<Input type="file" name="image" accept="image/jpeg, image/png"></Input>
				</div>
				{#if submitStatus === true}
					<Toast color="green" class="w-full" dismissable={false}>
						{#snippet icon()}
							<CheckCircleSolid class="h-5 w-5" />
							<span class="sr-only">Check icon</span>
						{/snippet}
						Produk berhasil dirubah.
					</Toast>
				{:else if submitStatus === false}
					<Toast color="red" class="w-full" dismissable={false}>
						{#snippet icon()}
							<CloseCircleSolid class="h-5 w-5" />
							<span class="sr-only">Close icon</span>
						{/snippet}
						Gagal merubah produk.
					</Toast>
				{/if}
				<Button type="submit" disabled={submitting || submitStatus === true}>
					{#if submitting}
						<div class="contents">
							<LoaderIcon class="animate-[spin_3s_linear_infinite]" /> Menambahkan produk...
						</div>
					{:else}
						Edit Produk
					{/if}
				</Button>
			</div>
		</form>
	{/if}
</Modal>
