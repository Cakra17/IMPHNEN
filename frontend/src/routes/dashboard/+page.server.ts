import { api } from '$lib/server/api';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	try {
		const summaryData = await api.get('/transactions/stats/days', cookies, {
			days: 1
		});
		const txTodayData = await api.get('/transactions/days', cookies, {
			days: 1
		});
		console.log(txTodayData);
		console.log(summaryData);
		return {
			summary: summaryData.data,
			transactions: txTodayData.data.transactions
		};
	} catch (error) {
		console.error('Failed to load transaction stats:', error);
		console.error('Error details:', {
			message: error instanceof Error ? error.message : 'Unknown',
			stack: error instanceof Error ? error.stack : undefined
		});

		// Return empty data instead of throwing so page still loads
		return {
			summary: null,
			error: error instanceof Error ? error.message : 'Failed to load data'
		};
	}
};

export const actions: Actions = {
	createManualTransaction: async ({ request, cookies }) => {
		// --- 1. Parse form submission ---
		const form = await request.formData();

		const body = {
			amount: Number(form.get('amount')),
			source: 'manual',
			transaction_date: String(form.get('transaction_date')),
			type: String(form.get('type'))
		};

		// --- 2. Call API using your helper ---
		// NOTE: No "default", this is a named action.
		const result = await api.post('/transactions', body, cookies);

		// --- 3. Return something to page ---
		return {
			success: true,
			result
		};
	}
};
