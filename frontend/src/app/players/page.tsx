import Link from "next/link";
import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";

export default async function PlayersPage() {
  const locale = await getServerLocale();

  return (
    <div className="py-8">
      <h1 className="text-3xl font-bold tracking-tight text-[var(--fw-ink)]">{translate(locale, "players.title")}</h1>
      <p className="mt-3 text-sm text-[var(--fw-muted)]">
        {translate(locale, "players.body")}
      </p>
      <Link
        href="/matches"
        className="fw-button fw-button--secondary mt-6"
      >
        {translate(locale, "home.browseMatches")}
      </Link>
    </div>
  );
}
