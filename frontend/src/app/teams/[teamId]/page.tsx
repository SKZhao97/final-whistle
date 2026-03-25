import Link from "next/link";
import { notFound } from "next/navigation";

import { ApiError, teamsApi } from "@/lib/api/client";
import type { TeamDetail } from "@/types/api";

type TeamDetailPageProps = {
  params: Promise<{ teamId: string }>;
};

async function getTeamDetail(teamId: string) {
  try {
    return await teamsApi.detail<TeamDetail>(teamId, { cache: "no-store" });
  } catch (error) {
    if (error instanceof ApiError && error.code === "NOT_FOUND") {
      notFound();
    }
    throw error;
  }
}

export default async function TeamDetailPage({ params }: TeamDetailPageProps) {
  const { teamId } = await params;
  const team = await getTeamDetail(teamId);

  return (
    <div className="py-8">
      <h1 className="text-3xl font-bold tracking-tight">{team.name}</h1>
      <p className="mt-2 text-sm text-neutral-600">
        {team.shortName ? `${team.shortName} · ` : ""}Average rating {team.ratingSummary.avgRating ?? "No samples"}
      </p>

      <section className="mt-8 rounded-xl border p-5">
        <h2 className="text-lg font-semibold">Recent Matches</h2>
        {team.recentMatches.length === 0 ? (
          <p className="mt-4 text-sm text-neutral-600">No related matches yet.</p>
        ) : (
          <div className="mt-4 grid gap-3">
            {team.recentMatches.map((match) => (
              <Link key={match.id} href={`/matches/${match.id}`} className="rounded-lg border p-4 hover:bg-neutral-50">
                <p className="font-medium">
                  {match.homeTeam.name} {typeof match.homeScore === "number" ? match.homeScore : "-"}:
                  {typeof match.awayScore === "number" ? match.awayScore : "-"} {match.awayTeam.name}
                </p>
                <p className="mt-1 text-sm text-neutral-600">
                  {new Date(match.kickoffAt).toLocaleString()}
                </p>
              </Link>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
