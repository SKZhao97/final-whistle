import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";

export default async function PlayerNotFound() {
  const locale = await getServerLocale();

  return (
    <div className="py-10">
      <h1 className="text-2xl font-semibold">{translate(locale, "player.notFound.title")}</h1>
      <p className="mt-2 text-sm text-neutral-600">{translate(locale, "player.notFound.body")}</p>
    </div>
  );
}
