import { useState, useEffect } from 'react'
import AccountManager from './components/AccountManager'
import ApiTester from './components/ApiTester'
import BatchImport from './components/BatchImport'
import VercelSync from './components/VercelSync'

const TABS = [
    { id: 'accounts', label: '🔑 账号管理' },
    { id: 'test', label: '🧪 API 测试' },
    { id: 'import', label: '📦 批量导入' },
    { id: 'vercel', label: '☁️ Vercel 同步' },
]

export default function App() {
    const [activeTab, setActiveTab] = useState('accounts')
    const [config, setConfig] = useState({ keys: [], accounts: [] })
    const [loading, setLoading] = useState(true)
    const [message, setMessage] = useState(null)

    const fetchConfig = async () => {
        try {
            setLoading(true)
            const res = await fetch('/admin/config')
            if (res.ok) {
                const data = await res.json()
                setConfig(data)
            }
        } catch (e) {
            console.error('获取配置失败:', e)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchConfig()
    }, [])

    const showMessage = (type, text) => {
        setMessage({ type, text })
        setTimeout(() => setMessage(null), 5000)
    }

    const renderTab = () => {
        switch (activeTab) {
            case 'accounts':
                return <AccountManager config={config} onRefresh={fetchConfig} onMessage={showMessage} />
            case 'test':
                return <ApiTester config={config} onMessage={showMessage} />
            case 'import':
                return <BatchImport onRefresh={fetchConfig} onMessage={showMessage} />
            case 'vercel':
                return <VercelSync onMessage={showMessage} />
            default:
                return null
        }
    }

    return (
        <div className="app">
            <header className="header">
                <h1>DS2API Admin</h1>
                <p>账号管理 · API 测试 · Vercel 部署</p>
            </header>

            {message && (
                <div className={`alert alert-${message.type}`}>
                    {message.text}
                </div>
            )}

            <div className="stats">
                <div className="stat">
                    <div className="stat-value">{config.keys?.length || 0}</div>
                    <div className="stat-label">API Keys</div>
                </div>
                <div className="stat">
                    <div className="stat-value">{config.accounts?.length || 0}</div>
                    <div className="stat-label">账号</div>
                </div>
            </div>

            <div className="tabs">
                {TABS.map(tab => (
                    <button
                        key={tab.id}
                        className={`tab ${activeTab === tab.id ? 'active' : ''}`}
                        onClick={() => setActiveTab(tab.id)}
                    >
                        {tab.label}
                    </button>
                ))}
            </div>

            {loading ? (
                <div className="card">
                    <div className="empty-state">
                        <span className="loading"></span> 加载中...
                    </div>
                </div>
            ) : (
                renderTab()
            )}
        </div>
    )
}
