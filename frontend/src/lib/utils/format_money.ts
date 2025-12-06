// src/lib/utils/currency.ts

/**
 * Format number to Indonesian Rupiah format
 * @param amount - The amount to format
 * @param includeDecimals - Whether to include decimal places (default: false)
 * @returns Formatted string like "Rp 1.234.567" or "Rp 1.234.567,89"
 */
export function formatRupiah(amount: number, includeDecimals: boolean = false): string {
	// Handle null/undefined
	if (amount == null) return 'Rp 0';

	// Format with Indonesian locale
	const formatted = new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		minimumFractionDigits: includeDecimals ? 2 : 0,
		maximumFractionDigits: includeDecimals ? 2 : 0
	}).format(amount);

	return formatted;
}

/**
 * Format number to compact Indonesian Rupiah format
 * @param amount - The amount to format
 * @returns Formatted string like "Rp 1,2 jt" or "Rp 1,5 M"
 */
export function formatRupiahCompact(amount: number): string {
	if (amount == null) return 'Rp 0';

	const abs = Math.abs(amount);
	const sign = amount < 0 ? '-' : '';

	if (abs >= 1_000_000_000) {
		return `${sign}Rp ${(abs / 1_000_000_000).toFixed(1)} M`;
	} else if (abs >= 1_000_000) {
		return `${sign}Rp ${(abs / 1_000_000).toFixed(1)} jt`;
	} else if (abs >= 1_000) {
		return `${sign}Rp ${(abs / 1_000).toFixed(1)} rb`;
	} else {
		return `${sign}Rp ${abs}`;
	}
}

/**
 * Parse Indonesian Rupiah string back to number
 * @param value - String like "Rp 1.234.567" or "1.234.567,89"
 * @returns Number value
 */
export function parseRupiah(value: string): number {
	// Remove Rp, spaces, and dots (thousand separators)
	const cleaned = value.replace(/Rp/g, '').replace(/\s/g, '').replace(/\./g, '').replace(/,/g, '.'); // Replace comma (decimal) with dot

	return parseFloat(cleaned) || 0;
}
