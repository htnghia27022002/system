import { env } from '@/config/env'
import { backendModuleNames } from '@/config/modules'
import { loadSystemMock } from '@/services/mock'

import { apiClient } from './api-client'

/** Public deployment info from GET /api/app/info (BE: public/handlers/app_handler.go). */
export type AppInfo = {
  /** Enabled BE modules (MODULES env on the BE). */
  modules: string[]
  /** Public runtime settings (key → value), e.g. `auth.registrationEnabled`. */
  settings: Record<string, string | boolean | number>
}

export async function fetchAppInfo(): Promise<AppInfo> {
  if (env.USE_MOCK_API) {
    return {
      modules: [...backendModuleNames],
      settings: (await loadSystemMock()).publicSettings(),
    }
  }

  const { data } = await apiClient.get<AppInfo>('/app/info')
  return { modules: data.modules ?? [], settings: data.settings ?? {} }
}
