// Generates src/types/api.gen.ts from the BE OpenAPI document.
//   pnpm gen:api                                  # monorepo: ../be/api/openapi.json
//   OPENAPI_SOURCE=https://api.example.com/api/openapi.json pnpm gen:api
// The output is committed, so the FE builds without the BE present.
import { execFileSync } from 'node:child_process'
import { existsSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)))
const source = process.env.OPENAPI_SOURCE || path.join(root, '..', 'be', 'api', 'openapi.json')
if (!/^https?:\/\//.test(source) && !existsSync(source)) {
  console.error(`OpenAPI source not found: ${source}\nSet OPENAPI_SOURCE to the BE spec URL or file.`)
  process.exit(1)
}
const out = path.join(root, 'src', 'types', 'api.gen.ts')
execFileSync(
  path.join(root, 'node_modules', '.bin', 'openapi-typescript'),
  [source, '-o', out, '--root-types', '--root-types-no-schema-prefix'],
  { stdio: 'inherit' },
)
