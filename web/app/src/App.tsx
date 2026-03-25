import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { TeamSelector } from './pages/TeamSelector'
import { HealthCheckView } from './pages/HealthCheckView'

export default function App() {
  return (
    <BrowserRouter>
      <div style={{
        maxWidth: '800px',
        margin: '0 auto',
        padding: '24px 16px',
        fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
        color: '#111827',
      }}>
        <Routes>
          <Route path="/" element={<TeamSelector />} />
          <Route path="/healthcheck/:id" element={<HealthCheckView />} />
        </Routes>
      </div>
    </BrowserRouter>
  )
}
