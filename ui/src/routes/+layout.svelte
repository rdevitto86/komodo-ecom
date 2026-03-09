<script lang="ts">
	import '$lib/styles/index.css';
	import Header from '$lib/components/navigation/Header.svelte';

	interface Props {
		children: import('svelte').Snippet;
	}

	let { children }: Props = $props();

	logger.info('Layout loaded');
</script>

<div class="dark">
	<Header />
	{@render children()}
</div>
