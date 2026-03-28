import Link from "next/link";

import { TeamCrest } from "@/components/experience/FootballPrimitives";
import { formatDateTime, formatNumber } from "@/lib/i18n/domain";
import { translate } from "@/lib/i18n/core";
import type { Locale } from "@/lib/i18n/config";
import type { MatchListItem } from "@/types/api";

export function MatchFixtureCard({
  match,
  locale,
}: {
  match: MatchListItem;
  locale: Locale;
}) {
  return (
    <Link
      href={`/matches/${match.id}`}
      className="match-shell transition-transform duration-150 hover:-translate-y-0.5 hover:border-[var(--fw-field-300)]"
    >
      <div className="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
        <div className="space-y-4">
          <div className="flex items-center gap-4">
            <div className="flex min-w-0 items-center gap-3">
              <TeamCrest team={match.homeTeam} size="sm" />
              <div className="min-w-0">
                <p className="text-xs uppercase tracking-[0.18em] text-[var(--fw-muted)]">
                  {translate(locale, "enum.supporterSide.home")}
                </p>
                <p className="truncate text-lg font-semibold text-[var(--fw-ink)]">
                  {match.homeTeam.name}
                </p>
              </div>
            </div>

            <div className="rounded-[1.1rem] border border-[var(--fw-line)] bg-[var(--fw-surface)]/9 px-4 py-2 text-center">
              <p className="text-2xl font-semibold tracking-[-0.08em] text-[var(--fw-score)]">
                {typeof match.homeScore === "number" ? match.homeScore : "-"}:
                {typeof match.awayScore === "number" ? match.awayScore : "-"}
              </p>
            </div>

            <div className="flex min-w-0 items-center gap-3">
              <div className="min-w-0 text-right">
                <p className="text-xs uppercase tracking-[0.18em] text-[var(--fw-muted)]">
                  {translate(locale, "enum.supporterSide.away")}
                </p>
                <p className="truncate text-lg font-semibold text-[var(--fw-ink)]">
                  {match.awayTeam.name}
                </p>
              </div>
              <TeamCrest team={match.awayTeam} size="sm" />
            </div>
          </div>

          <p className="text-sm text-[var(--fw-muted)]">
            {formatDateTime(match.kickoffAt, locale)}
          </p>
        </div>

        <div className="grid gap-2 text-sm text-[var(--fw-muted)] lg:text-right">
          <p>{translate(locale, "matches.checkIns", { count: match.aggregates.checkInCount })}</p>
          <p>{translate(locale, "matches.avgRating", { value: formatNumber(match.aggregates.matchRatingAvg, locale) })}</p>
        </div>
      </div>
    </Link>
  );
}
