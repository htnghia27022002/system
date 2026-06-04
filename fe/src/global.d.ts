// Allow CSS file imports (handled by Next.js webpack/turbopack)
declare module '*.css' {
  const styles: Record<string, string>
  export default styles
}
