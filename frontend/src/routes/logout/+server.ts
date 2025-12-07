import { redirect } from '@sveltejs/kit';

export const GET = ({ cookies }) => {
	// delete access_token cookie
	cookies.delete('access_token', {
		path: '/',
		httpOnly: true,
		secure: process.env.NODE_ENV === 'production',
		sameSite: 'strict'
	});
	cookies.delete('user_data', {
		path: '/',
		httpOnly: false,
		secure: process.env.NODE_ENV === 'production',
		sameSite: 'strict'
	});

	// redirect to login (or home)
	throw redirect(302, '/auth/login');
};
