const MODEL_POLICY_FAMILIES = ['flash', 'pro', 'vision']

export function buildRouteTargetOptions(t, family) {
    return MODEL_POLICY_FAMILIES
        .filter((candidate) => candidate !== family)
        .map((candidate) => ({
            value: candidate,
            label: t(`settings.modelPolicy.${candidate}`),
        }))
}

export function resolveRouteTarget(target, family) {
    const normalizedTarget = String(target || '').trim().toLowerCase()
    const options = buildRouteTargetOptions((key) => key, family)
    if (options.some((option) => option.value === normalizedTarget)) {
        return normalizedTarget
    }
    return options[0]?.value || ''
}
