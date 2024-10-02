import { useCallback } from 'react';

interface FetchOptions {
  method?: string;
  headers?: Record<string, string>;
  body?: any;
  timeout?: number; // ms
}

export default function useFetch() {
  const fetchData = useCallback(
    async (url: string, options?: FetchOptions) => {
      const controller = new AbortController();
      const { signal } = controller;

      const timeoutId = options?.timeout
        ? setTimeout(() => {
            controller.abort();
          }, options.timeout)
        : null;

      try {
        const response = await fetch(url, {
          method: options?.method || 'GET',
          headers: {
            'Content-Type': 'application/json',
            ...options?.headers,
          },
          body: options?.body ?
            typeof options.body === "object" ? 
              JSON.stringify(options.body)
              : options.body
            : undefined,
          signal,
        });

        if (!response.ok) {
          throw new Error(`Error: ${response.status}`);
        }

        const data = await response.json();

        return { data, error: null };
      } catch (error: any) {
        if (signal.aborted) {
          return { data: null, error: 'Request aborted due to timeout' };
        }
        return { data: null, error: error.message || 'Something went wrong' };
      } finally {
        if (timeoutId) {
          clearTimeout(timeoutId);
        }
      }
    }, []);

  return fetchData;
}
