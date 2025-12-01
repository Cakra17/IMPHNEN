<script lang="ts">
    type Direction = "Vertical" | "Horizontal"

    interface Props {
        direction: Direction;
        className?: string;
    }

    let { direction, className }: Props = $props();
</script>

<div
	class={`${direction === 'Horizontal' ? 'h-[1px] w-full' : 'w-[1px] h-full'} ${className}`}
></div>
