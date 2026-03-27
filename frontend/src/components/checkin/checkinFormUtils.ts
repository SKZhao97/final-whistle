import type {
  CheckInDetail,
  CheckInPlayerRatingInput,
  PlayerSummary,
  UpsertCheckInRequest,
} from "../../types/api";
import type { Locale } from "../../lib/i18n/config";
import { translate } from "../../lib/i18n/core";

export type CheckInFormState = {
  watchedType: UpsertCheckInRequest["watchedType"];
  supporterSide: UpsertCheckInRequest["supporterSide"];
  matchRating: string;
  homeTeamRating: string;
  awayTeamRating: string;
  shortReview: string;
  watchedAt: string;
  tags: number[];
  playerRatings: Array<{
    playerId: string;
    rating: string;
    note: string;
  }>;
};

export type CheckInFormErrors = {
  form?: string;
  watchedType?: string;
  supporterSide?: string;
  matchRating?: string;
  homeTeamRating?: string;
  awayTeamRating?: string;
  shortReview?: string;
  watchedAt?: string;
  playerRatings?: string;
};

export const DEFAULT_PLAYER_RATING = { playerId: "", rating: "", note: "" };
const MAX_SHORT_REVIEW_LENGTH = 280;
const MAX_PLAYER_NOTE_LENGTH = 80;

export function createDefaultFormState(): CheckInFormState {
  return {
    watchedType: "FULL",
    supporterSide: "NEUTRAL",
    matchRating: "8",
    homeTeamRating: "8",
    awayTeamRating: "8",
    shortReview: "",
    watchedAt: toDatetimeLocal(new Date()),
    tags: [],
    playerRatings: [],
  };
}

export function createFormStateFromCheckIn(checkIn: CheckInDetail): CheckInFormState {
  return {
    watchedType: checkIn.watchedType,
    supporterSide: checkIn.supporterSide,
    matchRating: String(checkIn.matchRating),
    homeTeamRating: String(checkIn.homeTeamRating),
    awayTeamRating: String(checkIn.awayTeamRating),
    shortReview: checkIn.shortReview ?? "",
    watchedAt: toDatetimeLocal(new Date(checkIn.watchedAt)),
    tags: checkIn.tags.map((tag) => tag.id),
    playerRatings: checkIn.playerRatings.map((rating) => ({
      playerId: String(rating.player.id),
      rating: String(rating.rating),
      note: rating.note ?? "",
    })),
  };
}

export function validateFormState(
  formState: CheckInFormState,
  roster: PlayerSummary[],
  locale: Locale = "en",
): CheckInFormErrors {
  const errors: CheckInFormErrors = {};
  const rosterIds = new Set(roster.map((player) => player.id));

  if (!formState.watchedType) {
    errors.watchedType = translate(locale, "validation.watchedTypeRequired");
  }
  if (!formState.supporterSide) {
    errors.supporterSide = translate(locale, "validation.supporterSideRequired");
  }

  validateRatingString(formState.matchRating, translate(locale, "checkin.matchRating"), errors, "matchRating", locale);
  validateRatingString(formState.homeTeamRating, translate(locale, "checkin.homeTeamRating"), errors, "homeTeamRating", locale);
  validateRatingString(formState.awayTeamRating, translate(locale, "checkin.awayTeamRating"), errors, "awayTeamRating", locale);

  if (formState.shortReview.length > MAX_SHORT_REVIEW_LENGTH) {
    errors.shortReview = translate(locale, "validation.shortReviewLength");
  }
  if (!formState.watchedAt) {
    errors.watchedAt = translate(locale, "validation.watchedAtRequired");
  } else if (Number.isNaN(new Date(formState.watchedAt).getTime())) {
    errors.watchedAt = translate(locale, "validation.watchedAtInvalid");
  }

  const selectedPlayers = new Set<number>();
  for (const entry of formState.playerRatings) {
    const playerRatingError = validatePlayerRatingEntry(entry, rosterIds, selectedPlayers, locale);
    if (playerRatingError) {
      errors.playerRatings = playerRatingError;
      break;
    }
  }

  return errors;
}

export function buildPayload(formState: CheckInFormState): UpsertCheckInRequest {
  const matchRating = parseNumericField(formState.matchRating, "matchRating");
  const homeTeamRating = parseNumericField(formState.homeTeamRating, "homeTeamRating");
  const awayTeamRating = parseNumericField(formState.awayTeamRating, "awayTeamRating");
  const watchedAt = new Date(formState.watchedAt);
  if (Number.isNaN(watchedAt.getTime())) {
    throw new Error("Invalid watchedAt value");
  }

  return {
    watchedType: formState.watchedType,
    supporterSide: formState.supporterSide,
    matchRating,
    homeTeamRating,
    awayTeamRating,
    shortReview: formState.shortReview.trim() ? formState.shortReview.trim() : undefined,
    watchedAt: watchedAt.toISOString(),
    tags: formState.tags,
    playerRatings: formState.playerRatings.map(
      (entry): CheckInPlayerRatingInput => ({
        playerId: parseNumericField(entry.playerId, "playerId"),
        rating: parseNumericField(entry.rating, "playerRating"),
        note: entry.note.trim() ? entry.note.trim() : undefined,
      }),
    ),
  };
}

export function toDatetimeLocal(date: Date) {
  const year = date.getFullYear();
  const month = `${date.getMonth() + 1}`.padStart(2, "0");
  const day = `${date.getDate()}`.padStart(2, "0");
  const hours = `${date.getHours()}`.padStart(2, "0");
  const minutes = `${date.getMinutes()}`.padStart(2, "0");
  return `${year}-${month}-${day}T${hours}:${minutes}`;
}

function validateRatingString(
  value: string,
  label: string,
  errors: CheckInFormErrors,
  key: "matchRating" | "homeTeamRating" | "awayTeamRating",
  locale: Locale,
) {
  const numericValue = parseRequiredNumber(value);
  if (numericValue === null || numericValue < 1 || numericValue > 10) {
    errors[key] = translate(locale, "validation.ratingRange", { label });
  }
}

function parseNumericField(value: string, field: string) {
  const parsed = parseRequiredNumber(value);
  if (parsed === null) {
    throw new Error(`Invalid ${field} value`);
  }
  return parsed;
}

function validatePlayerRatingEntry(
  entry: CheckInFormState["playerRatings"][number],
  rosterIds: Set<number>,
  selectedPlayers: Set<number>,
  locale: Locale,
) {
  const playerId = parseRequiredNumber(entry.playerId);
  if (playerId === null) {
    return translate(locale, "validation.playerRequired");
  }
  if (!rosterIds.has(playerId)) {
    return translate(locale, "validation.playerRoster");
  }
  if (selectedPlayers.has(playerId)) {
    return translate(locale, "validation.playerDuplicate");
  }
  selectedPlayers.add(playerId);

  const rating = parseRequiredNumber(entry.rating);
  if (rating === null || rating < 1 || rating > 10) {
    return translate(locale, "validation.playerRatingRange");
  }
  if (entry.note.length > MAX_PLAYER_NOTE_LENGTH) {
    return translate(locale, "validation.playerNoteLength");
  }

  return undefined;
}

function parseRequiredNumber(value: string) {
  const trimmedValue = value.trim();
  if (trimmedValue === "") {
    return null;
  }

  const parsed = Number(trimmedValue);
  if (Number.isNaN(parsed)) {
    return null;
  }

  return parsed;
}
