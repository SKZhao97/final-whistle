import test from "node:test";
import assert from "node:assert/strict";

import { groupMatchesBySeasonAndRound } from "../src/components/matches/matchGrouping.ts";
import type { MatchListItem } from "../src/types/api.ts";

function createMatch(
  id: number,
  season: string,
  round: string | undefined,
  kickoffAt: string,
): MatchListItem {
  return {
    id,
    competition: "Premier League",
    season,
    round,
    status: "FINISHED",
    kickoffAt,
    homeTeam: { id: id * 10 + 1, name: `Home ${id}`, slug: `home-${id}` },
    awayTeam: { id: id * 10 + 2, name: `Away ${id}`, slug: `away-${id}` },
    homeScore: 1,
    awayScore: 0,
    aggregates: {
      matchRatingAvg: null,
      homeTeamRatingAvg: null,
      awayTeamRatingAvg: null,
      checkInCount: 0,
    },
  };
}

test("groupMatchesBySeasonAndRound orders seasons newest first and rounds ascending", () => {
  const grouped = groupMatchesBySeasonAndRound(
    [
      createMatch(1, "2024-2025", "Matchday 2", "2026-03-28T10:00:00Z"),
      createMatch(2, "2025-2026", "Matchday 3", "2026-03-28T10:00:00Z"),
      createMatch(3, "2025-2026", "Matchday 1", "2026-03-26T10:00:00Z"),
    ],
    "Round TBD",
  );

  assert.equal(grouped[0].label, "2025-2026");
  assert.equal(grouped[0].rounds[0].label, "Matchday 1");
  assert.equal(grouped[0].rounds[1].label, "Matchday 3");
  assert.equal(grouped[1].label, "2024-2025");
});

test("groupMatchesBySeasonAndRound puts missing rounds into fallback group", () => {
  const grouped = groupMatchesBySeasonAndRound(
    [
      createMatch(1, "2024-2025", undefined, "2026-03-28T10:00:00Z"),
      createMatch(2, "2024-2025", "Matchday 1", "2026-03-26T10:00:00Z"),
    ],
    "Round TBD",
  );

  assert.equal(grouped[0].rounds[0].label, "Matchday 1");
  assert.equal(grouped[0].rounds[1].label, "Round TBD");
  assert.equal(grouped[0].rounds[1].isFallback, true);
});

test("groupMatchesBySeasonAndRound sorts fixtures inside a round by kickoff", () => {
  const grouped = groupMatchesBySeasonAndRound(
    [
      createMatch(1, "2024-2025", "Matchday 1", "2026-03-29T10:00:00Z"),
      createMatch(2, "2024-2025", "Matchday 1", "2026-03-28T10:00:00Z"),
    ],
    "Round TBD",
  );

  assert.equal(grouped[0].rounds[0].matches[0].id, 2);
  assert.equal(grouped[0].rounds[0].matches[1].id, 1);
});
