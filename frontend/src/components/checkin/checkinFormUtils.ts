import type {
  CheckInDetail,
  CheckInPlayerRatingInput,
  PlayerSummary,
  UpsertCheckInRequest,
} from "../../types/api";

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
): CheckInFormErrors {
  const errors: CheckInFormErrors = {};
  const rosterIds = new Set(roster.map((player) => player.id));

  if (!formState.watchedType) {
    errors.watchedType = "Watched type is required.";
  }
  if (!formState.supporterSide) {
    errors.supporterSide = "Supporter side is required.";
  }

  validateRatingString(formState.matchRating, "Match rating", errors, "matchRating");
  validateRatingString(formState.homeTeamRating, "Home team rating", errors, "homeTeamRating");
  validateRatingString(formState.awayTeamRating, "Away team rating", errors, "awayTeamRating");

  if (formState.shortReview.length > 280) {
    errors.shortReview = "Short review must be 280 characters or fewer.";
  }
  if (!formState.watchedAt) {
    errors.watchedAt = "Watched at is required.";
  } else if (Number.isNaN(new Date(formState.watchedAt).getTime())) {
    errors.watchedAt = "Watched at must be a valid date and time.";
  }

  const selectedPlayers = new Set<number>();
  for (const entry of formState.playerRatings) {
    if (!entry.playerId) {
      errors.playerRatings = "Each player rating needs a selected player.";
      break;
    }

    const playerId = Number(entry.playerId);
    if (!rosterIds.has(playerId)) {
      errors.playerRatings = "Player ratings must use players from this match roster.";
      break;
    }
    if (selectedPlayers.has(playerId)) {
      errors.playerRatings = "A player can only be rated once per record.";
      break;
    }
    selectedPlayers.add(playerId);

    const rating = Number(entry.rating);
    if (!entry.rating || Number.isNaN(rating) || rating < 1 || rating > 10) {
      errors.playerRatings = "Each player rating must be between 1 and 10.";
      break;
    }
    if (entry.note.length > 80) {
      errors.playerRatings = "Each player note must be 80 characters or fewer.";
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
) {
  const numericValue = Number(value);
  if (!value || Number.isNaN(numericValue) || numericValue < 1 || numericValue > 10) {
    errors[key] = `${label} must be between 1 and 10.`;
  }
}

function parseNumericField(value: string, field: string) {
  const parsed = Number(value);
  if (!value || Number.isNaN(parsed)) {
    throw new Error(`Invalid ${field} value`);
  }
  return parsed
}
