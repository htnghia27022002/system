import 'axios'

// Allow CSS side-effect imports under TypeScript 6 (noUncheckedSideEffectImports).
// Path-aliased CSS files also need a co-located `*.css.d.ts` (see src/styles/).
declare module '*.css' {
  const styles: Record<string, string>
  export default styles
}

declare module 'axios' {
  export interface AxiosRequestConfig {
    /** When true, request does not increment the top nav loading bar. */
    skipNavLoading?: boolean
  }
}
