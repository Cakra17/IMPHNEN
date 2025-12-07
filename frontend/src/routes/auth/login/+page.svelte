<script lang="ts">
	import Separator from '$lib/components/primitives/separator.svelte';
	import {
		MailIcon,
		EyeIcon,
		EyeOffIcon,
		KeyRoundIcon,
		ArrowRightIcon,
		InfoIcon
	} from '@lucide/svelte';
	import { redirect } from '@sveltejs/kit';
	import { Alert, Button, Heading, Input, Label, P } from 'flowbite-svelte';
	import type { ActionData, PageData } from './$types';
	import { onMount } from 'svelte';
	import { afterNavigate, goto } from '$app/navigation';
	import { page } from '$app/stores';

	const { form, data } = $props<{ form: ActionData; data: PageData }>();

	let email = $state(form?.email ?? '');
	let password = $state('');
	let passwordVisible = $state(false);

	let hasSuccessParam = $derived($page.url.searchParams.has('registered'));
</script>

<div class="flex flex-col gap-3">
	<Heading tag="h2">Selamat Datang</Heading>
	<P class="text-teal-800">Masuk. Biarkan resi dan chat berbicara.</P>
</div>
<form method="POST">
	<div class="grid gap-6 grid-cols-1">
		<div>
			<Label for="email" class="mb-2">Email</Label>
			<Input
				class={`ps-10 ${email === '' ? 'text-teal-500' : ''}`}
				type="text"
				id="email"
				name="email"
				placeholder="nama@bisnisku.id"
				bind:value={email}
				required
			>
				{#snippet left()}
					<MailIcon />
				{/snippet}
			</Input>
		</div>
		<div>
			<Label for="password" class="mb-2">Password</Label>
			<Input
				class={`ps-10 ${password === '' ? 'text-teal-500' : ''}`}
				type={passwordVisible ? 'text' : 'password'}
				id="password"
				name="password"
				placeholder={passwordVisible ? 'supersecret' : '•••••••••••'}
				bind:value={password}
				required
			>
				{#snippet left()}
					<KeyRoundIcon />
				{/snippet}
				{#snippet right()}
					<button
						type="button"
						class="outline-none"
						onclick={() => (passwordVisible = !passwordVisible)}
					>
						{#if passwordVisible}
							<EyeOffIcon />
						{:else}
							<EyeIcon />
						{/if}
					</button>
				{/snippet}
			</Input>
		</div>
		{#if hasSuccessParam && !form?.error}
			<Alert border color="green">
				{#snippet icon()}<InfoIcon class="h-5 w-5" />{/snippet}
				Registrasi berhasil! Silahkan login ke akunmu.
			</Alert>
		{/if}
		{#if form?.error}
			<Alert border color="red">
				{#snippet icon()}<InfoIcon class="h-5 w-5" />{/snippet}
				{form.error}
			</Alert>
		{/if}
		<Button type="submit" disabled={email === '' || password.length < 8}
			><div class="flex flex-row gap-2">Masuk <ArrowRightIcon /></div></Button
		>
	</div>
</form>
<div class="w-full flex flex-row items-center gap-4">
	<Separator direction="Horizontal" className="flex-1 bg-slate-400" />
	<P class="text-teal-800">Belum punya akun?</P>
	<Separator direction="Horizontal" className="flex-1 bg-slate-400" />
</div>
<a href="/auth/register"
	><Heading tag="h6" class="w-full text-center text-teal-700 hover:text-teal-900 cursor-pointer"
		>Daftar Gratis</Heading
	></a
>
