import test from "node:test";
import assert from "node:assert/strict";

import {
  buildPayload,
  createFormStateFromCheckIn,
  validateFormState,
} from "../src/components/checkin/checkinFormUtils.ts";
import type { CheckInDetail, PlayerSummary } from "../src/types/api.ts";

const roster: PlayerSummary[] = [
  {
    id: 1,
    name: "Kevin De Bruyne",
    slug: "kevin-de-bruyne",
    position: "Midfielder",
    team: { id: 1, name: "Manchester City", slug: "manchester-city", shortName: "MCI" },
  },
  {
    id: 2,
    name: "Erling Haaland",
    slug: "erling-haaland",
    position: "Forward",
    team: { id: 1, name: "Manchester City", slug: "manchester-city", shortName: "MCI" },
  },
];

test("validateFormState rejects duplicate players", () => {
  const errors = validateFormState(
    {
      watchedType: "FULL",
      supporterSide: "NEUTRAL",
      matchRating: "8",
      homeTeamRating: "8",
      awayTeamRating: "7",
      shortReview: "",
      watchedAt: "2026-03-26T14:00",
      tags: [],
      playerRatings: [
        { playerId: "1", rating: "8", note: "" },
        { playerId: "1", rating: "7", note: "" },
      ],
    },
    roster,
  );

  assert.equal(errors.playerRatings, "A player can only be rated once per record.");
});

test("validateFormState rejects missing player rating values with clear errors", () => {
  const errors = validateFormState(
    {
      watchedType: "FULL",
      supporterSide: "NEUTRAL",
      matchRating: "8",
      homeTeamRating: "8",
      awayTeamRating: "7",
      shortReview: "",
      watchedAt: "2026-03-26T14:00",
      tags: [],
      playerRatings: [{ playerId: "1", rating: "", note: "" }],
    },
    roster,
  );

  assert.equal(errors.playerRatings, "Each player rating must be between 1 and 10.");
});

test("validateFormState rejects players outside the match roster", () => {
  const errors = validateFormState(
    {
      watchedType: "FULL",
      supporterSide: "NEUTRAL",
      matchRating: "8",
      homeTeamRating: "8",
      awayTeamRating: "7",
      shortReview: "",
      watchedAt: "2026-03-26T14:00",
      tags: [],
      playerRatings: [{ playerId: "99", rating: "8", note: "" }],
    },
    roster,
  );

  assert.equal(errors.playerRatings, "Player ratings must use players from this match roster.");
});

test("validateFormState rejects invalid watchedAt format", () => {
  const errors = validateFormState(
    {
      watchedType: "FULL",
      supporterSide: "NEUTRAL",
      matchRating: "8",
      homeTeamRating: "8",
      awayTeamRating: "7",
      shortReview: "",
      watchedAt: "not-a-date",
      tags: [],
      playerRatings: [],
    },
    roster,
  );

  assert.equal(errors.watchedAt, "Watched at must be a valid date and time.");
});

test("buildPayload trims optional text and preserves all player ratings", () => {
  const payload = buildPayload({
    watchedType: "FULL",
    supporterSide: "HOME",
    matchRating: "9",
    homeTeamRating: "8",
    awayTeamRating: "7",
    shortReview: "  Sharp display  ",
    watchedAt: "2026-03-26T14:00",
    tags: [1, 4],
    playerRatings: [
      { playerId: "1", rating: "8", note: "  Engine  " },
      { playerId: "2", rating: "9", note: "" },
    ],
  });

  assert.equal(payload.shortReview, "Sharp display");
  assert.equal(payload.playerRatings.length, 2);
  assert.equal(payload.playerRatings[0].note, "Engine");
  assert.equal(payload.playerRatings[1].note, undefined);
});

test("buildPayload throws on invalid numeric values", () => {
  assert.throws(
    () =>
      buildPayload({
        watchedType: "FULL",
        supporterSide: "HOME",
        matchRating: "oops",
        homeTeamRating: "8",
        awayTeamRating: "7",
        shortReview: "",
        watchedAt: "2026-03-26T14:00",
        tags: [],
        playerRatings: [],
      }),
    /Invalid matchRating value/,
  );
});

test("createFormStateFromCheckIn maps existing check-in data into editable strings", () => {
  const checkIn: CheckInDetail = {
    id: 10,
    matchId: 1,
    watchedType: "PARTIAL",
    supporterSide: "AWAY",
    matchRating: 7,
    homeTeamRating: 6,
    awayTeamRating: 8,
    shortReview: "Late swing",
    watchedAt: "2026-03-26T14:00:00.000Z",
    tags: [{ id: 4, name: "经典", slug: "classic" }],
    playerRatings: [
      {
        id: 1,
        player: roster[0],
        rating: 8,
        note: "Creator",
      },
    ],
    createdAt: "2026-03-26T14:00:00.000Z",
    updatedAt: "2026-03-26T14:10:00.000Z",
  };

  const formState = createFormStateFromCheckIn(checkIn);
  assert.equal(formState.watchedType, "PARTIAL");
  assert.equal(formState.playerRatings[0].playerId, "1");
  assert.equal(formState.playerRatings[0].note, "Creator");
  assert.deepEqual(formState.tags, [4]);
});
