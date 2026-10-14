import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'

export default defineConfig([
  globalIgnores(['dist', '.next', 'node_modules', 'next-env.d.ts']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      tseslint.configs.recommended,
      reactHooks.configs.flat.recommended,
    ],
    languageOptions: {
      globals: globals.browser,
    },
    rules: {
      // React Compiler strict rule — setMounted/setState patterns in shadcn are intentional
      'react-hooks/set-state-in-effect': 'off',
      // Allow _-prefixed params to be unused (intentional no-op / backward-compat stubs)
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
      ],
    },
  },
  // Design-token guard (fe/DESIGN.md). shadcn primitives in components/ui and
  // marketing polish in features/landing are exempt.
  {
    files: ['src/**/*.{ts,tsx}'],
    ignores: ['src/components/ui/**', 'src/features/landing/**'],
    rules: {
      'no-restricted-syntax': ['error', ...designTokenRestrictions()],
    },
  },
])

function designTokenRestrictions() {
  const banned = [
    {
      pattern:
        '\\b(bg|text|border|ring|outline|fill|stroke|from|via|to|divide|decoration|accent|caret|shadow)-(slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose)-\\d{2,3}\\b|\\b(bg|text|border)-(black|white)\\b',
      message:
        'Use theme tokens (primary, muted, destructive, success, warning, chart-N, …) from src/styles/index.css instead of raw Tailwind palette colors.',
    },
    {
      pattern: '\\bz-\\[\\d+\\]',
      message:
        'Use the z-index scale from src/styles/index.css, e.g. z-(--z-modal) or z-(--z-floating).',
    },
    {
      pattern: '\\b(text|tracking|leading)-\\[[\\d.]+(px|rem|em)?\\]',
      message:
        'Use the Tailwind type scale (text-xs … text-2xl) instead of arbitrary font size, tracking or leading.',
    },
    {
      pattern: '\\brounded-(2xl|3xl)\\b',
      message:
        'Surfaces use rounded-xl, controls rounded-lg (see fe/DESIGN.md). Larger radii are for landing pages only.',
    },
  ]
  return banned.flatMap(({ pattern, message }) => [
    { selector: `Literal[value=/${pattern}/]`, message },
    { selector: `TemplateElement[value.raw=/${pattern}/]`, message },
  ])
}
