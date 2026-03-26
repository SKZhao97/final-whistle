import Link from "next/link";
import { notFound } from "next/navigation";

import MatchCheckInPanel from "@/components/checkin/MatchCheckInPanel";
import { ApiError, matchesApi } from "@/lib/api/client";
import type { MatchDetail } from "@/types/api";

type MatchDetailPageProps = {
  params: Promise<{ matchId: string }>;
};

async function getMatchDetail(matchId: string) {
  try {
    return await matchesApi.detail<MatchDetail>(matchId, { cache: "no-store" });
  } catch (error) {
    if (error instanceof ApiError && error.code === "NOT_FOUND") {
      notFound();
    }
    throw error;
  }
}

export default async function MatchDetailPage({ params }: MatchDetailPageProps) {
  const { matchId } = await params;
  const match = await getMatchDetail(matchId);

  return (
    <div className="py-8">
      <div className="mb-8">
        <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">
          {match.competition} · {match.season}
        </p>
        <h1 className="mt-2 text-3xl font-bold tracking-tight">
          {match.homeTeam.name} {typeof match.homeScore === "number" ? match.homeScore : "-"}:
          {typeof match.awayScore === "number" ? match.awayScore : "-"} {match.awayTeam.name}
        </h1>
        <p className="mt-2 text-sm text-neutral-600">
          {match.round ?? "Round TBD"} · {new Date(match.kickoffAt).toLocaleString()}
          {match.venue ? ` · ${match.venue}` : ""}
        </p>
        <div className="mt-4 flex flex-wrap gap-3 text-sm">
          <Link href={`/teams/${match.homeTeam.id}`} className="underline">
            {match.homeTeam.name}
          </Link>
          <Link href={`/teams/${match.awayTeam.id}`} className="underline">
            {match.awayTeam.name}
          </Link>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <MatchCheckInPanel match={match} />

        <section className="rounded-xl border p-5">
          <h2 className="text-lg font-semibold">Community Snapshot</h2>
          <dl className="mt-4 space-y-2 text-sm text-neutral-700">
            <div className="flex justify-between gap-4">
              <dt>Check-ins</dt>
              <dd>{match.aggregates.checkInCount}</dd>
            </div>
            <div className="flex justify-between gap-4">
              <dt>Match avg</dt>
              <dd>{match.aggregates.matchRatingAvg ?? "No samples"}</dd>
            </div>
            <div className="flex justify-between gap-4">
              <dt>{match.homeTeam.name}</dt>
              <dd>{match.aggregates.homeTeamRatingAvg ?? "No samples"}</dd>
            </div>
            <div className="flex justify-between gap-4">
              <dt>{match.awayTeam.name}</dt>
              <dd>{match.aggregates.awayTeamRatingAvg ?? "No samples"}</dd>
            </div>
          </dl>
        </section>

        <section className="rounded-xl border p-5 lg:col-span-2">
          <h2 className="text-lg font-semibold">Player Ratings</h2>
          {match.playerRatings.length === 0 ? (
            <p className="mt-4 text-sm text-neutral-600">No player ratings yet.</p>
          ) : (
            <div className="mt-4 space-y-3">
              {match.playerRatings.map((rating) => (
                <div key={rating.player.id} className="flex items-center justify-between gap-4">
                  <div>
                    <Link href={`/players/${rating.player.id}`} className="font-medium underline">
                      {rating.player.name}
                    </Link>
                    <p className="text-sm text-neutral-600">{rating.player.team.name}</p>
                  </div>
                  <div className="text-right text-sm">
                    <p>{rating.avgRating ?? "No avg"}</p>
                    <p className="text-neutral-500">{rating.ratingCount} ratings</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        <section className="rounded-xl border p-5 lg:col-span-3">
          <h2 className="text-lg font-semibold">Recent Reviews</h2>
          {match.recentReviews.length === 0 ? (
            <p className="mt-4 text-sm text-neutral-600">No reviews yet.</p>
          ) : (
            <div className="mt-4 grid gap-4">
              {match.recentReviews.map((review) => (
                <article key={review.id} className="rounded-lg border p-4">
                  <div className="flex items-center justify-between gap-4">
                    <p className="font-medium">{review.user.name}</p>
                    <p className="text-sm text-neutral-500">
                      Match rating {review.matchRating}
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
