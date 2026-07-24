/**
 * Exports the photo-matched stylized hero to public/models/hero-character.glb
 *
 * Run from fe/:  node --import tsx scripts/export-hero-character.mts
 */
import { writeFileSync, mkdirSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { Blob as NodeBlob } from 'node:buffer'
import * as THREE from 'three'
import { GLTFExporter } from 'three/examples/jsm/exporters/GLTFExporter.js'

import { buildHeroCharacterGroup } from '../src/features/landing/components/hero-character/build-hero-character.ts'

// GLTFExporter expects browser Blob/FileReader APIs.
if (typeof globalThis.Blob === 'undefined') {
  globalThis.Blob = NodeBlob as unknown as typeof Blob
}

class FileReaderPolyfill {
  result: ArrayBuffer | null = null
  onloadend: ((ev: { target: FileReaderPolyfill }) => void) | null = null
  onerror: ((ev: unknown) => void) | null = null

  readAsArrayBuffer(blob: Blob) {
    Promise.resolve(blob.arrayBuffer())
      .then((buffer) => {
        this.result = buffer
        this.onloadend?.({ target: this })
      })
      .catch((err) => this.onerror?.(err))
  }
}

globalThis.FileReader = FileReaderPolyfill as unknown as typeof FileReader

const __dirname = dirname(fileURLToPath(import.meta.url))
const outPath = resolve(__dirname, '../public/models/hero-character.glb')

const scene = new THREE.Scene()
scene.add(buildHeroCharacterGroup())

const exporter = new GLTFExporter()

await new Promise<void>((resolvePromise, reject) => {
  exporter.parse(
    scene,
    (result) => {
      if (!(result instanceof ArrayBuffer)) {
        reject(new Error('Expected binary GLB ArrayBuffer from GLTFExporter'))
        return
      }
      mkdirSync(dirname(outPath), { recursive: true })
      writeFileSync(outPath, Buffer.from(result))
      console.log(`Wrote ${outPath} (${result.byteLength} bytes)`)
      resolvePromise()
    },
    (error) => reject(error),
    { binary: true },
  )
})
