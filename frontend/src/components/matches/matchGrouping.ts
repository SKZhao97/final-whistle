import type { MatchListItem } from "@/types/api";

export type GroupedRound = {
  id: string;
  label: string;
  round: string | null;
  matches: MatchListItem[];
  isFallback: boolean;
};

export type GroupedSeason = {
  id: string;
  label: string;
  rounds: GroupedRound[];
};

function parseSeasonSortValue(season: string) {
  const explicitYears = season.match(/\d{4}/g);
  if (explicitYears && explicitYears.length > 0) {
    return Number(explicitYears[0]);
  }

  const compactYears = season.match(/\d+/g);
  if (compactYears && compactYears.length > 0) {
    return Number(compactYears[0]);
  }

  return -1;
}

function parseRoundSortValue(round: string | null) {
  if (!round) {
    return Number.POSITIVE_INFINITY;
  }

  const numeric = round.match(/\d+/);
  if (numeric) {
    return Number(numeric[0]);
  }

  return Number.POSITIVE_INFINITY - 1;
}

function compareKickoff(a: MatchListItem, b: MatchListItem) {
  return new Date(a.kickoffAt).getTime() - new Date(b.kickoffAt).getTime();
}

export function groupMatchesBySeasonAndRound(
  matches: MatchListItem[],
  fallbackRoundLabel: string,
): GroupedSeason[] {
  const seasons = new Map<string, Map<string, GroupedRound>>();

  for (const match of matches) {
    const seasonKey = match.season;
    const roundKey = match.round?.trim() || "__fallback_round__";
    const roundLabel = match.round?.trim() || fallbackRoundLabel;

    if (!seasons.has(seasonKey)) {
      seasons.set(seasonKey, new Map());
    }

    const rounds = seasons.get(seasonKey)!;
    if (!rounds.has(roundKey)) {
      rounds.set(roundKey, {
        id: `${seasonKey}-${roundKey}`,
        label: roundLabel,
        round: match.round ?? null,
        matches: [],
        isFallback: !match.round,
      });
    }

    rounds.get(roundKey)!.matches.push(match);
  }

  return Array.from(seasons.entries())
    .sort((a, b) => parseSeasonSortValue(b[0]) - parseSeasonSortValue(a[0]) || b[0].localeCompare(a[0]))
    .map(([season, rounds]) => ({
      id: season,
      label: season,
      rounds: Array.from(rounds.values())
        .map((group) => ({
          ...group,
          matches: [...group.matches].sort(compareKickoff),
        }))
        .sort((a, b) => {
          if (a.isFallback !== b.isFallback) {
            return a.isFallback ? 1 : -1;
          }
          const byRound = parseRoundSortValue(a.round) - parseRoundSortValue(b.round);
          if (byRound !== 0) {
            return byRound;
          }
          return a.label.localeCompare(b.label);
        }),
    }));
}
