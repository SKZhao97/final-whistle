import Link from "next/link";
import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";

export default async function Home() {
  const locale = await getServerLocale();

  return (
    <div className="flex flex-col items-center justify-center py-12">
      <div className="max-w-2xl mx-auto text-center space-y-8">
        <h1 className="text-4xl font-bold tracking-tight text-gray-900 dark:text-gray-100 sm:text-5xl">
          {translate(locale, "home.title")}
        </h1>
        <p className="text-lg text-gray-600 dark:text-gray-400">
          {translate(locale, "home.subtitle")}
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link
            href="/matches"
            className="inline-flex items-center justify-center rounded-md bg-primary px-6 py-3 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            {translate(locale, "home.browseMatches")}
          </Link>
          <Link
            href="/me"
            className="inline-flex items-center justify-center rounded-md border border-input bg-background px-6 py-3 text-sm font-medium shadow-sm transition-colors hover:bg-accent hover:text-accent-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            {translate(locale, "home.myProfile")}
          </Link>
        </div>
        <div className="pt-8 grid grid-cols-1 sm:grid-cols-2 gap-6 text-left">
          <div className="space-y-3">
            <h3 className="text-lg font-semibold">{translate(locale, "home.feature.record.title")}</h3>
            <p className="text-sm text-gray-500">
              {translate(locale, "home.feature.record.body")}
            </p>
          </div>
          <div className="space-y-3">
            <h3 className="text-lg font-semibold">{translate(locale, "home.feature.review.title")}</h3>
            <p className="text-sm text-gray-500">
              {translate(locale, "home.feature.review.body")}
            </p>
          </div>
          <div className="space-y-3">
            <h3 className="text-lg font-semibold">{translate(locale, "home.feature.archive.title")}</h3>
            <p className="text-sm text-gray-500">
              {translate(locale, "home.feature.archive.body")}
            </p>
          </div>
          <div className="space-y-3">
            <h3 className="text-lg font-semibold">{translate(locale, "home.feature.community.title")}</h3>
            <p className="text-sm text-gray-500">
              {translate(locale, "home.feature.community.body")}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
