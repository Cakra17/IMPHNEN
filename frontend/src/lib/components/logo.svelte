<script lang="ts">
	import { BookXIcon } from "@lucide/svelte";
    type Size = "sm" | "lg"
    
    const { size }: { size: Size } = $props()
</script>

<div class={`shrink-0 relative bg-red-400 rounded-lg ${size == 'sm' ? 'h-10 w-10' : 'h-12 w-12'}`}>
	<div class="absolute inset-0 flex justify-center items-center">
		<BookXIcon size={size === 'sm' ? 24 : 36} />
	</div>
</div>
