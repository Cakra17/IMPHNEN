import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals, url }) => {
	const { user } = locals;
	const isAuthRoute = url.pathname.startsWith('/auth');

	// User IS logged in
	if (user) {
		// Redirect away from auth pages to dashboard
		if (isAuthRoute) {
			throw redirect(302, '/dashboard');
		}

		// Redirect root to dashboard
		if (url.pathname === '/') {
			throw redirect(302, '/dashboard');
		}

		// Return user data for authenticated routes
		return { user };
	}

	// User IS NOT logged in
	else {
		// Allow access to auth routes (login/register)
		if (isAuthRoute) {
			return {};
		}

		// Redirect protected routes to login
		const redirectTo = url.pathname === '/' ? '' : `?redirectTo=${url.pathname}`;
		throw redirect(302, `/auth/login${redirectTo}`);
	}
};
