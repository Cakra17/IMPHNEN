<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { orderStatus, type Order, type OrderStatus } from '$lib/types/orders.js';
	import { formatDateTime } from '$lib/utils/date.js';
	import { CheckIcon, ChevronLeftIcon, MapPinIcon, PhoneCallIcon } from '@lucide/svelte';
	import {
		Badge,
		Button,
		Dropdown,
		DropdownItem,
		Heading,
		P,
		PaginationNav
	} from 'flowbite-svelte';
	import { ArrowLeftOutline, ArrowRightOutline, ChevronDownOutline } from 'flowbite-svelte-icons';

	let { data } = $props();
	let orders: Order[] = $derived(data.orders);
	let meta: PaginationMeta = $derived(data.meta);
	const statuses = Object.keys(orderStatus); // ['Pending', 'Diterima', 'Dicancel']

	let filters = $derived(data.filters);

	let statusFilter = $derived($page.url.searchParams.get('status'));

	// Filter handlers
	function applyFilters(status?: string) {
		const params = new URLSearchParams($page.url.searchParams);

		if (status) {
			params.set('status', status);
		} else {
			params.delete('status');
		}

		params.set('page', '1'); // Reset to page 1 when filtering
		goto(`?${params.toString()}`, { invalidate: ['/dashboard/orders'] });
	}

	function changePage(newPage: number) {
		const params = new URLSearchParams($page.url.searchParams);
		params.set('page', String(newPage));
		goto(`?${params.toString()}`);
	}
</script>

<a
	href="/dashboard"
	class="flex flex-row gap-2 font-bold text-teal-600 block lg:hidden hover:text-teal-800 hover:underline cursor-pointer"
	><ChevronLeftIcon />Kembali</a
>

<div class="flex flex-col md:flex-row justify-between items-start gap-4 pb-6">
	<div class="flex-1">
		<Heading class="text-lg md:text-2xl ml-8 lg:ml-0">Daftar Orderan</Heading>
	</div>
	<div class="flex-1 w-full flex flex-row items-center justify-start md:justify-end gap-4">
		<Button size="xs"
			>{statusFilter === null
				? 'Semua Order'
				: orderStatus[statusFilter as OrderStatus]}<ChevronDownOutline
				class="ms-2 h-6 w-6"
			/></Button
		>
		<Dropdown simple>
			<DropdownItem class={`flex flex-row gap-2`} onclick={() => applyFilters()}
				><CheckIcon class={statusFilter === null ? '' : 'opacity-0 text-transparent'} />Semua Order</DropdownItem
			>
			{#each statuses as status}
				<DropdownItem class={`flex flex-row gap-2`} onclick={() => applyFilters(status)}
					><CheckIcon
						class={statusFilter === status ? '' : 'opacity-0 text-transparent'}
					/>{orderStatus[status as OrderStatus]}</DropdownItem
				>
			{/each}
		</Dropdown>
		<PaginationNav
			currentPage={meta.page}
			totalPages={meta.total_page !== 0 ? meta.total_page : 1}
			onPageChange={changePage}
		>
			{#snippet prevContent()}
				<span class="sr-only">Previous</span>
				<ArrowLeftOutline class="h-5 w-5" />
			{/snippet}
			{#snippet nextContent()}
				<span class="sr-only">Next</span>
				<ArrowRightOutline class="h-5 w-5" />
			{/snippet}
		</PaginationNav>
	</div>
</div>

<!-- Orders List -->
{#if data.error}
	<div class="error text-red-600">{data.error}</div>
{/if}

{#if orders.length === 0}
	<div class="text-center py-12 text-gray-400">
		<p>No orders found</p>
	</div>
{:else}
	<div class="orders-list space-y-4">
		{#each orders as order}
			<div
				class="order-card flex flex-row bg-white border border-teal-200 rounded-xl overflow-hidden"
			>
				<div class="flex-1 flex flex-col bg-pattern p-4">
					<div class="flex flex-row justify-between">
						<p class="text-sm text-white">Order #{order.id.slice(0, 8)}</p>
						{#if order.status === 'pending'}
							<Badge large rounded color="yellow">Pending</Badge>
						{:else if order.status === 'confirmed'}
							<Badge large rounded color="green">Diterima</Badge>
						{:else}
							<Badge large rounded color="red">Dibatalkan</Badge>
						{/if}
					</div>
					<Heading tag="h5" class="my-2 text-white">{order.customer.name}</Heading>
					<P class="font-semibold text-white text-sm mb-1 flex gap-2 items-center"
						><PhoneCallIcon size={20} />{order.customer.phone}</P
					>
					<P class="font-semibold text-white text-sm flex gap-2 items-center mb-8"
						><MapPinIcon size={20} />{order.customer.address}</P
					>
					<P class="font-semibold text-sm text-white flex gap-2 items-center"
						>{formatDateTime(order.created_at)}</P
					>
				</div>
				<div class="flex-2 flex flex-col p-4">
					<div class="flex-1 text-sm">
						{#each order.order_items as item}
							<div class="flex justify-between">
								<Heading tag="h6">{item.product.name} x {item.quantity}</Heading>
								<Heading tag="h6">Rp {item.total_price.toLocaleString('id-ID')}</Heading>
							</div>
						{/each}
					</div>
					<div class="mt-3 pt-3 border-t border-teal-400 flex justify-between font-bold">
						<span>Total</span>
						<span>Rp {order.total_price.toLocaleString('id-ID')}</span>
					</div>
				</div>
			</div>
		{/each}
	</div>
{/if}
