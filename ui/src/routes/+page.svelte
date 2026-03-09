<script lang="ts">
	import Hero from '$lib/components/layout/hero.svelte';

	interface HeroItem {
		id: string;
		title: string;
		image: string;
		type: 'product' | 'service';
	}

	const heroItems: HeroItem[] = [
		{
			id: 'symmetrix',
			title: 'SYMMETRIX',
			image: 'https://images.unsplash.com/photo-1511499767150-a48a237f0083?w=800&q=80',
			type: 'product'
		},
		{
			id: 'matrix',
			title: 'MATRIX',
			image: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=800&q=80',
			type: 'product'
		},
		{
			id: 'maverix',
			title: 'MAVERIX',
			image: 'https://images.unsplash.com/photo-1509967419530-da38b4704bc6?w=800&q=80',
			type: 'service'
		},
		{
			id: 'semtex',
			title: 'SEMTEX',
			image: 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800&q=80',
			type: 'service'
		},
		{
			id: 'assentrix',
			title: 'ASSENTRIX',
			image: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&q=80',
			type: 'product'
		}
	];
</script>

<Hero items={heroItems} />
