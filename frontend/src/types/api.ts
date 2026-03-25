// Base response envelope
export interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}

// Error codes
export type ApiErrorCode =
  | 'VALIDATION_ERROR'
  | 'UNAUTHORIZED'
  | 'FORBIDDEN'
  | 'NOT_FOUND'
  | 'CONFLICT'
  | 'INTERNAL_ERROR';

// Pagination
export interface PaginatedResponse<T> {
  items: T[];
  page: number;
  pageSize: number;
  total: number;
}

// Common filters
export interface PaginationParams {
  page?: number;
  pageSize?: number;
}

export interface User {
  id: number;
  name: string;
  email: string;
  avatarUrl?: string;
  createdAt: string;
  updatedAt: string;
}

// Enums from backend
export enum MatchStatus {
  SCHEDULED = 'SCHEDULED',
  FINISHED = 'FINISHED',
}

export enum WatchedType {
  FULL = 'FULL',
  PARTIAL = 'PARTIAL',
  HIGHLIGHTS = 'HIGHLIGHTS',
}

export enum SupporterSide {
  HOME = 'HOME',
  AWAY = 'AWAY',
  NEUTRAL = 'NEUTRAL',
}

// Core entities
export interface Team {
  id: number;
  name: string;
  shortName: string;
  slug: string;
  logoUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Player {
  id: number;
  teamId: number;
  name: string;
  slug: string;
  position?: string;
  avatarUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Match {
  id: number;
  competition: string;
  season: string;
  round?: string;
  status: MatchStatus;
  kickoffAt: string;
  homeTeamId: number;
  awayTeamId: number;
  homeScore?: number;
  awayScore?: number;
  venue?: string;
  createdAt: string;
  updatedAt: string;
}

export interface MatchPlayer {
  id: number;
  matchId: number;
  playerId: number;
  teamId: number;
}

export interface Tag {
  id: number;
  name: string;
  slug: string;
  sortOrder: number;
  isActive: boolean;
}

export interface CheckInTag {
  id: number;
  checkInId: number;
  tagId: number;
}

export interface Session {
  id: number;
  userId: number;
  token: string;
  expiredAt: string;
  createdAt: string;
}
