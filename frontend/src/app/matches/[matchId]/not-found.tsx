import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";

export default async function MatchNotFound() {
  const locale = await getServerLocale();

  return (
    <div className="py-10">
      <h1 className="text-2xl font-semibold">{translate(locale, "matchDetail.notFound.title")}</h1>
      <p className="mt-2 text-sm text-neutral-600">
        {translate(locale, "matchDetail.notFound.body")}
      </p>
    </div>
  );
}
