import 'axios'

// Allow CSS file imports (handled by Next.js webpack/turbopack)
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
