import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import '@/config/i18n'
import '@/styles/index.css'

import App from '@/App'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
