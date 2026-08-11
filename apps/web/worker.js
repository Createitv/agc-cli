const removedDistributionPaths = new Set(['/install.sh', '/install.ps1']);

export default {
  fetch(request, env) {
    const url = new URL(request.url);

    if (removedDistributionPaths.has(url.pathname) || url.pathname.startsWith('/downloads/')) {
      return new Response('Not found', {
        status: 404,
        headers: {
          'content-type': 'text/plain; charset=utf-8',
          'x-content-type-options': 'nosniff',
        },
      });
    }

    return env.ASSETS.fetch(request);
  },
};
