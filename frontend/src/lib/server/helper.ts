import type { Cookies } from '@sveltejs/kit';
import { api } from './api';
import { formatDateForBackend } from '$lib/utils/date';

// API Helper to get transactions stats for range of day
export async function getStatsOnRange(cookies: Cookies, start_date: Date, end_date: Date) {
	return api.get('/transactions/stats', cookies, {
		start_date: formatDateForBackend(start_date),
		end_date: formatDateForBackend(end_date)
	});
}

// API Helper to get list of transactions for a day
export async function getTransactionsOnDay(cookies: Cookies, date: Date) {
	return api.get('/transactions/date', cookies, {
		date: formatDateForBackend(date)
	});
}
