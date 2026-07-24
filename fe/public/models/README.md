# Hero 3D character models

- `hero-character.glb` — stylized avatar matched to site-owner reference photos
  (grey oversized tee, light jeans, cream sneakers). Regenerated with:

```bash
cd fe && pnpm exec tsx scripts/export-hero-character.mts
```

Config: `src/features/landing/components/hero-character/hero-character-config.ts`
(`mode: 'gltf'` loads this file; idle + head look-at use named nodes `Head` / `LeftArm` / `RightArm` / `Torso`).

Ready Player Me is discontinued (2026). For a photo-rigged replacement, export GLB from Mixamo / Avatar SDK / a 3D artist and overwrite `hero-character.glb`.
