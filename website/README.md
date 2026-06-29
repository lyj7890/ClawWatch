# ClawWatch — Website

The official landing page for **ClawWatch**, real-time AI agent monitoring for
engineering teams.

## Stack

- **Vite** + **React 19** + **TypeScript**
- **Tailwind CSS v4** (via `@tailwindcss/vite` — no PostCSS config)
- **motion/react** for entrance + scroll-reveal animation
- **@phosphor-icons/react** for iconography
- **Geist Sans / Geist Mono**, self-hosted via `@fontsource`

## Design

Dark-tech minimalist. `zinc-950` base, single `emerald-400` brand accent,
left-aligned hero with a live terminal stream, bento feature grid, terminal-style
architecture + code blocks. All motion respects `prefers-reduced-motion`.

## Develop

```bash
npm install
npm run dev      # http://localhost:5173
```

## Build

```bash
npm run build    # type-check + production bundle → dist/
npm run preview  # serve the build locally
```

## Structure

```
src/
  App.tsx                  # layout shell + ambient glow/grain
  main.tsx                 # entry, font imports
  index.css                # Tailwind v4 theme tokens + utilities
  components/
    Nav.tsx
    Hero.tsx               # live terminal feed
    Features.tsx           # bento grid
    Architecture.tsx       # Agent → Hub → Console
    HowItWorks.tsx         # 3-step flow
    CodeExample.tsx        # SDK snippet, copy-to-clipboard
    Footer.tsx
    Logo.tsx
    Reveal.tsx             # scroll-reveal wrapper
```

## License

MIT
