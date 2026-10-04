// Allow CSS side-effect imports under TypeScript 6 (noUncheckedSideEffectImports),
// including package CSS such as `leaflet/dist/leaflet.css`. Must stay a non-module
// declaration file (no top-level import/export) so the wildcard applies globally.
declare module '*.css' {
  const styles: Record<string, string>
  export default styles
}
