// In a utils/fetch.js file (or similar)
export const customFetch = async (url: string, options = {}) => {
  // Merge provided options with default credentials: 'include'
  const fetchOptions = {
    ...options,
    credentials: 'include' as any,
  };
  
  return fetch(url, fetchOptions);
};