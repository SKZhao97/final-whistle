import Link from "next/link";
import { notFound } from "next/navigation";

import { ApiError, playersApi } from "@/lib/api/client";
import type { PlayerDetail } from "@/types/api";

type PlayerDetailPageProps = {
  params: Promise<{ playerId: string }>;
};

async function getPlayerDetail(playerId: string) {
  try {
    return await playersApi.detail<PlayerDetail>(playerId, { cache: "no-store" });
  } catch (error) {
    if (error instanceof ApiError && error.code === "NOT_FOUND") {
      notFound();
    }
    throw error;
  }
}

export default async function PlayerDetailPage({ params }: PlayerDetailPageProps) {
  const { playerId } = await params;
  const player = await getPlayerDetail(playerId);

  return (
    <div className="py-8">
      <h1 className="text-3xl font-bold tracking-tight">{player.name}</h1>
      <p className="mt-2 text-sm text-neutral-600">
        {player.position ?? "Position TBD"} · <Link href={`/teams/${player.team.id}`} className="underline">{player.team.name}</Link>
      </p>
      <p className="mt-2 text-sm text-neutral-600">
        Average rating {player.ratingSummary.avgRating ?? "No samples"} from {player.ratingSummary.ratingCount} ratings
      </p>

      <section className="mt-8 rounded-xl border p-5">
        <h2 className="text-lg font-semibold">Recent Rated Matches</h2>
        {player.recentMatches.length === 0 ? (
          <p className="mt-4 text-sm text-neutral-600">No recent rated matches yet.</p>
        ) : (
          <div className="mt-4 grid gap-3">
            {player.recentMatches.map((item) => (
              <Link key={item.match.id} href={`/matches/${item.match.id}`} className="rounded-lg border p-4 hover:bg-neutral-50">
                <p className="font-medium">
                  {item.match.homeTeam.name} {typeof item.match.homeScore === "number" ? item.match.homeScore : "-"}:
                  {typeof item.match.awayScore === "number" ? item.match.awayScore : "-"} {item.match.awayTeam.name}
                </p>
                <p className="mt-1 text-sm text-neutral-600">
                  Avg player rating {item.avgRating ?? "No avg"} · {item.ratingCount} ratings
                </p>
              </Link>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
