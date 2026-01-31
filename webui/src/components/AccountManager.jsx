import { useState } from 'react'

export default function AccountManager({ config, onRefresh, onMessage }) {
    const [showAddKey, setShowAddKey] = useState(false)
    const [showAddAccount, setShowAddAccount] = useState(false)
    const [newKey, setNewKey] = useState('')
    const [newAccount, setNewAccount] = useState({ email: '', mobile: '', password: '' })
    const [loading, setLoading] = useState(false)

    const addKey = async () => {
        if (!newKey.trim()) return
        setLoading(true)
        try {
            const res = await fetch('/admin/keys', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ key: newKey.trim() }),
            })
            if (res.ok) {
                onMessage('success', 'API Key 添加成功')
                setNewKey('')
                setShowAddKey(false)
                onRefresh()
            } else {
                const data = await res.json()
                onMessage('error', data.detail || '添加失败')
            }
        } catch (e) {
            onMessage('error', '网络错误')
        } finally {
            setLoading(false)
        }
    }

    const deleteKey = async (key) => {
        if (!confirm('确定删除此 API Key？')) return
        try {
            const res = await fetch(`/admin/keys/${encodeURIComponent(key)}`, { method: 'DELETE' })
            if (res.ok) {
                onMessage('success', '删除成功')
                onRefresh()
            } else {
                onMessage('error', '删除失败')
            }
        } catch (e) {
            onMessage('error', '网络错误')
        }
    }

    const addAccount = async () => {
        if (!newAccount.password || (!newAccount.email && !newAccount.mobile)) {
            onMessage('error', '请填写密码和邮箱/手机号')
            return
        }
        setLoading(true)
        try {
            const res = await fetch('/admin/accounts', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(newAccount),
            })
            if (res.ok) {
                onMessage('success', '账号添加成功')
                setNewAccount({ email: '', mobile: '', password: '' })
                setShowAddAccount(false)
                onRefresh()
            } else {
                const data = await res.json()
                onMessage('error', data.detail || '添加失败')
            }
        } catch (e) {
            onMessage('error', '网络错误')
        } finally {
            setLoading(false)
        }
    }

    const deleteAccount = async (id) => {
        if (!confirm('确定删除此账号？')) return
        try {
            const res = await fetch(`/admin/accounts/${encodeURIComponent(id)}`, { method: 'DELETE' })
            if (res.ok) {
                onMessage('success', '删除成功')
                onRefresh()
            } else {
                onMessage('error', '删除失败')
            }
        } catch (e) {
            onMessage('error', '网络错误')
        }
    }

    return (
        <div className="section">
            {/* API Keys */}
            <div className="card">
                <div className="card-header">
                    <span className="card-title">🔑 API Keys</span>
                    <button className="btn btn-primary" onClick={() => setShowAddKey(true)}>+ 添加</button>
                </div>

                {config.keys?.length > 0 ? (
                    <div className="list">
                        {config.keys.map((key, i) => (
                            <div key={i} className="list-item">
                                <span className="list-item-text">{key.slice(0, 16)}****</span>
                                <button className="btn btn-danger" onClick={() => deleteKey(key)}>删除</button>
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="empty-state">暂无 API Key</div>
                )}
            </div>

            {/* Accounts */}
            <div className="card">
                <div className="card-header">
                    <span className="card-title">👤 DeepSeek 账号</span>
                    <button className="btn btn-primary" onClick={() => setShowAddAccount(true)}>+ 添加</button>
                </div>

                {config.accounts?.length > 0 ? (
                    <div className="list">
                        {config.accounts.map((acc, i) => (
                            <div key={i} className="list-item">
                                <div className="list-item-info">
                                    <span className="list-item-text">{acc.email || acc.mobile}</span>
                                    <span className={`badge ${acc.has_token ? 'badge-success' : 'badge-warning'}`}>
                                        {acc.has_token ? '已登录' : '未登录'}
                                    </span>
                                </div>
                                <button className="btn btn-danger" onClick={() => deleteAccount(acc.email || acc.mobile)}>删除</button>
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="empty-state">暂无账号</div>
                )}
            </div>

            {/* Add Key Modal */}
            {showAddKey && (
                <div className="modal-overlay" onClick={() => setShowAddKey(false)}>
                    <div className="modal" onClick={e => e.stopPropagation()}>
                        <div className="modal-header">
                            <span className="modal-title">添加 API Key</span>
                            <button className="modal-close" onClick={() => setShowAddKey(false)}>&times;</button>
                        </div>
                        <div className="form-group">
                            <label className="form-label">API Key</label>
                            <input
                                type="text"
                                className="form-input"
                                placeholder="输入你自定义的 API Key"
                                value={newKey}
                                onChange={e => setNewKey(e.target.value)}
                            />
                        </div>
                        <div className="btn-group">
                            <button className="btn btn-secondary" onClick={() => setShowAddKey(false)}>取消</button>
                            <button className="btn btn-primary" onClick={addKey} disabled={loading}>
                                {loading ? <span className="loading"></span> : '添加'}
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Add Account Modal */}
            {showAddAccount && (
                <div className="modal-overlay" onClick={() => setShowAddAccount(false)}>
                    <div className="modal" onClick={e => e.stopPropagation()}>
                        <div className="modal-header">
                            <span className="modal-title">添加 DeepSeek 账号</span>
                            <button className="modal-close" onClick={() => setShowAddAccount(false)}>&times;</button>
                        </div>
                        <div className="form-group">
                            <label className="form-label">Email（可选）</label>
                            <input
                                type="email"
                                className="form-input"
                                placeholder="user@example.com"
                                value={newAccount.email}
                                onChange={e => setNewAccount({ ...newAccount, email: e.target.value })}
                            />
                        </div>
                        <div className="form-group">
                            <label className="form-label">手机号（可选）</label>
                            <input
                                type="text"
                                className="form-input"
                                placeholder="+86..."
                                value={newAccount.mobile}
                                onChange={e => setNewAccount({ ...newAccount, mobile: e.target.value })}
                            />
                        </div>
                        <div className="form-group">
                            <label className="form-label">密码（必填）</label>
                            <input
                                type="password"
                                className="form-input"
                                placeholder="DeepSeek 账号密码"
                                value={newAccount.password}
                                onChange={e => setNewAccount({ ...newAccount, password: e.target.value })}
                            />
                        </div>
                        <div className="btn-group">
                            <button className="btn btn-secondary" onClick={() => setShowAddAccount(false)}>取消</button>
                            <button className="btn btn-primary" onClick={addAccount} disabled={loading}>
                                {loading ? <span className="loading"></span> : '添加'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}
