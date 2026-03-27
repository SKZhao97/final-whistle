import Link from "next/link";

import { LeagueMark, TeamCrest } from "@/components/experience/FootballPrimitives";
import { matchesApi, withLocaleHeaders } from "@/lib/api/client";
import { formatDateTime, formatNumber } from "@/lib/i18n/domain";
import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";
import type { MatchListResponse } from "@/types/api";

export default async function MatchesPage() {
  const locale = await getServerLocale();
  const data = await matchesApi.list<MatchListResponse>(
    { page: 1, pageSize: 20 },
    withLocaleHeaders(locale, { cache: "no-store" }),
  );

  return (
    <div className="py-8">
      <div className="mb-8 flex items-end justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-[var(--fw-ink)]">{translate(locale, "matches.title")}</h1>
          <p className="mt-2 text-sm text-[var(--fw-muted)]">
            {translate(locale, "matches.subtitle")}
          </p>
        </div>
        <p className="text-sm text-[var(--fw-muted)]">{translate(locale, "matches.total", { total: data.total })}</p>
      </div>

      {data.items.length === 0 ? (
        <div className="rounded-[1.4rem] border border-dashed border-[var(--fw-line)] bg-[var(--fw-surface)]/75 p-8 text-sm text-[var(--fw-muted)]">
          {translate(locale, "matches.empty")}
        </div>
      ) : (
        <div className="grid gap-4">
          {data.items.map((match) => (
            <Link
              key={match.id}
              href={`/matches/${match.id}`}
              className="match-shell transition-transform duration-150 hover:-translate-y-0.5 hover:border-[var(--fw-field-300)]"
            >
              <div className="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
                <div className="space-y-4">
                  <div className="flex flex-wrap items-center gap-3">
                    <LeagueMark label={match.competition} />
                    <span className="inline-flex items-center rounded-full border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] px-3 py-1 text-xs font-medium text-[var(--fw-ink-soft)]">
                      {match.season}
                    </span>
                  </div>

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
                    {match.round ?? translate(locale, "matches.roundTbd")} · {formatDateTime(match.kickoffAt, locale)}
                  </p>
                </div>

                <div className="grid gap-2 text-sm text-[var(--fw-muted)] lg:text-right">
                  <p>{translate(locale, "matches.checkIns", { count: match.aggregates.checkInCount })}</p>
                  <p>{translate(locale, "matches.avgRating", { value: formatNumber(match.aggregates.matchRatingAvg, locale) })}</p>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
