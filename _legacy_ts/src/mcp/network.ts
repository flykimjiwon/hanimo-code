import { resolve } from 'node:dns/promises';

const DNS_TIMEOUT_MS = 3000;

export async function detectNetworkMode(
  configured: 'auto' | 'online' | 'offline'
): Promise<'online' | 'offline'> {
  if (configured === 'online' || configured === 'offline') {
    return configured;
  }

  // configured === 'auto': probe DNS to determine connectivity
  try {
    await Promise.race([
      resolve('dns.google'),
      new Promise<never>((_, reject) =>
        setTimeout(() => reject(new Error('DNS timeout')), DNS_TIMEOUT_MS)
      ),
    ]);
    return 'online';
  } catch {
    return 'offline';
  }
}
