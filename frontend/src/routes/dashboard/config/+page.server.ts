import { api } from '$lib/server/api';
import type { Actions, PageServerLoad } from './$types';
import { fail, redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies, url }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const perPage = Number(url.searchParams.get('per_page')) || 10;

	try {
		const userData = await api.get('/users/me', cookies);
		const products = await api.get('/products', cookies, {
			page,
			per_page: perPage
		});

		return {
			user: userData.data.user,
			products_meta: products.meta,
			products: products.data.products
		};
	} catch (error) {
		console.error('Failed to load transaction stats:', error);
		console.error('Error details:', {
			message: error instanceof Error ? error.message : 'Unknown',
			stack: error instanceof Error ? error.stack : undefined
		});

		// Return empty data instead of throwing so page still loads
		return {
			user: null,
			error: error instanceof Error ? error.message : 'Failed to load data'
		};
	}
};

export const actions = {
	updateProfile: async ({ request, cookies }) => {
		const formData = await request.formData();

		const firstname = formData.get('firstname') as string;
		const lastname = formData.get('lastname') as string;
		const storeName = formData.get('store_name') as string;

		// Validation
		if (!firstname || firstname.trim() === '') {
			return fail(400, { error: 'First name is required' });
		}
		if (!lastname || lastname.trim() === '') {
			return fail(400, { error: 'Last name is required' });
		}
		if (!storeName || storeName.trim() === '') {
			return fail(400, { error: 'Store name is required' });
		}

		try {
			const result = await api.put(
				'/users/me',
				{
					firstname: firstname.trim(),
					lastname: lastname.trim(),
					store_name: storeName.trim()
				},
				cookies
			);

			// Update the user_data cookie with new info
			const currentUserData = cookies.get('user_data');
			if (currentUserData) {
				try {
					const userData = JSON.parse(currentUserData);
					userData.firstname = firstname.trim();
					userData.lastname = lastname.trim();
					userData.store_name = storeName.trim();

					cookies.set('user_data', JSON.stringify(userData), {
						path: '/',
						httpOnly: true,
						secure: true,
						sameSite: 'strict',
						maxAge: 60 * 60 * 24 * 7 // 7 days
					});
				} catch (e) {
					console.error('Failed to update user_data cookie:', e);
				}
			}

			return {
				success: true,
				message: 'Profile updated successfully',
				data: result
			};
		} catch (error) {
			console.error('Profile update error:', error);
			return fail(500, {
				error: error instanceof Error ? error.message : 'Failed to update profile'
			});
		}
	},
	submitProduct: async ({ request, cookies }) => {
		const formData = await request.formData();

		// Get form fields
		const image = formData.get('image') as File;
		const name = formData.get('name') as string;
		const price = formData.get('price') as string;
		const stock = formData.get('stock') as string;

		// Validation
		if (!image || image.size === 0) {
			return fail(400, { error: 'Product image is required' });
		}
		if (!name || name.trim() === '') {
			return fail(400, { error: 'Product name is required' });
		}
		if (!price || isNaN(Number(price))) {
			return fail(400, { error: 'Valid price is required' });
		}
		if (!stock || isNaN(Number(stock))) {
			return fail(400, { error: 'Valid stock is required' });
		}

		// Create FormData for API (since it accepts file upload)
		const apiFormData = new FormData();
		apiFormData.append('image', image);
		apiFormData.append('name', name);
		apiFormData.append('price', price);
		apiFormData.append('stock', stock);

		try {
			const result = await api.post('/products', apiFormData, cookies);

			return {
				success: true,
				data: result
			};
		} catch (error) {
			console.error('Product submission error:', error);
			return fail(500, {
				error: error instanceof Error ? error.message : 'Failed to create product'
			});
		}
	},
	updateProduct: async ({ request, cookies }) => {
		const formData = await request.formData();

		// Get form fields
		const id = formData.get('id') as string;
		const image = formData.get('image') as File;
		const name = formData.get('name') as string;
		const price = formData.get('price') as string;
		const stock = formData.get('stock') as string;

		const apiFormData = new FormData();
		// Validation & Append
		if (!stock || isNaN(Number(stock))) {
			return fail(400, { error: 'Valid stock is required' });
		}
		apiFormData.append('stock', stock);
		if (image && image.size >= 0) {
			apiFormData.append('image', image);
		}
		if (name && name.trim() === '') {
			apiFormData.append('name', name);
		}
		if (price && !isNaN(Number(price))) {
			apiFormData.append('price', price);
		}

		// Create FormData for API (since it accepts file upload)

		try {
			const result = await api.put(`/products/${id}`, apiFormData, cookies);

			return {
				success: true,
				data: result
			};
		} catch (error) {
			console.error('Product submission error:', error);
			return fail(500, {
				error: error instanceof Error ? error.message : 'Failed to create product'
			});
		}
	},
	deleteProfile: async ({ cookies }) => {
		try {
			await api.delete('/users/me', cookies);

			// Clear auth cookies
			cookies.delete('access_token', { path: '/' });
			cookies.delete('user_data', { path: '/' });

			// Redirect to login with message
			throw redirect(302, '/auth/login?deleted=true');
		} catch (error) {
			// If the error is a redirect, let it through
			if (error instanceof Response && error.status === 302) {
				throw error;
			}

			return fail(500, {
				error: error instanceof Error ? error.message : 'Failed to delete account'
			});
		}
	}
} satisfies Actions;
