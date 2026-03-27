"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { type ReactNode, useEffect, useMemo, useState } from "react";

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
import {
  ArchivePill,
  SectionShell,
  TeamCrest,
} from "@/components/experience/FootballPrimitives";
import { useLocale } from "@/components/i18n/LocaleProvider";
import { ApiError, matchesApi } from "@/lib/api/client";
import { translate } from "@/lib/i18n/core";
import { formatDateTime, getSupporterSideLabel, getWatchedTypeLabel } from "@/lib/i18n/domain";
import type { CheckInDetail, MatchDetail } from "@/types/api";

type MatchCheckInPanelProps = {
  match: MatchDetail;
};

export default function MatchCheckInPanel({ match }: MatchCheckInPanelProps) {
  const router = useRouter();
  const { status, refresh } = useAuth();
  const { locale, t } = useLocale();
  const [myCheckIn, setMyCheckIn] = useState<CheckInDetail | null | undefined>(undefined);
  const [loadingRecord, setLoadingRecord] = useState(false);
  const [recordError, setRecordError] = useState<string | null>(null);
  const [editing, setEditing] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [formState, setFormState] = useState<CheckInFormState>(() => createDefaultFormState());
  const [formErrors, setFormErrors] = useState<CheckInFormErrors>({});

  const isFinished = match.status === "FINISHED";
  const roster = useMemo(() => match.matchPlayers ?? [], [match.matchPlayers]);
  const tagOptions = useMemo(() => match.availableTags ?? [], [match.availableTags]);

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
          if (error.code === "NOT_FOUND") {
            setRecordError(t("checkin.backendMissing"));
          } else if (error.code === "UNAUTHORIZED") {
            setMyCheckIn(undefined);
            setRecordError(null);
            void refresh();
          } else {
            setRecordError(error.message);
          }
        } else {
          setRecordError(t("checkin.failedLoad"));
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
  }, [match.id, refresh, status, t]);

  const availablePlayers = useMemo(
    () => roster.map((player) => ({ ...player, label: `${player.name} · ${player.team.name}` })),
    [roster],
  );

  function toggleTag(tagId: number) {
    setFormState((current) => {
      const isSelected = current.tags.includes(tagId);
      return {
        ...current,
        tags: isSelected
          ? current.tags.filter((value) => value !== tagId)
          : [...current.tags, tagId],
      };
    });
  }

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
    const validation = validateFormState(formState, roster, locale);
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
        setSubmitError(t("checkin.failedSave"));
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <SectionShell
      eyebrow={t("matchDetail.recordEyebrow")}
      title={t("checkin.title")}
      description={t("checkin.subtitle")}
      accent="field"
      className="space-y-6"
    >
      {status === "loading" ? (
        <StateMessage
          eyebrow={t("checkin.editingEyebrow")}
          title={t("checkin.loadingSession")}
          body={t("checkin.entryBody")}
        />
      ) : null}

      {status !== "loading" && status !== "authenticated" ? (
        <StateCard
          eyebrow={t("checkin.emptyEyebrow")}
          title={t("checkin.entryTitle")}
          body={t("checkin.signInPrompt")}
          cta={
            <Link href="/login" className={primaryButtonClass}>
              {t("profile.goToLogin")}
            </Link>
          }
        />
      ) : null}

      {status === "authenticated" && !isFinished ? (
        <StateCard
          eyebrow={t("checkin.emptyEyebrow")}
          title={t("checkin.title")}
          body={t("checkin.notFinished")}
          aside={<ArchivePill>{match.status}</ArchivePill>}
        />
      ) : null}

      {status === "authenticated" && isFinished && (roster.length === 0 || tagOptions.length === 0) ? (
        <StateCard
          eyebrow={t("checkin.emptyEyebrow")}
          title={t("checkin.title")}
          body={t("checkin.missingFormData")}
        />
      ) : null}

      {status === "authenticated" && isFinished && roster.length > 0 && tagOptions.length > 0 ? (
        <>
          {!editing && !loadingRecord && myCheckIn ? (
            <SavedCheckInSummary checkIn={myCheckIn} locale={locale} onEdit={openEdit} />
          ) : null}

          {!editing && !loadingRecord && !myCheckIn ? (
            <StateCard
              eyebrow={t("checkin.emptyEyebrow")}
              title={t("checkin.entryTitle")}
              body={t("checkin.empty")}
              cta={
                <button type="button" onClick={openCreate} className={primaryButtonClass}>
                  {t("checkin.create")}
                </button>
              }
            />
          ) : null}

          {loadingRecord ? (
            <StateMessage
              eyebrow={t("checkin.editingEyebrow")}
              title={t("checkin.loadingRecord")}
              body={t("checkin.entryBody")}
            />
          ) : null}

          {recordError ? <p className="text-sm text-red-600">{recordError}</p> : null}

          {editing ? (
            <form onSubmit={handleSubmit} className="space-y-6 rounded-[1.5rem] border border-[var(--fw-line)] bg-white/88 p-6">
              <div className="space-y-3">
                <p className="match-eyebrow">{t("checkin.editingEyebrow")}</p>
                <h3 className="text-2xl font-semibold tracking-tight text-[var(--fw-ink)]">
                  {t("checkin.entryTitle")}
                </h3>
                <p className="max-w-3xl text-sm leading-6 text-[var(--fw-muted)]">{t("checkin.entryBody")}</p>
              </div>

              <div className="grid gap-4 md:grid-cols-3">
                <RatingField
                  label={t("checkin.matchRating")}
                  value={formState.matchRating}
                  error={formErrors.matchRating}
                  onChange={(value) => setFormState((current) => ({ ...current, matchRating: value }))}
                />
                <label className="block text-sm">
                  <span className="mb-2 block font-medium text-[var(--fw-ink-soft)]">{t("checkin.supporterSide")}</span>
                  <select
                    value={formState.supporterSide}
                    onChange={(event) =>
                      setFormState((current) => ({
                        ...current,
                        supporterSide: event.target.value as CheckInFormState["supporterSide"],
                      }))
                    }
                    className={inputClass}
                  >
                    <option value="HOME">{match.homeTeam.name}</option>
                    <option value="AWAY">{match.awayTeam.name}</option>
                    <option value="NEUTRAL">{t("enum.supporterSide.neutral")}</option>
                  </select>
                  {formErrors.supporterSide ? <ErrorText text={formErrors.supporterSide} /> : null}
                </label>
                <label className="block text-sm">
                  <span className="mb-2 block font-medium text-[var(--fw-ink-soft)]">{t("checkin.watchedType")}</span>
                  <select
                    value={formState.watchedType}
                    onChange={(event) =>
                      setFormState((current) => ({
                        ...current,
                        watchedType: event.target.value as CheckInFormState["watchedType"],
                      }))
                    }
                    className={inputClass}
                  >
                    <option value="FULL">{t("enum.watchedType.full")}</option>
                    <option value="PARTIAL">{t("enum.watchedType.partial")}</option>
                    <option value="HIGHLIGHTS">{t("enum.watchedType.highlights")}</option>
                  </select>
                  {formErrors.watchedType ? <ErrorText text={formErrors.watchedType} /> : null}
                </label>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <RatingField
                  label={`${match.homeTeam.name} ${t("checkin.rating")}`}
                  value={formState.homeTeamRating}
                  error={formErrors.homeTeamRating}
                  onChange={(value) => setFormState((current) => ({ ...current, homeTeamRating: value }))}
                />
                <RatingField
                  label={`${match.awayTeam.name} ${t("checkin.rating")}`}
                  value={formState.awayTeamRating}
                  error={formErrors.awayTeamRating}
                  onChange={(value) => setFormState((current) => ({ ...current, awayTeamRating: value }))}
                />
              </div>

              <fieldset className="space-y-3">
                <legend className="text-sm font-medium text-[var(--fw-ink-soft)]">{t("checkin.tags")}</legend>
                <div className="flex flex-wrap gap-2">
                  {tagOptions.map((tag) => {
                    const selected = formState.tags.includes(tag.id);
                    return (
                      <button
                        key={tag.id}
                        type="button"
                        onClick={() => toggleTag(tag.id)}
                        aria-pressed={selected}
                        className={`rounded-full border px-3 py-2 text-sm transition-colors ${
                          selected
                            ? "border-[var(--fw-field-900)] bg-[var(--fw-field-900)] text-white"
                            : "border-[var(--fw-line)] bg-[var(--fw-paper-strong)] text-[var(--fw-ink-soft)] hover:bg-[var(--fw-field-100)]"
                        }`}
                      >
                        {tag.name}
                      </button>
                    );
                  })}
                </div>
              </fieldset>

              <label className="block text-sm">
                <span className="mb-2 block font-medium text-[var(--fw-ink-soft)]">{t("checkin.shortReview")}</span>
                <textarea
                  value={formState.shortReview}
                  onChange={(event) => setFormState((current) => ({ ...current, shortReview: event.target.value }))}
                  rows={4}
                  maxLength={280}
                  className={`${inputClass} resize-none`}
                  placeholder={t("checkin.shortReviewPlaceholder")}
                />
                <span className="mt-2 block text-xs text-[var(--fw-muted)]">{formState.shortReview.length}/280</span>
                {formErrors.shortReview ? <ErrorText text={formErrors.shortReview} /> : null}
              </label>

              <div className="space-y-3">
                <div className="flex items-center justify-between gap-4">
                  <h3 className="text-sm font-medium text-[var(--fw-ink-soft)]">{t("checkin.playerRatings")}</h3>
                  <button
                    type="button"
                    onClick={() =>
                      setFormState((current) => ({
                        ...current,
                        playerRatings: [...current.playerRatings, { ...DEFAULT_PLAYER_RATING }],
                      }))
                    }
                    className="text-sm font-medium text-[var(--fw-field-700)] underline-offset-4 hover:underline"
                  >
                    {t("checkin.addPlayerRating")}
                  </button>
                </div>

                {formState.playerRatings.length === 0 ? (
                  <p className="text-sm text-[var(--fw-muted)]">
                    {locale === "zh"
                      ? "你可以为本场名单中的任意球员评分。"
                      : "Rate as many players from this match roster as you want."}
                  </p>
                ) : (
                  <div className="space-y-4">
                    {formState.playerRatings.map((playerRating, index) => (
                      <div
                        key={`${index}-${playerRating.playerId}`}
                        className="rounded-[1.25rem] border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] p-4"
                      >
                        <div className="grid gap-4 md:grid-cols-[1.5fr_0.6fr]">
                          <label className="block text-sm">
                            <span className="mb-2 block font-medium text-[var(--fw-ink-soft)]">{t("checkin.player")}</span>
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
                              className={inputClass}
                            >
                              <option value="">{locale === "zh" ? "选择球员" : "Select a player"}</option>
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
                            label={t("checkin.rating")}
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
                          <span className="mb-2 block font-medium text-[var(--fw-ink-soft)]">{t("checkin.note")}</span>
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
                            className={inputClass}
                            placeholder={t("checkin.notePlaceholder")}
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
                          className="mt-4 text-sm text-red-700 underline-offset-4 hover:underline"
                        >
                          {t("checkin.remove")}
                        </button>
                      </div>
                    ))}
                  </div>
                )}
                {formErrors.playerRatings ? <ErrorText text={formErrors.playerRatings} /> : null}
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <label className="block text-sm">
                  <span className="mb-2 block font-medium text-[var(--fw-ink-soft)]">{t("checkin.watchedAt")}</span>
                  <input
                    type="datetime-local"
                    value={formState.watchedAt}
                    onChange={(event) => setFormState((current) => ({ ...current, watchedAt: event.target.value }))}
                    className={inputClass}
                  />
                  {formErrors.watchedAt ? <ErrorText text={formErrors.watchedAt} /> : null}
                </label>
              </div>

              {submitError ? <ErrorText text={submitError} /> : null}

              <div className="flex flex-wrap gap-3">
                <button type="submit" disabled={submitting} className={primaryButtonClass}>
                  {submitting ? t("checkin.saving") : t("checkin.save")}
                </button>
                <button type="button" onClick={cancelEdit} className={secondaryButtonClass}>
                  {t("checkin.cancel")}
                </button>
              </div>
            </form>
          ) : null}
        </>
      ) : null}
    </SectionShell>
  );
}

function SavedCheckInSummary({
  checkIn,
  locale,
  onEdit,
}: {
  checkIn: CheckInDetail;
  locale: "en" | "zh";
  onEdit: () => void;
}) {
  return (
    <div className="grid gap-5 xl:grid-cols-[0.9fr_1.1fr]">
      <StateCard
        eyebrow={translate(locale, "checkin.archived")}
        title={translate(locale, "checkin.savedTitle")}
        body={translate(locale, "checkin.savedBody")}
        aside={<ArchivePill>{translate(locale, "checkin.savedAt", { value: formatDateTime(checkIn.updatedAt, locale) })}</ArchivePill>}
        cta={
          <button type="button" onClick={onEdit} className={secondaryButtonClass}>
            {translate(locale, "checkin.edit")}
          </button>
        }
      >
        <dl className="grid gap-3 text-sm text-[var(--fw-ink-soft)]">
          <SummaryRow
            label={translate(locale, "checkin.watchedType")}
            value={getWatchedTypeLabel(checkIn.watchedType, locale)}
          />
          <SummaryRow
            label={translate(locale, "checkin.supporterSide")}
            value={getSupporterSideLabel(checkIn.supporterSide, locale)}
          />
          <SummaryRow label={translate(locale, "checkin.matchRating")} value={String(checkIn.matchRating)} />
        </dl>

        {checkIn.tags.length > 0 ? (
          <div className="mt-5 flex flex-wrap gap-2">
            {checkIn.tags.map((tag) => (
              <ArchivePill key={tag.id}>{tag.name}</ArchivePill>
            ))}
          </div>
        ) : null}
      </StateCard>

      <div className="rounded-[1.5rem] border border-[var(--fw-line)] bg-white/88 p-6">
        <div className="flex items-center justify-between gap-4">
          <div>
            <p className="match-eyebrow">{translate(locale, "checkin.playerNotes")}</p>
            <h3 className="mt-2 text-xl font-semibold tracking-tight text-[var(--fw-ink)]">
              {translate(locale, "profile.matchArchiveLabel")}
            </h3>
          </div>
          <ArchivePill>{translate(locale, "profile.savedRecord")}</ArchivePill>
        </div>

        {checkIn.playerRatings.length === 0 ? (
          <p className="mt-4 text-sm text-[var(--fw-muted)]">
            {translate(locale, "checkin.noSavedPlayerRatings")}
          </p>
        ) : (
          <div className="mt-5 space-y-3">
            {checkIn.playerRatings.map((rating) => (
              <div
                key={rating.id}
                className="rounded-[1.2rem] border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] p-4"
              >
                <div className="flex items-center justify-between gap-4">
                  <div className="flex items-center gap-3">
                    <TeamCrest team={rating.player.team} size="sm" />
                    <div>
                      <p className="font-medium text-[var(--fw-ink)]">{rating.player.name}</p>
                      <p className="text-sm text-[var(--fw-muted)]">{rating.player.team.name}</p>
                    </div>
                  </div>
                  <ArchivePill>
                    {translate(locale, "checkin.ratingValue", { value: rating.rating })}
                  </ArchivePill>
                </div>
                {rating.note ? <p className="mt-3 text-sm leading-6 text-[var(--fw-ink-soft)]">{rating.note}</p> : null}
              </div>
            ))}
          </div>
        )}

        {checkIn.shortReview ? (
          <div className="mt-5 rounded-[1.2rem] border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] p-4 text-sm leading-6 text-[var(--fw-ink-soft)]">
            {checkIn.shortReview}
          </div>
        ) : null}
      </div>
    </div>
  );
}

function StateCard({
  eyebrow,
  title,
  body,
  aside,
  cta,
  children,
}: {
  eyebrow: string;
  title: string;
  body: string;
  aside?: ReactNode;
  cta?: ReactNode;
  children?: ReactNode;
}) {
  return (
    <div className="rounded-[1.5rem] border border-[var(--fw-line)] bg-white/88 p-6">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="space-y-3">
          <p className="match-eyebrow">{eyebrow}</p>
          <div className="space-y-2">
            <h3 className="text-2xl font-semibold tracking-tight text-[var(--fw-ink)]">{title}</h3>
            <p className="max-w-3xl text-sm leading-6 text-[var(--fw-muted)]">{body}</p>
          </div>
        </div>
        {aside}
      </div>
      {children ? <div className="mt-6">{children}</div> : null}
      {cta ? <div className="mt-6">{cta}</div> : null}
    </div>
  );
}

function StateMessage({ eyebrow, title, body }: { eyebrow: string; title: string; body: string }) {
  return <StateCard eyebrow={eyebrow} title={title} body={body} />;
}

function SummaryRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-4">
      <dt>{label}</dt>
      <dd className="font-medium text-[var(--fw-ink)]">{value}</dd>
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
      <span className="mb-2 block font-medium text-[var(--fw-ink-soft)]">{label}</span>
      <input
        type="number"
        min={1}
        max={10}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        className={inputClass}
      />
      {error ? <ErrorText text={error} /> : null}
    </label>
  );
}

function ErrorText({ text }: { text: string }) {
  return <span className="mt-2 block text-sm text-red-700">{text}</span>;
}

const inputClass =
  "w-full rounded-[1rem] border border-[var(--fw-line)] bg-white px-4 py-3 text-[var(--fw-ink)] outline-none transition focus:border-[var(--fw-field-500)] focus:ring-2 focus:ring-[rgba(110,143,114,0.18)]";

const primaryButtonClass =
  "inline-flex items-center justify-center rounded-full bg-[var(--fw-field-900)] px-5 py-3 text-sm font-medium text-white transition-colors hover:bg-[var(--fw-field-700)] disabled:opacity-50";

const secondaryButtonClass =
  "inline-flex items-center justify-center rounded-full border border-[var(--fw-line)] bg-white px-5 py-3 text-sm font-medium text-[var(--fw-ink-soft)] transition-colors hover:bg-[var(--fw-field-100)]";
