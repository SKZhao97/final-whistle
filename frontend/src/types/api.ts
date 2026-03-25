export interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}

export type ApiErrorCode =
  | "VALIDATION_ERROR"
  | "UNAUTHORIZED"
  | "FORBIDDEN"
  | "NOT_FOUND"
  | "CONFLICT"
  | "INTERNAL_ERROR";

export interface PaginationParams {
  page?: number;
  pageSize?: number;
}

export interface TeamSummary {
  id: number;
  name: string;
  shortName?: string;
  slug: string;
  logoUrl?: string;
}

export interface UserSummary {
  id: number;
  name: string;
  avatarUrl?: string;
}

export interface LoginRequest {
  email: string;
  name: string;
}

export interface AuthUserResponse {
  user: UserSummary;
}

export interface LogoutResponse {
  ok: true;
}

export interface TagSummary {
  id: number;
  name: string;
  slug: string;
}

export interface MatchAggregateSummary {
  matchRatingAvg: number | null;
  homeTeamRatingAvg: number | null;
  awayTeamRatingAvg: number | null;
  checkInCount: number;
}

export interface MatchListItem {
  id: number;
  competition: string;
  season: string;
  round?: string;
  status: string;
  kickoffAt: string;
  homeTeam: TeamSummary;
  awayTeam: TeamSummary;
  homeScore?: number;
  awayScore?: number;
  aggregates: MatchAggregateSummary;
}

export interface MatchListResponse {
  items: MatchListItem[];
  page: number;
  pageSize: number;
  total: number;
}

export interface PlayerSummary {
  id: number;
  name: string;
  slug: string;
  position?: string;
  avatarUrl?: string;
  team: TeamSummary;
}

export interface MatchPlayerRatingSummary {
  player: PlayerSummary;
  avgRating: number | null;
  ratingCount: number;
}

export interface MatchRecentReview {
  id: number;
  user: UserSummary;
  matchRating: number;
  shortReview: string;
  tags: TagSummary[];
  createdAt: string;
}

export interface MatchDetail {
  id: number;
  competition: string;
  season: string;
  round?: string;
  status: string;
  kickoffAt: string;
  homeTeam: TeamSummary;
  awayTeam: TeamSummary;
  homeScore?: number;
  awayScore?: number;
  venue?: string;
  aggregates: MatchAggregateSummary;
  playerRatings: MatchPlayerRatingSummary[];
  recentReviews: MatchRecentReview[];
}

export interface TeamRatingSummary {
  avgRating: number | null;
  ratingCount: number;
}

export interface TeamDetail {
  id: number;
  name: string;
  shortName?: string;
  slug: string;
  logoUrl?: string;
  recentMatches: MatchListItem[];
  ratingSummary: TeamRatingSummary;
}

export interface PlayerRecentMatch {
  match: MatchListItem;
  avgRating: number | null;
  ratingCount: number;
}

export interface PlayerRatingSummary {
  avgRating: number | null;
  ratingCount: number;
}

export interface PlayerDetail {
  id: number;
  name: string;
  slug: string;
  position?: string;
  avatarUrl?: string;
  team: TeamSummary;
  recentMatches: PlayerRecentMatch[];
  ratingSummary: PlayerRatingSummary;
}
