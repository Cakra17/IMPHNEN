import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const accessToken = event.cookies.get('access_token');
	const userDataCookie = event.cookies.get('user_data');

	// Simply trust the cookies exist = user is logged in
	// Any invalid/expired tokens will be caught when actual API calls return 401
	if (accessToken && userDataCookie) {
		try {
			event.locals.user = JSON.parse(userDataCookie);
		} catch (error) {
			console.error('Failed to parse user_data cookie:', error);
			event.cookies.delete('access_token', { path: '/' });
			event.cookies.delete('user_data', { path: '/' });
		}
	}

	return resolve(event);
};
