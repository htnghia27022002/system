import { env } from '@/config/env'
import { backendModuleNames } from '@/config/modules'

import { apiClient } from './api-client'

/** Public deployment info from GET /api/app/info (BE: public/handlers/app_handler.go). */
export type AppInfo = {
  /** Enabled BE modules (MODULES env on the BE). */
  modules: string[]
}

export async function fetchAppInfo(): Promise<AppInfo> {
  if (env.USE_MOCK_API) {
    return { modules: [...backendModuleNames] }
  }

  const { data } = await apiClient.get<AppInfo>('/app/info')
  return data
}
