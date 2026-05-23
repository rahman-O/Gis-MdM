import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { AuthProvider } from '@/app/providers'
import { App } from '@/app/App'
import '@/i18n/config'
import './index.css'

const densityKey = 'ui-density'
const savedDensity = window.localStorage.getItem(densityKey)
const densityClass =
  savedDensity === 'comfortable' ? 'density-comfortable' : savedDensity === 'dense' ? 'density-dense' : 'density-compact'
document.body.classList.remove('density-comfortable', 'density-compact', 'density-dense')
document.body.classList.add(densityClass)

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <AuthProvider>
        <App />
      </AuthProvider>
    </BrowserRouter>
  </StrictMode>
)
