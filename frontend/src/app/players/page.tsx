import Link from "next/link";
import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";

export default async function PlayersPage() {
  const locale = await getServerLocale();

  return (
    <div className="py-8">
      <h1 className="text-3xl font-bold tracking-tight">{translate(locale, "players.title")}</h1>
      <p className="mt-3 text-sm text-neutral-600">
        {translate(locale, "players.body")}
      </p>
      <Link
        href="/matches"
        className="mt-6 inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
      >
        {translate(locale, "home.browseMatches")}
      </Link>
    </div>
  );
}
