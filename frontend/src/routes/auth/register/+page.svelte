<script lang="ts">
	import Separator from "$lib/components/primitives/separator.svelte";
    import { MailIcon, EyeIcon, KeyRoundIcon, ArrowRightIcon } from "@lucide/svelte";
    import { Button, Heading, Input, Label, P } from "flowbite-svelte";

    let { children } = $props();

    let firstName = $state('');
    let lastName = $state('');

    let email = $state('');
    let password = $state('');
    let passwordVisible = $state(false);

    let valid = $state(false);
    $effect(() => {
        if (firstName === '' || lastName === '' || email === '' || password === '') {
            valid = false;
        } else {
            valid = true;
        }
    })
</script>

<div class="flex flex-col gap-3">
	<Heading tag="h2">Buat Akun Baru</Heading>
	<P class="text-teal-800">Daftar sekarang. Santai pembukuan. Fokus bisnis.</P>
</div>
<form>
	<div class="grid gap-6 grid-cols-2">
		<div>
			<Label for="first_name" class="mb-2">Nama Depan</Label>
			<Input
				type="text"
				id="first_name"
				placeholder="John"
				class={firstName === '' ? 'text-teal-500' : ''}
				bind:value={firstName}
				required
			/>
		</div>
		<div>
			<Label for="last_name" class="mb-2">Nama Belakang</Label>
			<Input
				type="text"
				id="last_name"
				placeholder="Doe"
				class={lastName === '' ? 'text-teal-500' : ''}
				bind:value={lastName}
				required
			/>
		</div>
		<div class="col-span-2">
			<Label for="email" class="mb-2">Email</Label>
			<Input
				class={`ps-10 ${email === '' ? 'text-teal-500' : ''}`}
				type="text"
				id="email"
				placeholder="nama@bisnisumkmku.id"
				bind:value={email}
				required
			>
				{#snippet left()}
					<MailIcon />
				{/snippet}
			</Input>
		</div>
		<div class="col-span-2">
			<Label for="password" class="mb-2">Password</Label>
			<Input
				class={`ps-10 ${password === '' ? 'text-teal-500' : ''}`}
				type="password"
				id={passwordVisible ? 'text' : 'password'}
				placeholder="••••••••"
				bind:value={password}
				required
			>
				{#snippet left()}
					<KeyRoundIcon />
				{/snippet}
				{#snippet right()}
					<EyeIcon />
				{/snippet}
			</Input>
		</div>
		<Button type="submit" class="col-span-2 flex flex-row gap-3" disabled={!valid}
			>Daftar Sekarang<ArrowRightIcon /></Button
		>
	</div>
</form>
<div class="w-full flex flex-row items-center gap-4">
	<Separator direction="Horizontal" className="flex-1 bg-slate-400" />
	<P class="text-teal-800">Sudah punya akun?</P>
	<Separator direction="Horizontal" className="flex-1 bg-slate-400" />
</div>
<a href="/auth/login"
	><Heading tag="h6" class="w-full text-center text-teal-700 hover:text-teal-900 cursor-pointer"
		>Masuk ke Akun</Heading
	></a
>
