"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { ApiError, usersApi } from "@/lib/api/client";
import { useAuth } from "@/components/auth/AuthProvider";
import { useLocale } from "@/components/i18n/LocaleProvider";
import {
  buildHistorySummary,
  buildPaginationMeta,
  buildProfileStats,
} from "@/components/profile/profilePageUtils";
import { formatDateTime } from "@/lib/i18n/domain";
import type { UserCheckInHistoryResponse, UserProfileSummary } from "@/types/api";

/**
 * 个人资料页面组件。
 * 显示用户个人资料摘要和签到历史。
 * 需要用户认证，未认证用户会看到登录提示。
 */
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
      <div className="py-8">
        <h1 className="text-3xl font-bold">{t("profile.title")}</h1>
        <p className="mt-4 text-neutral-600">{t("profile.checkingSession")}</p>
      </div>
    );
  }

  if (status === "unauthenticated" || !user) {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">{t("profile.title")}</h1>
        <p className="mt-4 text-neutral-600">
          {t("profile.signInRequired")}
        </p>
        <Link
          href="/login"
          className="mt-6 inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
        >
          {t("profile.goToLogin")}
        </Link>
      </div>
    );
  }

  if (profileLoading && !profile) {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">{t("profile.title")}</h1>
        <p className="mt-4 text-neutral-600">{t("profile.loading")}</p>
      </div>
    );
  }

  if (profileError && !profile) {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">{t("profile.title")}</h1>
        <p className="mt-4 text-sm text-red-600">{profileError}</p>
        <p className="mt-2 text-sm text-neutral-600">
          {t("profile.error.restart")}
        </p>
      </div>
    );
  }

  if (!profile || !history) {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">{t("profile.title")}</h1>
        <p className="mt-4 text-neutral-600">{t("profile.unavailable")}</p>
      </div>
    );
  }

  const stats = buildProfileStats(profile, locale);
  const pagination = buildPaginationMeta(history.total, history.page, history.pageSize);

  return (
    <div className="py-8">
      <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">{t("profile.authenticated")}</p>
      <h1 className="mt-2 text-3xl font-bold tracking-tight">{t("profile.title")}</h1>

      <div className="mt-8 rounded-2xl border p-6">
        <p className="text-sm text-neutral-500">{t("profile.signedInAs")}</p>
        <p className="mt-2 text-xl font-semibold">{profile.user.name}</p>
      </div>

      <div className="mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-5">
        {stats.map((stat) => (
          <div key={stat.label} className="rounded-2xl border p-5">
            <p className="text-sm text-neutral-500">{stat.label}</p>
            <p className="mt-2 text-2xl font-semibold">{stat.value}</p>
          </div>
        ))}
      </div>

      <div className="mt-10">
        <div className="flex items-center justify-between gap-4">
          <div>
            <h2 className="text-xl font-semibold">{t("profile.history.title")}</h2>
            <p className="mt-1 text-sm text-neutral-600">
              {t("profile.history.subtitle")}
            </p>
          </div>
          <p className="text-sm text-neutral-500">{t("profile.total", { total: history.total })}</p>
        </div>

        {profileError ? <p className="mt-4 text-sm text-red-600">{profileError}</p> : null}
        {historyError ? <p className="mt-2 text-sm text-red-600">{historyError}</p> : null}
        {historyLoading && history ? (
          <p className="mt-2 text-sm text-neutral-500">{t("profile.refreshingHistory")}</p>
        ) : null}

        {history.items.length === 0 ? (
          <div className="mt-6 rounded-2xl border border-dashed p-6">
            <p className="text-sm text-neutral-600">
              {t("profile.empty")}
            </p>
            <Link
              href="/matches"
              className="mt-4 inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
            >
              {t("home.browseMatches")}
            </Link>
          </div>
        ) : (
          <div className="mt-6 space-y-4">
            {history.items.map((item) => (
              <div key={item.id} className="rounded-2xl border p-5">
                <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                  <div>
                    <p className="text-lg font-semibold">{buildHistorySummary(item, locale)}</p>
                    <p className="mt-1 text-sm text-neutral-500">
                      {item.match.competition} · {item.match.season}
                      {item.match.round ? ` · ${item.match.round}` : ""}
                    </p>
                    <p className="mt-2 text-sm text-neutral-600">
                      Match {item.match.homeScore ?? "-"}:{item.match.awayScore ?? "-"} · {t("profile.yourRating", { value: item.matchRating })}
                    </p>
                    <p className="mt-1 text-sm text-neutral-600">
                      {t("profile.watchedAt", { value: formatDateTime(item.watchedAt, locale) })}
                    </p>
                  </div>
                  <Link
                    href={`/matches/${item.matchId}`}
                    className="inline-flex items-center justify-center rounded-md border px-3 py-2 text-sm font-medium transition-colors hover:bg-neutral-50"
                  >
                    {t("profile.viewMatch")}
                  </Link>
                </div>

                {item.shortReview ? (
                  <p className="mt-4 text-sm text-neutral-700">{item.shortReview}</p>
                ) : null}

                {item.tags.length > 0 ? (
                  <div className="mt-4 flex flex-wrap gap-2">
                    {item.tags.map((tag) => (
                      <span
                        key={tag.id}
                        className="rounded-full bg-neutral-100 px-3 py-1 text-xs font-medium text-neutral-700"
                      >
                        {tag.name}
                      </span>
                    ))}
                  </div>
                ) : null}
              </div>
            ))}
          </div>
        )}

        <div className="mt-6 flex items-center justify-between">
          <button
            type="button"
            onClick={() => setPage((current) => Math.max(1, current - 1))}
            disabled={!pagination.canGoPrev}
            className="inline-flex items-center justify-center rounded-md border px-4 py-2 text-sm font-medium transition-colors hover:bg-neutral-50 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Previous
          </button>
          <p className="text-sm text-neutral-500">
            Page {history.page} of {pagination.totalPages}
          </p>
          <button
            type="button"
            onClick={() => setPage((current) => current + 1)}
            disabled={!pagination.canGoNext}
            className="inline-flex items-center justify-center rounded-md border px-4 py-2 text-sm font-medium transition-colors hover:bg-neutral-50 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  );
}
