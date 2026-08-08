# PWA icons

PNG icons are generated from `public/favicon.svg`.

```bash
pnpm generate:pwa-icons
```

Outputs:

- `apple-touch-icon.png` (180×180) — iOS Home Screen
- `manifest-icon-192.maskable.png` — Android install
- `manifest-icon-512.maskable.png` — Android splash / install

Re-run after changing the logo SVG.
