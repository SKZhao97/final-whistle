import Link from "next/link";
import { notFound } from "next/navigation";

import MatchCheckInPanel from "@/components/checkin/MatchCheckInPanel";
import {
  ArchivePill,
  LeagueMark,
  SectionShell,
  TeamCrest,
} from "@/components/experience/FootballPrimitives";
import { ApiError, matchesApi, withLocaleHeaders } from "@/lib/api/client";
import { translate } from "@/lib/i18n/core";
import { formatDateTime, formatNumber } from "@/lib/i18n/domain";
import { getServerLocale } from "@/lib/i18n/server";
import type { MatchDetail } from "@/types/api";

type MatchDetailPageProps = {
  params: Promise<{ matchId: string }>;
};

async function getMatchDetail(matchId: string, locale: "en" | "zh") {
  try {
    return await matchesApi.detail<MatchDetail>(
      matchId,
      withLocaleHeaders(locale, { cache: "no-store" }),
    );
  } catch (error) {
    if (error instanceof ApiError && error.code === "NOT_FOUND") {
      notFound();
    }
    throw error;
  }
}

export default async function MatchDetailPage({ params }: MatchDetailPageProps) {
  const { matchId } = await params;
  const locale = await getServerLocale();
  const match = await getMatchDetail(matchId, locale);
  const normalizedMatch: MatchDetail = {
    ...match,
    availableTags: match.availableTags ?? [],
    matchPlayers: match.matchPlayers ?? [],
    playerRatings: match.playerRatings ?? [],
    recentReviews: match.recentReviews ?? [],
  };

  const scoreline = `${typeof normalizedMatch.homeScore === "number" ? normalizedMatch.homeScore : "-"}:${typeof normalizedMatch.awayScore === "number" ? normalizedMatch.awayScore : "-"}`;
  const matchMeta = [
    normalizedMatch.round ?? translate(locale, "matches.roundTbd"),
    formatDateTime(normalizedMatch.kickoffAt, locale),
    normalizedMatch.venue,
  ].filter(Boolean);

  return (
    <div className="space-y-8 pb-10">
      <section className="match-shell match-shell--field overflow-hidden">
        <div className="flex flex-col gap-8">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="space-y-3">
              <p className="match-eyebrow">{translate(locale, "matchDetail.heroLabel")}</p>
              <div className="flex flex-wrap items-center gap-3">
                <LeagueMark label={normalizedMatch.competition} />
                <ArchivePill>{normalizedMatch.season}</ArchivePill>
                <ArchivePill>{normalizedMatch.status}</ArchivePill>
              </div>
            </div>
            <p className="text-sm text-[var(--fw-muted)]">{translate(locale, "matchDetail.heroMeta")}</p>
          </div>

          <div className="grid gap-6 lg:grid-cols-[1fr_auto_1fr] lg:items-center">
            <div className="flex items-center gap-4 lg:justify-start">
              <TeamCrest team={normalizedMatch.homeTeam} />
              <div className="space-y-2">
                <p className="text-xs uppercase tracking-[0.18em] text-[var(--fw-muted)]">
                  {translate(locale, "enum.supporterSide.home")}
                </p>
                <Link
                  href={`/teams/${normalizedMatch.homeTeam.id}`}
                  className="text-2xl font-semibold tracking-tight text-[var(--fw-ink)] transition-colors hover:text-[var(--fw-field-900)]"
                >
                  {normalizedMatch.homeTeam.name}
                </Link>
              </div>
            </div>

            <div className="score-card rounded-[1.8rem] border border-[var(--fw-line)] px-8 py-6 text-center shadow-[0_24px_50px_rgba(16,31,24,0.08)]">
              <p className="text-xs uppercase tracking-[0.18em] text-[var(--fw-muted)]">
                {translate(locale, "matchDetail.contextEyebrow")}
              </p>
              <p className="mt-2 text-5xl font-semibold tracking-[-0.08em] text-[var(--fw-score)]">
                {scoreline}
              </p>
              <p className="mt-3 text-sm text-[var(--fw-muted)]">{matchMeta.join(" · ")}</p>
            </div>

            <div className="flex items-center gap-4 lg:justify-end">
              <div className="space-y-2 text-right">
                <p className="text-xs uppercase tracking-[0.18em] text-[var(--fw-muted)]">
                  {translate(locale, "enum.supporterSide.away")}
                </p>
                <Link
                  href={`/teams/${normalizedMatch.awayTeam.id}`}
                  className="text-2xl font-semibold tracking-tight text-[var(--fw-ink)] transition-colors hover:text-[var(--fw-field-900)]"
                >
                  {normalizedMatch.awayTeam.name}
                </Link>
              </div>
              <TeamCrest team={normalizedMatch.awayTeam} />
            </div>
          </div>
        </div>
      </section>

      <MatchCheckInPanel match={normalizedMatch} />

      <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
        <SectionShell
          eyebrow={translate(locale, "matchDetail.pulseEyebrow")}
          title={translate(locale, "matchDetail.snapshot")}
          description={translate(locale, "matchDetail.communitySubtitle")}
          accent="paper"
        >
          <div className="mt-6 grid gap-4 sm:grid-cols-2">
            <SnapshotStat
              label={translate(locale, "matchDetail.checkIns")}
              value={String(normalizedMatch.aggregates.checkInCount)}
            />
            <SnapshotStat
              label={translate(locale, "matchDetail.matchAvg")}
              value={formatNumber(normalizedMatch.aggregates.matchRatingAvg, locale)}
            />
            <SnapshotStat
              label={normalizedMatch.homeTeam.name}
              value={formatNumber(normalizedMatch.aggregates.homeTeamRatingAvg, locale)}
            />
            <SnapshotStat
              label={normalizedMatch.awayTeam.name}
              value={formatNumber(normalizedMatch.aggregates.awayTeamRatingAvg, locale)}
            />
          </div>
        </SectionShell>

        <SectionShell
          eyebrow={translate(locale, "matchDetail.pulseEyebrow")}
          title={translate(locale, "matchDetail.playerRatings")}
          description={translate(locale, "matchDetail.communitySubtitle")}
          accent="paper"
        >
          {normalizedMatch.playerRatings.length === 0 ? (
            <p className="mt-6 text-sm text-[var(--fw-muted)]">
              {translate(locale, "matchDetail.noPlayerRatings")}
            </p>
          ) : (
            <div className="mt-6 space-y-3">
              {normalizedMatch.playerRatings.map((rating) => (
                <div
                  key={rating.player.id}
                  className="flex items-center justify-between gap-4 rounded-[1.2rem] border border-[var(--fw-line)] bg-[color-mix(in_srgb,var(--fw-surface)_92%,transparent)] px-4 py-4 shadow-[0_18px_35px_rgba(16,31,24,0.05)]"
                >
                  <div className="flex items-center gap-3">
                    <TeamCrest team={rating.player.team} size="sm" />
                    <div>
                      <Link
                        href={`/players/${rating.player.id}`}
                        className="font-medium text-[var(--fw-ink)] underline-offset-4 hover:underline"
                      >
                        {rating.player.name}
                      </Link>
                      <p className="text-sm text-[var(--fw-muted)]">{rating.player.team.name}</p>
                    </div>
                  </div>
                  <div className="text-right text-sm">
                    <p className="text-lg font-semibold text-[var(--fw-field-900)]">
                      {formatNumber(rating.avgRating, locale)}
                    </p>
                    <p className="text-[var(--fw-muted)]">
                      {translate(locale, "matchDetail.ratings", { count: rating.ratingCount })}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </SectionShell>
      </div>

      <SectionShell
        eyebrow={translate(locale, "matchDetail.pulseEyebrow")}
        title={translate(locale, "matchDetail.recentReviews")}
        description={translate(locale, "matchDetail.communitySubtitle")}
        accent="paper"
      >
        {normalizedMatch.recentReviews.length === 0 ? (
          <p className="mt-6 text-sm text-[var(--fw-muted)]">
            {translate(locale, "matchDetail.noReviews")}
          </p>
        ) : (
          <div className="mt-6 grid gap-4 lg:grid-cols-2">
            {normalizedMatch.recentReviews.map((review) => (
              <article
                key={review.id}
                className="rounded-[1.35rem] border border-[var(--fw-line)] bg-[color-mix(in_srgb,var(--fw-surface)_92%,transparent)] p-5 shadow-[0_18px_35px_rgba(16,31,24,0.05)]"
              >
                <div className="flex items-center justify-between gap-4">
                  <p className="font-medium text-[var(--fw-ink)]">{review.user.name}</p>
                  <ArchivePill>
                    {translate(locale, "matchDetail.matchRating", { value: review.matchRating })}
                  </ArchivePill>
                </div>
                <p className="mt-4 text-sm leading-6 text-[var(--fw-ink-soft)]">{review.shortReview}</p>
                {review.tags.length > 0 ? (
                  <div className="mt-4 flex flex-wrap gap-2">
                    {review.tags.map((tag) => (
                      <ArchivePill key={tag.id}>{tag.name}</ArchivePill>
                    ))}
                  </div>
                ) : null}
              </article>
            ))}
          </div>
        )}
      </SectionShell>
    </div>
  );
}

function SnapshotStat({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-[1.25rem] border border-[var(--fw-line)] bg-[color-mix(in_srgb,var(--fw-surface)_92%,transparent)] p-4 shadow-[0_18px_35px_rgba(16,31,24,0.05)]">
      <p className="text-xs uppercase tracking-[0.18em] text-[var(--fw-muted)]">{label}</p>
      <p className="mt-3 text-2xl font-semibold tracking-tight text-[var(--fw-field-900)]">{value}</p>
    </div>
  );
}
