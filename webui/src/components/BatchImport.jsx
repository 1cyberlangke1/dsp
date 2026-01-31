import { useState } from 'react'
import { FileCode, Download, Upload, Copy, Check, AlertTriangle } from 'lucide-react'
import clsx from 'clsx'

const TEMPLATES = {
    full: {
        name: 'Full Configuration',
        desc: 'Includes keys, accounts, and model mappings',
        config: {
            keys: ["your-api-key-1", "your-api-key-2"],
            accounts: [
                { email: "user1@example.com", password: "password1", token: "" },
                { email: "user2@example.com", password: "password2", token: "" },
                { mobile: "+8613800138001", password: "password3", token: "" }
            ],
            claude_model_mapping: {
                fast: "deepseek-chat",
                slow: "deepseek-reasoner"
            }
        }
    },
    email_only: {
        name: 'Email Only',
        desc: 'Batch import email accounts',
        config: {
            keys: ["your-api-key"],
            accounts: [
                { email: "account1@example.com", password: "pass1", token: "" },
                { email: "account2@example.com", password: "pass2", token: "" },
                { email: "account3@example.com", password: "pass3", token: "" }
            ]
        }
    },
    mobile_only: {
        name: 'Mobile Only',
        desc: 'Batch import mobile number accounts',
        config: {
            keys: ["your-api-key"],
            accounts: [
                { mobile: "+8613800000001", password: "pass1", token: "" },
                { mobile: "+8613800000002", password: "pass2", token: "" },
                { mobile: "+8613800000003", password: "pass3", token: "" }
            ]
        }
    },
    keys_only: {
        name: 'API Keys Only',
        desc: 'Just adding API access keys',
        config: {
            keys: ["key-1", "key-2", "key-3"]
        }
    }
}

export default function BatchImport({ onRefresh, onMessage, authFetch }) {
    const [jsonInput, setJsonInput] = useState('')
    const [loading, setLoading] = useState(false)
    const [result, setResult] = useState(null)
    const [copied, setCopied] = useState(false)

    const apiFetch = authFetch || fetch

    const handleImport = async () => {
        if (!jsonInput.trim()) {
            onMessage('error', 'Please enter JSON configuration')
            return
        }

        let config
        try {
            config = JSON.parse(jsonInput)
        } catch (e) {
            onMessage('error', 'Invalid JSON format')
            return
        }

        setLoading(true)
        setResult(null)
        try {
            const res = await apiFetch('/admin/import', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(config),
            })
            const data = await res.json()
            if (res.ok) {
                setResult(data)
                onMessage('success', `Imported: ${data.imported_keys} Keys, ${data.imported_accounts} Accounts`)
                onRefresh()
            } else {
                onMessage('error', data.detail || 'Import failed')
            }
        } catch (e) {
            onMessage('error', 'Network error')
        } finally {
            setLoading(false)
        }
    }

    const loadTemplate = (key) => {
        const tpl = TEMPLATES[key]
        if (tpl) {
            setJsonInput(JSON.stringify(tpl.config, null, 2))
            onMessage('info', `Loaded template: ${tpl.name}`)
        }
    }

    const handleExport = async () => {
        try {
            const res = await apiFetch('/admin/export')
            if (res.ok) {
                const data = await res.json()
                setJsonInput(JSON.stringify(JSON.parse(data.json), null, 2))
                onMessage('success', 'Configuration loaded')
            }
        } catch (e) {
            onMessage('error', 'Failed to fetch config')
        }
    }

    const copyBase64 = async () => {
        try {
            const res = await apiFetch('/admin/export')
            if (res.ok) {
                const data = await res.json()
                await navigator.clipboard.writeText(data.base64)
                setCopied(true)
                setTimeout(() => setCopied(false), 2000)
                onMessage('success', 'Base64 copied to clipboard')
            }
        } catch (e) {
            onMessage('error', 'Copy failed')
        }
    }

    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 h-[calc(100vh-140px)]">
            {/* Templates Panel */}
            <div className="md:col-span-1 space-y-4">
                <div className="bg-card border border-border rounded-xl p-5 shadow-sm">
                    <h3 className="font-semibold flex items-center gap-2 mb-4">
                        <FileCode className="w-4 h-4 text-primary" />
                        Quick Templates
                    </h3>
                    <div className="space-y-3">
                        {Object.entries(TEMPLATES).map(([key, tpl]) => (
                            <button
                                key={key}
                                onClick={() => loadTemplate(key)}
                                className="w-full text-left p-3 rounded-lg border border-border bg-secondary/20 hover:bg-secondary/50 hover:border-primary/50 transition-all custom-focus group"
                            >
                                <div className="font-medium text-sm group-hover:text-primary transition-colors">{tpl.name}</div>
                                <div className="text-xs text-muted-foreground mt-0.5">{tpl.desc}</div>
                            </button>
                        ))}
                    </div>
                </div>

                <div className="bg-linear-to-br from-primary/10 to-transparent border border-primary/20 rounded-xl p-5 shadow-sm">
                    <h3 className="font-semibold flex items-center gap-2 mb-2 text-primary">
                        <Download className="w-4 h-4" />
                        Export Data
                    </h3>
                    <p className="text-sm text-muted-foreground mb-4">
                        Get your configuration as a Base64 string for Vercel environment variables.
                    </p>
                    <button
                        onClick={copyBase64}
                        className="w-full btn btn-primary bg-primary/90 hover:bg-primary shadow-lg shadow-primary/20"
                    >
                        {copied ? <Check className="w-4 h-4 mr-2" /> : <Copy className="w-4 h-4 mr-2" />}
                        {copied ? 'Copied!' : 'Copy Base64 Config'}
                    </button>
                    <p className="text-[10px] text-muted-foreground mt-2 text-center">
                        Variable Name: <code className="bg-background px-1 py-0.5 rounded border border-border">DS2API_CONFIG_JSON</code>
                    </p>
                </div>
            </div>

            {/* Editor Panel */}
            <div className="md:col-span-2 flex flex-col bg-card border border-border rounded-xl shadow-sm overflow-hidden h-full">
                <div className="p-4 border-b border-border flex items-center justify-between bg-muted/20">
                    <h3 className="font-semibold flex items-center gap-2">
                        <Upload className="w-4 h-4 text-primary" />
                        JSON Editor
                    </h3>
                    <div className="flex gap-2">
                        <button onClick={handleExport} className="btn btn-secondary text-xs h-8">
                            Load Current
                        </button>
                        <button onClick={handleImport} disabled={loading} className="btn btn-primary text-xs h-8">
                            {loading ? 'Importing...' : 'Simulate Import'}
                        </button>
                    </div>
                </div>

                <div className="flex-1 relative">
                    <textarea
                        className="absolute inset-0 w-full h-full p-4 font-mono text-sm bg-secondary/10 resize-none focus:outline-none custom-scrollbar"
                        value={jsonInput}
                        onChange={e => setJsonInput(e.target.value)}
                        placeholder={'{\n  "keys": ["your-api-key"],\n  "accounts": [\n    {"email": "...", "password": "...", "token": ""}\n  ]\n}'}
                        spellCheck={false}
                    />
                </div>

                {result && (
                    <div className={clsx(
                        "p-4 border-t",
                        result.imported_keys || result.imported_accounts ? "bg-emerald-500/10 border-emerald-500/20" : "bg-destructive/10 border-destructive/20"
                    )}>
                        <div className="flex items-start gap-3">
                            {result.imported_keys || result.imported_accounts ? (
                                <Check className="w-5 h-5 text-emerald-500 mt-0.5" />
                            ) : (
                                <AlertTriangle className="w-5 h-5 text-destructive mt-0.5" />
                            )}
                            <div>
                                <h4 className={clsx("font-medium", result.imported_keys || result.imported_accounts ? "text-emerald-500" : "text-destructive")}>
                                    Import Operation Completed
                                </h4>
                                <p className="text-sm opacity-80 mt-1">
                                    Successfully imported {result.imported_keys} API keys and updated {result.imported_accounts} accounts.
                                </p>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}
