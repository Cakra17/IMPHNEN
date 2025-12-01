<script lang="ts">
	import Logo from "$lib/components/logo.svelte";
	import Separator from "$lib/components/primitives/separator.svelte";
    import { MailIcon, EyeIcon, KeyRoundIcon, ArrowRightIcon, BookXIcon } from "@lucide/svelte";
    import { Button, Heading, Input, Label, P } from "flowbite-svelte";

    let { children } = $props();
</script>

<div class="h-screen grid grid-cols-1 lg:grid-cols-2 font-jakarta">
	<div class="h-full flex flex-col py-16 px-8 md:px-12 lg:px-16 xl:px-24 gap-2 bg-pattern flex-1">
		<div class="flex flex-row gap-2 items-start">
			<Logo size="lg" />
			<div>
				<Heading tag="h4" class="text-white">IMPHNEN</Heading>
				<P italic class="text-white">Ingin Menjadi Pengusaha Handal Namun Enggan Ngebuku</P>
			</div>
		</div>
		<div class="flex-1 flex flex-col py-16 justify-center">
			<Heading tag="h1" class="text-white">Kelola Bisnis,</Heading>
			<Heading tag="h1" class="text-emerald-400">Tanpa Mager.</Heading>
		</div>
		<P class="text-white"
			><strong>IMPHNEN</strong> membantu memudahkan pencatatan nota dan orderan pelanggan kamu. Lihat
			untung-rugi kamu tanpa perlu nyatet dan hitung manual. Kamu jadi punya waktu lebih banyak untuk
			kembangkan bisnismu.</P
		>
	</div>
	<div class="h-full flex flex-col justify-center px-16 md:px-24 pt-16 pb-8 md:pt-8">
		<div class="flex-1 flex flex-col gap-8 justify-center w-full">{@render children()}</div>
		<P size="sm" class="text-teal-500 text-center pt-16"
			>Made by Cakra & Theerlan for IMPHNEN x Kolosal.Ai Hackathon</P
		>
	</div>
</div>
