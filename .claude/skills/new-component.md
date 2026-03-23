# Skill: /new-component

Scaffold a new Svelte 5 component in the UI following established conventions.

## Usage

```
/new-component <Name> <category> [--no-test]
```

- `<Name>` — PascalCase component name, e.g. `ProductCard`, `OrderSummary`, `PriceTag`.
- `<category>` — Subdirectory under `src/lib/components/`. Must be one of the established categories:
  `commerce` | `sections` | `forms` | `loading` | `overlays` | `navigation` | `user` | `feedback` | `primitives` | `display` | `data` | `animations` | `media`
- `--no-test` — Skip generating the test stub (use only for purely visual/presentational components).

**Must be run from inside `ui/`.**

---

## Before generating anything

1. List `src/lib/components/<category>/` — check if a component with a similar name or purpose already exists. If so, flag it and ask whether to extend it instead.
2. Read one or two existing components in the same category to understand the naming, prop interface, and Tailwind patterns already in use.
3. Read `ui/CLAUDE.md` for any styling or accessibility notes that apply.

---

## Component file

Create `src/lib/components/<category>/<Name>.svelte`.

**Svelte 5 rules — non-negotiable:**
- Props via `$props()`. Never `export let`.
- Reactive state via `$state()`. Never `let` with reactive assignments.
- Derived values via `$derived()`. Never `$:`.
- Event callbacks as typed props (e.g. `onclick?: () => void`). Never `createEventDispatcher`.
- Two-way binding via `$bindable()` only when the component is a genuine form control.

```svelte
<script lang="ts">
  interface Props {
    // TODO: define props
    class?: string;
  }

  let { class: className = '', ...props }: Props = $props();
</script>

<!-- TODO: implement component -->
<div class={className}>
  <slot />
</div>
```

**Tailwind v4 notes:**
- Use utility classes directly. No `@apply` in `<style>` blocks unless unavoidable.
- Responsive prefixes: `sm:`, `md:`, `lg:` (mobile-first).
- Animation: prefer `tw-animate-css` classes or GSAP for complex sequences. Do not reach for `transition-*` for entrance animations on first render.

**Accessibility — required on every component:**
- Interactive elements (`button`, `a`, custom click targets) must have visible focus styles and keyboard handlers.
- Images must have `alt`. Decorative images use `alt=""`.
- ARIA roles/labels where the element's semantic role is ambiguous.
- Color contrast must meet WCAG AA (4.5:1 text, 3:1 UI elements).
- Do not use `tabindex > 0`.

---

## Test stub

Unless `--no-test`, create `src/lib/components/<category>/__tests__/<Name>.test.ts`.

```typescript
import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import <Name> from '../<Name>.svelte';

describe('<Name>', () => {
  it('renders without crashing', () => {
    render(<Name>);
    // TODO: add meaningful assertions
  });

  // TODO: test prop variations, edge cases, interactions
});
```

---

## After generating

1. Print the file path created.
2. Remind the developer to:
   - Import and use the component: `import <Name> from '$lib/components/<category>/<Name>.svelte'`.
   - Fill in the prop interface and implementation before shipping.
   - Run `bun run test:unit` to confirm the stub passes.
   - Check contrast ratios if using custom colors.
