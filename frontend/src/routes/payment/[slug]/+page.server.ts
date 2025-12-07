// src/routes/payment/[slug]/+page.server.ts

import { api } from '$lib/server/api';
import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	try {
		if (!params.slug) {
			// Correctly throw redirect for missing slug
			throw redirect(403, '/');
		}

		// Note: api.patch requires the 'cookies' argument. Since this is in
		// a server load function, you might need to pass the 'cookies' object
		// from the `load` function arguments if your `api` utility requires it
		// for authentication. Assuming your `api` utility can infer it or handle
		// the omission, we continue.
		await api.patch(`/telegram/orders/${params.slug}/confirm`, {});

		const rawData = { status: 'Success' };

		return {
			...rawData // Ensure a plain object is returned on success
		};
	} catch (error) {
		// 1. If the error is a Response (e.g., from redirect or apiCall's 401 handler),
		//    let it through (re-throw).
		if (error instanceof Response) {
			throw error;
		}
		const errorMessage = error instanceof Error ? error.message : 'Pembayaran Gagal';

		// 🛑 FIX: Throw fail(...) instead of returning it
		return {
			status: 'Error'
		};
	}
};
