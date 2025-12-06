import { API_BASE_URL } from '$lib/server/api';
import type { Actions } from './$types';
import { fail, isRedirect, redirect } from '@sveltejs/kit';

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();
		const email = data.get('email')?.toString() || '';
		const firstname = data.get('firstname')?.toString() || '';
		const lastname = data.get('lastname')?.toString() || '';
		const password = data.get('password')?.toString() || '';
		const store_name = data.get('store_name')?.toString() || '';

		// Validate required fields
		if (!email || !firstname || !lastname || !password || !store_name) {
			return fail(400, {
				email,
				firstname,
				lastname,
				store_name,
				error: 'Semua field wajib diisi.'
			});
		}

		// Validate email format
		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		if (!emailRegex.test(email)) {
			return fail(400, {
				email,
				firstname,
				lastname,
				store_name,
				error: 'Format email tidak valid.'
			});
		}

		// Validate password length
		if (password.length < 8) {
			return fail(400, {
				email,
				firstname,
				lastname,
				store_name,
				error: 'Password minimal 8 karakter.'
			});
		}

		try {
			// Call registration API
			const apiResponse = await fetch(`${API_BASE_URL}/auth/register`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email, firstname, lastname, password, store_name })
			});

			if (apiResponse.ok) {
				const result = await apiResponse.json();
				console.log('Registration successful:', result.message);

				// Redirect to login page after successful registration
				redirect(303, '/auth/login?registered=true');
			} else if (apiResponse.status === 400) {
				return fail(400, {
					email,
					firstname,
					lastname,
					store_name,
					error: 'Data registrasi tidak valid. Periksa kembali input Anda.'
				});
			} else if (apiResponse.status === 409) {
				return fail(409, {
					email,
					firstname,
					lastname,
					store_name,
					error: 'Email sudah terdaftar. Silakan gunakan email lain atau login.'
				});
			} else if (apiResponse.status === 500) {
				console.error('API Error: Server error during registration.');
				return fail(500, {
					email,
					firstname,
					lastname,
					store_name,
					error: 'Terjadi masalah pada server. Silakan coba sebentar lagi.'
				});
			} else {
				return fail(apiResponse.status, {
					email,
					firstname,
					lastname,
					store_name,
					error: 'Registrasi gagal. Ada masalah yang tidak terduga.'
				});
			}
		} catch (error) {
			if (isRedirect(error)) {
				// If it is, re-throw it so SvelteKit can handle the redirect properly.
				throw error;
			}

			// Network error
			console.error('Registration error:', error);
			return fail(500, {
				email,
				firstname,
				lastname,
				store_name,
				error: 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda.'
			});
		}
	}
};
