<script lang="ts">
	interface NavItem {
		label: string;
		href: string;
	}

	const navItems: NavItem[] = [
		{ label: 'Products', href: '/products' },
		{ label: 'Services', href: '/services' }
	];
</script>

<header class="header-glass">
	<div class="container mx-auto px-6">
		<div class="flex items-center justify-between h-20">
			<!-- Left Nav -->
			<nav class="flex items-center gap-8">
				{#each navItems.slice(0, 1) as item}
					<a
						href={item.href}
						class="nav-link font-header text-sm tracking-widest uppercase"
					>
						{item.label}
					</a>
				{/each}
			</nav>

			<!-- Center Logo -->
			<div class="shrink-0">
				<a href="/" class="logo-link">
					<img src="/komodo-logo.png" alt="Komodo" class="h-12" />
				</a>
			</div>

			<!-- Right Nav -->
			<nav class="flex items-center gap-8">
				{#each navItems.slice(1) as item}
					<a
						href={item.href}
						class="nav-link font-header text-sm tracking-widest uppercase"
					>
						{item.label}
					</a>
				{/each}
			</nav>

			<!-- Action Icons -->
			<div class="flex items-center gap-6 ml-8">
				<button class="icon-button" aria-label="Search">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
				</button>
				<button class="icon-button" aria-label="Account">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
				</button>
				<button class="icon-button relative" aria-label="Cart">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="8" cy="21" r="1"/><circle cx="19" cy="21" r="1"/><path d="M2.05 2.05h2l2.66 12.42a2 2 0 0 0 2 1.58h9.78a2 2 0 0 0 1.95-1.57l1.65-7.43H5.12"/></svg>
					<span class="cart-badge">0</span>
				</button>
			</div>
		</div>
	</div>
</header>

<style>
	.header-glass {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 50;
		background: hsl(var(--komodo-cyber-blue) / 0.3);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-bottom: 1px solid hsl(var(--komodo-electric-cyan) / 0.2);
		box-shadow: 0 4px 30px hsl(var(--komodo-cyber-blue) / 0.1);
	}

	.nav-link {
		color: hsl(var(--komodo-ghost-white));
		text-decoration: none;
		padding: 0.5rem 0;
	}

	.nav-link:hover {
		color: hsl(var(--komodo-electric-cyan));
	}

	.logo-link {
		text-decoration: none;
		display: inline-block;
	}

	.icon-button {
		background: transparent;
		border: 1px solid hsl(var(--komodo-electric-cyan) / 0.3);
		color: hsl(var(--komodo-ghost-white));
		padding: 0.5rem;
		border-radius: 0.375rem;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.icon-button:hover {
		background: hsl(var(--komodo-electric-cyan) / 0.1);
		border-color: hsl(var(--komodo-electric-cyan));
		color: hsl(var(--komodo-electric-cyan));
	}

	.cart-badge {
		position: absolute;
		top: -4px;
		right: -4px;
		background: hsl(var(--komodo-fuchsia));
		color: hsl(var(--komodo-ghost-white));
		font-size: 0.625rem;
		font-weight: bold;
		padding: 0.125rem 0.375rem;
		border-radius: 9999px;
		min-width: 1.25rem;
		text-align: center;
	}

	@media (max-width: 768px) {
		.header-glass {
			padding: 0 1rem;
		}

		.nav-link,
		.icon-button {
			font-size: 0.75rem;
		}

		.logo-link img {
			height: 2.5rem;
		}
	}
</style>
