/**
 * Hero 3D character config.
 *
 * Default `procedural` = live Three.js mesh (best while tuning look).
 * Set `mode: 'gltf'` to load `public/models/hero-character.glb`
 * (regenerate: `pnpm export:hero-character`).
 */
export const heroCharacterConfig = {
  /** Prefer procedural while iterating look; GLB mirrors the same builder. */
  mode: 'gltf' as 'procedural' | 'gltf',
  gltfUrl: '/models/hero-character.glb',
} as const
