<script lang="ts">
  import { onMount, onDestroy } from "svelte";

  let now = new Date();
  let timer: ReturnType<typeof setInterval> | null = null;

  onMount(() => {
    timer = setInterval(() => {
      now = new Date();
    }, 2000); // Potato friendly
  });

  onDestroy(() => {
    if (timer) clearInterval(timer);
  });

  $: time = now.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
</script>

<span>{time}</span>
