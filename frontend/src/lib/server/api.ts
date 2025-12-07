// Centralized API utility for making authenticated requests

import { redirect } from '@sveltejs/kit';
import type { Cookies } from '@sveltejs/kit';

export const API_BASE_URL = 'http://202.155.95.111/api/v1';

interface ApiOptions {
	method?: string;
	body?: any;
	cookies: Cookies;
	params?: Record<string, string | number | boolean>;
}

export async function apiCall(endpoint: string, options: ApiOptions) {
	const { method = 'GET', body, cookies, params } = options;

	const accessToken = cookies.get('access_token');

	if (!accessToken) {
		throw redirect(302, '/auth/login');
	}

	// Build URL
	let url = `${API_BASE_URL}${endpoint}`;

	if (params) {
		const query = new URLSearchParams(Object.entries(params).map(([k, v]) => [k, String(v)]));
		if (query.toString()) url += `?${query}`;
	}

	/**
	 * ⚠️ IMPORTANT:
	 * Do NOT send Content-Type header for GET request
	 * SvelteKit + JSON server = errors if you send Content-Type without body
	 */
	const headers: Record<string, string> = {
		Authorization: `Bearer ${accessToken}`
	};

	if (method !== 'GET' && !(body instanceof FormData)) {
		headers['Content-Type'] = 'application/json';
	}

	let fetchOptions: RequestInit = {
		method,
		headers
	};

	if (body) {
		fetchOptions.body = body instanceof FormData ? body : JSON.stringify(body);
	}

	let response: Response;

	try {
		response = await fetch(url, fetchOptions);
	} catch (err) {
		if (err instanceof Error && err.name === 'AbortError') {
			throw new Error('Request timeout - server tidak merespons');
		}
		throw err;
	}

	// Handle expired auth
	if (response.status === 401) {
		cookies.delete('access_token', { path: '/' });
		cookies.delete('user_data', { path: '/' });

		throw redirect(302, '/auth/login?session=expired');
	}

	// Handle other error codes
	if (!response.ok) {
		let message = `API Error: ${response.status} ${response.statusText}`;
		try {
			const data = await response.json();
			if (data?.message) message = data.message;
		} catch {
			// Response wasn't JSON
		}
		throw new Error(message);
	}

	// Handle non-JSON (204, empty body, etc)
	if (response.status === 204) return null;

	try {
		return await response.json();
	} catch {
		return null;
	}
}

// Convenience methods
export const api = {
	get: (endpoint: string, cookies: Cookies, params?: Record<string, string | number | boolean>) =>
		apiCall(endpoint, { method: 'GET', cookies, params }),

	post: (
		endpoint: string,
		body: any,
		cookies: Cookies,
		params?: Record<string, string | number | boolean>
	) => apiCall(endpoint, { method: 'POST', body, cookies, params }),

	put: (
		endpoint: string,
		body: any,
		cookies: Cookies,
		params?: Record<string, string | number | boolean>
	) => apiCall(endpoint, { method: 'PUT', body, cookies, params }),

	delete: (
		endpoint: string,
		cookies: Cookies,
		params?: Record<string, string | number | boolean>
	) => apiCall(endpoint, { method: 'DELETE', cookies, params })
};
