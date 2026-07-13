type FrontendConfig = {
  apiBaseUrl: string;
};

export const config: FrontendConfig = {
  apiBaseUrl:
    import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, "") ??
    "http://localhost:8080",
};
