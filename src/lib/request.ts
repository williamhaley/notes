import { getToken } from './stores/auth';

export default function request(url: string, options: RequestInit = {}) {
  return fetch(url, {
    ...options,
    headers: {
      token: getToken(),
      ...options.headers,
      'content-type': 'application/json'
    }
  });
}
