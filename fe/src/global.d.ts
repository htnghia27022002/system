import 'axios'

declare module 'axios' {
  export interface AxiosRequestConfig {
    /** When true, request does not increment the top nav loading bar. */
    skipNavLoading?: boolean
  }
}
