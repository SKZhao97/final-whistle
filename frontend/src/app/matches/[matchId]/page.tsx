import Link from "next/link";
import { notFound } from "next/navigation";

import MatchCheckInPanel from "@/components/checkin/MatchCheckInPanel";
import { ApiError, matchesApi, withLocaleHeaders } from "@/lib/api/client";
import { formatDateTime, formatNumber } from "@/lib/i18n/domain";
import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";
import type { MatchDetail } from "@/types/api";

type MatchDetailPageProps = {
  params: Promise<{ matchId: string }>;
};

async function getMatchDetail(matchId: string, locale: "en" | "zh") {
  try {
    return await matchesApi.detail<MatchDetail>(matchId, withLocaleHeaders(locale, { cache: "no-store" }));
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

  return (
    <div className="py-8">
      <div className="mb-8">
        <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">
          {normalizedMatch.competition} · {normalizedMatch.season}
        </p>
        <h1 className="mt-2 text-3xl font-bold tracking-tight">
          {normalizedMatch.homeTeam.name} {typeof normalizedMatch.homeScore === "number" ? normalizedMatch.homeScore : "-"}:
          {typeof normalizedMatch.awayScore === "number" ? normalizedMatch.awayScore : "-"} {normalizedMatch.awayTeam.name}
        </h1>
        <p className="mt-2 text-sm text-neutral-600">
          {normalizedMatch.round ?? translate(locale, "matches.roundTbd")} · {formatDateTime(normalizedMatch.kickoffAt, locale)}
          {normalizedMatch.venue ? ` · ${normalizedMatch.venue}` : ""}
        </p>
        <div className="mt-4 flex flex-wrap gap-3 text-sm">
          <Link href={`/teams/${normalizedMatch.homeTeam.id}`} className="underline">
            {normalizedMatch.homeTeam.name}
          </Link>
          <Link href={`/teams/${normalizedMatch.awayTeam.id}`} className="underline">
            {normalizedMatch.awayTeam.name}
          </Link>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <MatchCheckInPanel match={normalizedMatch} />

        <section className="rounded-xl border p-5">
          <h2 className="text-lg font-semibold">{translate(locale, "matchDetail.snapshot")}</h2>
          <dl className="mt-4 space-y-2 text-sm text-neutral-700">
            <div className="flex justify-between gap-4">
              <dt>{translate(locale, "matchDetail.checkIns")}</dt>
              <dd>{normalizedMatch.aggregates.checkInCount}</dd>
            </div>
            <div className="flex justify-between gap-4">
              <dt>{translate(locale, "matchDetail.matchAvg")}</dt>
              <dd>{formatNumber(normalizedMatch.aggregates.matchRatingAvg, locale)}</dd>
            </div>
            <div className="flex justify-between gap-4">
              <dt>{normalizedMatch.homeTeam.name}</dt>
              <dd>{formatNumber(normalizedMatch.aggregates.homeTeamRatingAvg, locale)}</dd>
            </div>
            <div className="flex justify-between gap-4">
              <dt>{normalizedMatch.awayTeam.name}</dt>
              <dd>{formatNumber(normalizedMatch.aggregates.awayTeamRatingAvg, locale)}</dd>
            </div>
          </dl>
        </section>

        <section className="rounded-xl border p-5 lg:col-span-2">
          <h2 className="text-lg font-semibold">{translate(locale, "matchDetail.playerRatings")}</h2>
          {normalizedMatch.playerRatings.length === 0 ? (
            <p className="mt-4 text-sm text-neutral-600">{translate(locale, "matchDetail.noPlayerRatings")}</p>
          ) : (
            <div className="mt-4 space-y-3">
              {normalizedMatch.playerRatings.map((rating) => (
                <div key={rating.player.id} className="flex items-center justify-between gap-4">
                  <div>
                    <Link href={`/players/${rating.player.id}`} className="font-medium underline">
                      {rating.player.name}
                    </Link>
                    <p className="text-sm text-neutral-600">{rating.player.team.name}</p>
                  </div>
                  <div className="text-right text-sm">
                    <p>{formatNumber(rating.avgRating, locale)}</p>
                    <p className="text-neutral-500">{translate(locale, "matchDetail.ratings", { count: rating.ratingCount })}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        <section className="rounded-xl border p-5 lg:col-span-3">
          <h2 className="text-lg font-semibold">{translate(locale, "matchDetail.recentReviews")}</h2>
          {normalizedMatch.recentReviews.length === 0 ? (
            <p className="mt-4 text-sm text-neutral-600">{translate(locale, "matchDetail.noReviews")}</p>
          ) : (
            <div className="mt-4 grid gap-4">
              {normalizedMatch.recentReviews.map((review) => (
                <article key={review.id} className="rounded-lg border p-4">
                  <div className="flex items-center justify-between gap-4">
                    <p className="font-medium">{review.user.name}</p>
                    <p className="text-sm text-neutral-500">
                      {translate(locale, "matchDetail.matchRating", { value: review.matchRating })}
                    </p>
                  </div>
                  <p className="mt-2 text-sm text-neutral-700">{review.shortReview}</p>
                  {review.tags.length > 0 ? (
                    <div className="mt-3 flex flex-wrap gap-2">
                      {review.tags.map((tag) => (
                        <span key={tag.id} className="rounded-full bg-neutral-100 px-2 py-1 text-xs">
                          {tag.name}
                        </span>
                      ))}
                    </div>
                  ) : null}
                </article>
              ))}
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
