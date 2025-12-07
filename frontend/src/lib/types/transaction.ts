export type Transaction = {
	id: string;
	type: 'income' | 'expense';
	source: 'receipt' | 'bot' | 'manual';
	amount: number;
	transaction_date: string;
	created_at: string;
};

export const transactionsSources = {
	receipt: 'Scan Struk',
	bot: 'Penjualan Telegram',
	manual: 'Input Manual'
};
