import { getStatsOnRange, getTransactionsOnDay } from '$lib/server/helper.js';
import { parse } from 'path';

export const load = async ({ cookies }) => {
	const today = new Date();
	const from = new Date();
	from.setDate(today.getDate() - 7);

	const cashflowStats = await getStatsOnRange(cookies, from, today);

	const dailyTransactions = await getTransactionsOnDay(cookies, today);

	return {
		stats: cashflowStats.data,
		dailyTransactions: dailyTransactions.data.transactions
	};
};

export const actions = {
	getDaily: async ({ request, cookies }) => {
		const data = await request.formData();
		const dateString = data.get('date') as string;
		const date = new Date(dateString);

		const result = await getTransactionsOnDay(cookies, date);

		// If result is a string, parse it first:
		const parsedResult = typeof result === 'string' ? JSON.parse(result) : result;
		return {
			dailyTransactions: parsedResult.data.transactions // Return the actual data object, not the string
		};
	}
};
