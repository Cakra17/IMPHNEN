// src/lib/utils/date.ts

/**
 * Format date to Indonesian readable format
 * @param date - Date string (YYYY-MM-DD) or Date object
 * @param options - Formatting options
 * @returns Formatted string like "6 Desember 2025"
 */
export function formatDate(
	date: string | Date,
	options?: {
		includeDay?: boolean; // Include day name (Senin, Selasa, etc.)
		short?: boolean; // Use short month names (Des instead of Desember)
		includeTime?: boolean; // Include time (HH:MM)
	}
): string {
	const { includeDay = false, short = false, includeTime = false } = options || {};

	// Parse date string to Date object
	const dateObj = typeof date === 'string' ? new Date(date) : date;

	// Check if valid date
	if (isNaN(dateObj.getTime())) {
		return 'Invalid date';
	}

	// Format options
	const formatOptions: Intl.DateTimeFormatOptions = {
		day: 'numeric',
		month: short ? 'short' : 'long',
		year: 'numeric'
	};

	if (includeDay) {
		formatOptions.weekday = 'long';
	}

	if (includeTime) {
		formatOptions.hour = '2-digit';
		formatOptions.minute = '2-digit';
	}

	return new Intl.DateTimeFormat('id-ID', formatOptions).format(dateObj);
}

/**
 * Format date to short format
 * @param date - Date string or Date object
 * @returns Formatted string like "6 Des 2025"
 */
export function formatDateShort(date: string | Date): string {
	return formatDate(date, { short: true });
}

/**
 * Format date with day name
 * @param date - Date string or Date object
 * @returns Formatted string like "Jumat, 6 Desember 2025"
 */
export function formatDateWithDay(date: string | Date): string {
	return formatDate(date, { includeDay: true });
}

/**
 * Get relative time in days only with limit
 * @param date - Date string or Date object
 * @param limit - Maximum days to show relative time (default: 7)
 * @returns Relative time string or formatted date if exceeds limit
 */
export function formatRelativeTime(date: string | Date, limit: number = 7): string {
	const dateObj = typeof date === 'string' ? new Date(date) : date;
	const now = new Date();

	// Reset time to midnight for accurate day comparison
	const dateOnly = new Date(dateObj.getFullYear(), dateObj.getMonth(), dateObj.getDate());
	const nowOnly = new Date(now.getFullYear(), now.getMonth(), now.getDate());

	const diffMs = nowOnly.getTime() - dateOnly.getTime();
	const diffDay = Math.floor(diffMs / (1000 * 60 * 60 * 24));

	// If difference exceeds limit, fallback to formatted date
	if (Math.abs(diffDay) > limit) {
		return formatDate(date);
	}

	// Relative time within limit
	if (diffDay === 0) {
		return 'Hari ini';
	} else if (diffDay === 1) {
		return 'Kemarin';
	} else if (diffDay > 1) {
		return `${diffDay} hari yang lalu`;
	} else if (diffDay === -1) {
		return 'Besok';
	} else if (diffDay < -1) {
		return `${Math.abs(diffDay)} hari lagi`;
	} else {
		return formatDate(date);
	}
}

/**
 * Format to DD/MM/YYYY
 * @param date - Date string or Date object
 * @returns Formatted string like "06/12/2025"
 */
export function formatDateSlash(date: string | Date): string {
	const dateObj = typeof date === 'string' ? new Date(date) : date;

	if (isNaN(dateObj.getTime())) {
		return 'Invalid date';
	}

	const day = String(dateObj.getDate()).padStart(2, '0');
	const month = String(dateObj.getMonth() + 1).padStart(2, '0');
	const year = dateObj.getFullYear();

	return `${day}/${month}/${year}`;
}

/**
 * Parse DD/MM/YYYY back to Date
 * @param dateStr - Date string in DD/MM/YYYY format
 * @returns Date object
 */
export function parseDateSlash(dateStr: string): Date {
	const [day, month, year] = dateStr.split('/');
	return new Date(parseInt(year), parseInt(month) - 1, parseInt(day));
}
