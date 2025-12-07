<script lang="ts">
	import type { ApexOptions } from 'apexcharts';
	import { Chart } from '@flowbite-svelte-plugins/chart';
	import {
		Card,
		A,
		Button,
		Dropdown,
		DropdownItem,
		Popover,
		Heading,
		Datepicker,
		P
	} from 'flowbite-svelte';
	import {
		InfoCircleSolid,
		ChevronRightOutline,
		ChevronDownOutline,
		FileLinesSolid
	} from 'flowbite-svelte-icons';

	let dateRange: { from: Date | undefined; to: Date | undefined } = $state({
		from: new Date(),
		to: new Date()
	});
	let availableRange = $derived.by(() => {
		if (!dateRange.from) {
			return { from: undefined, to: undefined };
		}

		const today = new Date();
		today.setHours(0, 0, 0, 0); // normalize

		// --- FROM ---
		const from = new Date(dateRange.from);
		from.setDate(from.getDate() - 30);

		// --- TO ---
		const to = new Date(dateRange.from);
		to.setDate(to.getDate() + 30);

		// clamp "to" so it's never after today
		if (to > today) {
			to.setTime(today.getTime());
		}

		return { from, to };
	});

	let options: ApexOptions = {
		chart: {
			zoom: {
				enabled: false,
				allowMouseWheelZoom: false
			},
			height: '400px',
			type: 'line',
			fontFamily: 'Inter, sans-serif',
			dropShadow: {
				enabled: false
			},
			toolbar: {
				show: false
			}
		},
		tooltip: {
			enabled: true,
			x: {
				show: false
			}
		},
		dataLabels: {
			enabled: false
		},
		stroke: {
			width: 6,
			curve: 'straight'
		},
		grid: {
			show: true,
			strokeDashArray: 4,
			padding: {
				left: 2,
				right: 2,
				top: -26
			}
		},
		series: [
			{
				name: 'Pemasukan',
				data: [
					{
						x: new Date(2025, 11, 1),
						y: 50000
					},
					{
						x: new Date(2025, 11, 2),
						y: 40000
					},
					{
						x: new Date(2025, 11, 3),
						y: 80000
					},
					{
						x: new Date(2025, 11, 4),
						y: 100000
					},
					{
						x: new Date(2025, 11, 5),
						y: 80000
					},
					{
						x: new Date(2025, 11, 6),
						y: 150000
					},
					{
						x: new Date(2025, 11, 7),
						y: 40000
					}
				],
				color: '#4CAF50'
			},
			{
				name: 'Pengeluaran',
				data: [
					{
						x: new Date(2025, 11, 1),
						y: 10000
					},
					{
						x: new Date(2025, 11, 2),
						y: 10000
					},
					{
						x: new Date(2025, 11, 3),
						y: 16000
					},
					{
						x: new Date(2025, 11, 4),
						y: 20000
					},
					{
						x: new Date(2025, 11, 5),
						y: 10000
					},
					{
						x: new Date(2025, 11, 6),
						y: 50000
					},
					{
						x: new Date(2025, 11, 7),
						y: 15000
					}
				],
				color: '#de1a24'
			}
		],
		legend: {
			show: false
		},
		xaxis: {
			type: 'datetime',
			labels: {
				show: true,
				style: {
					fontFamily: 'Inter, sans-serif',
					cssClass: 'text-xs font-normal fill-gray-500 dark:fill-gray-400'
				}
			},
			axisBorder: {
				show: false
			},
			axisTicks: {
				show: false
			}
		},
		yaxis: {
			min: 0,
			show: true
		}
	};
</script>

<div class="flex flex-col bg-white border border-teal-200 rounded-2xl p-6">
	<div class="flex-1 mb-5 flex flex-col md:flex-row justify-between gap-4">
		<div class="flex flex-col justify-center md:justify-start">
			<Heading class="text-xl mb-1 text-center md:text-left">Arus Kas</Heading>
			<span class="text-sm mb-2 text-center md:text-left">Pemasukan vs. Pengeluaran</span>
		</div>
		<div class="flex-1 hidden md:block"></div>
		<div class="flex-1 flex flex-col justify-center md:justify-end">
			<Heading class="text-md mb-1 text-right hidden md:block">Rentang tanggal</Heading>
			<Datepicker
				range
				bind:rangeFrom={dateRange.from}
				bind:rangeTo={dateRange.to}
				availableFrom={availableRange.from}
				availableTo={availableRange.to}
				color="teal"
			/>
		</div>
	</div>
	<Chart {options} />
	<div class="flex flex-row gap-6 justify-center">
		<div class="flex flex-row items-center gap-2">
			<div class="w-3 h-3 rounded-full" style="background-color: #4CAF50;"></div>
			<span>Pemasukan</span>
		</div>
		<div class="flex flex-row items-center gap-2">
			<div class="w-3 h-3 rounded-full" style="background-color: #de1a24;"></div>
			<span>Pengeluaran</span>
		</div>
	</div>
</div>
