"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { ApiError, usersApi } from "@/lib/api/client";
import { useAuth } from "@/components/auth/AuthProvider";
import {
  buildHistorySummary,
  buildPaginationMeta,
  buildProfileStats,
} from "@/components/profile/profilePageUtils";
import type { UserCheckInHistoryResponse, UserProfileSummary } from "@/types/api";

export default function MePage() {
  const { status, user } = useAuth();
  const [profile, setProfile] = useState<UserProfileSummary | null>(null);
  const [history, setHistory] = useState<UserCheckInHistoryResponse | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (status !== "authenticated" || !user) {
      setProfile(null);
      setHistory(null);
      setError(null);
      return;
    }

    let cancelled = false;

    async function load() {
      setLoading(true);
      setError(null);
      try {
        const [profileResult, historyResult] = await Promise.all([
          usersApi.profile(),
          usersApi.checkins({ page, pageSize: 10 }),
        ]);
        if (cancelled) {
          return;
        }
        setProfile(profileResult);
        setHistory(historyResult);
      } catch (err) {
        if (cancelled) {
          return;
        }
        if (err instanceof ApiError) {
          setError(err.message);
        } else {
          setError("Failed to load your profile.");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void load();
    return () => {
      cancelled = true;
    };
  }, [page, status, user]);

  if (status === "loading") {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">My Profile</h1>
        <p className="mt-4 text-neutral-600">Checking your session...</p>
      </div>
    );
  }

  if (status === "unauthenticated" || !user) {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">My Profile</h1>
        <p className="mt-4 text-neutral-600">
          You need to sign in before your profile and future check-ins are available.
        </p>
        <Link
          href="/login"
          className="mt-6 inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
        >
          Go to Dev Login
        </Link>
      </div>
    );
  }

  if (loading && !profile) {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">My Profile</h1>
        <p className="mt-4 text-neutral-600">Loading your profile and check-in history...</p>
      </div>
    );
  }

  if (error && !profile) {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">My Profile</h1>
        <p className="mt-4 text-sm text-red-600">{error}</p>
      </div>
    );
  }

  if (!profile || !history) {
    return (
      <div className="py-8">
        <h1 className="text-3xl font-bold">My Profile</h1>
        <p className="mt-4 text-neutral-600">Profile data is not available yet.</p>
      </div>
    );
  }

  const stats = buildProfileStats(profile);
  const pagination = buildPaginationMeta(history.total, history.page, history.pageSize);

  return (
    <div className="py-8">
      <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">Authenticated</p>
      <h1 className="mt-2 text-3xl font-bold tracking-tight">My Profile</h1>

      <div className="mt-8 rounded-2xl border p-6">
        <p className="text-sm text-neutral-500">Signed in as</p>
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
            <h2 className="text-xl font-semibold">Check-in History</h2>
            <p className="mt-1 text-sm text-neutral-600">
              Your recent recorded matches and ratings.
            </p>
          </div>
          <p className="text-sm text-neutral-500">Total: {history.total}</p>
        </div>

        {error ? <p className="mt-4 text-sm text-red-600">{error}</p> : null}

        {history.items.length === 0 ? (
          <div className="mt-6 rounded-2xl border border-dashed p-6">
            <p className="text-sm text-neutral-600">
              You haven&apos;t recorded a match yet. Start from the matches page.
            </p>
            <Link
              href="/matches"
              className="mt-4 inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
            >
              Browse Matches
            </Link>
          </div>
        ) : (
          <div className="mt-6 space-y-4">
            {history.items.map((item) => (
              <div key={item.id} className="rounded-2xl border p-5">
                <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                  <div>
                    <p className="text-lg font-semibold">{buildHistorySummary(item)}</p>
                    <p className="mt-1 text-sm text-neutral-500">
                      {item.match.competition} · {item.match.season}
                      {item.match.round ? ` · ${item.match.round}` : ""}
                    </p>
                    <p className="mt-2 text-sm text-neutral-600">
                      Match {item.match.homeScore ?? "-"}:{item.match.awayScore ?? "-"} · Your rating {item.matchRating}/10
                    </p>
                    <p className="mt-1 text-sm text-neutral-600">
                      Watched {new Date(item.watchedAt).toLocaleString()}
                    </p>
                  </div>
                  <Link
                    href={`/matches/${item.matchId}`}
                    className="inline-flex items-center justify-center rounded-md border px-3 py-2 text-sm font-medium transition-colors hover:bg-neutral-50"
                  >
                    View Match
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
