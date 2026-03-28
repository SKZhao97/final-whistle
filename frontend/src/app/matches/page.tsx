import { ArchivePill, LeagueMark, SectionShell } from "@/components/experience/FootballPrimitives";
import { MatchFixtureCard } from "@/components/matches/MatchFixtureCard";
import { SeasonSelector } from "@/components/matches/SeasonSelector";
import { groupMatchesBySeasonAndRound } from "@/components/matches/matchGrouping";
import { matchesApi, withLocaleHeaders } from "@/lib/api/client";
import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";
import type { MatchListResponse } from "@/types/api";

type MatchesPageProps = {
  searchParams?: Promise<Record<string, string | string[] | undefined>>;
};

function readSeasonParam(raw: string | string[] | undefined) {
  if (Array.isArray(raw)) {
    return raw[0] ?? null;
  }
  return raw ?? null;
}

export default async function MatchesPage({ searchParams }: MatchesPageProps) {
  const locale = await getServerLocale();
  const resolvedSearchParams = searchParams ? await searchParams : {};
  const data = await matchesApi.list<MatchListResponse>(
    { page: 1, pageSize: 1000 },
    withLocaleHeaders(locale, { cache: "no-store" }),
  );
  const groupedSeasons = groupMatchesBySeasonAndRound(data.items, translate(locale, "matches.roundTbd"));
  const selectedSeasonParam = readSeasonParam(resolvedSearchParams.season);
  const selectedSeason =
    groupedSeasons.find((season) => season.id === selectedSeasonParam) ?? groupedSeasons[0] ?? null;

  return (
    <div className="space-y-8 py-8">
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
      ) : selectedSeason ? (
        <>
          <SectionShell
            eyebrow={translate(locale, "matches.seasonEyebrow")}
            title={selectedSeason.label}
            description={translate(locale, "matches.seasonDescription")}
            accent="field"
            headerAside={
              <SeasonSelector
                label={translate(locale, "matches.seasonSelector")}
                value={selectedSeason.id}
                options={groupedSeasons.map((season) => ({
                  id: season.id,
                  label: season.label,
                }))}
              />
            }
          >
            <div className="mt-6 space-y-6">
              {selectedSeason.rounds.map((round) => (
                <section key={round.id} className="space-y-4">
                  <div className="flex flex-wrap items-center justify-between gap-3">
                    <div className="flex flex-wrap items-center gap-3">
                      <LeagueMark label={round.label} />
                      {round.isFallback ? (
                        <ArchivePill>{translate(locale, "matches.roundFallback")}</ArchivePill>
                      ) : null}
                    </div>
                    <p className="text-sm text-[var(--fw-muted)]">
                      {translate(locale, "matches.roundCount", { count: round.matches.length })}
                    </p>
                  </div>

                  <div className="grid gap-4">
                    {round.matches.map((match) => (
                      <MatchFixtureCard key={match.id} match={match} locale={locale} />
                    ))}
                  </div>
                </section>
              ))}
            </div>
          </SectionShell>
        </>
      ) : (
        <div className="rounded-[1.4rem] border border-dashed border-[var(--fw-line)] bg-[var(--fw-surface)]/75 p-8 text-sm text-[var(--fw-muted)]">
          {translate(locale, "matches.empty")}
        </div>
      )}

      <div className="rounded-[1.4rem] border border-dashed border-[var(--fw-line)] bg-[var(--fw-surface)]/65 p-5 text-sm leading-6 text-[var(--fw-muted)]">
        {translate(locale, "matches.futureBrowseHint")}
      </div>
    </div>
  );
}
