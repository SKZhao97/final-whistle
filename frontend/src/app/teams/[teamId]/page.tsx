import Link from "next/link";
import { notFound } from "next/navigation";

import { ApiError, teamsApi, withLocaleHeaders } from "@/lib/api/client";
import { formatDateTime, formatNumber } from "@/lib/i18n/domain";
import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";
import type { TeamDetail } from "@/types/api";

type TeamDetailPageProps = {
  params: Promise<{ teamId: string }>;
};

async function getTeamDetail(teamId: string, locale: "en" | "zh") {
  try {
    return await teamsApi.detail<TeamDetail>(teamId, withLocaleHeaders(locale, { cache: "no-store" }));
  } catch (error) {
    if (error instanceof ApiError && error.code === "NOT_FOUND") {
      notFound();
    }
    throw error;
  }
}

export default async function TeamDetailPage({ params }: TeamDetailPageProps) {
  const { teamId } = await params;
  const locale = await getServerLocale();
  const team = await getTeamDetail(teamId, locale);

  return (
    <div className="py-8">
      <h1 className="text-3xl font-bold tracking-tight">{team.name}</h1>
      <p className="mt-2 text-sm text-neutral-600">
        {team.shortName ? `${team.shortName} · ` : ""}{translate(locale, "team.detail.averageRating", { value: formatNumber(team.ratingSummary.avgRating, locale) })}
      </p>

      <section className="mt-8 rounded-xl border p-5">
        <h2 className="text-lg font-semibold">{translate(locale, "team.detail.recentMatches")}</h2>
        {team.recentMatches.length === 0 ? (
          <p className="mt-4 text-sm text-neutral-600">{translate(locale, "team.detail.noMatches")}</p>
        ) : (
          <div className="mt-4 grid gap-3">
            {team.recentMatches.map((match) => (
              <Link key={match.id} href={`/matches/${match.id}`} className="rounded-lg border p-4 hover:bg-neutral-50">
                <p className="font-medium">
                  {match.homeTeam.name} {typeof match.homeScore === "number" ? match.homeScore : "-"}:
                  {typeof match.awayScore === "number" ? match.awayScore : "-"} {match.awayTeam.name}
                </p>
                <p className="mt-1 text-sm text-neutral-600">
                  {formatDateTime(match.kickoffAt, locale)}
                </p>
              </Link>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
