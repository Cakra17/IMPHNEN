<script>
	import CashflowChart from "$lib/components/cashflow_chart.svelte";
	import { ChevronLeftIcon } from "@lucide/svelte";
	import { Heading, P } from "flowbite-svelte";

</script>

<a
	href="/dashboard"
	class="flex flex-row gap-2 font-bold text-teal-600 hover:text-teal-800 hover:underline cursor-pointer"
	><ChevronLeftIcon />Kembali</a
>
<div class="flex-1">
	<Heading class="text-lg md:text-2xl ml-8 mb-8">Analisa & Statistik</Heading>
	<CashflowChart />
</div>
