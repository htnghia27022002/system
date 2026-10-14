#!/usr/bin/env node
// Turns a fresh clone of the base into a named project, in place:
//   - display name ("System App")        → --name
//   - database name (myapp_pg / _test)   → <slug>_db / <slug>_test
//   - dev hostname (system.local)        → --domain
//   - prod image (youruser/go-app)       → --image (optional)
//   - logo wordmark                      → --name
// and creates .env, be/.env, fe/.env from the examples with a fresh JWT secret.
//
// Usage (repo root, once, right after cloning):
//   make bootstrap name="Acme Portal" slug=acme domain=acme.local [image=acme/acme-api]
//   node scripts/bootstrap.mjs --name "Acme Portal" --slug acme --domain acme.local --dry-run
//
// Feature history under docs/features/ is left untouched. Re-running is safe:
// values that were already replaced are simply not found again.
import { randomBytes } from 'node:crypto'
import { copyFileSync, existsSync, readFileSync, writeFileSync } from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)))
const args = process.argv.slice(2)
const opts = {}
for (let i = 0; i < args.length; i++) {
  if (args[i] === '--dry-run') opts.dryRun = true
  else if (args[i].startsWith('--')) opts[args[i].slice(2)] = args[++i]
}

const name = opts.name?.trim()
const slug = opts.slug?.trim()
const domain = opts.domain?.trim() || (slug ? `${slug}.local` : '')
if (!name || !slug) {
  console.error('Usage: node scripts/bootstrap.mjs --name "Acme Portal" --slug acme [--domain acme.local] [--image org/acme-api] [--dry-run]')
  process.exit(1)
}
if (!/^[a-z][a-z0-9-]{1,30}$/.test(slug)) {
  console.error('--slug: lowercase letters, digits and dashes, starting with a letter (e.g. acme or acme-portal)')
  process.exit(1)
}
if (!/^[a-z0-9.-]+$/.test(domain)) {
  console.error('--domain: a hostname such as acme.local')
  process.exit(1)
}
const db = slug.replace(/-/g, '_')
const shortName = name.split(/\s+/)[0]
const jsonString = (s) => JSON.stringify(s).slice(1, -1) // escape for JSON/TS string bodies

/** [file, [[from, to], ...]] — literal replacements, applied to every occurrence. */
const plan = [
  // Display name
  ['be/.env.example', [['System App', name], ['myapp_pg', `${db}_db`], ['system.local', domain]]],
  ['be/config.yaml', [['myapp_pg', `${db}_db`], ['system.local', domain]]],
  ['be/internal/config/config.go', [['"System App"', JSON.stringify(name)], ['"myapp_pg"', `"${db}_db"`]]],
  ['be/test/unit/database/migrate_test.go', [['myapp_pg', `${db}_db`]]],
  ['be/public/routes/app.go', [['"System API"', JSON.stringify(`${name} API`)]]],
  ['be/api/openapi.json', [['"System API"', JSON.stringify(`${name} API`)]]],
  ['fe/.env.example', [['System App', name]]],
  ['fe/Dockerfile', [['System App', name]]],
  ['fe/src/config/env.ts', [["'System App'", `'${name.replace(/'/g, "\\'")}'`]]],
  ['fe/src/app/layout.tsx', [["'System App'", `'${name.replace(/'/g, "\\'")}'`], ["'%s | System App'", `'%s | ${name.replace(/'/g, "\\'")}'`], ["title: 'System',", `title: '${shortName.replace(/'/g, "\\'")}',`]]],
  ['fe/src/components/common/pwa-install-prompt.tsx', [['Install System App', `Install ${name}`]]],
  ['fe/src/components/common/admin-app-logo.tsx', [['alt="htnghia"', `alt=${JSON.stringify(name)}`]]],
  ['fe/src/assets/favicon/logo-lockup.svg', [['>htnghia<', `>${shortName.replace(/[<&>]/g, '')}<`]]],
  ['fe/src/features/landing/content.ts', [["brand: 'System',", `brand: '${shortName.replace(/'/g, "\\'")}',`]]],
  ['fe/public/manifest.webmanifest', [['"System App"', JSON.stringify(name)]]],
  ['fe/public/offline.html', [['System App', name]]],
  ['fe/src/locales/en/common.json', [['"appName": "System App"', `"appName": "${jsonString(name)}"`]]],
  ['fe/src/locales/vi/common.json', [['"appName": "System App"', `"appName": "${jsonString(name)}"`]]],
  // Infra
  ['.env.example', [['system.local', domain]]],
  ['docker-compose.yml', [['myapp_pg', `${db}_db`]]],
  ['Makefile', [['myapp_pg', `${db}_db`]]],
  ['docker/nginx/nginx.dev.conf', [['system.local', domain]]],
  ['docker/nginx/nginx.prod.conf', [['system.local', domain]]],
  ['.github/workflows/ci.yml', [['myapp_test', `${db}_test`]]],
  ['package.json', [['"name": "system"', `"name": "${slug}"`]]],
  // Docs that describe how to run this project
  ['README.md', [['# System — core base', `# ${name}`], ['system.local', domain]]],
  ['AGENTS.md', [['system.local', domain]]],
  ['docker/README.md', [['system.local', domain]]],
  ['.claude/rules/environment.md', [['system.local', domain]]],
  ['.cursor/rules/environment.mdc', [['system.local', domain]]],
]
if (opts.image) {
  for (const file of ['.env.example', 'docker-compose.yml', 'Makefile', 'docker/README.md']) {
    plan.push([file, [['youruser/go-app', opts.image]]])
  }
}

// Merge entries per file so each file is read and reported once.
const byFile = new Map()
for (const [rel, pairs] of plan) byFile.set(rel, [...(byFile.get(rel) ?? []), ...pairs])

let changed = 0
for (const [rel, pairs] of byFile) {
  const file = path.join(root, rel)
  if (!existsSync(file)) continue
  let text = readFileSync(file, 'utf8')
  const before = text
  for (const [from, to] of pairs) text = text.split(from).join(to)
  if (text !== before) {
    changed++
    console.log(`  ~ ${rel}`)
    if (!opts.dryRun) writeFileSync(file, text)
  }
}

// Env files: create from the examples (never overwrite), with a random JWT secret.
for (const [example, target] of [['.env.example', '.env'], ['be/.env.example', 'be/.env'], ['fe/.env.example', 'fe/.env']]) {
  const dest = path.join(root, target)
  if (existsSync(dest)) {
    console.log(`  = ${target} exists, left as is`)
    continue
  }
  console.log(`  + ${target}`)
  if (opts.dryRun) continue
  copyFileSync(path.join(root, example), dest)
  if (target === 'be/.env') {
    const text = readFileSync(dest, 'utf8').replace(/^JWT_SECRET=.*$/m, `JWT_SECRET=${randomBytes(32).toString('hex')}`)
    writeFileSync(dest, text)
  }
}

console.log(`
${opts.dryRun ? '[dry run] ' : ''}${changed} files updated for "${name}" (slug ${slug}, db ${db}_db, host ${domain}).
Next:
  1. echo "127.0.0.1 ${domain}" | sudo tee -a /etc/hosts
  2. Replace the logo: fe/src/assets/favicon/*, fe/public/icons/* (pnpm -C fe generate:pwa-icons)
  3. make up-d   →   http://${domain}:\${NGINX_HTTP_PORT:-80}
  4. Read docs/new-project.md for the rest of the checklist.`)
