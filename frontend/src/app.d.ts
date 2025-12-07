// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: User | null; // The user object you set in hooks.server.ts
			token: string | null; // The access token you set in hooks.server.ts
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
