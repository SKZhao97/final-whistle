const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export type ApiResponse<T = unknown> = {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
};

export class ApiError extends Error {
  constructor(
    public code: string,
    message: string,
    public details?: Record<string, unknown>
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  const contentType = response.headers.get('content-type');
  const isJson = contentType && contentType.includes('application/json');

  if (!response.ok) {
    if (isJson) {
      const errorData = await response.json();
      throw new ApiError(
        errorData.error?.code || "HTTP_ERROR",
        errorData.error?.message || `HTTP ${response.status}`,
        errorData.error?.details
      );
    } else {
      throw new ApiError("HTTP_ERROR", `HTTP ${response.status}: ${response.statusText}`);
    }
  }

  if (isJson) {
    const data: ApiResponse<T> = await response.json();
    if (!data.success) {
      throw new ApiError(
        data.error?.code || "API_ERROR",
        data.error?.message || "Unknown API error",
        data.error?.details
      );
    }
    return data.data as T;
  }

  return (await response.text()) as unknown as T;
}

export async function apiRequest<T = unknown>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  const headers = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  const config: RequestInit = {
    ...options,
    headers,
  };

  const response = await fetch(url, config);
  return handleResponse<T>(response);
}

// Convenience methods
export const api = {
  get: <T = unknown>(endpoint: string, options?: RequestInit) =>
    apiRequest<T>(endpoint, { ...options, method: "GET" }),

  post: <T = unknown>(endpoint: string, body?: unknown, options?: RequestInit) =>
    apiRequest<T>(endpoint, {
      ...options,
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    }),

  put: <T = unknown>(endpoint: string, body?: unknown, options?: RequestInit) =>
    apiRequest<T>(endpoint, {
      ...options,
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    }),

  delete: <T = unknown>(endpoint: string, options?: RequestInit) =>
    apiRequest<T>(endpoint, { ...options, method: "DELETE" }),
};
