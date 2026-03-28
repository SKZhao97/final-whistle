import Link from "next/link";

import { ArchivePill, LeagueMark, SectionShell, TeamCrest } from "@/components/experience/FootballPrimitives";
import { translate } from "@/lib/i18n/core";
import { getServerLocale } from "@/lib/i18n/server";

export default async function Home() {
  const locale = await getServerLocale();

  const demoTeams = [
    { id: 3, name: locale === "zh" ? "阿森纳" : "Arsenal", shortName: "ARS", slug: "arsenal" },
    { id: 1, name: locale === "zh" ? "曼城" : "Manchester City", shortName: "MCI", slug: "manchester-city" },
  ];

  return (
    <div className="space-y-8 pb-10">
      <section className="match-shell match-shell--field hero-pitch overflow-hidden">
        <div className="grid gap-8 lg:grid-cols-[1.2fr_0.8fr] lg:items-center">
          <div className="space-y-6">
            <div className="flex flex-wrap items-center gap-3">
              <LeagueMark label={locale === "zh" ? "英超赛后记录" : "Post-match football archive"} />
              <ArchivePill>{locale === "zh" ? "A-first 体验" : "A-first experience"}</ArchivePill>
            </div>
            <div className="space-y-4">
              <h1 className="max-w-4xl text-4xl font-semibold tracking-tight text-[var(--fw-ink)] sm:text-5xl lg:text-6xl">
                {locale === "zh" ? (
                  <>
                    欢迎来到 <span className="text-[var(--fw-field-700)]">Final Whistle</span>
                  </>
                ) : (
                  <>
                    Welcome to <span className="text-[var(--fw-field-700)]">Final Whistle</span>
                  </>
                )}
              </h1>
              <p className="max-w-3xl text-lg leading-8 text-[var(--fw-muted)]">
                {translate(locale, "home.subtitle")}
              </p>
            </div>
            <div className="flex flex-col gap-3 sm:flex-row">
              <Link href="/matches" className="fw-button fw-button--primary">
                {translate(locale, "home.browseMatches")}
              </Link>
              <Link href="/me" className="fw-button fw-button--secondary">
                {translate(locale, "home.myProfile")}
              </Link>
            </div>
          </div>

          <div className="rounded-[1.8rem] border border-[var(--fw-line)] bg-[color:var(--fw-surface)]/86 p-6 shadow-[0_24px_55px_rgba(16,31,24,0.08)]">
            <p className="match-eyebrow">{translate(locale, "home.demoCard.eyebrow")}</p>
            <div className="mt-5 grid gap-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <TeamCrest team={demoTeams[0]} size="sm" />
                  <span className="font-medium text-[var(--fw-ink)]">{demoTeams[0].name}</span>
                </div>
                <span className="text-sm text-[var(--fw-muted)]">2</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <TeamCrest team={demoTeams[1]} size="sm" />
                  <span className="font-medium text-[var(--fw-ink)]">{demoTeams[1].name}</span>
                </div>
                <span className="text-sm text-[var(--fw-muted)]">2</span>
              </div>
              <div className="rounded-[1.2rem] border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] p-4">
                <p className="text-sm leading-6 text-[var(--fw-ink-soft)]">
                  {translate(locale, "home.demoCard.body")}
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <div className="grid gap-6 lg:grid-cols-2">
        <SectionShell
          eyebrow={locale === "zh" ? "记录价值" : "Recording Value"}
          title={translate(locale, "home.feature.record.title")}
          description={translate(locale, "home.feature.record.body")}
        >
          <div className="mt-6 space-y-3 text-sm leading-6 text-[var(--fw-ink-soft)]">
            <p>{translate(locale, "home.feature.review.body")}</p>
            <p>{translate(locale, "home.feature.archive.body")}</p>
          </div>
        </SectionShell>

        <SectionShell
          eyebrow={locale === "zh" ? "社区陪衬" : "Community Layer"}
          title={translate(locale, "home.feature.community.title")}
          description={translate(locale, "home.feature.community.body")}
          accent="field"
        >
          <div className="mt-6 grid gap-3 sm:grid-cols-2">
            <ArchivePill>{locale === "zh" ? "聚合评分" : "Aggregate ratings"}</ArchivePill>
            <ArchivePill>{locale === "zh" ? "热门标签" : "Hot tags"}</ArchivePill>
            <ArchivePill>{locale === "zh" ? "球员评分榜" : "Player board"}</ArchivePill>
            <ArchivePill>{locale === "zh" ? "近期短评" : "Recent reactions"}</ArchivePill>
          </div>
        </SectionShell>
      </div>
    </div>
  );
}
