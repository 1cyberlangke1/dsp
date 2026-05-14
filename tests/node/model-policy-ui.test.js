'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');

async function loadUtils() {
  return import('../../webui/src/features/settings/modelPolicy.js');
}

test('route target helpers exclude the current family and fall back to a valid target', async () => {
  const {
    buildRouteTargetOptions,
    resolveRouteTarget,
  } = await loadUtils();
  const t = (key) => key;

  assert.deepEqual(
    buildRouteTargetOptions(t, 'flash'),
    [
      { value: 'pro', label: 'settings.modelPolicy.pro' },
      { value: 'vision', label: 'settings.modelPolicy.vision' },
    ],
  );
  assert.equal(resolveRouteTarget('flash', 'flash'), 'pro');
  assert.equal(resolveRouteTarget('', 'pro'), 'flash');
  assert.equal(resolveRouteTarget('vision', 'pro'), 'vision');
});
