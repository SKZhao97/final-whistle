"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { useAuth } from "@/components/auth/AuthProvider";
import { ArchivePill, ArchiveStat, SectionShell } from "@/components/experience/FootballPrimitives";
import {
  buildArchiveMemory,
  buildHistorySummary,
  buildPaginationMeta,
  buildProfileStats,
} from "@/components/profile/profilePageUtils";
import { ApiError, usersApi } from "@/lib/api/client";
import { formatDateTime } from "@/lib/i18n/domain";
import { useLocale } from "@/components/i18n/LocaleProvider";
import type { UserCheckInHistoryResponse, UserProfileSummary } from "@/types/api";

export default function MePage() {
  const { status, user, refresh } = useAuth();
  const { locale, t } = useLocale();
  const [profile, setProfile] = useState<UserProfileSummary | null>(null);
  const [history, setHistory] = useState<UserCheckInHistoryResponse | null>(null);
  const [page, setPage] = useState(1);
  const [profileLoading, setProfileLoading] = useState(false);
  const [historyLoading, setHistoryLoading] = useState(false);
  const [profileError, setProfileError] = useState<string | null>(null);
  const [historyError, setHistoryError] = useState<string | null>(null);

  useEffect(() => {
    if (status !== "authenticated" || !user) {
      setProfile(null);
      setProfileError(null);
      return;
    }

    let cancelled = false;

    async function loadProfile() {
      setProfileLoading(true);
      setProfileError(null);
      try {
        const profileResult = await usersApi.profile();
        if (cancelled) {
          return;
        }
        setProfile(profileResult);
      } catch (err) {
        if (cancelled) {
          return;
        }
        if (err instanceof ApiError) {
          if (err.code === "UNAUTHORIZED") {
            await refresh();
            return;
          }
          if (err.code === "NOT_FOUND") {
            setProfileError(t("profile.backendMissing"));
          } else {
            setProfileError(err.message);
          }
        } else {
          setProfileError(t("profile.unavailable"));
        }
      } finally {
        if (!cancelled) {
          setProfileLoading(false);
        }
      }
    }

    void loadProfile();
    return () => {
      cancelled = true;
    };
  }, [refresh, status, t, user]);

  useEffect(() => {
    if (status !== "authenticated" || !user) {
      setHistory(null);
      setHistoryError(null);
      return;
    }

    let cancelled = false;

    async function loadHistory() {
      setHistoryLoading(true);
      setHistoryError(null);
      try {
        const historyResult = await usersApi.checkins({ page, pageSize: 10 });
        if (cancelled) {
          return;
        }
        setHistory(historyResult);
      } catch (err) {
        if (cancelled) {
          return;
        }
        if (err instanceof ApiError) {
          if (err.code === "UNAUTHORIZED") {
            await refresh();
            return;
          }
          setHistoryError(err.message);
        } else {
          setHistoryError(t("profile.unavailable"));
        }
      } finally {
        if (!cancelled) {
          setHistoryLoading(false);
        }
      }
    }

    void loadHistory();
    return () => {
      cancelled = true;
    };
  }, [page, refresh, status, t, user]);

  if (status === "loading") {
    return (
      <div className="space-y-6 py-8">
        <p className="match-eyebrow">{t("profile.kicker")}</p>
        <h1 className="text-4xl font-semibold tracking-tight text-[var(--fw-ink)]">{t("profile.title")}</h1>
        <p className="text-[var(--fw-muted)]">{t("profile.checkingSession")}</p>
      </div>
    );
  }

  if (status === "unauthenticated" || !user) {
    return (
      <SectionShell
        eyebrow={t("profile.kicker")}
        title={t("profile.title")}
        description={t("profile.signInRequired")}
        accent="field"
      >
        <div className="mt-6">
            <Link
              href="/login"
              className="fw-button fw-button--primary"
            >
              {t("profile.goToLogin")}
            </Link>
        </div>
      </SectionShell>
    );
  }

  if (profileLoading && !profile) {
    return (
      <SectionShell eyebrow={t("profile.kicker")} title={t("profile.title")} description={t("profile.loading")} accent="field">
        <div className="mt-2 text-sm text-[var(--fw-muted)]">{t("profile.loading")}</div>
      </SectionShell>
    );
  }

  if (profileError && !profile) {
    return (
      <SectionShell eyebrow={t("profile.kicker")} title={t("profile.title")} description={profileError} accent="field">
        <p className="mt-2 text-sm text-[var(--fw-muted)]">{t("profile.error.restart")}</p>
      </SectionShell>
    );
  }

  if (!profile || !history) {
    return (
      <SectionShell eyebrow={t("profile.kicker")} title={t("profile.title")} description={t("profile.unavailable")} accent="field">
        <div className="mt-2 text-sm text-[var(--fw-muted)]">{t("profile.unavailable")}</div>
      </SectionShell>
    );
  }

  const stats = buildProfileStats(profile, locale);
  const memoryLines = buildArchiveMemory(profile, locale);
  const pagination = buildPaginationMeta(history.total, history.page, history.pageSize);

  return (
    <div className="space-y-8 pb-10">
      <SectionShell
        eyebrow={t("profile.kicker")}
        title={t("profile.title")}
        description={t("profile.subtitle")}
        accent="field"
      >
        <div className="mt-6 grid gap-5 xl:grid-cols-[0.9fr_1.1fr]">
          <div className="rounded-[1.6rem] border border-[var(--fw-line)] bg-[color-mix(in_srgb,var(--fw-surface)_94%,transparent)] p-6 shadow-[0_18px_40px_rgba(16,31,24,0.08)]">
            <p className="match-eyebrow">{t("profile.identity.title")}</p>
            <h2 className="mt-3 text-3xl font-semibold tracking-tight text-[var(--fw-ink)]">{profile.user.name}</h2>
            <p className="mt-3 text-sm leading-6 text-[var(--fw-muted)]">{t("profile.identity.subtitle")}</p>
          </div>

          <div className="rounded-[1.6rem] border border-[var(--fw-line)] bg-[color-mix(in_srgb,var(--fw-surface)_94%,transparent)] p-6 shadow-[0_18px_40px_rgba(16,31,24,0.08)]">
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="match-eyebrow">{t("profile.memory.title")}</p>
                <h2 className="mt-3 text-2xl font-semibold tracking-tight text-[var(--fw-ink)]">
                  {t("profile.memory.subtitle")}
                </h2>
              </div>
              <ArchivePill>{t("profile.authenticated")}</ArchivePill>
            </div>
            <div className="mt-5 space-y-3 text-sm leading-6 text-[var(--fw-ink-soft)]">
              {memoryLines.map((line) => (
                <p key={line}>{line}</p>
              ))}
              {profile.checkInCount === 0 ? <p>{t("profile.memory.empty")}</p> : null}
            </div>
          </div>
        </div>
      </SectionShell>

      <SectionShell
        eyebrow={t("profile.patterns.title")}
        title={t("profile.patterns.title")}
        description={t("profile.patterns.subtitle")}
      >
        <div className="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-5">
          {stats.map((stat) => (
            <ArchiveStat key={stat.label} label={stat.label} value={stat.value} />
          ))}
        </div>
      </SectionShell>

      <SectionShell
        eyebrow={t("profile.archive.title")}
        title={t("profile.archive.title")}
        description={t("profile.archive.subtitle")}
      >
        <div className="mt-2 flex items-center justify-between gap-4">
          <ArchivePill>{t("profile.total", { total: history.total })}</ArchivePill>
          {historyLoading ? <p className="text-sm text-[var(--fw-muted)]">{t("profile.refreshingHistory")}</p> : null}
        </div>

        {profileError ? <p className="mt-4 text-sm text-red-700">{profileError}</p> : null}
        {historyError ? <p className="mt-2 text-sm text-red-700">{historyError}</p> : null}

        {history.items.length === 0 ? (
          <div className="mt-6 rounded-[1.5rem] border border-dashed border-[var(--fw-line)] bg-[color-mix(in_srgb,var(--fw-surface)_94%,transparent)] p-6">
            <p className="text-sm leading-6 text-[var(--fw-muted)]">{t("profile.empty")}</p>
            <Link
              href="/matches"
              className="fw-button fw-button--primary mt-5"
            >
              {t("home.browseMatches")}
            </Link>
          </div>
        ) : (
          <div className="mt-6 space-y-4">
            {history.items.map((item) => (
              <article
                key={item.id}
                className="rounded-[1.5rem] border border-[var(--fw-line)] bg-[color-mix(in_srgb,var(--fw-surface)_94%,transparent)] p-6 shadow-[0_18px_40px_rgba(16,31,24,0.08)]"
              >
                <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
                  <div className="space-y-3">
                    <div className="flex flex-wrap items-center gap-2">
                      <ArchivePill>{t("profile.savedRecord")}</ArchivePill>
                      {item.tags.slice(0, 3).map((tag) => (
                        <ArchivePill key={tag.id}>{tag.name}</ArchivePill>
                      ))}
                    </div>
                    <div>
                      <p className="text-2xl font-semibold tracking-tight text-[var(--fw-ink)]">
                        {buildHistorySummary(item, locale)}
                      </p>
                      <p className="mt-2 text-sm text-[var(--fw-muted)]">
                        {item.match.competition} · {item.match.season}
                        {item.match.round ? ` · ${item.match.round}` : ""}
                      </p>
                    </div>
                    <div className="space-y-1 text-sm text-[var(--fw-ink-soft)]">
                      <p>
                        {t("profile.yourRating", { value: item.matchRating })} ·{" "}
                        {item.match.homeScore ?? "-"}:{item.match.awayScore ?? "-"}
                      </p>
                      <p>{t("profile.watchedAt", { value: formatDateTime(item.watchedAt, locale) })}</p>
                    </div>
                  </div>
                  <Link
                    href={`/matches/${item.matchId}`}
                    className="inline-flex items-center justify-center rounded-full border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] px-4 py-3 text-sm font-medium text-[var(--fw-ink-soft)] transition-colors hover:bg-[var(--fw-field-100)]"
                  >
                    {t("profile.openRecord")}
                  </Link>
                </div>

                {item.shortReview ? (
                  <div className="mt-5 rounded-[1.2rem] border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] p-4 text-sm leading-6 text-[var(--fw-ink-soft)]">
                    {item.shortReview}
                  </div>
                ) : null}
              </article>
            ))}
          </div>
        )}

        <div className="mt-6 flex items-center justify-between gap-4">
          <button
            type="button"
            onClick={() => setPage((current) => Math.max(1, current - 1))}
            disabled={!pagination.canGoPrev}
            className="inline-flex items-center justify-center rounded-full border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] px-4 py-3 text-sm font-medium text-[var(--fw-ink-soft)] transition-colors hover:bg-[var(--fw-field-100)] disabled:cursor-not-allowed disabled:opacity-50"
          >
            {t("profile.pagination.previous")}
          </button>
          <p className="text-sm text-[var(--fw-muted)]">
            {t("profile.pagination.page", { page: history.page, total: pagination.totalPages })}
          </p>
          <button
            type="button"
            onClick={() => setPage((current) => current + 1)}
            disabled={!pagination.canGoNext}
            className="inline-flex items-center justify-center rounded-full border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] px-4 py-3 text-sm font-medium text-[var(--fw-ink-soft)] transition-colors hover:bg-[var(--fw-field-100)] disabled:cursor-not-allowed disabled:opacity-50"
          >
            {t("profile.pagination.next")}
          </button>
        </div>
      </SectionShell>
    </div>
  );
}
