import { memo, useState } from 'react'

interface HeaderProps {
  activeTab: 'transactions' | 'stats'
  setActiveTab: (tab: 'transactions' | 'stats') => void
  onAddClick: () => void
  onBudgetClick: () => void
  onCategoryManagerClick: () => void
  onLogout: () => void
  username?: string
}

export const Header = memo(({
  activeTab,
  setActiveTab,
  onAddClick,
  onBudgetClick,
  onCategoryManagerClick,
  onLogout,
  username = 'Пользователь',
}: HeaderProps) => {
  const [showProfileMenu, setShowProfileMenu] = useState(false)

  const handleLogoutClick = () => {
    onLogout()
    setShowProfileMenu(false)
  }

  return (
    <header className="header">
      <div className="header-content">
        <h1>Финансовый учёт</h1>
        
        {/* Большая toggle-кнопка под заголовком */}
        <div className="tab-toggle">
          <button
            className={`tab-toggle-btn ${activeTab === 'transactions' ? 'active' : ''}`}
            onClick={() => setActiveTab('transactions')}
          >
            Транзакции
          </button>
          <button
            className={`tab-toggle-btn ${activeTab === 'stats' ? 'active' : ''}`}
            onClick={() => setActiveTab('stats')}
          >
            Инфографика
          </button>
        </div>

        <div className="header-actions">
          <div className="nav-buttons">
            <button className="btn btn-primary" onClick={onAddClick}>
              + Добавить
            </button>
            <button className="btn btn-secondary" onClick={onBudgetClick}>
              Бюджет
            </button>
            <button className="btn btn-secondary" onClick={onCategoryManagerClick}>
              Категории
            </button>
            <div className="profile-wrapper">
              <button
                className="profile-btn"
                onClick={() => setShowProfileMenu(v => !v)}
              >
                <span className="profile-icon">👤</span>
                <span className="profile-name">{username}</span>
              </button>
              {showProfileMenu && (
                <div className="profile-dropdown">
                  <button className="btn btn-secondary" onClick={handleLogoutClick}>
                    Выйти
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </header>
  )
})