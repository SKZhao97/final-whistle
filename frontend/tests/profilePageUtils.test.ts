import test from "node:test";
import assert from "node:assert/strict";

import {
  buildHistorySummary,
  buildPaginationMeta,
  buildProfileStats,
  formatAverageRating,
} from "../src/components/profile/profilePageUtils.ts";
import type { UserCheckInHistoryItem, UserProfileSummary } from "../src/types/api.ts";

test("formatAverageRating returns fallback when missing", () => {
  assert.equal(formatAverageRating(undefined), "No ratings yet");
});

test("buildProfileStats uses friendly labels", () => {
  const profile: UserProfileSummary = {
    user: { id: 1, name: "Demo User" },
    checkInCount: 12,
    avgMatchRating: 8.25,
    favoriteTeamId: 2,
    favoriteTeam: { id: 2, name: "Arsenal", slug: "arsenal" },
    mostUsedTagId: 3,
    mostUsedTag: { id: 3, name: "Electric", slug: "electric" },
    recentCheckInCount: 4,
  };

  const stats = buildProfileStats(profile);
  assert.equal(stats[0].value, "12");
  assert.equal(stats[1].value, "8.3");
  assert.equal(stats[3].value, "Arsenal");
  assert.equal(stats[4].value, "Electric");
});

test("buildPaginationMeta bounds next and previous buttons", () => {
  assert.deepEqual(buildPaginationMeta(25, 1, 10), {
    totalPages: 3,
    canGoPrev: false,
    canGoNext: true,
  });
});

test("buildHistorySummary returns readable match label", () => {
  const item = {
    id: 1,
    matchId: 2,
    watchedType: "FULL",
    supporterSide: "HOME",
    matchRating: 8,
    homeTeamRating: 9,
    awayTeamRating: 6,
    watchedAt: "2026-03-26T10:00:00Z",
    createdAt: "2026-03-26T10:00:00Z",
    updatedAt: "2026-03-26T10:00:00Z",
    tags: [],
    match: {
      id: 2,
      competition: "Premier League",
      season: "2025/26",
      status: "FINISHED",
      kickoffAt: "2026-03-25T18:30:00Z",
      homeTeam: { id: 1, name: "Liverpool", slug: "liverpool" },
      awayTeam: { id: 2, name: "Arsenal", slug: "arsenal" },
      aggregates: {
        matchRatingAvg: null,
        homeTeamRatingAvg: null,
        awayTeamRatingAvg: null,
        checkInCount: 0,
      },
    },
  } as UserCheckInHistoryItem;

  assert.equal(buildHistorySummary(item), "Liverpool vs Arsenal");
});
