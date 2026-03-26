"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

import { useAuth } from "@/components/auth/AuthProvider";
import {
  buildPayload,
  createDefaultFormState,
  createFormStateFromCheckIn,
  DEFAULT_PLAYER_RATING,
  type CheckInFormErrors,
  type CheckInFormState,
  validateFormState,
} from "@/components/checkin/checkinFormUtils";
import { ApiError, matchesApi } from "@/lib/api/client";
import type { CheckInDetail, MatchDetail } from "@/types/api";

type MatchCheckInPanelProps = {
  match: MatchDetail;
};

export default function MatchCheckInPanel({ match }: MatchCheckInPanelProps) {
  const router = useRouter();
  const { status } = useAuth();
  const [myCheckIn, setMyCheckIn] = useState<CheckInDetail | null | undefined>(undefined);
  const [loadingRecord, setLoadingRecord] = useState(false);
  const [recordError, setRecordError] = useState<string | null>(null);
  const [editing, setEditing] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [formState, setFormState] = useState<CheckInFormState>(() => createDefaultFormState());
  const [formErrors, setFormErrors] = useState<CheckInFormErrors>({});

  const isFinished = match.status === "FINISHED";
  const roster = match.matchPlayers;
  const tagOptions = match.availableTags;

  useEffect(() => {
    let cancelled = false;

    async function loadMyCheckIn() {
      if (status !== "authenticated") {
        setMyCheckIn(undefined);
        setEditing(false);
        setRecordError(null);
        setLoadingRecord(false);
        return;
      }

      setLoadingRecord(true);
      setRecordError(null);

      try {
        const result = await matchesApi.myCheckIn(match.id, { cache: "no-store" });
        if (cancelled) {
          return;
        }
        setMyCheckIn(result);
        setFormState(result ? createFormStateFromCheckIn(result) : createDefaultFormState());
      } catch (error) {
        if (cancelled) {
          return;
        }
        if (error instanceof ApiError) {
          setRecordError(error.message);
        } else {
          setRecordError("Failed to load your match record.");
        }
      } finally {
        if (!cancelled) {
          setLoadingRecord(false);
        }
      }
    }

    void loadMyCheckIn();
    return () => {
      cancelled = true;
    };
  }, [match.id, status]);

  const availablePlayers = useMemo(
    () => roster.map((player) => ({ ...player, label: `${player.name} · ${player.team.name}` })),
    [roster],
  );

  function openCreate() {
    setFormState(createDefaultFormState());
    setFormErrors({});
    setSubmitError(null);
    setEditing(true);
  }

  function openEdit() {
    if (!myCheckIn) {
      return;
    }
    setFormState(createFormStateFromCheckIn(myCheckIn));
    setFormErrors({});
    setSubmitError(null);
    setEditing(true);
  }

  function cancelEdit() {
    setEditing(false);
    setSubmitError(null);
    setFormErrors({});
    setFormState(myCheckIn ? createFormStateFromCheckIn(myCheckIn) : createDefaultFormState());
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const validation = validateFormState(formState, roster);
    setFormErrors(validation);
    setSubmitError(null);
    if (Object.keys(validation).length > 0) {
      return;
    }

    const payload = buildPayload(formState);
    setSubmitting(true);

    try {
      const result = myCheckIn
        ? await matchesApi.updateCheckIn(match.id, payload)
        : await matchesApi.createCheckIn(match.id, payload);
      setMyCheckIn(result);
      setFormState(createFormStateFromCheckIn(result));
      setEditing(false);
      router.refresh();
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      } else {
        setSubmitError("Failed to save your match record.");
      }
    } finally {
      setSubmitting(false);
    }
  }

  if (status === "loading") {
    return (
      <section className="rounded-xl border p-5 lg:col-span-3">
        <h2 className="text-lg font-semibold">My Match Record</h2>
        <p className="mt-3 text-sm text-neutral-600">Checking your session and match record...</p>
      </section>
    );
  }

  if (status !== "authenticated") {
    return (
      <section className="rounded-xl border p-5 lg:col-span-3">
        <h2 className="text-lg font-semibold">My Match Record</h2>
        <p className="mt-3 text-sm text-neutral-600">
          Sign in to record your reaction, ratings, and tags for this match.
        </p>
        <Link
          href="/login"
          className="mt-4 inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
        >
          Go to Dev Login
        </Link>
      </section>
    );
  }

  if (!isFinished) {
    return (
      <section className="rounded-xl border p-5 lg:col-span-3">
        <h2 className="text-lg font-semibold">My Match Record</h2>
        <p className="mt-3 text-sm text-neutral-600">
          Check-ins open after the match is finished. You can come back here once the final whistle
          blows.
        </p>
      </section>
    );
  }

  return (
    <section className="rounded-xl border p-5 lg:col-span-3">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h2 className="text-lg font-semibold">My Match Record</h2>
          <p className="mt-2 text-sm text-neutral-600">
            Save your own ratings, tags, and player notes for this match.
          </p>
        </div>
        {!editing && !loadingRecord ? (
          myCheckIn ? (
            <button
              type="button"
              onClick={openEdit}
              className="inline-flex items-center justify-center rounded-md border px-4 py-2 text-sm font-medium transition-colors hover:bg-neutral-50"
            >
              Edit My Record
            </button>
          ) : (
            <button
              type="button"
              onClick={openCreate}
              className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
            >
              Record This Match
            </button>
          )
        ) : null}
      </div>

      {loadingRecord ? <p className="mt-4 text-sm text-neutral-600">Loading your record...</p> : null}
      {recordError ? <p className="mt-4 text-sm text-red-600">{recordError}</p> : null}

      {!editing && myCheckIn ? <SavedCheckInSummary checkIn={myCheckIn} /> : null}

      {!editing && !loadingRecord && !myCheckIn ? (
        <p className="mt-5 text-sm text-neutral-600">
          You have not recorded this match yet. Use the button above to create your first entry.
        </p>
      ) : null}

      {editing ? (
        <form onSubmit={handleSubmit} className="mt-6 space-y-6 rounded-xl border border-neutral-200 p-5">
          <div className="grid gap-4 md:grid-cols-2">
            <label className="block text-sm">
              <span className="mb-2 block font-medium">Watched Type</span>
              <select
                value={formState.watchedType}
                onChange={(event) => setFormState((current) => ({ ...current, watchedType: event.target.value as CheckInFormState["watchedType"] }))}
                className="w-full rounded-md border px-3 py-2"
              >
                <option value="FULL">Full Match</option>
                <option value="PARTIAL">Partial</option>
                <option value="HIGHLIGHTS">Highlights</option>
              </select>
              {formErrors.watchedType ? <span className="mt-1 block text-red-600">{formErrors.watchedType}</span> : null}
            </label>

            <label className="block text-sm">
              <span className="mb-2 block font-medium">Supporter Side</span>
              <select
                value={formState.supporterSide}
                onChange={(event) => setFormState((current) => ({ ...current, supporterSide: event.target.value as CheckInFormState["supporterSide"] }))}
                className="w-full rounded-md border px-3 py-2"
              >
                <option value="HOME">{match.homeTeam.name}</option>
                <option value="AWAY">{match.awayTeam.name}</option>
                <option value="NEUTRAL">Neutral</option>
              </select>
              {formErrors.supporterSide ? <span className="mt-1 block text-red-600">{formErrors.supporterSide}</span> : null}
            </label>
          </div>

          <div className="grid gap-4 md:grid-cols-3">
            <RatingField
              label="Match Rating"
              value={formState.matchRating}
              error={formErrors.matchRating}
              onChange={(value) => setFormState((current) => ({ ...current, matchRating: value }))}
            />
            <RatingField
              label={`${match.homeTeam.name} Rating`}
              value={formState.homeTeamRating}
              error={formErrors.homeTeamRating}
              onChange={(value) => setFormState((current) => ({ ...current, homeTeamRating: value }))}
            />
            <RatingField
              label={`${match.awayTeam.name} Rating`}
              value={formState.awayTeamRating}
              error={formErrors.awayTeamRating}
              onChange={(value) => setFormState((current) => ({ ...current, awayTeamRating: value }))}
            />
          </div>

          <label className="block text-sm">
            <span className="mb-2 block font-medium">Watched At</span>
            <input
              type="datetime-local"
              value={formState.watchedAt}
              onChange={(event) => setFormState((current) => ({ ...current, watchedAt: event.target.value }))}
              className="w-full rounded-md border px-3 py-2"
            />
            {formErrors.watchedAt ? <span className="mt-1 block text-red-600">{formErrors.watchedAt}</span> : null}
          </label>

          <label className="block text-sm">
            <span className="mb-2 block font-medium">Short Review</span>
            <textarea
              value={formState.shortReview}
              onChange={(event) => setFormState((current) => ({ ...current, shortReview: event.target.value }))}
              rows={4}
              maxLength={280}
              className="w-full rounded-md border px-3 py-2"
              placeholder="What stood out after the final whistle?"
            />
            <span className="mt-1 block text-xs text-neutral-500">{formState.shortReview.length}/280</span>
            {formErrors.shortReview ? <span className="mt-1 block text-red-600">{formErrors.shortReview}</span> : null}
          </label>

          <fieldset className="space-y-3">
            <legend className="text-sm font-medium">Tags</legend>
            <div className="flex flex-wrap gap-2">
              {tagOptions.map((tag) => {
                const selected = formState.tags.includes(tag.id);
                return (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() =>
                      setFormState((current) => ({
                        ...current,
                        tags: selected
                          ? current.tags.filter((value) => value !== tag.id)
                          : [...current.tags, tag.id],
                      }))
                    }
                    className={`rounded-full border px-3 py-1 text-sm transition-colors ${
                      selected ? "border-primary bg-primary text-primary-foreground" : "border-neutral-300 hover:bg-neutral-50"
                    }`}
                  >
                    {tag.name}
                  </button>
                );
              })}
            </div>
          </fieldset>

          <div className="space-y-3">
            <div className="flex items-center justify-between gap-4">
              <h3 className="text-sm font-medium">Player Ratings</h3>
              <button
                type="button"
                onClick={() =>
                  setFormState((current) => ({
                    ...current,
                    playerRatings: [...current.playerRatings, { ...DEFAULT_PLAYER_RATING }],
                  }))
                }
                className="text-sm font-medium text-primary underline-offset-4 hover:underline"
              >
                Add Player Rating
              </button>
            </div>

            {formState.playerRatings.length === 0 ? (
              <p className="text-sm text-neutral-600">
                Rate as many players from this match roster as you want.
              </p>
            ) : (
              <div className="space-y-4">
                {formState.playerRatings.map((playerRating, index) => (
                  <div key={`${index}-${playerRating.playerId}`} className="rounded-lg border border-neutral-200 p-4">
                    <div className="grid gap-4 md:grid-cols-[1.4fr_0.5fr]">
                      <label className="block text-sm">
                        <span className="mb-2 block font-medium">Player</span>
                        <select
                          value={playerRating.playerId}
                          onChange={(event) =>
                            setFormState((current) => ({
                              ...current,
                              playerRatings: current.playerRatings.map((entry, entryIndex) =>
                                entryIndex === index ? { ...entry, playerId: event.target.value } : entry,
                              ),
                            }))
                          }
                          className="w-full rounded-md border px-3 py-2"
                        >
                          <option value="">Select a player</option>
                          {availablePlayers
                            .filter((player) => {
                              const selectedElsewhere = formState.playerRatings.some(
                                (entry, entryIndex) =>
                                  entryIndex !== index && entry.playerId === String(player.id),
                              );
                              return !selectedElsewhere || player.id === Number(playerRating.playerId);
                            })
                            .map((player) => (
                            <option key={player.id} value={String(player.id)}>
                              {player.label}
                            </option>
                            ))}
                        </select>
                      </label>

                      <RatingField
                        label="Rating"
                        value={playerRating.rating}
                        onChange={(value) =>
                          setFormState((current) => ({
                            ...current,
                            playerRatings: current.playerRatings.map((entry, entryIndex) =>
                              entryIndex === index ? { ...entry, rating: value } : entry,
                            ),
                          }))
                        }
                      />
                    </div>

                    <label className="mt-4 block text-sm">
                      <span className="mb-2 block font-medium">Note</span>
                      <input
                        type="text"
                        value={playerRating.note}
                        maxLength={80}
                        onChange={(event) =>
                          setFormState((current) => ({
                            ...current,
                            playerRatings: current.playerRatings.map((entry, entryIndex) =>
                              entryIndex === index ? { ...entry, note: event.target.value } : entry,
                            ),
                          }))
                        }
                        className="w-full rounded-md border px-3 py-2"
                        placeholder="Optional player note"
                      />
                    </label>

                    <button
                      type="button"
                      onClick={() =>
                        setFormState((current) => ({
                          ...current,
                          playerRatings: current.playerRatings.filter((_, entryIndex) => entryIndex !== index),
                        }))
                      }
                      className="mt-4 text-sm text-red-600 underline-offset-4 hover:underline"
                    >
                      Remove Player Rating
                    </button>
                  </div>
                ))}
              </div>
            )}
            {formErrors.playerRatings ? <p className="text-sm text-red-600">{formErrors.playerRatings}</p> : null}
          </div>

          {submitError ? <p className="text-sm text-red-600">{submitError}</p> : null}

          <div className="flex flex-wrap gap-3">
            <button
              type="submit"
              disabled={submitting}
              className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50"
            >
              {submitting ? "Saving..." : myCheckIn ? "Update Record" : "Create Record"}
            </button>
            <button
              type="button"
              onClick={cancelEdit}
              className="inline-flex items-center justify-center rounded-md border px-4 py-2 text-sm font-medium transition-colors hover:bg-neutral-50"
            >
              Cancel
            </button>
          </div>
        </form>
      ) : null}
    </section>
  );
}

function SavedCheckInSummary({ checkIn }: { checkIn: CheckInDetail }) {
  return (
    <div className="mt-6 grid gap-5 lg:grid-cols-[1fr_1.2fr]">
      <div className="rounded-lg border border-neutral-200 p-4">
        <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">Saved Record</p>
        <dl className="mt-4 space-y-2 text-sm">
          <div className="flex justify-between gap-4">
            <dt>Watched Type</dt>
            <dd>{checkIn.watchedType}</dd>
          </div>
          <div className="flex justify-between gap-4">
            <dt>Supporter Side</dt>
            <dd>{checkIn.supporterSide}</dd>
          </div>
          <div className="flex justify-between gap-4">
            <dt>Match Rating</dt>
            <dd>{checkIn.matchRating}</dd>
          </div>
          <div className="flex justify-between gap-4">
            <dt>Saved At</dt>
            <dd>{new Date(checkIn.updatedAt).toLocaleString()}</dd>
          </div>
        </dl>
        {checkIn.tags.length > 0 ? (
          <div className="mt-4 flex flex-wrap gap-2">
            {checkIn.tags.map((tag) => (
              <span key={tag.id} className="rounded-full bg-neutral-100 px-2 py-1 text-xs">
                {tag.name}
              </span>
            ))}
          </div>
        ) : null}
      </div>

      <div className="rounded-lg border border-neutral-200 p-4">
        <h3 className="text-sm font-medium">Player Notes</h3>
        {checkIn.playerRatings.length === 0 ? (
          <p className="mt-3 text-sm text-neutral-600">No player ratings saved yet.</p>
        ) : (
          <div className="mt-3 space-y-3">
            {checkIn.playerRatings.map((rating) => (
              <div key={rating.id} className="rounded-md border border-neutral-200 p-3">
                <div className="flex items-center justify-between gap-4">
                  <div>
                    <p className="font-medium">{rating.player.name}</p>
                    <p className="text-sm text-neutral-500">{rating.player.team.name}</p>
                  </div>
                  <p className="text-sm">Rating {rating.rating}</p>
                </div>
                {rating.note ? <p className="mt-2 text-sm text-neutral-700">{rating.note}</p> : null}
              </div>
            ))}
          </div>
        )}
        {checkIn.shortReview ? (
          <div className="mt-4 rounded-md bg-neutral-50 p-3 text-sm text-neutral-700">
            {checkIn.shortReview}
          </div>
        ) : null}
      </div>
    </div>
  );
}

function RatingField({
  label,
  value,
  onChange,
  error,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}) {
  return (
    <label className="block text-sm">
      <span className="mb-2 block font-medium">{label}</span>
      <input
        type="number"
        min={1}
        max={10}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        className="w-full rounded-md border px-3 py-2"
      />
      {error ? <span className="mt-1 block text-red-600">{error}</span> : null}
    </label>
  );
}
