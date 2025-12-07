// +page.server.ts
import { api } from '$lib/server/api';
import type { OrdersResponse } from '$lib/types/orders';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const perPage = Number(url.searchParams.get('per_page')) || 10;
	const customerId = url.searchParams.get('customer_id') || undefined;
	const status = url.searchParams.get('status') || undefined;

	try {
		// Build query params
		const params: Record<string, string | number> = {
			page,
			per_page: perPage
		};

		// Add optional filters if present
		if (status) params.status = status;

		const response: OrdersResponse = await api.get('/orders', cookies, params);

		return {
			orders: response.data || [],
			meta: response.meta || {
				page: 1,
				per_page: 10,
				total_data: 0,
				total_page: 0
			},
			filters: {
				customerId,
				status
			}
		};
	} catch (error) {
		console.error('Failed to load orders:', error);
		return {
			orders: [],
			meta: {
				page: 1,
				per_page: 10,
				total_data: 0,
				total_page: 0
			},
			filters: {
				customerId,
				status
			},
			error: error instanceof Error ? error.message : 'Failed to load orders'
		};
	}
};
