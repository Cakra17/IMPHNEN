<script lang="ts">
  import type { ApexOptions } from "apexcharts";
  import { Chart } from "@flowbite-svelte-plugins/chart";
  import { Card, A, Button, Dropdown, DropdownItem, Popover, Heading } from "flowbite-svelte";
  import { InfoCircleSolid, ChevronRightOutline, ChevronDownOutline, FileLinesSolid } from "flowbite-svelte-icons";

  let options: ApexOptions = {
    chart: {
      zoom: {
        enabled: false,
        allowMouseWheelZoom: false
      },
      height: "400px",
      type: "line",
      fontFamily: "Inter, sans-serif",
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
      curve: "straight"
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
        name: "Pemasukan",
        data: [{
            x: new Date(2025,11,1), y: 50000
        },{
            x: new Date(2025,11,2), y: 40000
        },{
            x: new Date(2025,11,3), y: 80000
        },{
            x: new Date(2025,11,4), y: 100000
        },{
            x: new Date(2025,11,5), y: 80000
        },{
            x: new Date(2025,11,6), y: 150000
        },{
            x: new Date(2025,11,7), y: 40000
        }],
        color: "#4CAF50"
      },
      {
        name: "Pengeluaran",
        data: [{
            x: new Date(2025,11,1), y: 10000
        },{
            x: new Date(2025,11,2), y: 10000
        },{
            x: new Date(2025,11,3), y: 16000
        },{
            x: new Date(2025,11,4), y: 20000
        },{
            x: new Date(2025,11,5), y: 10000
        },{
            x: new Date(2025,11,6), y: 50000
        },{
            x: new Date(2025,11,7), y: 15000
        }],
        color: "#de1a24"
      },
    ],
    legend: {
      show: false
    },
    xaxis: {
        type: "datetime",
      labels: {
        show: true,
        style: {
          fontFamily: "Inter, sans-serif",
          cssClass: "text-xs font-normal fill-gray-500 dark:fill-gray-400"
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
	<div class="mb-5 flex justify-between">
		<div class="flex flex-col">
			<Heading class="text-xl mb-1">Arus Kas</Heading>
			<span class="text-sm mb-2">Pemasukan vs. Pengeluaran</span>
		</div>
		<div>
			<Button color="light" class="px-3 py-2"
				>Last week<ChevronDownOutline class="ms-1.5 h-2.5 w-2.5" /></Button
			>
			<Dropdown simple class="w-40">
				<DropdownItem>Yesterday</DropdownItem>
				<DropdownItem>Today</DropdownItem>
				<DropdownItem>Last 7 days</DropdownItem>
				<DropdownItem>Last 30 days</DropdownItem>
				<DropdownItem>Last 90 days</DropdownItem>
			</Dropdown>
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
