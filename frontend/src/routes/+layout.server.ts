import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals, url }) => {
	const { user } = locals;
	const { pathname } = url;

	const isAuthRoute = pathname.startsWith('/auth');
	const isDashboardRoute = pathname.startsWith('/dashboard');

	// --- User IS logged in ---
	if (user) {
		// Redirect away from /auth pages to /dashboard
		if (isAuthRoute) {
			console.log('User logged in. Redirecting from /auth.');
			throw redirect(302, '/dashboard');
		}

		// Allow access to all other pages (including /payment and /dashboard)
		return { user };
	}

	// --- User IS NOT logged in ---
	else {
		// Only protected route is /dashboard. Redirect to login if user tries to access it.
		if (isDashboardRoute) {
			console.log('User not logged in. Redirecting from /dashboard.');

			// Capture the intended path for redirect back after login
			const redirectTo = pathname === '/' ? '' : `?redirectTo=${pathname}`;

			throw redirect(302, `/auth/login${redirectTo}`);
		}

		// Allow access to all other pages, including /auth and /payment
		return {};
	}
};
