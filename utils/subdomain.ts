const BASE_DOMAIN = (import.meta as any).env.VITE_BASE_DOMAIN || 'blytz.cloud';

export const getSubdomain = (): string | null => {
  const hostname = window.location.hostname;

  if (hostname === 'localhost' || hostname === '127.0.0.1') {
    return null;
  }

  const parts = hostname.split('.');

  if (parts.length < 2) {
    return null;
  }

  const baseParts = BASE_DOMAIN.split('.');

  if (parts.length <= baseParts.length) {
    return null;
  }

  const potentialSubdomain = parts.slice(0, parts.length - baseParts.length).join('.');

  if (potentialSubdomain === 'www' || potentialSubdomain.startsWith('192-168-') || potentialSubdomain.startsWith('10-0-')) {
    return null;
  }

  return potentialSubdomain || null;
};

export const isSubdomain = (): boolean => {
  return getSubdomain() !== null;
};

export const getBaseDomain = (): string => {
  return BASE_DOMAIN;
};
