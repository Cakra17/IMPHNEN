// src/routes/payment/[slug]/+page.server.ts

import { api } from '$lib/server/api';
import { fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const actions: Actions = {
	// We define a default action, but you could name it 'confirm' if your form
	// submission used: <form method="POST">
	default: async ({ params }) => {
		// NOTE: The 'cookies' object is needed here because your 'api.patch'
		// likely relies on it for authentication.

		try {
			if (!params.slug) {
				// For actions, if a crucial parameter is missing, we use 'fail'
				// instead of redirecting the user away, giving them feedback.
				return fail(400, {
					success: false,
					message: 'Missing required order identifier.'
				});
			}

			// 1. Call the API using the slug from params.
			//    (Assuming api.patch can infer/handle the cookies argument here,
			//    otherwise you must pass it if your utility requires it).
			await api.patch(`/telegram/orders/${params.slug}/confirm`, {});

			// 2. On success, return a success object.
			//    SvelteKit actions automatically return a plain object.
			return {
				success: true,
				status: 'Success',
				message: 'Pembayaran Berhasil!'
			};
		} catch (error) {
			// 3. Handle Errors

			// If the error is a Response (e.g., from redirect or apiCall's 401 handler),
			// we must re-throw it so SvelteKit can handle the redirect or error response.
			if (error instanceof Response) {
				// If it's a redirect, let SvelteKit handle the navigation.
				if (error.status === 302 || error.status === 401 || error.status === 403) {
					throw error;
				}
			}

			// For all other errors (fetch error, 500 from API, etc.):
			const errorMessage =
				error instanceof Error ? error.message : 'Pembayaran Gagal. Silakan coba lagi.';

			// 🛑 Use fail() to return an error response to the form.
			//    The status is 500 (Server Error) or another appropriate code.
			return fail(500, {
				success: false,
				message: errorMessage,
				status: 'Error'
			});
		}
	}
};
