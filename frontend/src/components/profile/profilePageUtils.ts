import type { UserCheckInHistoryItem, UserProfileSummary } from "@/types/api";
import type { Locale } from "../../lib/i18n/config";
import { translate } from "../../lib/i18n/core";

export function formatAverageRating(value: number | undefined, locale: Locale = "en") {
  if (value === undefined || Number.isNaN(value)) {
    return translate(locale, "stats.noRatingsYet");
  }
  return value.toFixed(1);
}

export function buildProfileStats(profile: UserProfileSummary, locale: Locale = "en") {
  return [
    { label: translate(locale, "stats.checkIns"), value: String(profile.checkInCount) },
    { label: translate(locale, "stats.avgMatchRating"), value: formatAverageRating(profile.avgMatchRating, locale) },
    { label: translate(locale, "stats.recent30Days"), value: String(profile.recentCheckInCount) },
    {
      label: translate(locale, "stats.favoriteTeam"),
      value: profile.favoriteTeam?.name ?? translate(locale, "stats.notEnoughData"),
    },
    {
      label: translate(locale, "stats.mostUsedTag"),
      value: profile.mostUsedTag?.name ?? translate(locale, "stats.notEnoughData"),
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

export function buildHistorySummary(item: UserCheckInHistoryItem, locale: Locale = "en") {
  return locale === "zh"
    ? `${item.match.homeTeam.name} 对 ${item.match.awayTeam.name}`
    : `${item.match.homeTeam.name} vs ${item.match.awayTeam.name}`;
}
