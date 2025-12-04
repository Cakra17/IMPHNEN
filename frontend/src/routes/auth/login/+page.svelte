<script lang="ts">
	import Separator from "$lib/components/primitives/separator.svelte";
    import { MailIcon, EyeIcon, EyeOffIcon, KeyRoundIcon, ArrowRightIcon } from "@lucide/svelte";
	import { redirect } from "@sveltejs/kit";
    import { Button, Heading, Input, Label, P } from "flowbite-svelte";

    let { children } = $props();
    let email = $state('');
    let password = $state('');
	let passwordVisible = $state(false);

	function handleSubmit() {
		console.log("Mock login:", { email, password });

		// simple redirect
		window.location.href = "/dashboard";
	}
</script>

<div class="flex flex-col gap-3">
	<Heading tag="h2">Selamat Datang</Heading>
	<P class="text-teal-800">Masuk. Biarkan resi dan chat berbicara.</P>
</div>
<form onsubmit={handleSubmit}>
	<div class="grid gap-6 grid-cols-1">
		<div>
			<Label for="email" class="mb-2">Email</Label>
			<Input
				class={`ps-10 ${email === '' ? 'text-teal-500' : ''}`}
				type="text"
				id="email"
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
