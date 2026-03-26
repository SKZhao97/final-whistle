import type { UserCheckInHistoryItem, UserProfileSummary } from "@/types/api";

export function formatAverageRating(value?: number) {
  if (value === undefined || Number.isNaN(value)) {
    return "No ratings yet";
  }
  return value.toFixed(1);
}

export function buildProfileStats(profile: UserProfileSummary) {
  return [
    { label: "Check-ins", value: String(profile.checkInCount) },
    { label: "Avg Match Rating", value: formatAverageRating(profile.avgMatchRating) },
    { label: "Recent 30 Days", value: String(profile.recentCheckInCount) },
    {
      label: "Favorite Team",
      value: profile.favoriteTeam?.name ?? "Not enough data",
    },
    {
      label: "Most Used Tag",
      value: profile.mostUsedTag?.name ?? "Not enough data",
    },
  ];
}

export function buildPaginationMeta(total: number, page: number, pageSize: number) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  return {
    totalPages,
    canGoPrev: page > 1,
    canGoNext: page < totalPages,
  };
}

export function buildHistorySummary(item: UserCheckInHistoryItem) {
  return `${item.match.homeTeam.name} vs ${item.match.awayTeam.name}`;
}
