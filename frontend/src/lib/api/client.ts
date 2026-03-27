import type { UserCheckInHistoryResponse, UserProfileSummary } from "@/types/api";
import { LOCALE_COOKIE_NAME, type Locale } from "@/lib/i18n/config";

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

function getClientLocale(): Locale | null {
  if (typeof document === "undefined") {
    return null;
  }

  const match = document.cookie
    .split("; ")
    .find((item) => item.startsWith(`${LOCALE_COOKIE_NAME}=`));
  if (!match) {
    return null;
  }

  const value = match.split("=")[1];
  return value === "zh" ? "zh" : "en";
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
  const headers = new Headers(options.headers);
  if (!headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  const clientLocale = getClientLocale();
  if (clientLocale && !headers.has("X-Final-Whistle-Locale")) {
    headers.set("X-Final-Whistle-Locale", clientLocale);
  }

  const config: RequestInit = {
    ...options,
    credentials: options.credentials ?? "include",
    headers,
  };

  const response = await fetch(url, config);
  return handleResponse<T>(response);
}

export function withLocaleHeaders(locale: Locale, options: RequestInit = {}): RequestInit {
  const headers = new Headers(options.headers);
  headers.set("X-Final-Whistle-Locale", locale);

  return {
    ...options,
    headers,
  };
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

export type AuthUserResponse = {
  user: {
    id: number;
    name: string;
    avatarUrl?: string;
  };
};

export type LoginRequest = {
  email: string;
  name: string;
};

export type MatchRosterPlayer = {
  id: number;
  name: string;
  slug: string;
  position?: string;
  avatarUrl?: string;
  team: {
    id: number;
    name: string;
    shortName?: string;
    slug: string;
    logoUrl?: string;
  };
};

export type CheckInPlayerRating = {
  id: number;
  player: MatchRosterPlayer;
  rating: number;
  note?: string;
};

export type CheckInDetailResponse = {
  id: number;
  matchId: number;
  watchedType: "FULL" | "PARTIAL" | "HIGHLIGHTS";
  supporterSide: "HOME" | "AWAY" | "NEUTRAL";
  matchRating: number;
  homeTeamRating: number;
  awayTeamRating: number;
  shortReview?: string;
  watchedAt: string;
  tags: Array<{ id: number; name: string; slug: string }>;
  playerRatings: CheckInPlayerRating[];
  createdAt: string;
  updatedAt: string;
};

export type MatchTagOption = {
  id: number;
  name: string;
  slug: string;
};

export type UpsertCheckInRequest = {
  watchedType: "FULL" | "PARTIAL" | "HIGHLIGHTS";
  supporterSide: "HOME" | "AWAY" | "NEUTRAL";
  matchRating: number;
  homeTeamRating: number;
  awayTeamRating: number;
  shortReview?: string;
  watchedAt: string;
  tags: number[];
  playerRatings: Array<{
    playerId: number;
    rating: number;
    note?: string;
  }>;
};

type MatchListQuery = {
  competition?: string;
  season?: string;
  page?: number;
  pageSize?: number;
};

function withQuery(path: string, query?: Record<string, string | number | undefined>) {
  const params = new URLSearchParams();
  Object.entries(query ?? {}).forEach(([key, value]) => {
    if (value !== undefined && value !== "") {
      params.set(key, String(value));
    }
  });
  const queryString = params.toString();
  return queryString ? `${path}?${queryString}` : path;
}

export const matchesApi = {
  list: <T = unknown>(query?: MatchListQuery, options?: RequestInit) =>
    api.get<T>(withQuery("/matches", query), options),
  detail: <T = unknown>(matchId: string | number, options?: RequestInit) =>
    api.get<T>(`/matches/${matchId}`, options),
  myCheckIn: (matchId: string | number, options?: RequestInit) =>
    api.get<CheckInDetailResponse | null>(`/matches/${matchId}/my-checkin`, options),
  createCheckIn: (matchId: string | number, body: UpsertCheckInRequest, options?: RequestInit) =>
    api.post<CheckInDetailResponse>(`/matches/${matchId}/checkin`, body, options),
  updateCheckIn: (matchId: string | number, body: UpsertCheckInRequest, options?: RequestInit) =>
    api.put<CheckInDetailResponse>(`/matches/${matchId}/checkin`, body, options),
};

export const teamsApi = {
  detail: <T = unknown>(teamId: string | number, options?: RequestInit) =>
    api.get<T>(`/teams/${teamId}`, options),
};

export const playersApi = {
  detail: <T = unknown>(playerId: string | number, options?: RequestInit) =>
    api.get<T>(`/players/${playerId}`, options),
};

export const authApi = {
  login: (body: LoginRequest, options?: RequestInit) =>
    api.post<AuthUserResponse>("/auth/login", body, options),
  logout: (options?: RequestInit) =>
    api.post<{ ok: true }>("/auth/logout", undefined, options),
  me: (options?: RequestInit) =>
    api.get<AuthUserResponse>("/auth/me", options),
};

export const usersApi = {
  profile: (options?: RequestInit) =>
    api.get<UserProfileSummary>("/me/profile", options),
  checkins: (
    query?: { page?: number; pageSize?: number },
    options?: RequestInit,
  ) => api.get<UserCheckInHistoryResponse>(withQuery("/me/checkins", query), options),
};
