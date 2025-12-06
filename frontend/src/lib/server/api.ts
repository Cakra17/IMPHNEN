// src/lib/server/api.ts
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
		// No token = redirect to login
		throw redirect(302, '/auth/login');
	}

	// Build URL with query parameters
	let url = `${API_BASE_URL}${endpoint}`;
	if (params) {
		const queryString = new URLSearchParams(
			Object.entries(params).map(([key, value]) => [key, String(value)])
		).toString();
		url += `?${queryString}`;
	}

	// Create abort controller for timeout
	const controller = new AbortController();
	const timeoutId = setTimeout(() => controller.abort(), 10000); // 10 second timeout

	try {
		const response = await fetch(url, {
			method,
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${accessToken}`
			},
			body: body ? JSON.stringify(body) : undefined,
			signal: controller.signal
		});

		clearTimeout(timeoutId);

		// Handle 401 Unauthorized - token expired or invalid
		if (response.status === 401) {
			// Clear cookies and redirect to login
			cookies.delete('access_token', { path: '/' });
			cookies.delete('user_data', { path: '/' });
			throw redirect(302, '/auth/login?session=expired');
		}

		// Handle other error status codes
		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
			console.error('API Error Details:', {
				status: response.status,
				statusText: response.statusText,
				url: url,
				errorData
			});
			throw new Error(errorData.message || `API Error: ${response.status} ${response.statusText}`);
		}

		return response.json();
	} catch (error) {
		clearTimeout(timeoutId);

		if (error instanceof Error && error.name === 'AbortError') {
			console.error('API Request Timeout:', url);
			throw new Error('Request timeout - API server tidak merespons');
		}

		throw error;
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
