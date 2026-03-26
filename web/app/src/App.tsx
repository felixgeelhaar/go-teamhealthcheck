import { useState } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { useParticipant } from './hooks/useParticipant'
import { Navbar } from './components/Navbar'
import { ParticipantSetup } from './components/ParticipantSetup'
import { TeamSelector } from './pages/TeamSelector'
import { HealthCheckView } from './pages/HealthCheckView'
import { VotingPage } from './pages/VotingPage'
import { CreateTemplate } from './pages/CreateTemplate'
import { CreateHealthCheck } from './pages/CreateHealthCheck'

export default function App() {
  const { isSet } = useParticipant()
  const [showSetup, setShowSetup] = useState(!isSet)

  const handleChangeNameClick = () => {
    setShowSetup(true)
  }

  const handleSetupComplete = () => {
    setShowSetup(false)
  }

  return (
    <BrowserRouter>
      <div className="app-shell">
        <Navbar onChangeNameClick={handleChangeNameClick} />

        <main className="app-content">
          <Routes>
            <Route path="/" element={<TeamSelector />} />
            <Route path="/healthcheck/:id" element={<HealthCheckView />} />
            <Route path="/healthcheck/:id/vote" element={<VotingPage />} />
            <Route path="/templates/new" element={<CreateTemplate />} />
            <Route path="/healthcheck/new/:teamId" element={<CreateHealthCheck />} />
          </Routes>
        </main>

        {showSetup && (
          <ParticipantSetup onComplete={handleSetupComplete} />
        )}
      </div>
    </BrowserRouter>
  )
}
