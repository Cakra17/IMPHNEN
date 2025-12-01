<script lang="ts">
	import Separator from "$lib/components/primitives/separator.svelte";
    import { MailIcon, EyeIcon, KeyRoundIcon, ArrowRightIcon } from "@lucide/svelte";
    import { Button, Heading, Input, Label, P } from "flowbite-svelte";

    let { children } = $props();
</script>

<div class="h-screen grid grid-cols-2 font-jakarta">
	<div class="h-full flex flex-col py-16 px-16 md:px-24 bg-pattern flex-1">
		<Heading tag="h4" class="text-white">IMPHNEN</Heading>
		<P italic class="text-white">Ingin Menjadi Pengusaha Handal Namun Enggan Ngebuku</P>
		<div class="flex-1 flex flex-col justify-center">
			<Heading tag="h1" class="text-white">Kelola Bisnis,</Heading>
			<Heading tag="h1" class="text-emerald-400">Tanpa Mager.</Heading>
			<P class="text-white py-8"
				><strong>IMPHNEN</strong> membantu memudahkan pencatatan nota dan orderan pelanggan kamu. Lihat
				untung-rugi kamu tanpa perlu nyatet dan hitung manual. Kamu jadi punya waktu lebih banyak untuk
				kembangkan bisnismu.</P
			>
		</div>
		<P size="sm" class="text-teal-300"
			>Made by Cakra & Theerlan for IMPHNEN x Kolosal.Ai Hackathon</P
		>
	</div>
	<div class="h-full flex flex-col justify-center px-16 md:px-24 gap-8">
		{@render children()}
	</div>
</div>
