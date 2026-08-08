/**
 * Generates PWA PNG icons from public/favicon.svg.
 *
 * Run from fe/:  pnpm generate:pwa-icons
 */
import { mkdirSync, readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import sharp from 'sharp'

const __dirname = dirname(fileURLToPath(import.meta.url))
const root = resolve(__dirname, '..')
const svgPath = resolve(root, 'public/favicon.svg')
const outDir = resolve(root, 'public/icons')

const sizes = [
  { name: 'apple-touch-icon.png', size: 180 },
  { name: 'manifest-icon-192.maskable.png', size: 192 },
  { name: 'manifest-icon-512.maskable.png', size: 512 },
]

async function main() {
  mkdirSync(outDir, { recursive: true })
  const svg = readFileSync(svgPath)

  for (const { name, size } of sizes) {
    await sharp(svg)
      .resize(size, size, { fit: 'contain', background: '#3D72E0' })
      .png()
      .toFile(resolve(outDir, name))
    console.log(`Wrote public/icons/${name}`)
  }
}

main().catch((error) => {
  console.error(error)
  process.exit(1)
})
