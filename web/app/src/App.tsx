import { useState } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { useParticipant } from './hooks/useParticipant'
import { useApi } from './hooks/useApi'
import { Navbar } from './components/Navbar'
import { ParticipantSetup } from './components/ParticipantSetup'
import { TeamSelector } from './pages/TeamSelector'
import { HealthCheckView } from './pages/HealthCheckView'
import { VotingPage } from './pages/VotingPage'
import { CreateTemplate } from './pages/CreateTemplate'
import { CreateHealthCheck } from './pages/CreateHealthCheck'
import { CompareTeams } from './pages/CompareTeams'
import { RetroPage } from './pages/RetroPage'
import type { PluginEntry } from './types'

export default function App() {
  const { isSet } = useParticipant()
  const [showSetup, setShowSetup] = useState(!isSet)
  const { data: plugins } = useApi<PluginEntry[]>('/api/plugins')

  const handleChangeNameClick = () => {
    setShowSetup(true)
  }

  const handleSetupComplete = () => {
    setShowSetup(false)
  }

  return (
    <BrowserRouter>
      <div className="app-shell">
        <Navbar onChangeNameClick={handleChangeNameClick} plugins={plugins || []} />

        <main className="app-content">
          <Routes>
            <Route path="/" element={<TeamSelector />} />
            <Route path="/healthcheck/:id" element={<HealthCheckView />} />
            <Route path="/healthcheck/:id/vote" element={<VotingPage />} />
            <Route path="/templates/new" element={<CreateTemplate />} />
            <Route path="/compare" element={<CompareTeams />} />
            <Route path="/healthcheck/new/:teamId" element={<CreateHealthCheck />} />
            <Route path="/retro/:hcId" element={<RetroPage />} />
            {(plugins || []).map(plugin => {
              if (plugin.route === '/retro/:hcId') return null
              return (
                <Route
                  key={plugin.name}
                  path={plugin.route}
                  element={<div>Plugin: {plugin.label}</div>}
                />
              )
            })}
          </Routes>
        </main>

        {showSetup && (
          <ParticipantSetup onComplete={handleSetupComplete} />
        )}
      </div>
    </BrowserRouter>
  )
}
