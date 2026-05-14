import { buildRouteTargetOptions, resolveRouteTarget } from './modelPolicy'

export default function ModelSection({ t, form, setForm }) {
    const families = ['flash', 'pro', 'vision']
    const updateFamilyRule = (family, patch) => setForm((prev) => ({
        ...prev,
        model_family_policy: {
            ...prev.model_family_policy,
            [family]: {
                ...prev.model_family_policy[family],
                ...patch,
            },
        },
    }))
    const updateToolRule = (family, enabled) => setForm((prev) => ({
        ...prev,
        model_tool_policy: {
            ...prev.model_tool_policy,
            [family]: {
                ...prev.model_tool_policy[family],
                enabled,
            },
        },
    }))
    return (
        <div className="bg-card border border-border rounded-xl p-5 space-y-4">
            <h3 className="font-semibold">{t('settings.modelTitle')}</h3>
            <div className="space-y-3">
                <div className="text-sm text-muted-foreground">{t('settings.modelPolicyTitle')}</div>
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
                    {families.map((family) => {
                        const rule = form.model_family_policy?.[family] || { mode: 'allow', target: '' }
                        const routeTargets = buildRouteTargetOptions(t, family)
                        const targetValue = rule.mode === 'route'
                            ? resolveRouteTarget(rule.target, family)
                            : ''
                        return (
                            <div key={family} className="rounded-lg border border-border bg-background/60 p-4 space-y-3">
                                <div className="text-sm font-medium">{t(`settings.modelPolicy.${family}`)}</div>
                                <label className="text-xs space-y-1 block">
                                    <span className="text-muted-foreground">{t('settings.modelPolicyMode')}</span>
                                    <select
                                        value={rule.mode}
                                        onChange={(e) => updateFamilyRule(family, {
                                            mode: e.target.value,
                                            target: e.target.value === 'route'
                                                ? resolveRouteTarget(form.model_family_policy?.[family]?.target, family)
                                                : '',
                                        })}
                                        className="w-full bg-background border border-border rounded-lg px-3 py-2 text-sm"
                                    >
                                        <option value="allow">{t('settings.modelPolicyModeAllow')}</option>
                                        <option value="disable">{t('settings.modelPolicyModeDisable')}</option>
                                        <option value="route">{t('settings.modelPolicyModeRoute')}</option>
                                    </select>
                                </label>
                                <label className="text-xs space-y-1 block">
                                    <span className="text-muted-foreground">{t('settings.modelPolicyTarget')}</span>
                                    <select
                                        value={targetValue}
                                        disabled={rule.mode !== 'route'}
                                        onChange={(e) => updateFamilyRule(family, { target: e.target.value })}
                                        className="w-full bg-background border border-border rounded-lg px-3 py-2 text-sm disabled:opacity-50"
                                    >
                                        {rule.mode !== 'route' && (
                                            <option value="">{t('settings.modelPolicyTargetUnset')}</option>
                                        )}
                                        {routeTargets.map((target) => (
                                            <option key={target.value} value={target.value}>{target.label}</option>
                                        ))}
                                    </select>
                                </label>
                                <label className="text-xs space-y-1 block">
                                    <span className="text-muted-foreground">{t('settings.modelToolPolicyEnabled')}</span>
                                    <select
                                        value={form.model_tool_policy?.[family]?.enabled === false ? 'false' : 'true'}
                                        onChange={(e) => updateToolRule(family, e.target.value === 'true')}
                                        className="w-full bg-background border border-border rounded-lg px-3 py-2 text-sm"
                                    >
                                        <option value="true">{t('settings.modelToolPolicyOn')}</option>
                                        <option value="false">{t('settings.modelToolPolicyOff')}</option>
                                    </select>
                                </label>
                            </div>
                        )
                    })}
                </div>
            </div>
            <label className="text-sm space-y-2 block">
                <span className="text-muted-foreground">{t('settings.modelAliases')}</span>
                <textarea
                    value={form.model_aliases_text}
                    onChange={(e) => setForm((prev) => ({ ...prev, model_aliases_text: e.target.value }))}
                    rows={12}
                    className="w-full bg-background border border-border rounded-lg px-3 py-2 font-mono text-xs"
                />
            </label>
        </div>
    )
}
