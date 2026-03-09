# UX Design

## Design System

### Styling Stack
- **Tailwind CSS v4** — utility-first, configured via vite plugin (no config file required)
- **`tw-animate-css`** — Tailwind-compatible CSS animation utilities
- **`class-variance-authority` (CVA)** — typed component variant system
- **`tailwind-merge`** — safe Tailwind class merging (prevents conflicts)
- **`clsx`** — conditional class composition

### Icons
`lucide-svelte` — consistent icon set, tree-shaken per import.

### Typography
Defined as Tailwind utilities in `src/lib/styles/`. Use the established type scale — do not introduce ad-hoc font sizes.

### Color & Spacing
Tailwind design tokens. Do not add arbitrary values (`w-[143px]`) unless there is no token equivalent. Prefer the spacing scale.

---

## Animation

### GSAP
Used for page transitions, scroll-triggered reveals, and complex sequenced animations. GSAP timelines live in `src/lib/animations/`. Do not add GSAP animations unless asked — they have a significant bundle cost.

### Lenis (Smooth Scroll)
`lenis` wraps the root scroll container for smooth inertia scrolling. Initialized in `+layout.svelte`. Do not disable without understanding downstream effects on scroll-linked animations.

### Threlte (Three.js)
`@threlte/core` + `@threlte/extras` for 3D scenes. Used for hero/product showcase sections. Three.js scenes are heavyweight — lazy-load unless they appear above the fold.

### `postprocessing`
WebGL post-processing effects for Three.js renders. Used sparingly.

---

## Component Principles

1. **Extend before creating.** Check `src/lib/components/` exhaustively. If a variant of an existing component covers the use case, use CVA to add the variant.

2. **Accessibility is not optional.**
   - All interactive elements must have accessible names (label, `aria-label`, or visible text)
   - Keyboard navigation must work: focus visible, tab order logical
   - Color contrast must meet WCAG AA (4.5:1 for text, 3:1 for UI components)
   - Modals and drawers must trap focus and respond to Escape

3. **No gratuitous animation.** Animations must serve a purpose (guide attention, indicate state, provide continuity). Respect `prefers-reduced-motion`.

4. **No layout side-effects.** A component's style scope is its own box. Never use negative margins or absolute positioning that bleeds into parent layout without explicit intent.

---

## Component Variants (CVA Pattern)

New components that need style variants should use CVA:

```ts
import { cva, type VariantProps } from 'class-variance-authority';

const button = cva('base-classes', {
  variants: {
    intent: { primary: '...', secondary: '...' },
    size: { sm: '...', md: '...', lg: '...' },
  },
  defaultVariants: { intent: 'primary', size: 'md' },
});
```

---

## Responsive Design

Mobile-first. Breakpoints follow Tailwind defaults. All new pages and components must be tested at mobile (375px), tablet (768px), and desktop (1280px+) widths.
