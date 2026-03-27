import Link from "next/link";

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
          <h1 className="text-3xl font-bold tracking-tight">{translate(locale, "matches.title")}</h1>
          <p className="mt-2 text-sm text-neutral-600">
            {translate(locale, "matches.subtitle")}
          </p>
        </div>
        <p className="text-sm text-neutral-500">{translate(locale, "matches.total", { total: data.total })}</p>
      </div>

      {data.items.length === 0 ? (
        <div className="rounded-xl border border-dashed p-8 text-sm text-neutral-600">
          {translate(locale, "matches.empty")}
        </div>
      ) : (
        <div className="grid gap-4">
          {data.items.map((match) => (
            <Link
              key={match.id}
              href={`/matches/${match.id}`}
              className="rounded-xl border p-5 transition-colors hover:bg-neutral-50"
            >
              <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
                <div>
                  <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">
                    {match.competition} · {match.season}
                  </p>
                  <h2 className="mt-1 text-xl font-semibold">
                    {match.homeTeam.name} {typeof match.homeScore === "number" ? match.homeScore : "-"}:
                    {typeof match.awayScore === "number" ? match.awayScore : "-"} {match.awayTeam.name}
                  </h2>
                  <p className="mt-1 text-sm text-neutral-600">
                    {match.round ?? translate(locale, "matches.roundTbd")} · {formatDateTime(match.kickoffAt, locale)}
                  </p>
                </div>
                <div className="text-sm text-neutral-600">
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
