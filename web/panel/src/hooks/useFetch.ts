import { useCallback } from 'react';

interface FetchOptions {
  headers?: Record<string, string>;
  signal?: AbortSignal;
  timeout?: number;
  method?: string;
  body?: any;
}

export default function useFetch() {
  const fetchData = useCallback(
    async <T>(url: string, options?: FetchOptions) => {
      const controller = new AbortController();
      const { signal } = controller;

      const timeoutId = options?.timeout
        ? setTimeout(() => {
            controller.abort();
          }, options.timeout)
        : null;

      try {
        const isObject = typeof options?.body === 'object';
        const response = await fetch(url, {
          method: options?.method || 'GET',
          credentials: 'include',
          headers: {
            'Content-Type': isObject ? 'application/json' : 'text/plain',
            ...options?.headers,
          },
          body: options?.body ?
            isObject ? JSON.stringify(options.body)
              : options.body
            : undefined,
          signal : options?.signal || signal,
        });

        const data = await response.json();
        return { data: data as T, code: response.status };
      } catch (error: any) {
        if (signal.aborted) return { error: 'Request aborted due to timeout' };
        return { error: error.message || 'Something went wrong' };
      } finally {
        if (timeoutId) clearTimeout(timeoutId);
      }
    }, []);

  return fetchData;
}
