import { error } from 'console';
import type { Actions } from './$types';
import { fail, redirect } from '@sveltejs/kit';
import { API_BASE_URL } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const actions: Actions = {
	default: async ({ request, cookies, fetch }) => {
		const data = await request.formData();
		const email = data.get('email');
		const password = data.get('password');

		if (!email || !password) {
			return fail(400, {
				email: email,
				error: 'Email dan password wajib diisi.'
			});
		}

		const apiResponse = await fetch(`${API_BASE_URL}/auth/login`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ email, password })
		});

		if (apiResponse.ok) {
			const result = await apiResponse.json();
			const { access_token } = result.data.token;
			const userData = result.data.user;

			cookies.set('access_token', access_token, {
				path: '/',
				httpOnly: true, // Prevent client-side access
				secure: process.env.NODE_ENV === 'production',
				sameSite: 'strict'
			});

			cookies.set('user_data', JSON.stringify(userData), {
				path: '/',
				httpOnly: false, // Can be read by the client
				secure: process.env.NODE_ENV === 'production',
				sameSite: 'strict'
			});

			throw redirect(303, '/dashboard');
		} else if (apiResponse.status === 400) {
			return fail(400, { error: 'Kata sandi atau email salah. Silahkan coba lagi.' });
		} else if (apiResponse.status === 404) {
			return fail(404, { error: 'Akun tidak ditemukan.' });
		} else if (apiResponse.status === 500) {
			console.error('API Error: Failed to generate token on server.');
			return fail(500, {
				error: 'Terjadi masalah pada server. Silakan coba sebentar lagi.'
			});
		} else {
			return fail(apiResponse.status, {
				error: 'Login gagal. Ada masalah yang tidak terduga.'
			});
		}
	}
};
